package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type LoginSettingsRepo struct {
	domain.LoginSettingsRepository
	mock *MockLoginSettingsRepository
}

func NewLoginSettingsRepo(ctrl *gomock.Controller) *LoginSettingsRepo {
	return &LoginSettingsRepo{
		mock:                    NewMockLoginSettingsRepository(ctrl),
		LoginSettingsRepository: repository.LoginSettingsRepository(),
	}
}

func (lsr *LoginSettingsRepo) EXPECT() *MockLoginSettingsRepositoryMockRecorder {
	return lsr.mock.EXPECT()
}

func (lsr *LoginSettingsRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LoginSettings, error) {
	return lsr.mock.Get(ctx, client, opts...)
}

func (lsr *LoginSettingsRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.LoginSettings, error) {
	return lsr.mock.List(ctx, client, opts...)
}

func (lsr *LoginSettingsRepo) Set(ctx context.Context, client database.QueryExecutor, settings *domain.LoginSettings) error {
	return lsr.mock.Set(ctx, client, settings)
}

func (lsr *LoginSettingsRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return lsr.mock.Delete(ctx, client, condition)
}

func (lsr *LoginSettingsRepo) SetColumns(ctx context.Context, client database.QueryExecutor, settings *domain.Settings, changes ...database.Change) error {
	return lsr.mock.SetColumns(ctx, client, settings, changes...)
}
