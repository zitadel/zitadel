package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate enumer -type AccessTokenType -transform lower -trimprefix AccessTokenType
type AccessTokenType uint8

const (
	AccessTokenTypeUnspecified AccessTokenType = iota
	AccessTokenTypeBearer
	AccessTokenTypeJWT
)

type Machine struct {
	Name            string          `json:"name,omitempty" db:"name"`
	Description     *string         `json:"description,omitempty" db:"description"`
	Secret          *string         `json:"secret,omitempty" db:"secret"`
	AccessTokenType AccessTokenType `json:"accessTokenType,omitempty" db:"access_token_type"`
}

//go:generate mockgen -typed -package domainmock -destination ./mock/user_machine.mock.go . MachineUserRepository
type MachineUserRepository interface {
	machineColumns
	machineConditions
	machineChanges

	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)

	// Keys returns the [MachineKeyRepository] for machine users
	Keys() MachineKeyRepository
	// LoadKeys enables fetching of machine keys when getting or listing machine users
	LoadKeys() MachineUserRepository
	// PersonalAccessTokens returns the [PersonalAccessTokenRepository] for machine users
	PersonalAccessTokens() PersonalAccessTokenRepository
	// LoadPersonalAccessTokens enables fetching of personal access tokens when getting or listing machine users
	LoadPersonalAccessTokens() MachineUserRepository
}

type machineColumns interface {
	userColumns
	NameColumn() database.Column
	DescriptionColumn() database.Column
	SecretColumn() database.Column
	AccessTokenTypeColumn() database.Column
}

type machineConditions interface {
	userConditions
	NameCondition(op database.TextOperation, name string) database.Condition
	DescriptionCondition(op database.TextOperation, description string) database.Condition
	AccessTokenTypeCondition(accessTokenType AccessTokenType) database.Condition
}

type machineChanges interface {
	userChanges
	SetName(name string) database.Change
	// SetDescription sets the description, nil to clear it
	SetDescription(description *string) database.Change
	// SetSecret sets the secret, nil to clear it
	SetSecret(secret *string) database.Change
	SetAccessTokenType(accessTokenType AccessTokenType) database.Change
}
