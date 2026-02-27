package saml

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/beevik/etree"
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
func (s *Session) GetAuth(ctx context.Context) (idp.Auth, error) {
	url, err := url.Parse(s.state)
	if err != nil {
		return nil, err
	}
	request := &http.Request{
		URL: url,
	}
	return s.auth(request.WithContext(ctx))
}

// PersistentParameters implements the [idp.Session] interface.
func (s *Session) PersistentParameters() map[string]any {
	return nil
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

	userMapper := NewUser()
	// nameID is required, but at least in ADFS it will not be sent unless explicitly configured
	if s.Assertion.Subject == nil || s.Assertion.Subject.NameID == nil {
		if strings.TrimSpace(s.TransientMappingAttributeName) == "" {
			return nil, zerrors.ThrowInvalidArgument(err, "SAML-EFG32", "Errors.Intent.MissingTransientMappingAttributeName")
		}
		// workaround to use the transient mapping attribute when the subject / nameID are missing (e.g. in ADFS, Shibboleth)
		mappingID, err := s.transientMappingID()
		if err != nil {
			return nil, err
		}
		userMapper.SetID(mappingID)
	} else {
		nameID := s.Assertion.Subject.NameID
		// use the nameID as default mapping id
		userMapper.SetID(nameID.Value)
		if nameID.Format == string(saml.TransientNameIDFormat) {
			mappingID, err := s.transientMappingID()
			if err != nil {
				return nil, err
			}
			userMapper.SetID(mappingID)
		}
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

func (s *Session) ExpiresAt() time.Time {
	if s.Assertion == nil || s.Assertion.Conditions == nil {
		return time.Time{}
	}
	return s.Assertion.Conditions.NotOnOrAfter
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

// auth is a modified copy of the [samlsp.Middleware.HandleStartAuthFlow] method.
// Instead of writing the response to the http.ResponseWriter, it returns the auth request as an [idp.Auth].
// In case of an error, it returns the error directly and does not write to the response.
func (s *Session) auth(r *http.Request) (idp.Auth, error) {
	if r.URL.Path == s.ServiceProvider.ServiceProvider.AcsURL.Path {
		// should never occur, but was handled in the original method, so we keep it here
		return nil, zerrors.ThrowInvalidArgument(nil, "SAML-Eoi24", "don't wrap Middleware with RequireAccount")
	}

	var binding, bindingLocation string
	if s.ServiceProvider.Binding != "" {
		binding = s.ServiceProvider.Binding
		bindingLocation = s.ServiceProvider.ServiceProvider.GetSSOBindingLocation(binding)
	} else {
		binding = saml.HTTPRedirectBinding
		bindingLocation = s.ServiceProvider.ServiceProvider.GetSSOBindingLocation(binding)
		if bindingLocation == "" {
			binding = saml.HTTPPostBinding
			bindingLocation = s.ServiceProvider.ServiceProvider.GetSSOBindingLocation(binding)
		}
	}

	authReq, err := s.ServiceProvider.ServiceProvider.MakeAuthenticationRequest(bindingLocation, binding, s.ServiceProvider.ResponseBinding)
	if err != nil {
		return nil, err
	}
	relayState, err := s.ServiceProvider.RequestTracker.TrackRequest(nil, r, authReq.ID)
	if err != nil {
		return nil, err
	}

	if binding == saml.HTTPRedirectBinding {
		redirectURL, err := authReq.Redirect(relayState, &s.ServiceProvider.ServiceProvider)
		if err != nil {
			return nil, err
		}
		return idp.Redirect(redirectURL.String())
	}
	if binding == saml.HTTPPostBinding {
		doc := etree.NewDocument()
		doc.SetRoot(authReq.Element())
		reqBuf, err := doc.WriteToBytes()
		if err != nil {
			return nil, err
		}
		encodedReqBuf := base64.StdEncoding.EncodeToString(reqBuf)
		return idp.Form(authReq.Destination,
			map[string]string{
				"SAMLRequest": encodedReqBuf,
				"RelayState":  relayState,
			})
	}
	return nil, zerrors.ThrowInvalidArgument(nil, "SAML-Eoi24", "Errors.Intent.Invalid")
}
