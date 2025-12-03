package domainmock

import (
	context "context"
	time "time"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type BrandingSettingsRepo struct {
	domain.BrandingSettingsRepository
	mock *MockBrandingSettingsRepository
}

func NewBrandingSettingsRepo(ctrl *gomock.Controller) *BrandingSettingsRepo {
	return &BrandingSettingsRepo{
		mock:                       NewMockBrandingSettingsRepository(ctrl),
		BrandingSettingsRepository: repository.BrandingSettingsRepository(),
	}
}

func (bsr *BrandingSettingsRepo) EXPECT() *MockBrandingSettingsRepositoryMockRecorder {
	return bsr.mock.EXPECT()
}

func (bsr *BrandingSettingsRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.BrandingSettings, error) {
	return bsr.mock.Get(ctx, client, opts...)
}

func (bsr *BrandingSettingsRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.BrandingSettings, error) {
	return bsr.mock.List(ctx, client, opts...)
}

func (bsr *BrandingSettingsRepo) Set(ctx context.Context, client database.QueryExecutor, settings *domain.BrandingSettings) error {
	return bsr.mock.Set(ctx, client, settings)
}

func (bsr *BrandingSettingsRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return bsr.mock.Delete(ctx, client, condition)
}

func (bsr *BrandingSettingsRepo) Activate(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return bsr.mock.Activate(ctx, client, condition)
}

func (bsr *BrandingSettingsRepo) ActivateAt(ctx context.Context, client database.QueryExecutor, condition database.Condition, updatedAt time.Time) (int64, error) {
	return bsr.mock.ActivateAt(ctx, client, condition, updatedAt)
}

func (bsr *BrandingSettingsRepo) SetColumns(ctx context.Context, client database.QueryExecutor, settings *domain.Settings, changes ...database.Change) error {
	return bsr.mock.SetColumns(ctx, client, settings, changes...)
}
