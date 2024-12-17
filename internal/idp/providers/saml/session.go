package saml

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ idp.Session = (*Session)(nil)

// Session is the [idp.Session] implementation for the SAML provider.
type Session struct {
	ServiceProvider               *samlsp.Middleware
	state                         string
	TransientMappingAttributeName string

	RequestID string
	Request   *http.Request

	Assertion *saml.Assertion
}

func NewSession(provider *Provider, requestID string, request *http.Request) (*Session, error) {
	sp, err := provider.GetSP()
	if err != nil {
		return nil, err
	}
	return &Session{
		ServiceProvider:               sp,
		TransientMappingAttributeName: provider.TransientMappingAttributeName(),
		RequestID:                     requestID,
		Request:                       request,
	}, nil
}

// GetAuth implements the [idp.Session] interface.
func (s *Session) GetAuth(ctx context.Context) (string, bool) {
	url, _ := url.Parse(s.state)
	resp := NewTempResponseWriter()

	request := &http.Request{
		URL: url,
	}
	s.ServiceProvider.HandleStartAuthFlow(
		resp,
		request.WithContext(ctx),
	)

	if location := resp.Header().Get("Location"); location != "" {
		return idp.Redirect(location)
	}
	return idp.Form(resp.content.String())
}

// FetchUser implements the [idp.Session] interface.
func (s *Session) FetchUser(ctx context.Context) (user idp.User, err error) {
	if s.RequestID == "" || s.Request == nil {
		return nil, zerrors.ThrowInvalidArgument(nil, "SAML-d09hy0wkex", "Errors.Intent.ResponseInvalid")
	}

	s.Assertion, err = s.ServiceProvider.ServiceProvider.ParseResponse(s.Request, []string{s.RequestID})
	if err != nil {
		invalidRespErr := new(saml.InvalidResponseError)
		if errors.As(err, &invalidRespErr) {
			return nil, zerrors.ThrowInvalidArgument(invalidRespErr.PrivateErr, "SAML-ajl3irfs", "Errors.Intent.ResponseInvalid")
		}
		return nil, zerrors.ThrowInvalidArgument(err, "SAML-nuo0vphhh9", "Errors.Intent.ResponseInvalid")
	}

	// nameID is required, but at least in ADFS it will not be sent unless explicitly configured
	if s.Assertion.Subject == nil || s.Assertion.Subject.NameID == nil {
		return nil, zerrors.ThrowInvalidArgument(err, "SAML-EFG32", "Errors.Intent.ResponseInvalid")
	}
	nameID := s.Assertion.Subject.NameID
	userMapper := NewUser()
	// use the nameID as default mapping id
	userMapper.SetID(nameID.Value)
	if nameID.Format == string(saml.TransientNameIDFormat) {
		mappingID, err := s.transientMappingID()
		if err != nil {
			return nil, err
		}
		userMapper.SetID(mappingID)
	}
	for _, statement := range s.Assertion.AttributeStatements {
		for _, attribute := range statement.Attributes {
			values := make([]string, len(attribute.Values))
			for i := range attribute.Values {
				values[i] = attribute.Values[i].Value
			}
			userMapper.Attributes[attribute.Name] = values
		}
	}
	return userMapper, nil
}

func (s *Session) transientMappingID() (string, error) {
	for _, statement := range s.Assertion.AttributeStatements {
		for _, attribute := range statement.Attributes {
			if attribute.Name != s.TransientMappingAttributeName {
				continue
			}
			if len(attribute.Values) != 1 {
				return "", zerrors.ThrowInvalidArgument(nil, "SAML-Soij4", "Errors.Intent.MissingSingleMappingAttribute")
			}
			return attribute.Values[0].Value, nil
		}
	}
	return "", zerrors.ThrowInvalidArgument(nil, "SAML-swwg2", "Errors.Intent.MissingSingleMappingAttribute")
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

func (w *TempResponseWriter) WriteHeader(statusCode int) {}

func NewTempResponseWriter() *TempResponseWriter {
	return &TempResponseWriter{
		header:  map[string][]string{},
		content: bytes.NewBuffer([]byte{}),
	}
}
