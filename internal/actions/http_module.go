package actions

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/zitadel/logging"
)

func WithHTTP(client *http.Client) runOpt {
	return func(c *runConfig) {
		c.modules["zitadel/http"] = func(runtime *goja.Runtime, module *goja.Object) {
			requireHTTP(client, c.end, runtime, module)
		}
	}
}

type HTTP struct {
	runtime *goja.Runtime
	client  *http.Client
	maxEnd  time.Time
}

func requireHTTP(client *http.Client, maxEnd time.Time, runtime *goja.Runtime, module *goja.Object) {
	c := &HTTP{
		client:  client,
		runtime: runtime,
		maxEnd:  maxEnd,
	}
	o := module.Get("exports").(*goja.Object)
	logging.OnError(
		o.Set("fetch", c.fetch),
	).Warn("somesing happened")
}

type fetchConfig struct {
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Body    io.Reader   `json:"body"`
}

var defaultFetchConfig = fetchConfig{
	Method: http.MethodGet,
	Headers: http.Header{
		"Content-Type": []string{"application/json"},
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
	Body       interface{} `json:"body"`
	StatusCode int         `json:"status"`
	//TODO: add headers
}

func (c *HTTP) fetch(call goja.FunctionCall) goja.Value {
	req, err := c.buildHTTPRequest(call.Arguments)
	if err != nil {
		// handle error
		logging.WithError(err).Warn("new req failed")
	}

	c.client.Timeout = time.Until(c.maxEnd)
	res, err := c.client.Do(req)
	if err != nil {
		logging.WithError(err).Warn("call failed")
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// TODO: do something meaningfull with the body
		logging.WithError(err).Warn("unable to parse body")
	}
	return c.runtime.ToValue(response{StatusCode: res.StatusCode, Body: string(body)})
}

func (c *HTTP) buildHTTPRequest(args []goja.Value) (req *http.Request, err error) {
	if len(args) > 2 {
		//TODO: thow error too many args
		logging.WithFields("count", len(args)).Error("more than 2 args provided")
	}

	if len(args) < 1 {
		//TODO: thow error no url provided
		logging.Error("no args provided")
	}

	url, ok := args[0].Export().(string)
	if !ok {
		logging.Error("url was not a string")
	}

	config := defaultFetchConfig
	if len(args) == 2 {
		config, err = c.fetchConfigFromArg(args[1].ToObject(c.runtime))
		if err != nil {
			return nil, err
		}
	}

	req, err = http.NewRequest(config.Method, url, config.Body)
	if err != nil {
		return nil, err
	}
	req.Header = config.Headers

	return req, nil
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
