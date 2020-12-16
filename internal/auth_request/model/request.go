package model

type Request interface {
	Type() AuthRequestType
	IsValid() bool
}

type AuthRequestType int32

var (
	authRequestTypeMapping = map[AuthRequestType]Request{
		AuthRequestTypeOIDC: &AuthRequestOIDC{},
	}
)

const (
	AuthRequestTypeOIDC AuthRequestType = iota
	AuthRequestTypeSAML
)

const (
	OrgDomainPrimaryScope = "urn:zitadel:iam:org:domain:primary:"
	OrgDomainPrimaryClaim = "urn:zitadel:iam:org:domain:primary"
	ProjectIDScope        = "urn:zitadel:iam:org:project:id:"
	AudSuffix             = ":aud"
)

type AuthRequestOIDC struct {
	Scopes        []string
	ResponseType  OIDCResponseType
	Nonce         string
	CodeChallenge *OIDCCodeChallenge
}

func (a *AuthRequestOIDC) Type() AuthRequestType {
	return AuthRequestTypeOIDC
}

func (a *AuthRequestOIDC) IsValid() bool {
	return len(a.Scopes) > 0 &&
		a.CodeChallenge == nil || a.CodeChallenge != nil && a.CodeChallenge.IsValid()
}

type AuthRequestSAML struct {
}

func (a *AuthRequestSAML) Type() AuthRequestType {
	return AuthRequestTypeSAML
}

func (a *AuthRequestSAML) IsValid() bool {
	return true
}

type OIDCResponseType int32

const (
	OIDCResponseTypeCode OIDCResponseType = iota
	OIDCResponseTypeIdToken
	OIDCResponseTypeIdTokenToken
)
