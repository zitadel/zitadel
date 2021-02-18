package repository

import (
	"context"

	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/project/model"
)

type ProjectRepository interface {
	ProjectByID(ctx context.Context, id string) (*model.ProjectView, error)
	SearchProjects(ctx context.Context, request *model.ProjectViewSearchRequest) (*model.ProjectViewSearchResponse, error)

	ProjectGrantsByProjectIDAndRoleKey(ctx context.Context, projectID, roleKey string) ([]*model.ProjectGrantView, error)
	SearchProjectGrants(ctx context.Context, request *model.ProjectGrantViewSearchRequest) (*model.ProjectGrantViewSearchResponse, error)
	SearchGrantedProjects(ctx context.Context, request *model.ProjectGrantViewSearchRequest) (*model.ProjectGrantViewSearchResponse, error)
	ProjectGrantViewByID(ctx context.Context, grantID string) (*model.ProjectGrantView, error)

	ProjectMemberByID(ctx context.Context, projectID, userID string) (*model.ProjectMemberView, error)
	SearchProjectMembers(ctx context.Context, request *model.ProjectMemberSearchRequest) (*model.ProjectMemberSearchResponse, error)
	GetProjectMemberRoles(ctx context.Context) ([]string, error)

	SearchProjectRoles(ctx context.Context, projectId string, request *model.ProjectRoleSearchRequest) (*model.ProjectRoleSearchResponse, error)
	ProjectChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool) (*model.ProjectChanges, error)

	ApplicationByID(ctx context.Context, projectID, appID string) (*model.ApplicationView, error)
	AddApplication(ctx context.Context, app *model.Application) (*model.Application, error)
	ChangeApplication(ctx context.Context, app *model.Application) (*model.Application, error)
	DeactivateApplication(ctx context.Context, projectID, appID string) (*model.Application, error)
	ReactivateApplication(ctx context.Context, projectID, appID string) (*model.Application, error)
	RemoveApplication(ctx context.Context, projectID, appID string) error
	ChangeOIDCConfig(ctx context.Context, config *model.OIDCConfig) (*model.OIDCConfig, error)
	ChangeAPIConfig(ctx context.Context, config *model.APIConfig) (*model.APIConfig, error)
	ChangeOIDConfigSecret(ctx context.Context, projectID, appID string) (*model.OIDCConfig, error)
	ChangeAPIConfigSecret(ctx context.Context, projectID, appID string) (*model.APIConfig, error)
	SearchApplications(ctx context.Context, request *model.ApplicationSearchRequest) (*model.ApplicationSearchResponse, error)
	ApplicationChanges(ctx context.Context, id string, secId string, lastSequence uint64, limit uint64, sortAscending bool) (*model.ApplicationChanges, error)
	SearchClientKeys(ctx context.Context, request *key_model.AuthNKeySearchRequest) (*key_model.AuthNKeySearchResponse, error)
	GetClientKey(ctx context.Context, projectID, applicationID, keyID string) (*key_model.AuthNKeyView, error)
	AddClientKey(ctx context.Context, key *model.ClientKey) (*model.ClientKey, error)
	RemoveClientKey(ctx context.Context, projectID, applicationID, keyID string) error

	ProjectGrantByID(ctx context.Context, grantID string) (*model.ProjectGrantView, error)
	SearchProjectGrantMembers(ctx context.Context, request *model.ProjectGrantMemberSearchRequest) (*model.ProjectGrantMemberSearchResponse, error)

	ProjectGrantMemberByID(ctx context.Context, projectID, userID string) (*model.ProjectGrantMemberView, error)
	GetProjectGrantMemberRoles() []string
}
