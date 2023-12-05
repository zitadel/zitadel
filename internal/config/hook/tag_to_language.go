package hook

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
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

		lang, err := domain.ParseLanguage(data.(string))
		return lang[0], err
	}
}
