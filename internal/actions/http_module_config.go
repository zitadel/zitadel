package actions

import (
	"net"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/zitadel/zitadel/internal/zerrors"
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
		DenyList: make([]AddressChecker, len(config.DenyList)),
	}

	for i, entry := range config.DenyList {
		if c.DenyList[i], err = parseDenyListEntry(entry); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func parseDenyListEntry(entry string) (AddressChecker, error) {
	if checker, err := NewIPChecker(entry); err == nil {
		return checker, nil
	}
	return &DomainChecker{Domain: entry}, nil
}

func NewIPChecker(i string) (AddressChecker, error) {
	_, network, err := net.ParseCIDR(i)
	if err == nil {
		return &IPChecker{Net: network}, nil
	}
	if ip := net.ParseIP(i); ip != nil {
		return &IPChecker{IP: ip}, nil
	}
	return nil, zerrors.ThrowInvalidArgument(nil, "ACTIO-ddJ7h", "invalid ip")
}

type IPChecker struct {
	Net *net.IPNet
	IP  net.IP
}

func (c *IPChecker) Matches(address string) bool {
	ip := net.ParseIP(address)
	if ip == nil {
		return false
	}

	if c.IP != nil {
		return c.IP.Equal(ip)
	}
	return c.Net.Contains(ip)
}

type DomainChecker struct {
	Domain string
}

func (c *DomainChecker) Matches(domain string) bool {
	//TODO: allow wild cards
	return c.Domain == domain
}
