package cacheUtils

import (
	"context"
	"encoding/json"

	"github.com/gomodule/redigo/redis"
)

func Set(ctx context.Context, pool *redis.Pool, key string, data interface{}) error {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	reply, err := redis.String(conn.Do("SET", key, dataJSON))
	if err != nil || reply != "OK" {
		return err
	}

	return nil
}
