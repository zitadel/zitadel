package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dop251/goja"

	"github.com/zitadel/logging"
)

func WithHTTP(ctx context.Context, client *http.Client) Option {
	return func(c *runConfig) {
		c.modules["zitadel/http"] = func(runtime *goja.Runtime, module *goja.Object) {
			requireHTTP(ctx, client, runtime, module)
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
	Body       string
	StatusCode int
	Headers    map[string][]string
	runtime    *goja.Runtime
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
		return c.runtime.ToValue(&response{StatusCode: res.StatusCode, Body: string(body), runtime: c.runtime})
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
