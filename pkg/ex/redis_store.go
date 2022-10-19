package ex

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"

	"github.com/happyxhw/iself/pkg/util"
)

/*
此代码修改自：https://github.com/boj/redistore
源代码使用的是 redigo，为了统一这里修改为了 go-redis
*/

// 没有设置 maxAge 情况下默认的 maxAge
var sessionExpire = 86400 * 30

// session 长度
const sessionLen = 32

// SessionSerializer 序列化、反序列化接口
type SessionSerializer interface {
	Deserialize(d []byte, ss *sessions.Session) error
	Serialize(ss *sessions.Session) ([]byte, error)
}

var _ SessionSerializer = (*JSONSerializer)(nil)

// JSONSerializer json 序列化、反序列化
type JSONSerializer struct{}

// Serialize 序列化
func (s JSONSerializer) Serialize(ss *sessions.Session) ([]byte, error) {
	m := make(map[string]interface{}, len(ss.Values))
	for k, v := range ss.Values {
		ks, ok := k.(string)
		if !ok {
			err := fmt.Errorf("non-string key value, cannot serialize session to json: %v", k)
			fmt.Printf("redistore.JSONSerializer.serialize() Error: %v", err)
			return nil, err
		}
		m[ks] = v
	}
	return json.Marshal(m)
}

// Deserialize 反序列化
func (s JSONSerializer) Deserialize(d []byte, ss *sessions.Session) error {
	m := make(map[string]interface{})
	err := json.Unmarshal(d, &m)
	if err != nil {
		fmt.Printf("redistore.JSONSerializer.deserialize() Error: %v", err)
		return err
	}
	for k, v := range m {
		ss.Values[k] = v
	}
	return nil
}

// GobSerializer gob 序列化、反序列化
type GobSerializer struct{}

// Serialize 序列化
func (s GobSerializer) Serialize(ss *sessions.Session) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(ss.Values)
	if err == nil {
		return buf.Bytes(), nil
	}
	return nil, err
}

// Deserialize 反序列化
func (s GobSerializer) Deserialize(d []byte, ss *sessions.Session) error {
	dec := gob.NewDecoder(bytes.NewBuffer(d))
	return dec.Decode(&ss.Values)
}

// RedisStore session redis 存储
type RedisStore struct {
	rdb        *redis.Client
	Codecs     []securecookie.Codec
	Options    *sessions.Options // default configuration
	maxLen     int
	keyPrefix  string
	serializer SessionSerializer
}

// NewRedisStore instantiates a RedisStore
func NewRedisStore(rdb *redis.Client, prefix string, keyPairs []byte) *RedisStore {
	rs := &RedisStore{
		rdb:    rdb,
		Codecs: securecookie.CodecsFromPairs(keyPairs),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: sessionExpire,
		},
		maxLen:     4096,
		keyPrefix:  prefix,
		serializer: JSONSerializer{},
	}
	return rs
}

// SetMaxLength sets RedisStore.maxLen if the `l` argument is greater or equal 0
// maxLen restricts the maximum length of new sessions to l.
// If l is 0 there is no limit to the size of a session, use with caution.
// The default for a new RedisStore is 4096. Redis allows for max.
// value sizes of up to 512MB (http://redis.io/topics/data-types)
// Default: 4096,
func (s *RedisStore) SetMaxLength(l int) {
	if l >= 0 {
		s.maxLen = l
	}
}

// SetKeyPrefix set the prefix
func (s *RedisStore) SetKeyPrefix(p string) {
	s.keyPrefix = p
}

// SetSerializer sets the serializer
func (s *RedisStore) SetSerializer(ss SessionSerializer) {
	s.serializer = ss
}

// SetMaxAge restricts the maximum age, in seconds, of the session record
// both in database and a browser. This is to change session storage configuration.
// If you want just to remove session use your session `s` object and change it's
// `Options.MaxAge` to -1, as specified in
//
//	http://godoc.org/github.com/gorilla/sessions#Options
//
// Default is the one provided by this package value - `sessionExpire`.
// Set it to 0 for no restriction.
// Because we use `MaxAge` also in SecureCookie crypting algorithm you should
// use this function to change `MaxAge` value.
func (s *RedisStore) SetMaxAge(v int) {
	var c *securecookie.SecureCookie
	var ok bool
	s.Options.MaxAge = v
	for i := range s.Codecs {
		if c, ok = s.Codecs[i].(*securecookie.SecureCookie); ok {
			c.MaxAge(v)
		} else {
			fmt.Printf("Can't change MaxAge on codec %v\n", s.Codecs[i])
		}
	}
}

// Get returns a session for the given name after adding it to the registry.
//
// See gorilla/sessions FilesystemStore.Get().
func (s *RedisStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// See gorilla/sessions FilesystemStore.New().
func (s *RedisStore) New(r *http.Request, name string) (*sessions.Session, error) {
	var (
		err error
		ok  bool
	)
	session := sessions.NewSession(s, name)
	// make a copy
	options := *s.Options
	session.Options = &options
	session.IsNew = true
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
		if err == nil {
			ok, err = s.load(session)
			session.IsNew = !(err == nil && ok) // not new if no error and data available
		}
	}
	return session, err
}

// Save adds a single session to the response.
func (s *RedisStore) Save(_ *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Marked for deletion.
	if session.Options.MaxAge < 0 {
		if err := s.delete(session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		// Build an alphanumeric key for the redis store.
		if session.ID == "" {
			session.ID = util.NanoID(sessionLen)
		}
		if err := s.save(session); err != nil {
			return err
		}
		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
		if err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	}
	return nil
}

// save stores the session in redis.
func (s *RedisStore) save(session *sessions.Session) error {
	b, err := s.serializer.Serialize(session)
	if err != nil {
		return err
	}
	if s.maxLen != 0 && len(b) > s.maxLen {
		return errors.New("SessionStore: the value to store is too big")
	}

	age := session.Options.MaxAge
	if age == 0 {
		age = 60 * 20
	}
	err = s.rdb.Set(context.TODO(), s.keyPrefix+session.ID, b, time.Duration(age)*time.Second).Err()
	return err
}

// load reads the session from redis.
// returns true if there is a session data in DB
func (s *RedisStore) load(session *sessions.Session) (bool, error) {
	data, err := s.rdb.Get(context.TODO(), s.keyPrefix+session.ID).Bytes()
	if err != nil && err != redis.Nil {
		return false, err
	}
	if len(data) == 0 {
		return false, nil // no data was associated with this key
	}

	return true, s.serializer.Deserialize(data, session)
}

// delete removes keys from redis if MaxAge<0
func (s *RedisStore) delete(session *sessions.Session) error {
	err := s.rdb.Del(context.TODO(), s.keyPrefix+session.ID).Err()
	return err
}

type store struct {
	*RedisStore
}

// NewStore 返回 redis store
func NewStore(rdb *redis.Client, prefix string, keyPairs []byte) *store {
	s := NewRedisStore(rdb, prefix, keyPairs)
	return &store{s}
}
