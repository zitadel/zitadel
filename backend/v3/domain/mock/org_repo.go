package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func NewOrgRepo(ctrl *gomock.Controller) *OrgRepo {
	return &OrgRepo{
		mock:                   NewMockOrganizationRepository(ctrl),
		OrganizationRepository: repository.OrganizationRepository(),
	}
}

type OrgRepo struct {
	domain.OrganizationRepository
	mock       *MockOrganizationRepository
	domainRepo *OrgDomainRepo

	existDomain database.Condition
}

func (r *OrgRepo) SetDomains(repo *OrgDomainRepo) {
	r.domainRepo = repo
}

// SetExistsDomain is a helper function that allows to skip calling EXISTS() when ExistsDomain() is called.
// You will not be able to mock the ExistsDomain() call though
func (r *OrgRepo) SetExistsDomain(existsDomainCondition database.Condition) {
	r.existDomain = existsDomainCondition

}

func (r *OrgRepo) ExistsDomain(_ database.Condition) database.Condition {
	return r.existDomain
}

func (r *OrgRepo) EXPECT() *MockOrganizationRepositoryMockRecorder {
	return r.mock.EXPECT()
}

func (r *OrgRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Organization, error) {
	return r.mock.Get(ctx, client, opts...)
}

func (r *OrgRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Organization, error) {
	return r.mock.List(ctx, client, opts...)
}

func (r *OrgRepo) Create(ctx context.Context, client database.QueryExecutor, org *domain.Organization) error {
	return r.mock.Create(ctx, client, org)
}

func (r *OrgRepo) Update(ctx context.Context, client database.QueryExecutor, conditions database.Condition, changes ...database.Change) (int64, error) {
	return r.mock.Update(ctx, client, conditions, changes...)
}

func (r *OrgRepo) Delete(ctx context.Context, client database.QueryExecutor, conditions database.Condition) (int64, error) {
	return r.mock.Delete(ctx, client, conditions)
}

func (r *OrgRepo) LoadDomains() domain.OrganizationRepository {
	return r.mock.LoadDomains()
}
