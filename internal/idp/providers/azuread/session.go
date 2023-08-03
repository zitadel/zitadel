package azuread

import (
	"net/http"

	httphelper "github.com/zitadel/oidc/v2/pkg/http"
	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

// Session extends the [oauth.Session] to extend it with the [idp.SessionSupportsMigration] functionality
type Session struct {
	*oauth.Session
}

// RetrievePreviousID implements the [idp.SessionSupportsMigration] interface by returning the `sub` from the userinfo endpoint
func (s *Session) RetrievePreviousID() (string, error) {
	req, err := http.NewRequest("GET", userinfoEndpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("authorization", s.Tokens.TokenType+" "+s.Tokens.AccessToken)
	userinfo := new(oidc.UserInfo)
	if err := httphelper.HttpRequest(s.Provider.HttpClient(), req, &userinfo); err != nil {
		return "", err
	}
	return userinfo.Subject, nil
}
