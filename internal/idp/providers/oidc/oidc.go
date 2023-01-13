package oidc

import (
	"context"
	"errors"

	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Provider = (*Provider)(nil)

var ErrCodeMissing = errors.New("no auth code provided")

type Provider struct {
	name string
	rp.RelyingParty
}

func New(issuer, clientID, clientSecret, redirectURI string) (*Provider, error) {
	provider := &Provider{}
	relyingParty, err := rp.NewRelyingPartyOIDC(issuer, clientID, clientSecret, redirectURI, []string{oidc.ScopeOpenID})
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
	oidcSession, _ := session.(*Session)
	if oidcSession.Tokens == nil {
		if err = p.authorize(oidcSession); err != nil {
			return idp.User{}, err
		}
	}
	info, err := rp.Userinfo(
		oidcSession.Tokens.AccessToken,
		oidcSession.Tokens.TokenType,
		oidcSession.Tokens.IDTokenClaims.GetSubject(),
		p.RelyingParty,
	)
	if err != nil {
		return idp.User{}, err
	}
	userFromClaims(info, &user)
	return user, nil
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

func userFromClaims(info oidc.UserInfo, user *idp.User) {
	user.ID = info.GetSubject()
	user.FirstName = info.GetGivenName()
	user.LastName = info.GetFamilyName()
	user.DisplayName = info.GetName()
	user.NickName = info.GetNickname()
	user.PreferredUsername = info.GetPreferredUsername()
	user.Email = info.GetEmail()
	user.IsEmailVerified = info.IsEmailVerified()
	user.Phone = info.GetPhoneNumber()
	user.IsPhoneVerified = info.IsPhoneNumberVerified()
	user.PreferredLanguage = info.GetLocale()
	user.AvatarURL = info.GetPicture()
	user.Profile = info.GetProfile()
}
