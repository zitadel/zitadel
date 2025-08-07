package test

import (
	"reflect"

	"github.com/stretchr/testify/assert"
)

func AssertMapContains[M ~map[K]V, K comparable, V any](t assert.TestingT, m M, key K, expectedValue V) {
	val, exists := m[key]
	assert.True(t, exists, "Key '%s' should exist in the map", key)
	if !exists {
		return
	}

	assert.Equal(t, expectedValue, val, "Key '%s' should have value '%d'", key, expectedValue)
}

// PartiallyDeepEqual is similar to reflect.DeepEqual,
// but only compares exported non-zero fields of the expectedValue
func PartiallyDeepEqual(expected, actual interface{}) bool {
	if expected == nil {
		return actual == nil
	}

	if actual == nil {
		return false
	}

	return partiallyDeepEqual(reflect.ValueOf(expected), reflect.ValueOf(actual))
}

func partiallyDeepEqual(expected, actual reflect.Value) bool {
	// Dereference pointers if needed
	if expected.Kind() == reflect.Ptr {
		if expected.IsNil() {
			return true
		}

		expected = expected.Elem()
	}

	if actual.Kind() == reflect.Ptr {
		if actual.IsNil() {
			return false
		}

		actual = actual.Elem()
	}

	if expected.Type() != actual.Type() {
		return false
	}

	switch expected.Kind() { //nolint:exhaustive
	case reflect.Struct:
		for i := 0; i < expected.NumField(); i++ {
			field := expected.Type().Field(i)
			if field.PkgPath != "" { // Skip unexported fields
				continue
			}

			expectedField := expected.Field(i)
			actualField := actual.Field(i)

			// Skip zero-value fields in expected
			if reflect.DeepEqual(expectedField.Interface(), reflect.Zero(expectedField.Type()).Interface()) {
				continue
			}

			// Compare fields recursively
			if !partiallyDeepEqual(expectedField, actualField) {
				return false
			}
		}
		return true

	case reflect.Slice, reflect.Array:
		if expected.Len() > actual.Len() {
			return false
		}

		for i := 0; i < expected.Len(); i++ {
			if !partiallyDeepEqual(expected.Index(i), actual.Index(i)) {
				return false
			}
		}

		return true

	default:
		// Compare primitive types
		return reflect.DeepEqual(expected.Interface(), actual.Interface())
	}
}

func Must[T any](result T, error error) T {
	if error != nil {
		panic(error)
	}

	return result
}
