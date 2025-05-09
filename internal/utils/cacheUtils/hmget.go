package cacheUtils

import (
	"context"
	"fmt"
	"pr-reviewer/internal/utils/cacheUtils/utils"
	"reflect"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

type HMGetOptions struct {
	ForceRefresh bool
	LoadFunc     LoadFunc
}

func HMGet(ctx context.Context, pool *redis.Pool, key string, dest interface{}) error {
	return GetWithOptions(ctx, pool, key, dest, GetOptions{})
}

func HMGetWithOptions(ctx context.Context, pool *redis.Pool, key string, dest interface{}, options GetOptions) error {
	return HMGetWithLoadAndExpire(ctx, pool, key, options.ForceRefresh, dest, options.LoadFunc)
}

func HMGetWithLoadAndExpire(ctx context.Context, pool *redis.Pool, key string, forceRefresh bool, dest interface{}, loadFunc LoadFunc) error {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	fields, err := utils.StructToFieldsArray(dest)
	if err != nil {
		return err
	}
	args := append([]interface{}{key}, fields...)

	reply, err := redis.Values(conn.Do("HMGET", args...))
	if err != nil && err != redis.ErrNil {
		return err
	}

	if ((len(reply) != 0 && reply[0] == nil) || forceRefresh) && loadFunc != nil {
		data, expiresIn, err := loadFunc()
		if err != nil {
			return err
		}

		if err := utils.SetDist(dest, data); err != nil {
			return errors.Wrap(err, "failed to set dist")
		}
		args, err := utils.StructToArgs(data)
		if err != nil {
			return errors.Wrap(err, "failed to convert struct to args")
		}
		args = append([]interface{}{key}, args...)

		reply, err := redis.Int64(conn.Do("HSET", args...))
		if err != nil || reply == 0 {
			return err
		}

		if expiresIn > 0 {
			_, err = conn.Do("EXPIRE", key, expiresIn)
			if err != nil {
				return errors.Wrap(err, "failed to set expiration")
			}
		}

		return nil
	}

	if err := scanForHash(fields, reply, dest); err != nil {
		return errors.Wrap(err, "failed to scan for hash")
	}

	return nil
}

// scanForHash maps Redis HMGET result to struct fields based on Redis tags
func scanForHash(fields []interface{}, data []interface{}, dest interface{}) error {
	// Ensure dest is a pointer to a struct
	val := reflect.ValueOf(dest)
	typ := reflect.TypeOf(dest)

	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dest must be a pointer to a struct")
	}

	// Dereference the pointer
	val = val.Elem()

	if len(fields) != len(data) {
		return fmt.Errorf("mismatched fields and data length")
	}

	// Define a map of handlers for different types
	typeHandlers := map[reflect.Kind]func(interface{}) (interface{}, error){
		reflect.String: func(v interface{}) (interface{}, error) {
			return ParseRedisReply[string](v, nil)
		},
		reflect.Int: func(v interface{}) (interface{}, error) {
			return ParseRedisReply[int](v, nil)
		},
		reflect.Int64: func(v interface{}) (interface{}, error) {
			return ParseRedisReply[int64](v, nil)
		},
		reflect.Uint64: func(v interface{}) (interface{}, error) {
			return ParseRedisReply[uint64](v, nil)
		},
		reflect.Float64: func(v interface{}) (interface{}, error) {
			return ParseRedisReply[float64](v, nil)
		},
		reflect.Bool: func(v interface{}) (interface{}, error) {
			return ParseRedisReply[bool](v, nil)
		},
	}

	// Iterate over fields
	for i, field := range fields {
		fieldName := field.(string)
		fieldValue := data[i]

		// Find the struct field by its tag (Redis tag)
		structField := val.FieldByNameFunc(func(name string) bool {
			field, ok := typ.Elem().FieldByName(name)
			if !ok {
				return false
			}
			redisTag := field.Tag.Get("redis")
			return redisTag == fieldName
		})

		if !structField.IsValid() {
			return fmt.Errorf("field with tag '%s' not found in struct", fieldName)
		}
		if !structField.CanSet() {
			return fmt.Errorf("field %s is not settable", fieldName)
		}

		// Get field type and find the corresponding parser function
		fieldType := structField.Type()
		handler, exists := typeHandlers[fieldType.Kind()]
		if !exists {
			return fmt.Errorf("unsupported field type: %v", fieldType)
		}

		// Parse the value using the handler function
		parsedValue, err := handler(fieldValue)
		if err != nil {
			return fmt.Errorf("error parsing field '%s': %v", fieldName, err)
		}

		// Set the parsed value to the struct field
		structField.Set(reflect.ValueOf(parsedValue))
	}

	return nil
}
