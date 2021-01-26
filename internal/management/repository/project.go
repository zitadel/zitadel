package repository

import (
	"context"

	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/project/model"
)

type ProjectRepository interface {
	ProjectByID(ctx context.Context, id string) (*model.ProjectView, error)
	CreateProject(ctx context.Context, project *model.Project) (*model.Project, error)
	UpdateProject(ctx context.Context, project *model.Project) (*model.Project, error)
	DeactivateProject(ctx context.Context, id string) (*model.Project, error)
	ReactivateProject(ctx context.Context, id string) (*model.Project, error)
	RemoveProject(ctx context.Context, id string) error
	SearchProjects(ctx context.Context, request *model.ProjectViewSearchRequest) (*model.ProjectViewSearchResponse, error)
	SearchProjectGrants(ctx context.Context, request *model.ProjectGrantViewSearchRequest) (*model.ProjectGrantViewSearchResponse, error)
	SearchGrantedProjects(ctx context.Context, request *model.ProjectGrantViewSearchRequest) (*model.ProjectGrantViewSearchResponse, error)
	ProjectGrantViewByID(ctx context.Context, grantID string) (*model.ProjectGrantView, error)

	ProjectMemberByID(ctx context.Context, projectID, userID string) (*model.ProjectMemberView, error)
	AddProjectMember(ctx context.Context, member *model.ProjectMember) (*model.ProjectMember, error)
	ChangeProjectMember(ctx context.Context, member *model.ProjectMember) (*model.ProjectMember, error)
	RemoveProjectMember(ctx context.Context, projectID, userID string) error
	SearchProjectMembers(ctx context.Context, request *model.ProjectMemberSearchRequest) (*model.ProjectMemberSearchResponse, error)
	GetProjectMemberRoles(ctx context.Context) ([]string, error)

	AddProjectRole(ctx context.Context, role *model.ProjectRole) (*model.ProjectRole, error)
	ChangeProjectRole(ctx context.Context, role *model.ProjectRole) (*model.ProjectRole, error)
	RemoveProjectRole(ctx context.Context, projectID, key string) error
	SearchProjectRoles(ctx context.Context, projectId string, request *model.ProjectRoleSearchRequest) (*model.ProjectRoleSearchResponse, error)
	ProjectChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool) (*model.ProjectChanges, error)
	BulkAddProjectRole(ctx context.Context, role []*model.ProjectRole) error

	ApplicationByID(ctx context.Context, projectID, appID string) (*model.ApplicationView, error)
	AddApplication(ctx context.Context, app *model.Application) (*model.Application, error)
	ChangeApplication(ctx context.Context, app *model.Application) (*model.Application, error)
	DeactivateApplication(ctx context.Context, projectID, appID string) (*model.Application, error)
	ReactivateApplication(ctx context.Context, projectID, appID string) (*model.Application, error)
	RemoveApplication(ctx context.Context, projectID, appID string) error
	ChangeOIDCConfig(ctx context.Context, config *model.OIDCConfig) (*model.OIDCConfig, error)
	ChangeOIDConfigSecret(ctx context.Context, projectID, appID string) (*model.OIDCConfig, error)
	SearchApplications(ctx context.Context, request *model.ApplicationSearchRequest) (*model.ApplicationSearchResponse, error)
	ApplicationChanges(ctx context.Context, id string, secId string, lastSequence uint64, limit uint64, sortAscending bool) (*model.ApplicationChanges, error)
	SearchApplicationKeys(ctx context.Context, request *key_model.AuthNKeySearchRequest) (*key_model.AuthNKeySearchResponse, error)
	GetApplicationKey(ctx context.Context, projectID, applicationID, keyID string) (*key_model.AuthNKeyView, error)
	AddApplicationKey(ctx context.Context, key *model.ApplicationKey) (*model.ApplicationKey, error)
	RemoveApplicationKey(ctx context.Context, projectID, applicationID, keyID string) error

	ProjectGrantByID(ctx context.Context, grantID string) (*model.ProjectGrantView, error)
	AddProjectGrant(ctx context.Context, grant *model.ProjectGrant) (*model.ProjectGrant, error)
	ChangeProjectGrant(ctx context.Context, grant *model.ProjectGrant) (*model.ProjectGrant, error)
	DeactivateProjectGrant(ctx context.Context, projectID, grantID string) (*model.ProjectGrant, error)
	ReactivateProjectGrant(ctx context.Context, projectID, grantID string) (*model.ProjectGrant, error)
	RemoveProjectGrant(ctx context.Context, projectID, grantID string) error
	SearchProjectGrantMembers(ctx context.Context, request *model.ProjectGrantMemberSearchRequest) (*model.ProjectGrantMemberSearchResponse, error)

	ProjectGrantMemberByID(ctx context.Context, projectID, userID string) (*model.ProjectGrantMemberView, error)
	AddProjectGrantMember(ctx context.Context, member *model.ProjectGrantMember) (*model.ProjectGrantMember, error)
	ChangeProjectGrantMember(ctx context.Context, member *model.ProjectGrantMember) (*model.ProjectGrantMember, error)
	RemoveProjectGrantMember(ctx context.Context, projectID, grantID, userID string) error
	GetProjectGrantMemberRoles() []string
}
