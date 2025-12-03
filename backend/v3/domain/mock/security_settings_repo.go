package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type SecuritySettingsRepo struct {
	domain.SecuritySettingsRepository
	mock *MockSecuritySettingsRepository
}

func NewSecuritySettingsRepo(ctrl *gomock.Controller) *SecuritySettingsRepo {
	return &SecuritySettingsRepo{
		mock:                       NewMockSecuritySettingsRepository(ctrl),
		SecuritySettingsRepository: repository.SecuritySettingsRepository(),
	}
}

func (ssr *SecuritySettingsRepo) EXPECT() *MockSecuritySettingsRepositoryMockRecorder {
	return ssr.mock.EXPECT()
}

func (ssr *SecuritySettingsRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.SecuritySettings, error) {
	return ssr.mock.Get(ctx, client, opts...)
}

func (ssr *SecuritySettingsRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.SecuritySettings, error) {
	return ssr.mock.List(ctx, client, opts...)
}

func (ssr *SecuritySettingsRepo) Set(ctx context.Context, client database.QueryExecutor, settings *domain.SecuritySettings) error {
	return ssr.mock.Set(ctx, client, settings)
}

func (ssr *SecuritySettingsRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return ssr.mock.Delete(ctx, client, condition)
}
