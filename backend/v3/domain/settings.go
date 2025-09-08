package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/domain"
)

//go:generate enumer -type SettingType -transform snake -trimprefix SettingType -sql
type SettingType uint8

const (
	SettingTypeUnspecified SettingType = iota
	SettingTypeLogin
	SettingTypePasswordComplexity
	SettingTypePasswordExpiry
	SettingTypeBranding
	SettingTypeDomain
	SettingTypeLegalAndSupport
	SettingTypeLockout
	SettingTypeGeneral
	SettingTypeSecurity
)

type Setting struct {
	ID         string          `json:"id,omitempty" db:"id"`
	InstanceID string          `json:"instanceId,omitempty" db:"instance_id"`
	OrgID      *string         `json:"orgId,omitempty" db:"org_id"`
	Type       SettingType     `json:"type,omitempty" db:"type"`
	Settings   json.RawMessage `json:"settings,omitempty" db:"settings"`
	CreatedAt  time.Time       `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt  time.Time       `json:"updatedAt,omitzero" db:"updated_at"`
}

type LoginSetting struct {
	AllowUserNamePassword      *bool                    `json:"allowUsernamePassword,omitempty"`
	AllowRegister              *bool                    `json:"allowRegister,omitempty"`
	AllowExternalSetting       *bool                    `json:"allowExternalIdp,omitempty"`
	ForceMFA                   *bool                    `json:"forceMFA,omitempty"`
	ForceMFALocalOnly          *bool                    `json:"forceMFALocalOnly,omitempty"`
	HidePasswordReset          *bool                    `json:"hidePasswordReset,omitempty"`
	IgnoreUnknownUsernames     *bool                    `json:"ignoreUnknownUsernames,omitempty"`
	AllowDomainDiscovery       *bool                    `json:"allowDomainDiscovery,omitempty"`
	DisableLoginWithEmail      *bool                    `json:"disableLoginWithEmail,omitempty"`
	DisableLoginWithPhone      *bool                    `json:"disableLoginWithPhone,omitempty"`
	PasswordlessType           *domain.PasswordlessType `json:"passwordlessType,omitempty"`
	IsDefault                  *bool                    `json:"isDefault,omitempty"`
	DefaultRedirectURI         *string                  `json:"defaultRedirectURI,omitempty"`
	PasswordCheckLifetime      *time.Duration           `json:"passwordCheckLifetime,omitempty"`
	ExternalLoginCheckLifetime *time.Duration           `json:"externalLoginCheckLifetime,omitempty"`
	MFAInitSkipLifetime        *time.Duration           `json:"mfaInitSkipLifetime,omitempty"`
	SecondFactorCheckLifetime  *time.Duration           `json:"secondFactorCheckLifetime,omitempty"`
	MultiFactorCheckLifetime   *time.Duration           `json:"multiFactorCheckLifetime,omitempty"`
}

type settingsColumns interface {
	IDColumn() database.Column
	InstanceIDColumn() database.Column
	OrgIDColumn() database.Column
	TypeColumn() database.Column
	SettingsColumn() database.Column
	CreatedAtColumn() database.Column
	UpdatedAtColumn() database.Column
}

type settingsConditions interface {
	InstanceIDCondition(id string) database.Condition
	OrgIDCondition(id *string) database.Condition
	IDCondition(id string) database.Condition
	TypeCondition(typ SettingType) database.Condition
}

type settingsChanges interface {
	SetType(state SettingType) database.Change
	SetSettings(settings string) database.Change
}

type SettingsRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	Get(ctx context.Context, id string, instanceID string, orgID *string) (*Setting, error)
	List(ctx context.Context, conditions ...database.Condition) ([]*Setting, error)

	Create(ctx context.Context, setting *Setting) error
	Update(ctx context.Context, id string, instanceID string, orgID *string, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, id string, instanceID string, orgID *string) (int64, error)
}
