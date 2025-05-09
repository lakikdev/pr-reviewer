package cacheUtils

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

// function to set redis list with new data
func RPush[T any](ctx context.Context, pool *redis.Pool, key string, expireIn int64, data []T) error {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	args := []interface{}{key}
	for _, item := range data {

		args = append(args, item)
	}

	//insert all the data
	reply, err := redis.Int(conn.Do("RPUSH", args...))
	if err != nil || reply == 0 {
		return err
	}

	//call EXPIRE to set expire time
	if expireIn > 0 {
		reply, err := redis.Int(conn.Do("EXPIRE", key, expireIn))
		if err != nil || reply == 0 {
			return err
		}
	}

	return nil
}
