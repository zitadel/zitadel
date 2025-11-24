package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func NewUserRepo(ctrl *gomock.Controller) *UserRepo {
	return &UserRepo{
		mock:           NewMockUserRepository(ctrl),
		UserRepository: repository.UserRepository(),
	}
}

type UserRepo struct {
	domain.UserRepository
	mock *MockUserRepository
}

func (u *UserRepo) EXPECT() *MockUserRepositoryMockRecorder {
	return u.mock.EXPECT()
}

func (u *UserRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.User, error) {
	return u.mock.Get(ctx, client, opts...)
}

func (u *UserRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.User, error) {
	return u.mock.List(ctx, client, opts...)
}

func (u *UserRepo) Create(ctx context.Context, client database.QueryExecutor, user *domain.User) error {
	return u.mock.Create(ctx, client, user)
}

func (u *UserRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) error {
	return u.mock.Delete(ctx, client, condition)
}

func (u *UserRepo) Human() domain.HumanRepository {
	return u.mock.Human()
}

func (u *UserRepo) Machine() domain.MachineRepository {
	return u.mock.Machine()
}
