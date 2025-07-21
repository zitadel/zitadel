package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate enumer -type IDPType -transform lower -trimprefix IDPType
type IDPType uint8

const (
	IDPTypeOIDC IDPType = iota
	IDPTypeOAUTH
	IDPTypeSAML
	IDPTypeLDAP
	IDPTypeGithub
	IDPTypeGoogle
	IDPTypeMicrosoft
	IDPTypeApple
)

//go:generate enumer -type IDPState -transform lower -trimprefix IDPState
type IDPState uint8

const (
	IDPStateActive IDPState = iota
	IDPStateInactive
)

type IdentityProvider struct {
	InstanceID        string    `json:"instanceId,omitempty" db:"instance_id"`
	OrgID             string    `json:"orgId,omitempty" db:"org_id"`
	ID                string    `json:"id,omitempty" db:"id"`
	State             string    `json:"state,omitempty" db:"state"`
	Name              string    `json:"name,omitempty" db:"name"`
	Type              string    `json:"type,omitempty" db:"type"`
	AllowCreation     bool      `json:"allowCreation,omitempty" db:"allow_creation"`
	AllowAutoCreation bool      `json:"allowAutoCreation,omitempty" db:"allow_auto_creation"`
	AllowAutoUpdate   bool      `json:"allowAutoUpdate,omitempty" db:"allow_auto_update"`
	AllowLinking      bool      `json:"allowLinking,omitempty" db:"allow_linking"`
	StylingType       int16     `json:"stylingType,omitempty" db:"styling_type"`
	Payload           string    `json:"payload,omitempty" db:"payload"`
	CreatedAt         time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

// IDPIdentifierCondition is used to help specify a single identity_provider,
// it will either be used as the  identity_provider ID or identity_provider name,
// as identity_provider can be identified either using (instnaceID + OrgID + ID) OR (instanceID + OrgID + name)
type IDPIdentifierCondition interface {
	database.Condition
}

type idProviderColumns interface {
	InstanceIDColumn() database.Column
	OrgIDColumn() database.Column
	IDColumn() database.Column
	StateColumn() database.Column
	NameColumn() database.Column
	TypeColumn() database.Column
	AllowCreationColumn() database.Column
	AllowAutoCreationColumn() database.Column
	AllowAutoUpdateColumn() database.Column
	AllowLinkingColumn() database.Column
	StylingTypeColumn() database.Column
	PayloadColumn() database.Column
	CreatedAtColumn() database.Column
	UpdatedAtColumn() database.Column
}

type idProviderConditions interface {
	InstanceIDCondition(id string) database.Condition
	OrgIDCondition(id string) database.Condition
	IDCondition(id string) IDPIdentifierCondition
	StateCondition(state IDPState) database.Condition
	NameCondition(name string) IDPIdentifierCondition
	TypeCondition(typee IDPType) database.Condition
	AllowCreationCondition(allow bool) database.Condition
	AllowAutoCreationCondition(allow bool) database.Condition
	AllowAutoUpdateCondition(allow bool) database.Condition
	AllowLinkingCondition(allow bool) database.Condition
	StylingTypeCondition(style int16) database.Condition
	PayloadCondition(payload string) database.Condition
}

type idProviderChanges interface {
	SetName(name string) database.Change
	SetState(state IDPState) database.Change
	SetAllowCreation(allow bool) database.Change
	SetAllowAutoCreation(allow bool) database.Change
	SetAllowAutoUpdate(allow bool) database.Change
	SetAllowLinking(allow bool) database.Change
	SetStylingType(stylingType int16) database.Change
	SetPayload(payload string) database.Change
}

type IDProviderRepository interface {
	idProviderColumns
	idProviderConditions
	idProviderChanges

	Get(ctx context.Context, id IDPIdentifierCondition, instnaceID string, orgID string) (*IdentityProvider, error)
	List(ctx context.Context, conditions ...database.Condition) ([]*IdentityProvider, error)

	Create(ctx context.Context, idp *IdentityProvider) error
	Update(ctx context.Context, id IDPIdentifierCondition, instnaceID string, orgID string, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, id IDPIdentifierCondition) (int64, error)
}
