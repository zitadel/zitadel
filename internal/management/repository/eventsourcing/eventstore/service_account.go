package eventstore

import (
	"context"

	svcacc_model "github.com/caos/zitadel/internal/user/model"
	svcacc_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type ServiceAccountRepo struct {
	UserEvents *svcacc_event.UserEventstore
}

func (repo *ServiceAccountRepo) CreateServiceAccount(ctx context.Context, account *svcacc_model.Machine) (*svcacc_model.Machine, error) {
	//TODO: create logic
	return repo.UserEvents.CreateServiceAccount(ctx, account)
}

func (repo *ServiceAccountRepo) UpdateServiceAccount(ctx context.Context, account *svcacc_model.Machine) (*svcacc_model.Machine, error) {
	//TODO: update logic
	return repo.UserEvents.UpdateServiceAccount(ctx, account)
}

func (repo *ServiceAccountRepo) DeactivateServiceAccount(ctx context.Context, id string) (*svcacc_model.Machine, error) {
	//TODO: deactivate logic
	return repo.UserEvents.DeactivateServiceAccount(ctx, id)
}

func (repo *ServiceAccountRepo) ReactivateServiceAccount(ctx context.Context, id string) (*svcacc_model.Machine, error) {
	//TODO: reactivate logic
	return repo.UserEvents.ReactivateServiceAccount(ctx, id)
}

func (repo *ServiceAccountRepo) LockServiceAccount(ctx context.Context, id string) (*svcacc_model.Machine, error) {
	//TODO: lock logic
	return repo.UserEvents.LockServiceAccount(ctx, id)
}

func (repo *ServiceAccountRepo) UnlockServiceAccount(ctx context.Context, id string) (*svcacc_model.Machine, error) {
	//TODO: unlock logic
	return repo.UserEvents.UnlockServiceAccount(ctx, id)
}

func (repo *ServiceAccountRepo) DeleteServiceAccount(ctx context.Context, id string) error {
	//TODO: delete logic
	return repo.UserEvents.DeleteServiceAccount(ctx, id)
}
