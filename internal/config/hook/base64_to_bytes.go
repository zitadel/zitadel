package hook

import (
	"encoding/base64"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

func Base64ToBytesHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf([]byte{}) {
			return data, nil
		}

		return base64.StdEncoding.DecodeString(data.(string))
	}
}
