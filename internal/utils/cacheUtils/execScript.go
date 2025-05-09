package cacheUtils

import (
	"context"
	"encoding/json"

	"github.com/gomodule/redigo/redis"
)

func ExecScript(ctx context.Context, pool *redis.Pool, scriptStr string, keys []string, args []interface{}, dest interface{}) error {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	script := redis.NewScript(len(keys), scriptStr)

	argsList := []interface{}{}
	for _, key := range keys {
		argsList = append(argsList, key)
	}
	argsList = append(argsList, args...)

	rawReplay, err := script.Do(conn, argsList...)
	if err != nil {
		return err
	}

	if dest != nil {
		reply, err := redis.String(rawReplay, err)
		if err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(reply), &dest); err != nil {
			return err
		}
	}

	return nil

}
