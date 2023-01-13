package github

import (
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

const (
	AuthURL    = "https://github.com/login/oauth/authorize"
	TokenURL   = "https://github.com/login/oauth/access_token"
	ProfileURL = "https://api.github.com/user"
)

func New(clientID, secret, callbackURL string, scopes ...string) (*Provider, error) {
	provider, err := oauth.New(
		newConfig(clientID, secret, callbackURL, AuthURL, TokenURL, scopes),
		ProfileURL,
		func() oauth.UserInfoMapper {
			return new(User)
		},
	)
	if err != nil {
		return nil, err
	}
	p := &Provider{
		Provider: provider,
	}
	return p, nil
}

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

type User struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	//NodeId            string      `json:"node_id"`
	AvatarUrl  string `json:"avatar_url"`
	GravatarId string `json:"gravatar_id"`
	//Url        string `json:"url"`
	HtmlUrl string `json:"html_url"`
	//FollowersUrl      string      `json:"followers_url"`
	//FollowingUrl      string      `json:"following_url"`
	//GistsUrl          string      `json:"gists_url"`
	//StarredUrl        string      `json:"starred_url"`
	//SubscriptionsUrl  string      `json:"subscriptions_url"`
	//OrganizationsUrl  string      `json:"organizations_url"`
	//ReposUrl          string      `json:"repos_url"`
	//EventsUrl         string      `json:"events_url"`
	//ReceivedEventsUrl string      `json:"received_events_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
	Name      string `json:"name"`
	//Company           string      `json:"company"`
	//Blog              string      `json:"blog"`
	//Location          string      `json:"location"`
	Email string `json:"email"`
	//Hireable          interface{} `json:"hireable"`
	//Bio               string      `json:"bio"`
	//TwitterUsername   string      `json:"twitter_username"`
	//PublicRepos       int         `json:"public_repos"`
	//PublicGists       int         `json:"public_gists"`
	//Followers         int         `json:"followers"`
	//Following         int         `json:"following"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) GetPreferredUsername() string {
	return u.Login
}

func (u *User) GetPhone() string {
	return ""
}

func (u *User) IsPhoneVerified() bool {
	return false
}

func (u *User) GetPreferredLanguange() language.Tag {
	return language.Und
}

func (u *User) GetProfile() string {
	return u.HtmlUrl
}

func (u *User) IsEmailVerified() bool {
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
	return ""
}

func (u *User) GetLastName() string {
	return ""
}

func (u *User) RawData() any {
	return u
}
