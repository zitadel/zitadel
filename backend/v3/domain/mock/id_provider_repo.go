package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type IDProviderRepo struct {
	domain.IDProviderRepository
	mock *MockIDProviderRepository
}

func NewIDProviderRepo(ctrl *gomock.Controller) *IDProviderRepo {
	return &IDProviderRepo{
		mock:                 NewMockIDProviderRepository(ctrl),
		IDProviderRepository: repository.IDProviderRepository(),
	}
}

func (r *IDProviderRepo) EXPECT() *MockIDProviderRepositoryMockRecorder {
	return r.mock.EXPECT()
}

func (r *IDProviderRepo) Get(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IdentityProvider, error) {
	return r.mock.Get(ctx, client, id, instanceID, orgID)
}

func (r *IDProviderRepo) List(ctx context.Context, client database.QueryExecutor, conditions ...database.Condition) ([]*domain.IdentityProvider, error) {
	return r.mock.List(ctx, client, conditions...)
}

func (r *IDProviderRepo) Create(ctx context.Context, client database.QueryExecutor, idp *domain.IdentityProvider) error {
	return r.mock.Create(ctx, client, idp)
}

func (r *IDProviderRepo) Update(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string, changes ...database.Change) (int64, error) {
	return r.mock.Update(ctx, client, id, instanceID, orgID, changes...)
}

func (r *IDProviderRepo) Delete(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (int64, error) {
	return r.mock.Delete(ctx, client, id, instanceID, orgID)
}

func (r *IDProviderRepo) GetOIDC(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPOIDC, error) {
	return r.mock.GetOIDC(ctx, client, id, instanceID, orgID)
}

func (r *IDProviderRepo) GetJWT(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPJWT, error) {
	return r.mock.GetJWT(ctx, client, id, instanceID, orgID)
}

func (r *IDProviderRepo) GetOAuth(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPOAuth, error) {
	return r.mock.GetOAuth(ctx, client, id, instanceID, orgID)
}

func (r *IDProviderRepo) GetAzureAD(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPAzureAD, error) {
	return r.mock.GetAzureAD(ctx, client, id, instanceID, orgID)
}

func (r *IDProviderRepo) GetGoogle(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPGoogle, error) {
	return r.mock.GetGoogle(ctx, client, id, instanceID, orgID)
}

func (r *IDProviderRepo) GetGithub(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPGithub, error) {
	return r.mock.GetGithub(ctx, client, id, instanceID, orgID)
}

func (r *IDProviderRepo) GetGithubEnterprise(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPGithubEnterprise, error) {
	return r.mock.GetGithubEnterprise(ctx, client, id, instanceID, orgID)
}

func (r *IDProviderRepo) GetGitlab(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPGitlab, error) {
	return r.mock.GetGitlab(ctx, client, id, instanceID, orgID)
}

func (r *IDProviderRepo) GetGitlabSelfHosting(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPGitlabSelfHosting, error) {
	return r.mock.GetGitlabSelfHosting(ctx, client, id, instanceID, orgID)
}

func (r *IDProviderRepo) GetLDAP(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPLDAP, error) {
	return r.mock.GetLDAP(ctx, client, id, instanceID, orgID)
}

func (r *IDProviderRepo) GetApple(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPApple, error) {
	return r.mock.GetApple(ctx, client, id, instanceID, orgID)
}

func (r *IDProviderRepo) GetSAML(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPSAML, error) {
	return r.mock.GetSAML(ctx, client, id, instanceID, orgID)
}
