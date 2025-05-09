package dbHelper

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func GetDBFields(i interface{}) []string {
	return GetDBFieldsWithIgnore(i, []string{})
}

// e.g. []string{"id", "name", "description"}
func GetDBFieldsWithIgnore(i interface{}, ignoreFields []string) []string {
	ignoreFieldsMap := make(map[string]string)
	for _, field := range ignoreFields {
		ignoreFieldsMap[field] = field
	}

	return getDBFieldsWithIgnoreRecurse(reflect.TypeOf(i), &ignoreFieldsMap)
}

// e.g. "id, name, description"
func GetDBFieldsCSV(fields []string) string {
	return strings.Join(fields, ", ")
}

// e.g. ":id, :name, :description"
func GetDBFieldsCSVColons(fields []string) string {
	for i := 0; i < len(fields); i++ {
		fields[i] = ":" + fields[i]
	}
	return strings.Join(fields, ", ")
}

func GetDBFieldsUpdate(fields []string) string {
	return GetDBFieldsUpdateWithExtra(fields, []string{})
}

// apply alias to fields
func ApplyAliasToFields(fields []string, alias string) []string {
	for i := 0; i < len(fields); i++ {
		fields[i] = fmt.Sprintf("%s.%s", alias, fields[i])
	}
	return fields
}

// e.g. "id=:id, name=:name, description=:description, {extra}, {extra}"
func GetDBFieldsUpdateWithExtra(fields []string, extra []string) string {
	var query []string
	for i := 0; i < len(fields); i++ {
		query = append(query, fmt.Sprintf("%s=:%s", fields[i], fields[i]))
	}
	query = append(query, extra...)
	return strings.Join(query, ", ")
}

func GetDBFieldNameByJsonName(i interface{}, jsonName string) (fieldName string, fieldType reflect.Type) {
	return GetDBFieldNameByJsonNameRecurse(reflect.TypeOf(i), jsonName)
}

func ParseValueToFieldType(fieldType reflect.Type, value string) (correctValue interface{}, err error) {
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	switch fieldType.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.ParseInt(value, 10, 64)
	case reflect.Float32, reflect.Float64:
		return strconv.ParseFloat(value, 64)
	case reflect.Bool:
		return strconv.ParseBool(value)

	default:
		return strings.ToLower(value), err
	}
}

// e.g. []string{"id", "name", "description"}
func getDBFieldsWithIgnoreRecurse(goType reflect.Type, ignoreFieldsMap *map[string]string) []string {

	// Look inside slices and pointers.
	goType = unpackElemType(goType)

	fields := make([]string, 0)
	for i := 0; i < goType.NumField(); i++ {
		field := goType.Field(i)
		dbTag := field.Tag.Get("db")

		// Skip ignored or disabled fields.
		if _, found := (*ignoreFieldsMap)[dbTag]; found || dbTag == "-" {
			continue
		}

		// Look inside slices and pointers.
		fieldType := unpackElemType(field.Type)

		// Recurse inside structs with no name - if it does have a name, we
		// instead assume it's a json column and don't look inside. For any
		// other type, we skip if it has no name.
		if fieldType.Kind() == reflect.Struct && dbTag == "" {
			fields = append(fields, getDBFieldsWithIgnoreRecurse(fieldType,
				ignoreFieldsMap)...)
		} else if dbTag != "" {
			fields = append(fields, dbTag)
		}
	}

	return fields
}

func GetDBFieldNameByJsonNameRecurse(goType reflect.Type, jsonName string) (fieldName string, fieldType reflect.Type) {
	goType = unpackElemType(goType)

	if goType.Kind() == reflect.Struct {
		for i := 0; i < goType.NumField(); i++ {
			member := goType.Field(i)
			dbTag := member.Tag.Get("db")

			// Look inside slices and pointers.
			memberType := unpackElemType(member.Type)

			// Recurse inside structs with no name.
			if memberType.Kind() == reflect.Struct && dbTag == "" {

				fieldName, fieldType = GetDBFieldNameByJsonNameRecurse(memberType, jsonName)
				if fieldType != nil {
					return fieldName, fieldType
				}

			} else if strings.Split(member.Tag.Get("json"), ",")[0] == jsonName {

				//json tag can contain ,omitempty flag which will break the code so we split by , (comma)  and get first part which is field name
				return member.Tag.Get("db"), member.Type
			}
		}
	}

	return fieldName, fieldType
}

// Peels off pointer and slice types to get at the contained type.
func unpackElemType(t reflect.Type) reflect.Type {

	for t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	return t
}
