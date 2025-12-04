package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
)

func NewSessionRepo(ctrl *gomock.Controller) *SessionRepo {
	return &SessionRepo{
		mock: NewMockSessionRepository(ctrl),
		// TODO(IAM-Marco): Uncomment when session repository is complete (see https://github.com/zitadel/zitadel/issues/10212)
		// SessionRepository: repository.SessionRepository(),
	}
}

type SessionRepo struct {
	domain.SessionRepository
	mock *MockSessionRepository
}

func (s *SessionRepo) EXPECT() *MockSessionRepositoryMockRecorder {
	s.mock.ctrl.T.Helper()
	return s.mock.EXPECT()
}

func (s *SessionRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Session, error) {
	s.mock.ctrl.T.Helper()
	return s.mock.Get(ctx, client, opts...)
}

func (s *SessionRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Session, error) {
	s.mock.ctrl.T.Helper()
	return s.mock.List(ctx, client, opts...)
}

func (s *SessionRepo) Create(ctx context.Context, client database.QueryExecutor, user *domain.Session) error {
	s.mock.ctrl.T.Helper()
	return s.mock.Create(ctx, client, user)
}

func (s *SessionRepo) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	s.mock.ctrl.T.Helper()
	return s.mock.Update(ctx, client, condition, changes...)
}

func (s *SessionRepo) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) error {
	s.mock.ctrl.T.Helper()
	return s.mock.Delete(ctx, client, condition)
}
