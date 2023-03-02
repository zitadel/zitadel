package domain

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type IDPConfig struct {
	es_models.ObjectRoot
	IDPConfigID  string
	Type         IDPConfigType
	Name         string
	StylingType  IDPConfigStylingType
	State        IDPConfigState
	OIDCConfig   *OIDCIDPConfig
	JWTConfig    *JWTIDPConfig
	AutoRegister bool
}

type IDPConfigView struct {
	AggregateID     string
	IDPConfigID     string
	Name            string
	StylingType     IDPConfigStylingType
	State           IDPConfigState
	CreationDate    time.Time
	ChangeDate      time.Time
	Sequence        uint64
	IDPProviderType IdentityProviderType
	AutoRegister    bool

	IsOIDC                     bool
	OIDCClientID               string
	OIDCClientSecret           *crypto.CryptoValue
	OIDCIssuer                 string
	OIDCScopes                 []string
	OIDCIDPDisplayNameMapping  OIDCMappingField
	OIDCUsernameMapping        OIDCMappingField
	OAuthAuthorizationEndpoint string
	OAuthTokenEndpoint         string

	JWTEndpoint     string
	JWTIssuer       string
	JWTKeysEndpoint string
}

type OIDCIDPConfig struct {
	es_models.ObjectRoot
	IDPConfigID           string
	ClientID              string
	ClientSecret          *crypto.CryptoValue
	ClientSecretString    string
	Issuer                string
	AuthorizationEndpoint string
	TokenEndpoint         string
	Scopes                []string
	IDPDisplayNameMapping OIDCMappingField
	UsernameMapping       OIDCMappingField
}

type JWTIDPConfig struct {
	es_models.ObjectRoot
	IDPConfigID  string
	JWTEndpoint  string
	Issuer       string
	KeysEndpoint string
	HeaderName   string
}

// IDPConfigType
// Deprecated: use [IDPType]
type IDPConfigType int32

const (
	IDPConfigTypeOIDC IDPConfigType = iota
	IDPConfigTypeSAML
	IDPConfigTypeJWT

	//count is for validation
	idpConfigTypeCount
	IDPConfigTypeUnspecified IDPConfigType = -1
)

func (f IDPConfigType) Valid() bool {
	return f >= 0 && f < idpConfigTypeCount
}

// IDPConfigState
// Deprecated: use [IDPStateType]
type IDPConfigState int32

const (
	IDPConfigStateUnspecified IDPConfigState = iota
	IDPConfigStateActive
	IDPConfigStateInactive
	IDPConfigStateRemoved

	idpConfigStateCount
)

func (s IDPConfigState) Valid() bool {
	return s >= 0 && s < idpConfigStateCount
}

func (s IDPConfigState) Exists() bool {
	return s != IDPConfigStateUnspecified && s != IDPConfigStateRemoved
}

// IDPConfigStylingType
// Deprecated: use a concrete provider
type IDPConfigStylingType int32

const (
	IDPConfigStylingTypeUnspecified IDPConfigStylingType = iota
	IDPConfigStylingTypeGoogle

	idpConfigStylingTypeCount
)

func (f IDPConfigStylingType) Valid() bool {
	return f >= 0 && f < idpConfigStylingTypeCount
}

func (st IDPConfigStylingType) GetCSSClass() string {
	switch st {
	case IDPConfigStylingTypeGoogle:
		return "google"
	default:
		return ""
	}
}
