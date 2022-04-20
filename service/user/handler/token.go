package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"git.happyxhw.cn/happyxhw/iself/pkg/log"

	"git.happyxhw.cn/happyxhw/iself/pkg/em"
)

var (
	// ErrGetToken 获取token错误
	ErrGetToken = em.NewError(http.StatusBadRequest, 40201, "get oauth2 token")
)

// Token srv
type Token struct {
	RDB *redis.Client
}

// SaveToken 保存到redis
func (t *Token) SaveToken(ctx context.Context, token *oauth2.Token, source string, id int64) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("oauth2:%s:%d", source, id)
	if err := t.RDB.Set(context.TODO(), key, data, 0).Err(); err != nil {
		log.Error("save token to redis", zap.Error(err), log.Ctx(ctx))
		return err
	}
	return nil
}

// GetToken 获取 access token，自动刷新
func (t *Token) GetToken(ctx context.Context, source string, id int64, conf *oauth2.Config) (*oauth2.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	key := fmt.Sprintf("oauth2:%s:%d", source, id)
	result, err := t.RDB.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	err = json.Unmarshal([]byte(result), &token)
	if err != nil {
		return nil, err
	}
	newToken, err := conf.TokenSource(ctx, &token).Token()
	if err != nil {
		return nil, err
	}
	if newToken.AccessToken != token.AccessToken {
		_ = t.SaveToken(ctx, newToken, source, id)
	}

	return newToken, nil
}