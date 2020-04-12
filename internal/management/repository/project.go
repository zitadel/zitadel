package repository

import (
	"context"
	"github.com/caos/zitadel/internal/project/model"
)

type ProjectRepository interface {
	ProjectByID(ctx context.Context, id string) (*model.Project, error)
	CreateProject(ctx context.Context, name string) (*model.Project, error)
	UpdateProject(ctx context.Context, project *model.Project) (*model.Project, error)
	DeactivateProject(ctx context.Context, id string) (*model.Project, error)
	ReactivateProject(ctx context.Context, id string) (*model.Project, error)

	ProjectMemberByID(ctx context.Context, projectID, userID string) (*model.ProjectMember, error)
	AddProjectMember(ctx context.Context, member *model.ProjectMember) (*model.ProjectMember, error)
	ChangeProjectMember(ctx context.Context, member *model.ProjectMember) (*model.ProjectMember, error)
	RemoveProjectMember(ctx context.Context, projectID, userID string) error

	AddProjectRole(ctx context.Context, role *model.ProjectRole) (*model.ProjectRole, error)
	ChangeProjectRole(ctx context.Context, role *model.ProjectRole) (*model.ProjectRole, error)
	RemoveProjectRole(ctx context.Context, projectID, key string) error
}
