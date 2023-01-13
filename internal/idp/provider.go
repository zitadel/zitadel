package idp

import "golang.org/x/text/language"

type Provider interface {
	Name() string
	//SetName(name string)
	BeginAuth(state string) (Session, error)
	//Authorize(session Session) (Session, error)
	//UnmarshalSession(string) (Session, error)
	FetchUser(Session) (User, error)
	//Debug(bool)
	//RefreshToken(refreshToken string) (*oauth2.Token, error) // Get new access token based on the refresh token
	//RefreshTokenAvailable() bool                             // Refresh token is provided by auth provider or not
}

type User struct {
	ID                string
	FirstName         string
	LastName          string
	DisplayName       string
	NickName          string
	PreferredUsername string
	Email             string
	IsEmailVerified   bool
	Phone             string
	IsPhoneVerified   bool
	PreferredLanguage language.Tag
	AvatarURL         string
	Profile           string
	RawData           any
}

/*
login

click on provider (e.g. 1 = google)
getprovider(1) -> return google template

provider.beginauth(state) -> return session (@max does this make sense)

session.authurl()
redirect to url

handle callback
	getprovider(1) -> return google template

	getsession(state) -> session
	getprovider(session.provider) -> return google template
	provider.FetchUser(session) -> user / error
		session.authorize(provider)
		provider.FetchUser(session) -> user / error



mgmt

getprovider(1) -> return google template
	type google
	client id
	client secret


*/
