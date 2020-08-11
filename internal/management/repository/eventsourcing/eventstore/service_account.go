package eventstore

import (
	"context"

	svcacc_model "github.com/caos/zitadel/internal/service_account/model"
	svcacc_event "github.com/caos/zitadel/internal/service_account/repository/eventsourcing"
)

type ServiceAccountRepo struct {
	ServiceAccountEvents *svcacc_event.ServiceAccountEventstore
}

func (repo *ServiceAccountRepo) CreateServiceAccount(ctx context.Context, account *svcacc_model.ServiceAccount) (*svcacc_model.ServiceAccount, error) {
	//TODO: create logic
	return repo.ServiceAccountEvents.CreateServiceAccount(ctx, account)
}

func (repo *ServiceAccountRepo) UpdateServiceAccount(ctx context.Context, account *svcacc_model.ServiceAccount) (*svcacc_model.ServiceAccount, error) {
	//TODO: update logic
	return repo.ServiceAccountEvents.UpdateServiceAccount(ctx, account)
}

func (repo *ServiceAccountRepo) DeactivateServiceAccount(ctx context.Context, id string) (*svcacc_model.ServiceAccount, error) {
	//TODO: deactivate logic
	return repo.ServiceAccountEvents.DeactivateServiceAccount(ctx, id)
}

func (repo *ServiceAccountRepo) ReactivateServiceAccount(ctx context.Context, id string) (*svcacc_model.ServiceAccount, error) {
	//TODO: reactivate logic
	return repo.ServiceAccountEvents.ReactivateServiceAccount(ctx, id)
}

func (repo *ServiceAccountRepo) LockServiceAccount(ctx context.Context, id string) (*svcacc_model.ServiceAccount, error) {
	//TODO: lock logic
	return repo.ServiceAccountEvents.LockServiceAccount(ctx, id)
}

func (repo *ServiceAccountRepo) UnlockServiceAccount(ctx context.Context, id string) (*svcacc_model.ServiceAccount, error) {
	//TODO: unlock logic
	return repo.ServiceAccountEvents.UnlockServiceAccount(ctx, id)
}

func (repo *ServiceAccountRepo) DeleteServiceAccount(ctx context.Context, id string) error {
	//TODO: delete logic
	return repo.ServiceAccountEvents.DeleteServiceAccount(ctx, id)
}
