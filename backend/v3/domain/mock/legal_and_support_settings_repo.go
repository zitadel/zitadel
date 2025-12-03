package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type LegalAndSupportSettingsRepo struct {
	domain.LegalAndSupportSettingsRepository
	mock *MockLegalAndSupportSettingsRepository
}

func NewLegalAndSupportSettingsRepo(ctrl *gomock.Controller) *LegalAndSupportSettingsRepo {
	return &LegalAndSupportSettingsRepo{
		mock:                              NewMockLegalAndSupportSettingsRepository(ctrl),
		LegalAndSupportSettingsRepository: repository.LegalAndSupportSettingsRepository(),
	}
}

func (lssr *LegalAndSupportSettingsRepo) EXPECT() *MockLegalAndSupportSettingsRepositoryMockRecorder {
	return lssr.mock.EXPECT()
}

func (lssr *LegalAndSupportSettingsRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LegalAndSupportSettings, error) {
	return lssr.mock.Get(ctx, client, opts...)
}

func (lssr *LegalAndSupportSettingsRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.LegalAndSupportSettings, error) {
	return lssr.mock.List(ctx, client, opts...)
}

func (lssr *LegalAndSupportSettingsRepo) Set(ctx context.Context, client database.QueryExecutor, settings *domain.LegalAndSupportSettings) error {
	return lssr.mock.Set(ctx, client, settings)
}

func (lssr *LegalAndSupportSettingsRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return lssr.mock.Delete(ctx, client, condition)
}
