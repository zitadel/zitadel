package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type InstancesDomainRepo struct {
	domain.InstanceDomainRepository
	mock *MockInstanceDomainRepository
}

func NewInstancesDomainRepo(ctrl *gomock.Controller) *InstancesDomainRepo {
	return &InstancesDomainRepo{
		mock:                     NewMockInstanceDomainRepository(ctrl),
		InstanceDomainRepository: repository.InstanceDomainRepository(),
	}
}

func (r *InstancesDomainRepo) EXPECT() *MockInstanceDomainRepositoryMockRecorder {
	return r.mock.EXPECT()
}

func (r *InstancesDomainRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.InstanceDomain, error) {
	return r.mock.Get(ctx, client, opts...)

}

func (r *InstancesDomainRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.InstanceDomain, error) {
	return r.mock.List(ctx, client, opts...)

}

func (r *InstancesDomainRepo) Add(ctx context.Context, client database.QueryExecutor, domain *domain.AddInstanceDomain) error {
	return r.mock.Add(ctx, client, domain)

}

func (r *InstancesDomainRepo) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return r.mock.Update(ctx, client, condition, changes...)

}

func (r *InstancesDomainRepo) Remove(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return r.mock.Remove(ctx, client, condition)

}
