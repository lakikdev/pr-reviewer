package utils

import (
	"fmt"
	"reflect"
)

func StructToFieldsArray(data interface{}) ([]interface{}, error) {
	// Create a slice to hold the converted arguments
	var args []interface{}

	// Reflect on the struct to get its fields and values
	v := reflect.ValueOf(data)
	//if got a pointer, unwrap it

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Check if the value is a struct
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", v.Kind())
	}

	// Iterate over each field in the struct
	for i := 0; i < v.NumField(); i++ {
		// Get the field name and value
		fieldTagName := v.Type().Field(i).Tag.Get("redis")

		// Append field name and value to the arguments slice
		args = append(args, fieldTagName)
	}

	return args, nil
}
