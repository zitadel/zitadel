package repository

import (
	"context"

	"github.com/caos/zitadel/internal/service_account/model"
)

type ServiceAccountRepository interface {
	ServiceAccountByID(ctx context.Context, id string) (*model.ServiceAccountView, error)
	CreateServiceAccount(ctx context.Context, user *model.ServiceAccount) (*model.ServiceAccount, error)
	DeactivateServiceAccount(ctx context.Context, id string) (*model.ServiceAccount, error)
	ReactivateServiceAccount(ctx context.Context, id string) (*model.ServiceAccount, error)
	LockServiceAccount(ctx context.Context, id string) (*model.ServiceAccount, error)
	UnlockServiceAccount(ctx context.Context, id string) (*model.ServiceAccount, error)
	SearchServiceAccounts(ctx context.Context, request *model.ServiceAccountSearchRequest) (*model.ServiceAccountSearchResult, error)
	ServiceAccountChanges(ctx context.Context, id string, lastSequence, limit uint64, sortAscending bool) (*model.ServiceAccountChanges, error)
	IsServiceAccountUnique(ctx context.Context) (bool, error)
}
