package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type instance struct {
	database.QueryExecutor
}

func Instance(client database.QueryExecutor) domain.InstanceRepository {
	return &instance{QueryExecutor: client}
}

func (i *instance) ByID(ctx context.Context, id string) (*domain.Instance, error) {
	var instance domain.Instance
	err := i.QueryExecutor.QueryRow(ctx, `SELECT id, name, created_at, updated_at, deleted_at FROM instances WHERE id = $1`, id).Scan(
		&instance.ID,
		&instance.Name,
		&instance.CreatedAt,
		&instance.UpdatedAt,
		&instance.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

const createInstanceStmt = `INSERT INTO instances (id, name) VALUES ($1, $2) RETURNING created_at, updated_at`

// Create implements [domain.InstanceRepository].
func (i *instance) Create(ctx context.Context, instance *domain.Instance) error {
	return i.QueryExecutor.QueryRow(ctx, createInstanceStmt,
		instance.ID,
		instance.Name,
	).Scan(
		&instance.CreatedAt,
		&instance.UpdatedAt,
	)
}

// On implements [domain.InstanceRepository].
func (i *instance) On(id string) domain.InstanceOperation {
	return &instanceOperation{
		QueryExecutor: i.QueryExecutor,
		id:            id,
	}
}

var _ domain.InstanceRepository = (*instance)(nil)
