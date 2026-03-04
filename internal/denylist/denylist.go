package denylist

import (
	"net"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type AddressChecker interface {
	IsDenied([]net.IP, string) error
}

func ParseDenyList(inputList []string) ([]AddressChecker, error) {
	var toReturn []AddressChecker
	for _, input := range inputList {
		parsed, parseErr := NewHostChecker(input)
		if parseErr != nil {
			return nil, parseErr
		}
		if parsed != nil {
			toReturn = append(toReturn, parsed)
		}
	}
	return toReturn, nil
}

func DenyListDecodeHook(from, to reflect.Value) (interface{}, error) {
	if to.Type() != reflect.TypeFor[[]AddressChecker]() {
		return from.Interface(), nil
	}
	var inputDenyList []string
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToSliceHookFunc(","),
		),
		WeaklyTypedInput: true,
		Result:           &inputDenyList,
	})
	if err != nil {
		return nil, err
	}
	if err = decoder.Decode(from.Interface()); err != nil {
		return nil, err
	}
	return ParseDenyList(inputDenyList)
}
