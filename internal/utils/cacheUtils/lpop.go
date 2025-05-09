package cacheUtils

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

func LPop[T any](ctx context.Context, pool *redis.Pool, key string) (*T, error) {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := conn.Do("LPOP", key)
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if reply == nil {
		return nil, ErrNoData
	}

	v, err := ParseRedisReply[T](reply, nil)
	if err != nil {
		return nil, err
	}

	return &v, nil
}
