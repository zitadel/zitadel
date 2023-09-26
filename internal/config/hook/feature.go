package hook

import (
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/zitadel/zitadel/internal/domain"
)

func StringToFeatureHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(domain.FeatureUnspecified) {
			return data, nil
		}

		return domain.FeatureString(data.(string))
	}
}
