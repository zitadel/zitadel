package domainmock

import (
	context "context"

	gomock "go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	database "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func NewProjectRepo(ctrl *gomock.Controller) *ProjectRepo {
	return &ProjectRepo{
		mock:              NewMockProjectRepository(ctrl),
		ProjectRepository: repository.ProjectRepository(),
	}
}

type ProjectRepo struct {
	domain.ProjectRepository
	mock *MockProjectRepository
}

func (r *ProjectRepo) EXPECT() *MockProjectRepositoryMockRecorder {
	return r.mock.EXPECT()
}

func (r *ProjectRepo) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Project, error) {
	return r.mock.Get(ctx, client, opts...)
}

func (r *ProjectRepo) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Project, error) {
	return r.mock.List(ctx, client, opts...)
}

func (r *ProjectRepo) Create(ctx context.Context, client database.QueryExecutor, org *domain.Project) error {
	return r.mock.Create(ctx, client, org)
}

func (r *ProjectRepo) Update(ctx context.Context, client database.QueryExecutor, conditions database.Condition, changes ...database.Change) (int64, error) {
	return r.mock.Update(ctx, client, conditions, changes...)
}

func (r *ProjectRepo) Delete(ctx context.Context, client database.QueryExecutor, conditions database.Condition) (int64, error) {
	return r.mock.Delete(ctx, client, conditions)
}
