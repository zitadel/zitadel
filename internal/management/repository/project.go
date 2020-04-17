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

	ApplicationByID(ctx context.Context, projectID, appID string) (*model.Application, error)
	AddApplication(ctx context.Context, app *model.Application) (*model.Application, error)
	ChangeApplication(ctx context.Context, app *model.Application) (*model.Application, error)
	DeactivateApplication(ctx context.Context, projectID, appID string) (*model.Application, error)
	ReactivateApplication(ctx context.Context, projectID, appID string) (*model.Application, error)
	RemoveApplication(ctx context.Context, projectID, appID string) error
	ChangeOIDCConfig(ctx context.Context, config *model.OIDCConfig) (*model.OIDCConfig, error)
	ChangeOIDConfigSecret(ctx context.Context, projectID, appID string) (*model.OIDCConfig, error)

	ProjectGrantByID(ctx context.Context, projectID, appID string) (*model.ProjectGrant, error)
	AddProjectGrant(ctx context.Context, app *model.ProjectGrant) (*model.ProjectGrant, error)
	ChangeProjectGrant(ctx context.Context, app *model.ProjectGrant) (*model.ProjectGrant, error)
	DeactivateProjectGrant(ctx context.Context, projectID, appID string) (*model.ProjectGrant, error)
	ReactivateProjectGrant(ctx context.Context, projectID, appID string) (*model.ProjectGrant, error)
	RemoveProjectGrant(ctx context.Context, projectID, appID string) error
}
