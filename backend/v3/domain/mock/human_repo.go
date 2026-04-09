package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func NewHumanRepo(ctrl *gomock.Controller) *HumanRepo {
	return &HumanRepo{
		mock:                NewMockHumanUserRepository(ctrl),
		HumanUserRepository: repository.HumanUserRepository(),
	}
}

type HumanRepo struct {
	domain.HumanUserRepository
	mock *MockHumanUserRepository
}

func (u *HumanRepo) EXPECT() *MockHumanUserRepositoryMockRecorder {
	return u.mock.EXPECT()
}

func (u *HumanRepo) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return u.mock.Update(ctx, client, condition, changes...)
}
