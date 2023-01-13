package idp

type Session interface {
	GetAuthURL() string
}
