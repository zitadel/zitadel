package mock

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

var _ driver.ValueConverter = (*TypeConverter)(nil)

type TypeConverter struct{}

// ConvertValue converts a value to a driver Value.
func (s TypeConverter) ConvertValue(v any) (driver.Value, error) {
	if driver.IsValue(v) {
		return v, nil
	}
	value := reflect.ValueOf(v)

	if rawMessage, ok := v.(json.RawMessage); ok {
		return convertBytes(rawMessage), nil
	}

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

func convertBytes(array []byte) string {
	var builder strings.Builder
	builder.Grow(hex.EncodedLen(len(array)) + 4)
	builder.WriteString(`\x`)
	builder.Write(hex.AppendEncode(nil, array))
	return builder.String()
}
