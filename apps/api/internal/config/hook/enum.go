package hook

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/constraints"
)

func EnumHookFunc[T constraints.Integer](resolve func(string) (T, error)) mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(T(0)) {
			return data, nil
		}
		return resolve(data.(string))
	}
}
