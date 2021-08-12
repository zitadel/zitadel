package twilio

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/providers"
	"github.com/kevinburke/twilio-go"
	"golang.org/x/net/http/httpproxy"
)

type Twilio struct {
	client *twilio.Client
}

func InitTwilioProvider(config TwilioConfig) *Twilio {
	httpClient := http.DefaultClient
	if config.Proxy != nil {
		httpClient.Transport = transport(config.Proxy)
	}
	return &Twilio{
		client: twilio.NewClient(config.SID, config.Token, httpClient),
	}
}

func (t *Twilio) CanHandleMessage(message providers.Message) bool {
	twilioMsg, ok := message.(*TwilioMessage)
	if !ok {
		return false
	}
	return twilioMsg.Content != "" && twilioMsg.RecipientPhoneNumber != "" && twilioMsg.SenderPhoneNumber != ""
}

func (t *Twilio) HandleMessage(message providers.Message) error {
	twilioMsg, ok := message.(*TwilioMessage)
	if !ok {
		return caos_errs.ThrowInternal(nil, "TWILI-s0pLc", "message is not TwilioMessage")
	}
	m, err := t.client.Messages.SendMessage(twilioMsg.SenderPhoneNumber, twilioMsg.RecipientPhoneNumber, twilioMsg.GetContent(), nil)
	if err != nil {
		return caos_errs.ThrowInternal(err, "TWILI-osk3S", "could not send message")
	}
	logging.LogWithFields("SMS_-f335c523", "message_sid", m.Sid, "status", m.Status).Debug("sms sent")
	return nil
}

func transport(proxy *Proxy) http.RoundTripper {
	tr := &http.Transport{
		Proxy: proxyFromConfig(proxy),
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if proxy.CertPath != "" {
		rootCAs := mustGetSystemCertPool()
		data, err := ioutil.ReadFile(proxy.CertPath)
		if err == nil {
			rootCAs.AppendCertsFromPEM(data)
		}
		tr.TLSClientConfig = &tls.Config{
			RootCAs: rootCAs,
		}
	}
	return tr
}

func proxyFromConfig(proxy *Proxy) func(req *http.Request) (*url.URL, error) {
	conf := &httpproxy.Config{HTTPProxy: proxy.HTTP, HTTPSProxy: proxy.HTTPS}
	return func(req *http.Request) (*url.URL, error) {
		return conf.ProxyFunc()(req.URL)
	}
}

func mustGetSystemCertPool() *x509.CertPool {
	pool, err := x509.SystemCertPool()
	if err != nil {
		return x509.NewCertPool()
	}
	return pool
}
