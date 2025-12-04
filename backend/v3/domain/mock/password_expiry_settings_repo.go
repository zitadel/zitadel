package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type PasswordExpirySettingsRepo struct {
	domain.PasswordExpirySettingsRepository
	mock *MockPasswordExpirySettingsRepository
}

func NewPasswordExpirySettingsRepo(ctrl *gomock.Controller) *PasswordExpirySettingsRepo {
	return &PasswordExpirySettingsRepo{
		mock:                             NewMockPasswordExpirySettingsRepository(ctrl),
		PasswordExpirySettingsRepository: repository.PasswordExpirySettingsRepository(),
	}
}

func (pesr *PasswordExpirySettingsRepo) EXPECT() *MockPasswordExpirySettingsRepositoryMockRecorder {
	pesr.mock.ctrl.T.Helper()
	return pesr.mock.EXPECT()
}

func (pesr *PasswordExpirySettingsRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.PasswordExpirySettings, error) {
	pesr.mock.ctrl.T.Helper()
	return pesr.mock.Get(ctx, client, opts...)
}

func (pesr *PasswordExpirySettingsRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.PasswordExpirySettings, error) {
	pesr.mock.ctrl.T.Helper()
	return pesr.mock.List(ctx, client, opts...)
}

func (pesr *PasswordExpirySettingsRepo) Set(ctx context.Context, client database.QueryExecutor, settings *domain.PasswordExpirySettings) error {
	pesr.mock.ctrl.T.Helper()
	return pesr.mock.Set(ctx, client, settings)
}

func (pesr *PasswordExpirySettingsRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	pesr.mock.ctrl.T.Helper()
	return pesr.mock.Delete(ctx, client, condition)
}
