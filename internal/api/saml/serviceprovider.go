package saml

import (
	"strings"

	"github.com/zitadel/saml/pkg/provider/serviceprovider"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
)

const (
	LoginSamlRequestParam = "samlRequest"
	LoginPath             = "/login"
)

type ServiceProvider struct {
	SP                *query.SAMLServiceProvider
	defaultLoginURL   string
	defaultLoginURLV2 string
}

func ServiceProviderFromBusiness(spQuery *query.SAMLServiceProvider, defaultLoginURL, defaultLoginURLV2 string) (*serviceprovider.ServiceProvider, error) {
	sp := &ServiceProvider{
		SP:                spQuery,
		defaultLoginURL:   defaultLoginURL,
		defaultLoginURLV2: defaultLoginURLV2,
	}

	return serviceprovider.NewServiceProvider(
		spQuery.AppID,
		&serviceprovider.Config{Metadata: spQuery.Metadata},
		sp.LoginURL,
	)
}

func (s *ServiceProvider) LoginURL(id string) string {
	// if the authRequest does not have the v2 prefix, it was created for login V1
	if !strings.HasPrefix(id, command.IDPrefixV2) {
		return s.defaultLoginURL + id
	}
	// any v2 login without a specific base uri will be sent to the configured login v2 UI
	// this way we're also backwards compatible
	if s.SP.LoginBaseURI == nil || s.SP.LoginBaseURI.String() == "" {
		return s.defaultLoginURLV2 + id
	}
	// for clients with a specific URI (internal or external) we only need to add the auth request id
	uri := s.SP.LoginBaseURI.JoinPath(LoginPath)
	q := uri.Query()
	q.Set(LoginSamlRequestParam, id)
	uri.RawQuery = q.Encode()
	return uri.String()
}
