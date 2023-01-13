package oauth

import (
	"context"
	"errors"
	"net/http"

	"github.com/zitadel/oidc/v2/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v2/pkg/http"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Provider = (*Provider)(nil)

var ErrCodeMissing = errors.New("no auth code provided")

type Provider struct {
	name string
	rp.RelyingParty
	userEndpoint string
	userMapper   func() UserInfoMapper
}

func New(config *oauth2.Config, userEndpoint string, userMapper func() UserInfoMapper) (*Provider, error) {
	provider := &Provider{
		userEndpoint: userEndpoint,
		userMapper:   userMapper,
	}

	relyingParty, err := rp.NewRelyingPartyOAuth(config)
	if err != nil {
		return nil, err
	}
	provider.RelyingParty = relyingParty
	return provider, nil
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) BeginAuth(state string) (idp.Session, error) {
	url := rp.AuthURL(state, p.RelyingParty)
	return &Session{AuthURL: url}, nil
}

func (p *Provider) FetchUser(session idp.Session) (user idp.User, err error) {
	oauthSession, _ := session.(*Session)
	if oauthSession.Tokens == nil {
		if err = p.authorize(oauthSession); err != nil {
			return idp.User{}, err
		}
	}
	req, err := http.NewRequest("GET", p.userEndpoint, nil)
	if err != nil {
		return idp.User{}, err
	}
	req.Header.Set("authorization", oauthSession.Tokens.TokenType+" "+oauthSession.Tokens.AccessToken)
	mapper := p.userMapper()
	if err := httphelper.HttpRequest(p.RelyingParty.HttpClient(), req, &mapper); err != nil {
		return idp.User{}, err
	}
	err = mapClaims(mapper, &user)
	return user, err
}

func (p *Provider) authorize(session *Session) error {
	if session.Code == "" {
		return ErrCodeMissing
	}
	tokens, err := rp.CodeExchange(context.TODO(), session.Code, p.RelyingParty)
	if err != nil {
		return err
	}
	session.Tokens = tokens
	return nil
}

func mapClaims(mapper UserInfoMapper, user *idp.User) error {
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
	user.PreferredLanguage = mapper.GetPreferredLanguange()
	user.AvatarURL = mapper.GetAvatarURL()
	user.Profile = mapper.GetProfile()
	user.RawData = mapper.RawData()
	return nil
}
