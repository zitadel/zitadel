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
		OrganizationRepository: repository.OrganizationRepository(nil),
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

func (r *OrgRepo) SetExistsDomain(existsDomainCondition database.Condition) {
	r.existDomain = existsDomainCondition

}

func (r *OrgRepo) ExistsDomain(_ database.Condition) database.Condition {
	return r.existDomain
}

func (r *OrgRepo) EXPECT() *MockOrganizationRepositoryMockRecorder {
	return r.mock.EXPECT()
}

func (r *OrgRepo) Get(ctx context.Context, opts ...database.QueryOption) (*domain.Organization, error) {
	return r.mock.Get(ctx, opts...)
}

func (r *OrgRepo) List(ctx context.Context, opts ...database.QueryOption) ([]*domain.Organization, error) {
	return r.mock.List(ctx, opts...)
}

func (r *OrgRepo) Create(ctx context.Context, organization *domain.Organization) error {
	return r.mock.Create(ctx, organization)
}

func (r *OrgRepo) Update(ctx context.Context, id domain.OrgIdentifierCondition, instance_id string, changes ...database.Change) (int64, error) {
	return r.mock.Update(ctx, id, instance_id, changes...)
}

func (r *OrgRepo) Delete(ctx context.Context, id domain.OrgIdentifierCondition, instance_id string) (int64, error) {
	return r.mock.Delete(ctx, id, instance_id)
}

func (r *OrgRepo) Domains(shouldLoad bool) domain.OrganizationDomainRepository {
	return r.domainRepo
}
