package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type MachineKeyType uint8

const (
	MachineKeyTypeUnspecified MachineKeyType = iota
	MachineKeyTypeNone
	MachineKeyTypeJSON
)

type MachineKey struct {
	ID        string
	Type      MachineKeyType
	PublicKey []byte
	CreatedAt time.Time
	ExpiresAt *time.Time
}

//go:generate mockgen -typed -package domainmock -destination ./mock/machine_key.mock.go . MachineKeyRepository
type MachineKeyRepository interface {
	Repository
	machineKeyColumns
	machineKeyConditions

	Get(ctx context.Context, client database.QueryExecutor, condition database.Condition, opts ...database.QueryOption) (*MachineKey, error)
	List(ctx context.Context, client database.QueryExecutor, condition database.Condition, opts ...database.QueryOption) ([]*MachineKey, error)

	Add(ctx context.Context, client database.QueryExecutor, key *MachineKey) error
	Remove(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}

type machineKeyColumns interface{}

type machineKeyConditions interface {
	PrimaryKeyCondition(instanceID, id string) database.Condition
}
