package repository

import (
	"context"
	"time"

	iam_model "github.com/caos/zitadel/internal/iam/model"

	"github.com/caos/zitadel/internal/project/model"
)

type ProjectRepository interface {
	GetProjectMemberRoles(ctx context.Context) ([]string, error)

	ProjectChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*model.ProjectChanges, error)

	ApplicationChanges(ctx context.Context, projectID string, appID string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*model.ApplicationChanges, error)

	GetProjectGrantMemberRoles() []string

	GetIAMByID(ctx context.Context) (*iam_model.IAM, error)
}
