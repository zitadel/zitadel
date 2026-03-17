package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func NewIDPIntentRepo(ctrl *gomock.Controller) *IDPIntentRepo {
	return &IDPIntentRepo{
		mock:                NewMockIDPIntentRepository(ctrl),
		IDPIntentRepository: repository.IDPIntentRepository(),
	}
}

type IDPIntentRepo struct {
	domain.IDPIntentRepository
	mock *MockIDPIntentRepository
}

func (u *IDPIntentRepo) EXPECT() *MockIDPIntentRepositoryMockRecorder {
	return u.mock.EXPECT()
}

func (i *IDPIntentRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPIntent, error) {
	return i.mock.Get(ctx, client, opts...)
}

func (i *IDPIntentRepo) Create(ctx context.Context, client database.QueryExecutor, intent *domain.IDPIntent) error {
	return i.mock.Create(ctx, client, intent)
}

func (i *IDPIntentRepo) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return i.mock.Update(ctx, client, condition, changes...)
}

func (i *IDPIntentRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return i.mock.Delete(ctx, client, condition)
}
