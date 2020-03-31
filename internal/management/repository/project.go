package repository

import (
	"context"
	"github.com/caos/zitadel/internal/project/model"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, name string) (*model.Project, error)
	UpdateProject(ctx context.Context, org *model.Project) (*model.Project, error)
	DeactivateProject(ctx context.Context, id string) (*model.Project, error)
	ReactivateProject(ctx context.Context, id string) (*model.Project, error)
}
