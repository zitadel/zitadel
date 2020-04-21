package model

type Request interface {
	Type() AuthSessionType
	IsValid() bool
}

type AuthSessionType int32

const (
	AuthSessionTypeOIDC AuthSessionType = iota
	AuthSessionTypeSAML
)

type AuthSessionOIDC struct {
	Scopes        []string
	ResponseType  OIDCResponseType
	Nonce         string
	CodeChallenge *OIDCCodeChallenge
}

func (a *AuthSessionOIDC) Type() AuthSessionType {
	return AuthSessionTypeOIDC
}

func (a *AuthSessionOIDC) IsValid() bool {
	return true
}

type AuthSessionSAML struct {
}

func (a *AuthSessionSAML) Type() AuthSessionType {
	return AuthSessionTypeSAML
}

func (a *AuthSessionSAML) IsValid() bool {
	return true
}
