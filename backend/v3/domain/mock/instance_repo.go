package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
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
	r.mock.ctrl.T.Helper()
	return r.mock.EXPECT()
}

func (r *InstanceRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Instance, error) {
	r.mock.ctrl.T.Helper()
	return r.mock.Get(ctx, client, opts...)

}

func (r *InstanceRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Instance, error) {
	r.mock.ctrl.T.Helper()
	return r.mock.List(ctx, client, opts...)

}

func (r *InstanceRepo) Create(ctx context.Context, client database.QueryExecutor, instance *domain.Instance) error {
	r.mock.ctrl.T.Helper()
	return r.mock.Create(ctx, client, instance)

}

func (r *InstanceRepo) Update(ctx context.Context, client database.QueryExecutor, instanceID string, changes ...database.Change) (int64, error) {
	r.mock.ctrl.T.Helper()
	return r.mock.Update(ctx, client, instanceID, changes...)

}

func (r *InstanceRepo) Delete(ctx context.Context, client database.QueryExecutor, instanceID string) (int64, error) {
	r.mock.ctrl.T.Helper()
	return r.mock.Delete(ctx, client, instanceID)

}

func (r *InstanceRepo) LoadDomains() domain.InstanceRepository {
	r.mock.ctrl.T.Helper()
	return r.mock.LoadDomains()
}
