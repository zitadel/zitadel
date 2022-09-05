package http

import (
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/http/httpproxy"
)

type ProxyConfig struct {
	HTTP    string
	HTTPS   string
	NoProxy []string
}

func ClientWithProxy(config ProxyConfig) *http.Client {
	client := http.DefaultClient
	if config.HTTP == "" &&
		config.HTTPS == "" &&
		len(config.NoProxy) == 0 {
		return client
	}
	proxy := httpproxy.Config{
		HTTPProxy:  config.HTTP,
		HTTPSProxy: config.HTTPS,
		NoProxy:    strings.Join(config.NoProxy, ","),
	}
	client.Transport = &http.Transport{
		Proxy: func(r *http.Request) (*url.URL, error) {
			return proxy.ProxyFunc()(r.URL)
		},
	}
	return client
}
