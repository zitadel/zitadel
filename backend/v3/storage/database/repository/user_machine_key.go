package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userMachineKey struct{}

func MachineKeyRepository() domain.MachineKeyRepository {
	return new(userMachineKey)
}

// Add implements [domain.MachineKeyRepository].
func (u userMachineKey) Add(ctx context.Context, client database.QueryExecutor, key *domain.MachineKey) error {
	panic("unimplemented")
}

// Get implements [domain.MachineKeyRepository].
func (u userMachineKey) Get(ctx context.Context, client database.QueryExecutor, condition database.Condition, opts ...database.QueryOption) (*domain.MachineKey, error) {
	panic("unimplemented")
}

// List implements [domain.MachineKeyRepository].
func (u userMachineKey) List(ctx context.Context, client database.QueryExecutor, condition database.Condition, opts ...database.QueryOption) ([]*domain.MachineKey, error) {
	panic("unimplemented")
}

// PrimaryKeyColumns implements [domain.MachineKeyRepository].
func (u userMachineKey) PrimaryKeyColumns() []database.Column {
	panic("unimplemented")
}

// PrimaryKeyCondition implements [domain.MachineKeyRepository].
func (u userMachineKey) PrimaryKeyCondition(instanceID string, id string) database.Condition {
	panic("unimplemented")
}

// Remove implements [domain.MachineKeyRepository].
func (u userMachineKey) Remove(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	panic("unimplemented")
}

var _ domain.MachineKeyRepository = (*userMachineKey)(nil)
