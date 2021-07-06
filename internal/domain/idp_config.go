package domain

import (
	"time"

	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

type IDPConfig interface {
	ObjectDetails() es_models.ObjectRoot
	ID() string
	IDPConfigType() IDPConfigType
	IDPConfigName() string
	IDPConfigStylingType() IDPConfigStylingType
	IDPConfigState() IDPConfigState
	//Details() *ObjectDetails
}

type CommonIDPConfig struct {
	es_models.ObjectRoot
	IDPConfigID     string
	Type            IDPConfigType
	Name            string
	StylingType     IDPConfigStylingType
	State           IDPConfigState
	IDPProviderType IdentityProviderType
	//OIDCConfig  *OIDCIDPConfig
}

func (c *CommonIDPConfig) ObjectDetails() es_models.ObjectRoot {
	return c.ObjectRoot
}

func (c *CommonIDPConfig) ID() string {
	return c.IDPConfigID
}

func (c *CommonIDPConfig) IDPConfigType() IDPConfigType {
	return c.Type
}

func (c *CommonIDPConfig) IDPConfigName() string {
	return c.Name
}

func (c *CommonIDPConfig) IDPConfigStylingType() IDPConfigStylingType {
	return c.StylingType
}

func (c *CommonIDPConfig) IDPConfigState() IDPConfigState {
	return c.State
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

	IsOIDC                     bool
	OIDCClientID               string
	OIDCClientSecret           *crypto.CryptoValue
	OIDCIssuer                 string
	OIDCScopes                 []string
	OIDCIDPDisplayNameMapping  OIDCMappingField
	OIDCUsernameMapping        OIDCMappingField
	OAuthAuthorizationEndpoint string
	OAuthTokenEndpoint         string
}

type OIDCIDPConfig struct {
	es_models.ObjectRoot
	CommonIDPConfig
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

func (c *OIDCIDPConfig) ID() string {
	return c.AggregateID
}

func (c *OIDCIDPConfig) Details() *ObjectDetails {
	return &ObjectDetails{
		Sequence:      c.Sequence,
		EventDate:     c.ChangeDate,
		ResourceOwner: c.ResourceOwner,
	}
}

type AuthConnectorIDPConfig struct {
	es_models.ObjectRoot
	CommonIDPConfig
	BaseURL     string
	ProviderID  string
	MachineID   string
	MachineName string
}

func (c *AuthConnectorIDPConfig) ID() string {
	return c.AggregateID
}

func (c *AuthConnectorIDPConfig) Details() *ObjectDetails {
	return &ObjectDetails{
		Sequence:      c.Sequence,
		EventDate:     c.ChangeDate,
		ResourceOwner: c.ResourceOwner,
	}
}

type IDPConfigType int32

const (
	IDPConfigTypeOIDC IDPConfigType = iota
	IDPConfigTypeSAML
	IDPConfigTypeAuthConnector

	//count is for validation
	idpConfigTypeCount
)

func (f IDPConfigType) Valid() bool {
	return f >= 0 && f < idpConfigTypeCount
}

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
	return s != IDPConfigStateUnspecified || s == IDPConfigStateRemoved
}

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
