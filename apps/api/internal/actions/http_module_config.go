package actions

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

func SetHTTPConfig(config *HTTPConfig) {
	httpConfig = config
}

var httpConfig *HTTPConfig

type HTTPConfig struct {
	DenyList []AddressChecker
}

func HTTPConfigDecodeHook(from, to reflect.Value) (interface{}, error) {
	if to.Type() != reflect.TypeOf(HTTPConfig{}) {
		return from.Interface(), nil
	}

	config := struct {
		DenyList []string
	}{}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
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
		DenyList: make([]AddressChecker, 0),
	}

	for _, unsplit := range config.DenyList {
		for _, split := range strings.Split(unsplit, ",") {
			parsed, parseErr := NewHostChecker(split)
			if parseErr != nil {
				return nil, parseErr
			}
			if parsed != nil {
				c.DenyList = append(c.DenyList, parsed)
			}
		}
	}

	return c, nil
}

func NewHostChecker(entry string) (AddressChecker, error) {
	if entry == "" {
		return nil, nil
	}
	_, network, err := net.ParseCIDR(entry)
	if err == nil {
		return &HostChecker{Net: network}, nil
	}
	if ip := net.ParseIP(entry); ip != nil {
		return &HostChecker{IP: ip}, nil
	}
	return &HostChecker{Domain: entry}, nil
}

type HostChecker struct {
	Net    *net.IPNet
	IP     net.IP
	Domain string
}

type AddressDeniedError struct {
	deniedBy string
}

func NewAddressDeniedError(deniedBy string) *AddressDeniedError {
	return &AddressDeniedError{deniedBy: deniedBy}
}

func (e *AddressDeniedError) Error() string {
	return fmt.Sprintf("address is denied by '%s'", e.deniedBy)
}

func (e *AddressDeniedError) Is(target error) bool {
	var addressDeniedErr *AddressDeniedError
	if !errors.As(target, &addressDeniedErr) {
		return false
	}
	return e.deniedBy == addressDeniedErr.deniedBy
}

func (c *HostChecker) IsDenied(ips []net.IP, address string) error {
	// if the address matches the domain, no additional checks as needed
	if c.Domain == address {
		return NewAddressDeniedError(c.Domain)
	}
	// otherwise we need to check on ips (incl. the resolved ips of the host)
	for _, ip := range ips {
		if c.Net != nil && c.Net.Contains(ip) {
			return NewAddressDeniedError(c.Net.String())
		}
		if c.IP != nil && c.IP.Equal(ip) {
			return NewAddressDeniedError(c.IP.String())
		}
	}
	return nil
}
