package cacheUtils

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func Delete(ctx context.Context, pool *redis.Pool, key string) error {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	deletedAmount, err := redis.Int(conn.Do("DEL", key))

	if deletedAmount > 0 {
		fmt.Printf("Removed %d keys from cache\n", deletedAmount)
	}

	return err
}
