package domainmock

import (
	"context"

	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type InstanceRepo struct {
	domain.InstanceRepository
	mock *MockInstanceRepository
}

func NewInstanceRepo(ctrl *gomock.Controller) *InstanceRepo {
	return &InstanceRepo{
		mock:               NewMockInstanceRepository(ctrl),
		InstanceRepository: repository.InstanceRepository(),
	}
}

func (r *InstanceRepo) EXPECT() *MockInstanceRepositoryMockRecorder {
	return r.mock.EXPECT()
}

func (r *InstanceRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Instance, error) {
	return r.mock.Get(ctx, client, opts...)

}

func (r *InstanceRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Instance, error) {
	return r.mock.List(ctx, client, opts...)

}

func (r *InstanceRepo) Create(ctx context.Context, client database.QueryExecutor, instance *domain.Instance) error {
	return r.mock.Create(ctx, client, instance)

}

func (r *InstanceRepo) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return r.mock.Update(ctx, client, condition, changes...)

}

func (r *InstanceRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return r.mock.Delete(ctx, client, condition)

}

func (r *InstanceRepo) LoadDomains() domain.InstanceRepository {
	return r.mock.LoadDomains()
}
