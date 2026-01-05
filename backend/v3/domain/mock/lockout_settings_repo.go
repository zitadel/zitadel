package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type LockoutSettingsRepo struct {
	domain.LockoutSettingsRepository
	mock *MockLockoutSettingsRepository
}

func NewLockoutSettingsRepo(ctrl *gomock.Controller) *LockoutSettingsRepo {
	return &LockoutSettingsRepo{
		mock:                      NewMockLockoutSettingsRepository(ctrl),
		LockoutSettingsRepository: repository.LockoutSettingsRepository(),
	}
}

func (lsr *LockoutSettingsRepo) EXPECT() *MockLockoutSettingsRepositoryMockRecorder {
	lsr.mock.ctrl.T.Helper()
	return lsr.mock.EXPECT()
}

func (lsr *LockoutSettingsRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LockoutSettings, error) {
	lsr.mock.ctrl.T.Helper()
	return lsr.mock.Get(ctx, client, opts...)
}

func (lsr *LockoutSettingsRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.LockoutSettings, error) {
	lsr.mock.ctrl.T.Helper()
	return lsr.mock.List(ctx, client, opts...)
}

func (lsr *LockoutSettingsRepo) Set(ctx context.Context, client database.QueryExecutor, settings *domain.LockoutSettings) error {
	lsr.mock.ctrl.T.Helper()
	return lsr.mock.Set(ctx, client, settings)
}

func (lsr *LockoutSettingsRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	lsr.mock.ctrl.T.Helper()
	return lsr.mock.Delete(ctx, client, condition)
}
