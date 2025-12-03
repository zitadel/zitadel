package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type DomainSettingsRepo struct {
	domain.DomainSettingsRepository
	mock *MockDomainSettingsRepository
}

func NewDomainSettingsRepo(ctrl *gomock.Controller) *DomainSettingsRepo {
	return &DomainSettingsRepo{
		mock:                     NewMockDomainSettingsRepository(ctrl),
		DomainSettingsRepository: repository.DomainSettingsRepository(),
	}
}

func (dsr *DomainSettingsRepo) EXPECT() *MockDomainSettingsRepositoryMockRecorder {
	return dsr.mock.EXPECT()
}

func (dsr *DomainSettingsRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.DomainSettings, error) {
	return dsr.mock.Get(ctx, client, opts...)
}

func (dsr *DomainSettingsRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.DomainSettings, error) {
	return dsr.mock.List(ctx, client, opts...)
}

func (dsr *DomainSettingsRepo) Set(ctx context.Context, client database.QueryExecutor, settings *domain.DomainSettings) error {
	return dsr.mock.Set(ctx, client, settings)
}

func (dsr *DomainSettingsRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return dsr.mock.Delete(ctx, client, condition)
}
