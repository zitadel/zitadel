package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type OrganizationSettingsRepo struct {
	domain.OrganizationSettingsRepository
	mock *MockOrganizationSettingsRepository
}

func NewOrganizationSettingsRepo(ctrl *gomock.Controller) *OrganizationSettingsRepo {
	return &OrganizationSettingsRepo{
		mock:                           NewMockOrganizationSettingsRepository(ctrl),
		OrganizationSettingsRepository: repository.OrganizationSettingsRepository(),
	}
}

func (osr *OrganizationSettingsRepo) EXPECT() *MockOrganizationSettingsRepositoryMockRecorder {
	return osr.mock.EXPECT()
}

func (osr *OrganizationSettingsRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.OrganizationSettings, error) {
	return osr.mock.Get(ctx, client, opts...)
}

func (osr *OrganizationSettingsRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.OrganizationSettings, error) {
	return osr.mock.List(ctx, client, opts...)
}

func (osr *OrganizationSettingsRepo) Set(ctx context.Context, client database.QueryExecutor, settings *domain.OrganizationSettings) error {
	return osr.mock.Set(ctx, client, settings)
}

func (osr *OrganizationSettingsRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return osr.mock.Delete(ctx, client, condition)
}
