package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type OrgDomainRepo struct {
	domain.OrganizationDomainRepository
	mock *MockOrganizationDomainRepository
}

func NewOrgDomainRepo(ctrl *gomock.Controller) *OrgDomainRepo {
	return &OrgDomainRepo{
		mock:                         NewMockOrganizationDomainRepository(ctrl),
		OrganizationDomainRepository: repository.OrganizationDomainRepository(),
	}
}

func (r *OrgDomainRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.OrganizationDomain, error) {
	return r.mock.Get(ctx, client, opts...)

}

func (r *OrgDomainRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.OrganizationDomain, error) {
	return r.mock.List(ctx, client, opts...)

}

func (r *OrgDomainRepo) Add(ctx context.Context, client database.QueryExecutor, domain *domain.AddOrganizationDomain) error {
	return r.mock.Add(ctx, client, domain)

}

func (r *OrgDomainRepo) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return r.mock.Update(ctx, client, condition, changes...)

}

func (r *OrgDomainRepo) Remove(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return r.mock.Remove(ctx, client, condition)

}
