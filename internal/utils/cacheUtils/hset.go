package cacheUtils

import (
	"context"
	"pr-reviewer/internal/utils/cacheUtils/utils"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

func HSet(ctx context.Context, pool *redis.Pool, key string, expiresIn int64, data interface{}) error {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get redis connection")
	}
	defer conn.Close()

	args, err := utils.StructToArgs(data)
	if err != nil {
		return errors.Wrap(err, "failed to convert struct to args")
	}

	args = append([]interface{}{key}, args...)

	reply, err := redis.Int(conn.Do("HSET", args...))
	if err != nil || reply == 0 {
		return errors.Wrap(err, "failed to set hash")
	}

	if expiresIn > 0 {
		_, err = conn.Do("EXPIRE", key, expiresIn)
		if err != nil {
			return errors.Wrap(err, "failed to set expiration")
		}
	}

	return nil
}
