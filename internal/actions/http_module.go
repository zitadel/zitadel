package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/mitchellh/mapstructure"
	"github.com/zitadel/logging"

	z_errs "github.com/zitadel/zitadel/internal/errors"
)

func WithHTTP(ctx context.Context) Option {
	return func(c *runConfig) {
		c.modules["zitadel/http"] = func(runtime *goja.Runtime, module *goja.Object) {
			requireHTTP(ctx, &http.Client{Transport: new(transport)}, runtime, module)
		}
	}
}

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

type transport struct{}

func (*transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if httpConfig == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	if isHostBlocked(httpConfig.DenyList, req.URL) {
		return nil, z_errs.ThrowInvalidArgument(nil, "ACTIO-N72d0", "host is denied")
	}
	return http.DefaultTransport.RoundTrip(req)
}

type HTTP struct {
	runtime *goja.Runtime
	client  *http.Client
}

func requireHTTP(ctx context.Context, client *http.Client, runtime *goja.Runtime, module *goja.Object) {
	c := &HTTP{
		client:  client,
		runtime: runtime,
	}
	o := module.Get("exports").(*goja.Object)
	logging.OnError(o.Set("fetch", c.fetch(ctx))).Warn("unable to set module")
}

type fetchConfig struct {
	Method  string
	Headers http.Header
	Body    io.Reader
}

var defaultFetchConfig = fetchConfig{
	Method: http.MethodGet,
	Headers: http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/json"},
	},
}

func (c *HTTP) fetchConfigFromArg(arg *goja.Object) (config fetchConfig, err error) {
	for _, key := range arg.Keys() {
		switch key {
		case "headers":
			config.Headers = parseHeaders(arg.Get(key).ToObject(c.runtime))
		case "method":
			config.Method = arg.Get(key).String()
		case "body":
			body, err := arg.Get(key).ToObject(c.runtime).MarshalJSON()
			if err != nil {
				return config, err
			}
			config.Body = bytes.NewReader(body)
		default:
			return config, errors.New("unimplemented")
		}
	}
	return config, nil
}

type response struct {
	Body    string
	Status  int
	Headers map[string][]string
	runtime *goja.Runtime
}

func (r *response) Json() goja.Value {
	var val interface{}

	if err := json.Unmarshal([]byte(r.Body), &val); err != nil {
		panic(err)
	}

	return r.runtime.ToValue(val)
}

func (r *response) Text() goja.Value {
	return r.runtime.ToValue(r.Body)
}

func (c *HTTP) fetch(ctx context.Context) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		req := c.buildHTTPRequest(call.Arguments)
		if deadline, ok := ctx.Deadline(); ok {
			c.client.Timeout = time.Until(deadline)
		}

		res, err := c.client.Do(req)
		if err != nil {
			logging.WithError(err).Debug("call failed")
			panic(err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			logging.WithError(err).Warn("unable to parse body")
			panic("unable to read response body")
		}
		return c.runtime.ToValue(&response{Status: res.StatusCode, Body: string(body), runtime: c.runtime})
	}
}

func (c *HTTP) buildHTTPRequest(args []goja.Value) (req *http.Request) {
	if len(args) > 2 {
		logging.WithFields("count", len(args)).Debug("more than 2 args provided")
		panic("too many args")
	}

	if len(args) < 1 {
		panic("no url provided")
	}

	config := defaultFetchConfig
	var err error
	if len(args) == 2 {
		config, err = c.fetchConfigFromArg(args[1].ToObject(c.runtime))
		if err != nil {
			panic(err)
		}
	}

	req, err = http.NewRequest(config.Method, args[0].Export().(string), config.Body)
	if err != nil {
		panic(err)
	}
	req.Header = config.Headers

	return req
}

func parseHeaders(headers *goja.Object) http.Header {
	h := make(http.Header, len(headers.Keys()))
	for _, k := range headers.Keys() {
		header := headers.Get(k).Export()
		values := []string{}

		switch headerValue := header.(type) {
		case string:
			values = strings.Split(headerValue, ",")
		case []interface{}:
			for _, v := range headerValue {
				values = append(values, v.(string))
			}
		}

		for _, v := range values {
			h.Add(k, strings.TrimSpace(v))
		}
	}
	return h
}

func parseDenyListEntry(entry string) (AddressChecker, error) {
	if checker, err := NewIPChecker(entry); err == nil {
		return checker, nil
	}
	return &DomainChecker{Domain: entry}, nil
}

func isHostBlocked(denyList []AddressChecker, address *url.URL) bool {
	for _, blocked := range denyList {
		if blocked.Matches(address.Hostname()) {
			return true
		}
	}
	return false
}

type AddressChecker interface {
	Matches(string) bool
}

func NewIPChecker(i string) (AddressChecker, error) {
	_, network, err := net.ParseCIDR(i)
	if err == nil {
		return &IPChecker{Net: network}, nil
	}
	if ip := net.ParseIP(i); ip != nil {
		return &IPChecker{IP: ip}, nil
	}
	return nil, z_errs.ThrowInvalidArgument(nil, "ACTIO-ddJ7h", "invalid ip")

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
	return c.Domain == domain
}
