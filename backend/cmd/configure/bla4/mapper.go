package bla4

import (
	"reflect"
	"strconv"
	"sync"
	"time"
)

var writeMu sync.RWMutex

func SetTypeMapper(typ reflect.Type, mapper func(input string) (any, error)) {
	writeMu.Lock()
	defer writeMu.Unlock()

	typeMappers[typ] = mapper
}

func SetTypeMapperFor[T any](mapper func(input string) (any, error)) {
	writeMu.Lock()
	defer writeMu.Unlock()

	typeMappers[reflect.TypeFor[T]()] = mapper
}

func typeMapper(typ reflect.Type) func(input string) (any, error) {
	writeMu.RLock()
	defer writeMu.RUnlock()

	return typeMappers[typ]
}

var typeMappers = map[reflect.Type]func(string) (any, error){
	reflect.TypeFor[time.Duration](): func(input string) (any, error) {
		return time.ParseDuration(input)
	},
}

// SetKindMapper overwrites the mapper for the given kind.
func SetKindMapper(kind reflect.Kind, mapper func(input string) (any, error)) {
	writeMu.Lock()
	defer writeMu.Unlock()

	kindMappers[kind] = mapper
}

func kindMapper(kind reflect.Kind) func(input string) (any, error) {
	writeMu.RLock()
	defer writeMu.RUnlock()

	return kindMappers[kind]
}

var kindMappers = map[reflect.Kind]func(input string) (any, error){
	reflect.String: func(input string) (any, error) {
		return input, nil
	},
	reflect.Bool: func(input string) (any, error) {
		return strconv.ParseBool(input)
	},
	reflect.Int: func(input string) (any, error) {
		return strconv.Atoi(input)
	},
	reflect.Int8: func(input string) (val any, err error) {
		val, err = strconv.ParseInt(input, 10, 8)
		val = int8(val.(int64))
		return val, err
	},
	reflect.Int16: func(input string) (val any, err error) {
		val, err = strconv.ParseInt(input, 10, 16)
		val = int16(val.(int64))
		return val, err
	},
	reflect.Int32: func(input string) (val any, err error) {
		val, err = strconv.ParseInt(input, 10, 32)
		val = int32(val.(int64))
		return val, err
	},
	reflect.Int64: func(input string) (any, error) {
		return strconv.ParseInt(input, 10, 64)
	},
	reflect.Uint: func(input string) (val any, err error) {
		val, err = strconv.ParseUint(input, 10, 0)
		val = uint(val.(uint64))
		return val, err
	},
	reflect.Uint8: func(input string) (val any, err error) {
		val, err = strconv.ParseUint(input, 10, 8)
		val = uint8(val.(uint64))
		return val, err
	},
	reflect.Uint16: func(input string) (val any, err error) {
		val, err = strconv.ParseUint(input, 10, 16)
		val = uint16(val.(uint64))
		return val, err
	},
	reflect.Uint32: func(input string) (val any, err error) {
		val, err = strconv.ParseUint(input, 10, 32)
		val = uint32(val.(uint64))
		return val, err
	},
	reflect.Uint64: func(input string) (any, error) {
		return strconv.ParseUint(input, 10, 64)
	},
	reflect.Float32: func(input string) (val any, err error) {
		val, err = strconv.ParseFloat(input, 32)
		val = float32(val.(float64))
		return val, err
	},
	reflect.Float64: func(input string) (any, error) {
		return strconv.ParseFloat(input, 64)
	},
	reflect.Complex64: func(input string) (val any, err error) {
		val, err = strconv.ParseComplex(input, 64)
		val = complex64(val.(complex128))
		return val, err
	},
	reflect.Complex128: func(input string) (any, error) {
		return strconv.ParseComplex(input, 128)
	},
}
