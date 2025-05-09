package cacheUtils

import (
	"context"
	"encoding/json"
	"pr-reviewer/internal/utils/cacheUtils/utils"

	"github.com/gomodule/redigo/redis"
)

// LoadFunc is a function type that defines a loader function which returns
// an interface{}, an int64, and an error. The interface{} represents the
// loaded data, the int64 represents the expiration time in seconds,
// and the error indicates if there was an issue during the loading
// process.
type LoadFunc func() (interface{}, int64, error)

type GetOptions struct {
	ForceRefresh bool
	LoadFunc     LoadFunc
}

func Get(ctx context.Context, pool *redis.Pool, key string, dest interface{}) error {
	return GetWithOptions(ctx, pool, key, dest, GetOptions{})
}

func GetWithOptions(ctx context.Context, pool *redis.Pool, key string, dest interface{}, options GetOptions) error {
	return getWithLoadAndExpire(ctx, pool, key, options.ForceRefresh, dest, options.LoadFunc)
}

func getWithLoadAndExpire(ctx context.Context, pool *redis.Pool, key string, forceRefresh bool, dest interface{}, loadFunc LoadFunc) error {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return err
	}

	if (reply == "" || forceRefresh) && loadFunc != nil {
		data, expiresIn, err := loadFunc()
		if err != nil {
			return err
		}

		if err := utils.SetDist(dest, data); err != nil {
			return err
		}

		dataJson, err := json.Marshal(dest)
		if err != nil {
			return err
		}

		if expiresIn == 0 {
			reply, err := redis.String(conn.Do("SET", key, dataJson))
			if err != nil || reply != "OK" {
				return err
			}
			return nil
		}

		reply, err := redis.String(conn.Do("SETEX", key, expiresIn, dataJson))
		if err != nil || reply != "OK" {
			return err
		}
		return nil
	}

	if reply != "" {
		if err := json.Unmarshal([]byte(reply), &dest); err != nil {
			return err
		}
	}

	return nil
}
