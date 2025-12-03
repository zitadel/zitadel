package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type PasswordComplexitySettingsRepo struct {
	domain.PasswordComplexitySettingsRepository
	mock *MockPasswordComplexitySettingsRepository
}

func NewPasswordComplexitySettingsRepo(ctrl *gomock.Controller) *PasswordComplexitySettingsRepo {
	return &PasswordComplexitySettingsRepo{
		mock:                                 NewMockPasswordComplexitySettingsRepository(ctrl),
		PasswordComplexitySettingsRepository: repository.PasswordComplexitySettingsRepository(),
	}
}

func (pcsr *PasswordComplexitySettingsRepo) EXPECT() *MockPasswordComplexitySettingsRepositoryMockRecorder {
	return pcsr.mock.EXPECT()
}

func (pcsr *PasswordComplexitySettingsRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.PasswordComplexitySettings, error) {
	return pcsr.mock.Get(ctx, client, opts...)
}

func (pcsr *PasswordComplexitySettingsRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.PasswordComplexitySettings, error) {
	return pcsr.mock.List(ctx, client, opts...)
}

func (pcsr *PasswordComplexitySettingsRepo) Set(ctx context.Context, client database.QueryExecutor, settings *domain.PasswordComplexitySettings) error {
	return pcsr.mock.Set(ctx, client, settings)
}

func (pcsr *PasswordComplexitySettingsRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return pcsr.mock.Delete(ctx, client, condition)
}
