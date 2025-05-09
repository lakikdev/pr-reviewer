package utils

import (
	"fmt"
	"reflect"
)

func SetDist(dest, data interface{}) error {
	destValue := reflect.ValueOf(dest)
	dataValue := reflect.ValueOf(data)

	// Ensure dest is a pointer
	if destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}

	// Get the actual element of the pointer
	destElem := destValue.Elem()

	// Ensure dest is settable
	if !destElem.CanSet() {
		return fmt.Errorf("dest is not settable")
	}

	// Check if types match
	if destElem.Type() == dataValue.Type() {
		destElem.Set(dataValue)
		return nil
	}

	// Handle struct assignment (if applicable)
	if destElem.Kind() == reflect.Struct && dataValue.Kind() == reflect.Struct {
		dataType := dataValue.Type()

		for i := 0; i < dataType.NumField(); i++ {
			field := dataType.Field(i)
			dataFieldValue := dataValue.Field(i)

			// Find the corresponding field in dest
			destField := destElem.FieldByName(field.Name)
			if destField.IsValid() && destField.CanSet() && destField.Type() == dataFieldValue.Type() {
				destField.Set(dataFieldValue)
			}
		}
		return nil
	}

	return fmt.Errorf("type mismatch: cannot assign %v to %v", dataValue.Type(), destElem.Type())
}
