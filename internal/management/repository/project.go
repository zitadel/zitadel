package repository

import (
	"context"
	"time"

	iam_model "github.com/caos/zitadel/internal/iam/model"

	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/project/model"
)

type ProjectRepository interface {
	GetProjectMemberRoles(ctx context.Context) ([]string, error)

	ProjectChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*model.ProjectChanges, error)

	ApplicationChanges(ctx context.Context, projectID string, appID string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*model.ApplicationChanges, error)
	SearchClientKeys(ctx context.Context, request *key_model.AuthNKeySearchRequest) (*key_model.AuthNKeySearchResponse, error)
	GetClientKey(ctx context.Context, projectID, applicationID, keyID string) (*key_model.AuthNKeyView, error)

	SearchProjectGrantMembers(ctx context.Context, request *model.ProjectGrantMemberSearchRequest) (*model.ProjectGrantMemberSearchResponse, error)

	ProjectGrantMemberByID(ctx context.Context, projectID, userID string) (*model.ProjectGrantMemberView, error)
	GetProjectGrantMemberRoles() []string

	GetIAMByID(ctx context.Context) (*iam_model.IAM, error)
}
