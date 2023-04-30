package domain

const (
	OrgDomainPrimaryScope = "urn:zitadel:iam:org:domain:primary:"
	OrgIDScope            = "urn:zitadel:iam:org:id:"
	OrgDomainPrimaryClaim = "urn:zitadel:iam:org:domain:primary"
	OrgIDClaim            = "urn:zitadel:iam:org:id"
	ProjectIDScope        = "urn:zitadel:iam:org:project:id:"
	ProjectIDScopeZITADEL = "zitadel"
	AudSuffix             = ":aud"
	SelectIDPScope        = "urn:zitadel:iam:org:idp:id:"
)

// TODO: Change AuthRequest to interface and let oidcauthreqesut implement it
type Request interface {
	Type() AuthRequestType
	IsValid() bool
}

type AuthRequestType int32

const (
	AuthRequestTypeOIDC AuthRequestType = iota
	AuthRequestTypeSAML
	AuthRequestTypeDevice
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
	ID          string
	BindingType string
	Code        string
	Issuer      string
	IssuerName  string
	Destination string
}

func (a *AuthRequestSAML) Type() AuthRequestType {
	return AuthRequestTypeSAML
}

func (a *AuthRequestSAML) IsValid() bool {
	return true
}

type AuthRequestDevice struct {
	ID         string
	DeviceCode string
	UserCode   string
	Scopes     []string
}

func (*AuthRequestDevice) Type() AuthRequestType {
	return AuthRequestTypeDevice
}

func (a *AuthRequestDevice) IsValid() bool {
	return a.DeviceCode != "" && a.UserCode != "" && len(a.Scopes) > 0
}
