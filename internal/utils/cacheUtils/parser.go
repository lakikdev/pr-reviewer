package cacheUtils

import (
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// ParseRedisReply tries to parse a Redis reply using the correct redigo helper function.
func ParseRedisReply[T any](reply interface{}, err error) (T, error) {
	var zero T // Zero value of T

	if err != nil {
		return zero, err
	}

	switch any(zero).(type) {
	case int:
		v, err := redis.Int(reply, nil)
		return any(v).(T), err
	case int64:
		v, err := redis.Int64(reply, nil)
		return any(v).(T), err
	case uint64:
		v, err := redis.Uint64(reply, nil)
		return any(v).(T), err
	case float64:
		v, err := redis.Float64(reply, nil)
		return any(v).(T), err
	case string:
		v, err := redis.String(reply, nil)
		return any(v).(T), err
	case []byte:
		v, err := redis.Bytes(reply, nil)
		return any(v).(T), err
	case bool:
		v, err := redis.Bool(reply, nil)
		return any(v).(T), err
	case []string:
		v, err := redis.Strings(reply, nil)
		return any(v).(T), err
	case []int:
		v, err := redis.Ints(reply, nil)
		return any(v).(T), err
	case []int64:
		v, err := redis.Int64s(reply, nil)
		return any(v).(T), err
	case []uint64:
		v, err := redis.Uint64s(reply, nil)
		return any(v).(T), err
	case [][]byte:
		v, err := redis.ByteSlices(reply, nil)
		return any(v).(T), err
	case map[string]string:
		v, err := redis.StringMap(reply, nil)
		return any(v).(T), err
	case map[string]int:
		v, err := redis.IntMap(reply, nil)
		return any(v).(T), err
	case map[string]int64:
		v, err := redis.Int64Map(reply, nil)
		return any(v).(T), err
	case map[string]uint64:
		v, err := redis.Uint64Map(reply, nil)
		return any(v).(T), err
	default:
		return zero, errors.New(fmt.Sprintf("unsupported type: %T", zero))
	}
}

