package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/storage/database"
)

type userOrchestrator interface {
	Create(ctx context.Context, tx database.Transaction, user *repository.User) (*repository.User, error)
	ByID(ctx context.Context, querier database.Querier, id string) (*repository.User, error)
}
