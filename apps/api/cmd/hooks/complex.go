package hooks

import (
	"encoding/json"
	"net/http"
	"reflect"
)

func SliceTypeStringDecode[T any](from, to reflect.Value) (any, error) {
	into := make([]T, 0)
	return complexTypeStringDecodeHook(from, to, into)
}

func MapTypeStringDecode[K ~string | ~int, V any](from, to reflect.Value) (any, error) {
	into := make(map[K]V, 0)
	return complexTypeStringDecodeHook(from, to, into)
}

func MapHTTPHeaderStringDecode(from, to reflect.Value) (any, error) {
	into := http.Header{}
	return complexTypeStringDecodeHook(from, to, into)
}

func complexTypeStringDecodeHook(from, to reflect.Value, out any) (any, error) {
	fromInterface := from.Interface()
	if to.Type() != reflect.TypeOf(out) {
		return fromInterface, nil
	}
	data, ok := fromInterface.(string)
	if !ok {
		return fromInterface, nil
	}
	err := json.Unmarshal([]byte(data), &out)
	return out, err
}
