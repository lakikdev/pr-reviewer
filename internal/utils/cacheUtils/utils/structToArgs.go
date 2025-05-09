package utils

import (
	"fmt"
	"reflect"
)

func StructToArgs(data interface{}) ([]interface{}, error) {
	// Create a slice to hold the converted arguments
	var args []interface{}

	// Reflect on the struct to get its fields and values
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", v.Kind())
	}

	// Iterate over each field in the struct
	for i := 0; i < v.NumField(); i++ {
		// Get the field name and value
		fieldTagName := v.Type().Field(i).Tag.Get("redis")
		fieldValue := v.Field(i)

		// Append field name and value to the arguments slice
		args = append(args, fieldTagName, fieldValue.Interface())
	}

	return args, nil
}
