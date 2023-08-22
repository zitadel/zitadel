package saml

import (
	"bytes"
	"context"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Session = (*Session)(nil)

// Session is the [idp.Session] implementation for the SAML provider.
type Session struct {
	serviceProvider *samlsp.Middleware
	state           string

	RequestID string
	Request   *http.Request
}

// GetAuthURL implements the [idp.Session] interface.
func (s *Session) GetAuth() (http.Header, []byte) {
	url, _ := url.Parse(s.state)
	resp := NewTempResponseWriter()

	s.serviceProvider.HandleStartAuthFlow(
		resp,
		&http.Request{
			URL: url,
		},
	)
	return resp.header, resp.content.Bytes()
}

// FetchUser implements the [idp.Session] interface.
func (s *Session) FetchUser(ctx context.Context) (user idp.User, err error) {
	if s.RequestID == "" && s.Request == nil {
		return nil, errors.ThrowInvalidArgument(nil, "SAML-tzb2sj", "Errors.*")
	}

	assertion, err := s.serviceProvider.ServiceProvider.ParseResponse(s.Request, []string{s.RequestID})
	if err != nil {
		return nil, err
	}

	userMapper := &UserMapper{}
	userMapper.SetID(assertion.Subject.NameID)
	return userMapper, nil
}

type TempResponseWriter struct {
	header  http.Header
	content *bytes.Buffer
}

func (w *TempResponseWriter) Header() http.Header {
	return w.header
}

func (w *TempResponseWriter) Write(content []byte) (int, error) {
	return w.content.Write(content)
}

func (w *TempResponseWriter) WriteHeader(statusCode int) {
	return
}

func NewTempResponseWriter() *TempResponseWriter {
	return &TempResponseWriter{
		header:  map[string][]string{},
		content: bytes.NewBuffer([]byte{}),
	}
}
