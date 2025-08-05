package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func WithHTTP(ctx context.Context) Option {
	return func(c *runConfig) {
		c.modules["zitadel/http"] = func(runtime *goja.Runtime, module *goja.Object) {
			requireHTTP(ctx, &http.Client{Transport: &transport{lookup: net.LookupIP}}, runtime, module)
		}
	}
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

func (c *HTTP) fetchConfigFromArg(arg *goja.Object, config *fetchConfig) (err error) {
	for _, key := range arg.Keys() {
		switch key {
		case "headers":
			config.Headers = parseHeaders(arg.Get(key).ToObject(c.runtime))
		case "method":
			config.Method = arg.Get(key).String()
		case "body":
			body, err := arg.Get(key).ToObject(c.runtime).MarshalJSON()
			if err != nil {
				return err
			}
			config.Body = bytes.NewReader(body)
		default:
			return zerrors.ThrowInvalidArgument(nil, "ACTIO-OfUeA", "key is invalid")
		}
	}
	return nil
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
		req := c.buildHTTPRequest(ctx, call.Arguments)
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

// the first argument has to be a string and is required
// the second agrument is optional and an object with the following fields possible:
// - `Headers`: map with string key and value of type string or string array
// - `Body`: json body of the request
// - `Method`: http method type
func (c *HTTP) buildHTTPRequest(ctx context.Context, args []goja.Value) (req *http.Request) {
	if len(args) > 2 {
		logging.WithFields("count", len(args)).Debug("more than 2 args provided")
		panic("too many args")
	}

	if len(args) == 0 {
		panic("no url provided")
	}

	config := defaultFetchConfig
	var err error
	if len(args) == 2 {
		if err = c.fetchConfigFromArg(args[1].ToObject(c.runtime), &config); err != nil {
			panic(err)
		}
	}

	req, err = http.NewRequestWithContext(ctx, config.Method, args[0].Export().(string), config.Body)
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
		var values []string

		switch headerValue := header.(type) {
		case string:
			values = strings.Split(headerValue, ",")
		case []any:
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

type transport struct {
	lookup func(string) ([]net.IP, error)
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if httpConfig == nil || len(httpConfig.DenyList) == 0 {
		return http.DefaultTransport.RoundTrip(req)
	}
	if err := t.isHostBlocked(httpConfig.DenyList, req.URL); err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "ACTIO-N72d0", "host is denied")
	}
	return http.DefaultTransport.RoundTrip(req)
}

func (t *transport) isHostBlocked(denyList []AddressChecker, address *url.URL) error {
	host := address.Hostname()
	ip := net.ParseIP(host)
	ips := []net.IP{ip}
	// if the hostname is a domain, we need to check resolve the ip(s), since it might be denied
	if ip == nil {
		var err error
		ips, err = t.lookup(host)
		if err != nil {
			return zerrors.ThrowInternal(err, "ACTIO-4m9s2", "lookup failed")
		}
	}
	for _, denied := range denyList {
		if err := denied.IsDenied(ips, host); err != nil {
			return err
		}
	}
	return nil
}

type AddressChecker interface {
	IsDenied([]net.IP, string) error
}
