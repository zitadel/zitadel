package idp

type Session interface {
	GetAuthURL() string
	FetchUser() (User, error)
}
