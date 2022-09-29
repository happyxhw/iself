package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cacher struct {
	rdb *redis.Client
}

func NewCacher(rdb *redis.Client) *Cacher {
	return &Cacher{
		rdb: rdb,
	}
}

func (cr *Cacher) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return cr.rdb.Get(ctx, key).Bytes()
}

func (cr *Cacher) GetString(ctx context.Context, key string) (string, error) {
	return cr.rdb.Get(ctx, key).Result()
}

func (cr *Cacher) GetObject(ctx context.Context, key string, out any) error {
	data, err := cr.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

func (cr *Cacher) Set(ctx context.Context, key string, val any, ex time.Duration) error {
	return cr.rdb.Set(ctx, key, val, ex).Err()
}

func (cr *Cacher) SetObject(ctx context.Context, key string, val any, ex time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return cr.rdb.Set(ctx, key, data, ex).Err()
}

func (cr *Cacher) SetNX(ctx context.Context, key string, val any, ex time.Duration) (bool, error) {
	return cr.rdb.SetNX(ctx, key, val, ex).Result()
}

func (cr *Cacher) SetNXObject(ctx context.Context, key string, val any, ex time.Duration) (bool, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return false, err
	}
	return cr.rdb.SetNX(ctx, key, data, ex).Result()
}

func (cr *Cacher) Del(ctx context.Context, key string) (int64, error) {
	return cr.rdb.Del(ctx, key).Result()
}
