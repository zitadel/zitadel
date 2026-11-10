package hook

import (
	"net/url"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

func StringToURLHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(&url.URL{}) {
			return data, nil
		}

		u, err := url.Parse(data.(string))
		return u, err
	}
}
