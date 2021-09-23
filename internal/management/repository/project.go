package repository

import (
	"context"
	"time"

	iam_model "github.com/caos/zitadel/internal/iam/model"

	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/project/model"
)

type ProjectRepository interface {
	ProjectMemberByID(ctx context.Context, projectID, userID string) (*model.ProjectMemberView, error)
	SearchProjectMembers(ctx context.Context, request *model.ProjectMemberSearchRequest) (*model.ProjectMemberSearchResponse, error)
	GetProjectMemberRoles(ctx context.Context) ([]string, error)

	ProjectChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*model.ProjectChanges, error)

	ApplicationByID(ctx context.Context, projectID, appID string) (*model.ApplicationView, error)
	SearchApplications(ctx context.Context, request *model.ApplicationSearchRequest) (*model.ApplicationSearchResponse, error)
	ApplicationChanges(ctx context.Context, projectID string, appID string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*model.ApplicationChanges, error)
	SearchClientKeys(ctx context.Context, request *key_model.AuthNKeySearchRequest) (*key_model.AuthNKeySearchResponse, error)
	GetClientKey(ctx context.Context, projectID, applicationID, keyID string) (*key_model.AuthNKeyView, error)

	ProjectGrantByID(ctx context.Context, grantID string) (*model.ProjectGrantView, error)
	SearchProjectGrantMembers(ctx context.Context, request *model.ProjectGrantMemberSearchRequest) (*model.ProjectGrantMemberSearchResponse, error)
	SearchProjectGrantRoles(ctx context.Context, projectID, grantID string, request *model.ProjectRoleSearchRequest) (*model.ProjectRoleSearchResponse, error)

	ProjectGrantMemberByID(ctx context.Context, projectID, userID string) (*model.ProjectGrantMemberView, error)
	GetProjectGrantMemberRoles() []string

	GetIAMByID(ctx context.Context) (*iam_model.IAM, error)
}
