package domainmock

import (
	"context"

	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type SecretGeneratorSettingsRepo struct {
	domain.SecretGeneratorSettingsRepository
	mock *MockSecretGeneratorSettingsRepository
}

func NewSecretGeneratorSettingsRepo(ctrl *gomock.Controller) *SecretGeneratorSettingsRepo {
	return &SecretGeneratorSettingsRepo{
		mock:                              NewMockSecretGeneratorSettingsRepository(ctrl),
		SecretGeneratorSettingsRepository: repository.SecretGeneratorSettingsRepository(),
	}
}

func (r *SecretGeneratorSettingsRepo) EXPECT() *MockSecretGeneratorSettingsRepositoryMockRecorder {
	r.mock.ctrl.T.Helper()
	return r.mock.EXPECT()
}

func (r *SecretGeneratorSettingsRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.SecretGeneratorSettings, error) {
	r.mock.ctrl.T.Helper()
	return r.mock.Get(ctx, client, opts...)
}

func (r *SecretGeneratorSettingsRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.SecretGeneratorSettings, error) {
	r.mock.ctrl.T.Helper()
	return r.mock.List(ctx, client, opts...)
}

func (r *SecretGeneratorSettingsRepo) Set(ctx context.Context, client database.QueryExecutor, settings *domain.SecretGeneratorSettings) error {
	r.mock.ctrl.T.Helper()
	return r.mock.Set(ctx, client, settings)
}

func (r *SecretGeneratorSettingsRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	r.mock.ctrl.T.Helper()
	return r.mock.Delete(ctx, client, condition)
}
