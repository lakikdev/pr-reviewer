package helper

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/copier"
)

func StructToMap(in interface{}, tag string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" {
			out[tagv] = v.Field(i).Interface()
		}
	}
	return out, nil
}

func ToFieldSlice[T any](slice interface{}, filedName string) []T {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("ToFieldSlice only accepts slice")
	}

	var result []T
	for i := 0; i < s.Len(); i++ {
		item := s.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		filed := item.FieldByName(filedName)
		if filed.Kind() == reflect.Ptr {
			result = append(result, filed.Elem().Interface().(T))
			continue
		}
		result = append(result, filed.Interface().(T))
	}
	return result
}

func ToInterfaceSlice[T any](items []T) []interface{} {
	result := make([]interface{}, len(items))
	for i, v := range items {
		result[i] = v
	}
	return result
}

func DeepCopy[T interface{}](src T) (copy T, err error) {
	err = copier.CopyWithOption(&copy, &src, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	return copy, err
}

// Pointer is a helper function to create a pointer of a value
func Pointer[T any](value T) *T {
	return &value
}

// CastPointer is a helper function to cast an interface to a pointer of a type
func CastPointer[T any](value interface{}) *T {
	v := value.(T)
	return &v
}
