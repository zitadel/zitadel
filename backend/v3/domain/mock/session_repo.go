package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func NewSessionRepo(ctrl *gomock.Controller) *SessionRepo {
	return &SessionRepo{
		mock:              NewMockSessionRepository(ctrl),
		SessionRepository: repository.SessionRepository(),
	}
}

type SessionRepo struct {
	domain.SessionRepository
	mock *MockSessionRepository
}

func (s *SessionRepo) EXPECT() *MockSessionRepositoryMockRecorder {
	return s.mock.EXPECT()
}

func (s *SessionRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Session, error) {
	return s.mock.Get(ctx, client, opts...)
}

func (s *SessionRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Session, error) {
	return s.mock.List(ctx, client, opts...)
}

func (s *SessionRepo) Create(ctx context.Context, client database.QueryExecutor, user *domain.Session) error {
	return s.mock.Create(ctx, client, user)
}

func (s *SessionRepo) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.mock.Update(ctx, client, condition, changes...)
}

func (s *SessionRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return s.mock.Delete(ctx, client, condition)
}
