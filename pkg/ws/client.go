package ws

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"git.happyxhw.cn/happyxhw/iself/pkg/log"
)

const (
	defaultWriteWait = 60 * time.Second

	defaultPongTimeout = 60 * time.Second
	defaultPingPeriod  = 30 * time.Second

	defaultReadLimit = 1024
)

// Client websocket client, held by stream
type Client struct {
	conn    *websocket.Conn // websocket conn
	stopCh  chan struct{}   // signal to stop client
	onClose func(*Client)   // do func when conn close

	stream *Stream
	once   sync.Once
	ch     chan *Msg
	errCh  chan error

	id     string // client id
	userID int64  // client user id
	token  string // user token
}

// NewClient return client instance
func NewClient(conn *websocket.Conn, stream *Stream, id, token string, userID int64, readFunc func([]byte)) *Client {
	cli := Client{
		conn:   conn,
		stream: stream,

		stopCh: make(chan struct{}, 1),
		ch:     make(chan *Msg),
		errCh:  make(chan error),

		id:     id,
		userID: userID,
		token:  token,
	}
	go cli.startReading(readFunc)
	go cli.startPingHandler()
	return &cli
}

func (c *Client) SetCloseFn(onClose func(*Client)) {
	c.onClose = onClose
}

// Close client
func (c *Client) Close() {
	c.close()
}

// close conn
func (c *Client) close() {
	c.once.Do(func() {
		c.stopCh <- struct{}{}
		_ = c.conn.Close()
		if c.onClose != nil {
			c.onClose(c)
		}
		c.stream.Remove(c)
		log.Info("client closed", zap.String("id", c.id), zap.Int64("user_id", c.userID))
	})
}

// StartReading starts listening on the Client connection.
// As we do not need anything from the Client,
// we ignore incoming messages. Leaves the loop on errors.
func (c *Client) startReading(readFunc func([]byte)) {
	defer c.close()
	c.conn.SetReadLimit(defaultReadLimit)
	_ = c.conn.SetReadDeadline(time.Now().Add(defaultPongTimeout))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(defaultPongTimeout))
		return nil
	})
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "close") {
				log.Info("ws close", zap.String("id", c.id), zap.Int64("user_id", c.userID), zap.Error(err))
			} else {
				log.Error("ws read", zap.String("id", c.id), zap.Int64("user_id", c.userID), zap.Error(err))
			}
			break
		}
		readFunc(msg)
	}
}

// ping loop, quit on error or stop signal
func (c *Client) startPingHandler() {
	pingTicker := time.NewTicker(defaultPingPeriod)
	defer func() {
		c.close()
		pingTicker.Stop()
	}()
	for {
		select {
		case <-pingTicker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(defaultWriteWait))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				if strings.Contains(err.Error(), "close") {
					log.Info("ping", zap.String("id", c.id), zap.Int64("user_id", c.userID), zap.Error(err))
				} else {
					log.Error("ping", zap.String("id", c.id), zap.Int64("user_id", c.userID), zap.Error(err))
				}
				return
			}
		case msg, ok := <-c.ch:
			if !ok {
				c.errCh <- nil
				return
			}
			err := c.conn.WriteJSON(msg)
			c.errCh <- err
		case <-c.stopCh:
			log.Info("stop ws client", zap.String("id", c.id), zap.Int64("user_id", c.userID))
			return
		}
	}
}

func (c *Client) send(msg *Msg) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) //nolint:gomnd
	defer cancel()
	select {
	case c.ch <- msg:
		err = <-c.errCh
		if err != nil {
			log.Error("send", zap.String("id", c.id), zap.Int64("user_id", c.userID), zap.Error(err))
		}
	case <-ctx.Done():
		log.Error("send timeout", zap.String("id", c.id), zap.Int64("user_id", c.userID))
		close(c.ch)
		err = errors.New("timeout")
	}

	return err
}
