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
	OrgID             string    `json:"org_id,omitempty" db:"org_id"`
	ID                string    `json:"id,omitempty" db:"id"`
	State             string    `json:"state,omitempty" db:"state"`
	Name              string    `json:"name,omitempty" db:"name"`
	Type              string    `json:"type,omitempty" db:"type"`
	AllowCreation     bool      `json:"allow_creation,omitempty" db:"allow_creation"`
	AllowAutoCreation bool      `json:"allow_auto_creation,omitempty" db:"allow_auto_creation"`
	AllowAutoUpdate   bool      `json:"allow_auto_update,omitempty" db:"allow_auto_update"`
	AllowLinking      bool      `json:"allow_linking,omitempty" db:"allow_linking"`
	StylingType       int16     `json:"styling_type,omitempty" db:"styling_type"`
	Payload           string    `json:"payload,omitempty" db:"payload"`
	CreatedAt         time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt,omitempty" db:"updated_at"`
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
	IDCondition(id string) database.Condition
	StateCondition(state IDPState) database.Condition
	NameCondition(name string) database.Condition
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

	Get(ctx context.Context, id string) (*IdentityProvider, error)
	List(ctx context.Context, conditions ...database.Condition) ([]*IdentityProvider, error)

	Create(ctx context.Context, idp *IdentityProvider) error
	Update(ctx context.Context, id string, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}
