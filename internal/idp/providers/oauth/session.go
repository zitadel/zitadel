package oauth

import (
	"context"
	"net/http"

	"github.com/zitadel/oidc/v2/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v2/pkg/http"
	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Session = (*Session)(nil)

// Session is the idp.Session implementation for the OAuth2.0 provider
type Session struct {
	AuthURL string
	Code    string
	Tokens  *oidc.Tokens

	Provider *Provider
}

func (s *Session) GetAuthURL() string {
	return s.AuthURL
}

func (s *Session) FetchUser() (user idp.User, err error) {
	if s.Tokens == nil {
		if err = s.authorize(); err != nil {
			return idp.User{}, err
		}
	}
	req, err := http.NewRequest("GET", s.Provider.userEndpoint, nil)
	if err != nil {
		return idp.User{}, err
	}
	req.Header.Set("authorization", s.Tokens.TokenType+" "+s.Tokens.AccessToken)
	mapper := s.Provider.userMapper()
	if err := httphelper.HttpRequest(s.Provider.RelyingParty.HttpClient(), req, &mapper); err != nil {
		return idp.User{}, err
	}
	err = mapUser(mapper, &user)
	return user, err
}

func (s *Session) authorize() error {
	if s.Code == "" {
		return ErrCodeMissing
	}
	tokens, err := rp.CodeExchange(context.TODO(), s.Code, s.Provider.RelyingParty)
	if err != nil {
		return err
	}
	s.Tokens = tokens
	return nil
}

func mapUser(mapper UserInfoMapper, user *idp.User) error {
	user.ID = mapper.GetID()
	user.FirstName = mapper.GetFirstName()
	user.LastName = mapper.GetLastName()
	user.DisplayName = mapper.GetDisplayName()
	user.NickName = mapper.GetNickName()
	user.PreferredUsername = mapper.GetPreferredUsername()
	user.Email = mapper.GetEmail()
	user.IsEmailVerified = mapper.IsEmailVerified()
	user.Phone = mapper.GetPhone()
	user.IsPhoneVerified = mapper.IsPhoneVerified()
	user.PreferredLanguage = mapper.GetPreferredLanguage()
	user.AvatarURL = mapper.GetAvatarURL()
	user.Profile = mapper.GetProfile()
	user.RawData = mapper.RawData()
	return nil
}
