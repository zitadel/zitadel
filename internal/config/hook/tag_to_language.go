package hook

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/text/language"
)

func TagToLanguageHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(language.Tag{}) {
			return data, nil
		}

		return language.Parse(data.(string))
	}
}
