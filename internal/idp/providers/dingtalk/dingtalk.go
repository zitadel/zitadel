package dingtalk

import (
	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

const (
	authURL    = "https://login.dingtalk.com/oauth2/auth"
	tokenURL   = "https://api.dingtalk.com/v1.0/oauth2/userAccessToken"
	profileURL = "https://api.dingtalk.com/v1.0/contact/users/me"
	name       = "DingTalk"
)

var _ idp.Provider = (*Provider)(nil)

// New creates a DingTalk provider using the [oauth.Provider] (OAuth 2.0 generic provider)
func New(clientID, secret, callbackURL string, scopes []string, options ...oauth.ProviderOpts) (*Provider, error) {
	return NewCustomURL(name, clientID, secret, callbackURL, authURL, tokenURL, profileURL, scopes, options...)
}

// NewCustomURL creates a DingTalk provider using the [oauth.Provider] (OAuth 2.0 generic provider)
// with custom endpoints
func NewCustomURL(name, clientID, secret, callbackURL, authURL, tokenURL, profileURL string, scopes []string, options ...oauth.ProviderOpts) (*Provider, error) {
	rp, err := oauth.New(
		newConfig(clientID, secret, callbackURL, authURL, tokenURL, scopes),
		name,
		profileURL,
		func() idp.User {
			return new(User)
		},
		options...,
	)
	if err != nil {
		return nil, err
	}
	return &Provider{
		Provider: rp,
	}, nil
}

// Provider is the [idp.Provider] implementation for DingTalk
type Provider struct {
	*oauth.Provider
}

func newConfig(clientID, secret, callbackURL, authURL, tokenURL string, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		RedirectURL:  callbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: scopes,
	}

	return c
}

// User is a representation of the authenticated DingTalk user and implements the [idp.User] interface
// https://developers.dingtalk.com/document/app/obtain-user-information-based-on-sns-temporary-authorization
type User struct {
	UnionID   string              `json:"unionId"`
	OpenID    string              `json:"openId"`
	Nick      string              `json:"nick"`
	AvatarURL string              `json:"avatarUrl"`
	Mobile    string              `json:"mobile"`
	StateCode string              `json:"stateCode"`
	Email     domain.EmailAddress `json:"email"`
	DingID    string              `json:"dingId"`
}

// GetID is an implementation of the [idp.User] interface.
func (u *User) GetID() string {
	return u.UnionID
}

// GetFirstName is an implementation of the [idp.User] interface.
// It returns an empty string because DingTalk does not provide separate first/last names.
func (u *User) GetFirstName() string {
	return ""
}

// GetLastName is an implementation of the [idp.User] interface.
// It returns an empty string because DingTalk does not provide separate first/last names.
func (u *User) GetLastName() string {
	return ""
}

// GetDisplayName is an implementation of the [idp.User] interface.
func (u *User) GetDisplayName() string {
	return u.Nick
}

// GetNickname is an implementation of the [idp.User] interface
// returning the nick name of the DingTalk user.
func (u *User) GetNickname() string {
	return u.Nick
}

// GetPreferredUsername is an implementation of the [idp.User] interface
// returning the nick name of the DingTalk user.
func (u *User) GetPreferredUsername() string {
	return u.Nick
}

// GetEmail is an implementation of the [idp.User] interface.
func (u *User) GetEmail() domain.EmailAddress {
	return u.Email
}

// IsEmailVerified is an implementation of the [idp.User] interface.
// It returns true because DingTalk validates emails themselves.
func (u *User) IsEmailVerified() bool {
	return true
}

// GetPhone is an implementation of the [idp.User] interface.
func (u *User) GetPhone() domain.PhoneNumber {
	return domain.PhoneNumber(u.Mobile)
}

// IsPhoneVerified is an implementation of the [idp.User] interface
// it returns true because DingTalk validates phone numbers themselves
func (u *User) IsPhoneVerified() bool {
	return u.Mobile != ""
}

// GetPreferredLanguage is an implementation of the [idp.User] interface.
// It returns [language.Chinese] as the default language for DingTalk users.
func (u *User) GetPreferredLanguage() language.Tag {
	return language.Chinese
}

// GetProfile is an implementation of the [idp.User] interface.
// DingTalk does not provide a profile URL, so we return an empty string.
func (u *User) GetProfile() string {
	return ""
}

// GetAvatarURL is an implementation of the [idp.User] interface.
func (u *User) GetAvatarURL() string {
	return u.AvatarURL
}
