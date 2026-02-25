package actions

import (
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/zitadel/zitadel/internal/denylist"
)

func SetHTTPConfig(config *HTTPConfig) {
	httpConfig = config
}

var httpConfig *HTTPConfig

type HTTPConfig struct {
	DenyList []denylist.AddressChecker
}

func HTTPConfigDecodeHook(from, to reflect.Value) (any, error) {
	if to.Type() != reflect.TypeFor[HTTPConfig]() {
		return from.Interface(), nil
	}

	config := struct {
		DenyList []string
	}{}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
		WeaklyTypedInput: true,
		Result:           &config,
	})
	if err != nil {
		return nil, err
	}

	if err = decoder.Decode(from.Interface()); err != nil {
		return nil, err
	}

	c := HTTPConfig{
		DenyList: make([]denylist.AddressChecker, 0),
	}

	c.DenyList, err = denylist.ParseDenyList(config.DenyList)
	if err != nil {
		return nil, err
	}

	return c, nil
}
