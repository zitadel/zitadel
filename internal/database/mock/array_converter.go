package mock

import (
	"database/sql/driver"
	"reflect"
	"strconv"
	"strings"
)

var _ driver.ValueConverter = (*ArrayConverter)(nil)

type ArrayConverter struct{}

// ConvertValue converts a value to a driver Value.
func (s ArrayConverter) ConvertValue(v any) (driver.Value, error) {
	if driver.IsValue(v) {
		return v, nil
	}

	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Slice {
		//nolint: exhaustive
		// only defined types
		switch value.Type().Elem().Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			return convertSigned(value), nil
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			return convertUnsigned(value), nil
		case reflect.String:
			return convertText(value), nil
		}
	}
	return v, nil
}

// converts a text array to valid pgx v5 representation
func convertSigned(array reflect.Value) string {
	slice := make([]string, array.Len())
	for i := 0; i < array.Len(); i++ {
		slice[i] = strconv.FormatInt(array.Index(i).Int(), 10)
	}

	return "{" + strings.Join(slice, ",") + "}"
}

// converts a text array to valid pgx v5 representation
func convertUnsigned(array reflect.Value) string {
	slice := make([]string, array.Len())
	for i := 0; i < array.Len(); i++ {
		slice[i] = strconv.FormatUint(array.Index(i).Uint(), 10)
	}

	return "{" + strings.Join(slice, ",") + "}"
}

// converts a text array to valid pgx v5 representation
func convertText(array reflect.Value) string {
	slice := make([]string, array.Len())
	for i := 0; i < array.Len(); i++ {
		slice[i] = array.Index(i).String()
	}

	return "{" + strings.Join(slice, ",") + "}"
}
