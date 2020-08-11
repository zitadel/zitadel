package repository

import (
	"context"

	"github.com/caos/zitadel/internal/service_account/model"
)

type ServiceAccountRepository interface {
	CreateServiceAccount(ctx context.Context, user *model.ServiceAccount) (*model.ServiceAccount, error)
	UpdateServiceAccount(ctx context.Context, user *model.ServiceAccount) (*model.ServiceAccount, error)
	DeactivateServiceAccount(ctx context.Context, id string) (*model.ServiceAccount, error)
	ReactivateServiceAccount(ctx context.Context, id string) (*model.ServiceAccount, error)
	LockServiceAccount(ctx context.Context, id string) (*model.ServiceAccount, error)
	UnlockServiceAccount(ctx context.Context, id string) (*model.ServiceAccount, error)
}
