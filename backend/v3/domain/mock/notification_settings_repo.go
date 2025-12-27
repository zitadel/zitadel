package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type NotificationSettingsRepo struct {
	domain.NotificationSettingsRepository
	mock *MockNotificationSettingsRepository
}

func NewNotificationSettingsRepo(ctrl *gomock.Controller) *NotificationSettingsRepo {
	return &NotificationSettingsRepo{
		mock:                           NewMockNotificationSettingsRepository(ctrl),
		NotificationSettingsRepository: repository.NotificationSettingsRepository(),
	}
}

func (nsr *NotificationSettingsRepo) EXPECT() *MockNotificationSettingsRepositoryMockRecorder {
	nsr.mock.ctrl.T.Helper()
	return nsr.mock.EXPECT()
}

func (nsr *NotificationSettingsRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.NotificationSettings, error) {
	nsr.mock.ctrl.T.Helper()
	return nsr.mock.Get(ctx, client, opts...)
}

func (nsr *NotificationSettingsRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.NotificationSettings, error) {
	nsr.mock.ctrl.T.Helper()
	return nsr.mock.List(ctx, client, opts...)
}

func (nsr *NotificationSettingsRepo) Set(ctx context.Context, client database.QueryExecutor, settings *domain.NotificationSettings) error {
	nsr.mock.ctrl.T.Helper()
	return nsr.mock.Set(ctx, client, settings)
}

func (nsr *NotificationSettingsRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	nsr.mock.ctrl.T.Helper()
	return nsr.mock.Delete(ctx, client, condition)
}
