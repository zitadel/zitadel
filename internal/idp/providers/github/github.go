package github

import (
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

const (
	authURL    = "https://github.com/login/oauth/authorize"
	tokenURL   = "https://github.com/login/oauth/access_token"
	profileURL = "https://api.github.com/user"
	name       = "GitHub"
)

var _ idp.Provider = (*Provider)(nil)

func New(clientID, secret, callbackURL string, scopes []string, options ...oauth.ProviderOpts) (*Provider, error) {
	return NewCustomURL(name, clientID, secret, callbackURL, authURL, tokenURL, profileURL, scopes, options...)
}

func NewCustomURL(name, clientID, secret, callbackURL, authURL, tokenURL, profileURL string, scopes []string, options ...oauth.ProviderOpts) (*Provider, error) {
	rp, err := oauth.New(
		newConfig(clientID, secret, callbackURL, authURL, tokenURL, scopes),
		name,
		profileURL,
		func() oauth.UserInfoMapper {
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

// Provider is the idp.Provider implementation for GitHub
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

// User is a representation of the authenticated GH user
// https://docs.github.com/en/rest/users/users?apiVersion=2022-11-28#get-the-authenticated-user
type User struct {
	Login                   string      `json:"login"`
	ID                      int         `json:"id"`
	NodeId                  string      `json:"node_id"`
	AvatarUrl               string      `json:"avatar_url"`
	GravatarId              string      `json:"gravatar_id"`
	Url                     string      `json:"url"`
	HtmlUrl                 string      `json:"html_url"`
	FollowersUrl            string      `json:"followers_url"`
	FollowingUrl            string      `json:"following_url"`
	GistsUrl                string      `json:"gists_url"`
	StarredUrl              string      `json:"starred_url"`
	SubscriptionsUrl        string      `json:"subscriptions_url"`
	OrganizationsUrl        string      `json:"organizations_url"`
	ReposUrl                string      `json:"repos_url"`
	EventsUrl               string      `json:"events_url"`
	ReceivedEventsUrl       string      `json:"received_events_url"`
	Type                    string      `json:"type"`
	SiteAdmin               bool        `json:"site_admin"`
	Name                    string      `json:"name"`
	Company                 string      `json:"company"`
	Blog                    string      `json:"blog"`
	Location                string      `json:"location"`
	Email                   string      `json:"email"`
	Hireable                interface{} `json:"hireable"`
	Bio                     string      `json:"bio"`
	TwitterUsername         string      `json:"twitter_username"`
	PublicRepos             int         `json:"public_repos"`
	PublicGists             int         `json:"public_gists"`
	Followers               int         `json:"followers"`
	Following               int         `json:"following"`
	CreatedAt               time.Time   `json:"created_at"`
	UpdatedAt               time.Time   `json:"updated_at"`
	PrivateGists            int         `json:"private_gists"`
	TotalPrivateRepos       int         `json:"total_private_repos"`
	OwnedPrivateRepos       int         `json:"owned_private_repos"`
	DiskUsage               int         `json:"disk_usage"`
	Collaborators           int         `json:"collaborators"`
	TwoFactorAuthentication bool        `json:"two_factor_authentication"`
	Plan                    struct {
		Name          string `json:"name"`
		Space         int    `json:"space"`
		PrivateRepos  int    `json:"private_repos"`
		Collaborators int    `json:"collaborators"`
	} `json:"plan"`
}

func (u *User) GetPreferredUsername() string {
	return u.Login
}

func (u *User) GetPhone() string {
	// GitHub does not provide the user's phone
	return ""
}

func (u *User) IsPhoneVerified() bool {
	// GitHub does not provide the user's phone
	return false
}

func (u *User) GetPreferredLanguage() language.Tag {
	// GitHub does not provide the user's language
	return language.Und
}

func (u *User) GetProfile() string {
	return u.HtmlUrl
}

func (u *User) IsEmailVerified() bool {
	// GitHub validates emails themself
	return true
}

func (u *User) GetID() string {
	return strconv.Itoa(u.ID)
}

func (u *User) GetDisplayName() string {
	return u.Name
}

func (u *User) GetNickName() string {
	return u.Login
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetAvatarURL() string {
	return u.AvatarUrl
}

func (u *User) GetFirstName() string {
	// GitHub does not provide the user's firstname
	return ""
}

func (u *User) GetLastName() string {
	// GitHub does not provide the user's lastname
	return ""
}

func (u *User) RawData() any {
	return u
}
