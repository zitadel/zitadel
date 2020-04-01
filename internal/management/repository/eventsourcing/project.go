package eventsourcing

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/eventstore"
	proj_model "github.com/caos/zitadel/internal/project/model"
)

type projectRepo struct {
	projectEvents *eventstore.ProjectEventstore
	//view      *view.View
}

func (repo *projectRepo) ProjectByID(ctx context.Context, id string) (project *proj_model.Project, err error) {
	//viewProject, err := repo.view.OrgByID(id)
	//if err != nil && !caos_errs.IsNotFound(err) {
	//	return nil, err
	//}
	//if viewProject != nil {
	//	project = org_view.OrgToModel(viewProject)
	//} else {
	project = proj_model.NewProject(id)
	//}
	return project, repo.projectEvents.ProjectByID(ctx, project)
}

func (repo *projectRepo) CreateProject(ctx context.Context, name string) (*proj_model.Project, error) {
	id, err := repo.projectEvents.CreateProject(ctx, name)
	if err != nil {
		return nil, err
	}

	return repo.ProjectByID(ctx, id)
}

func (repo *projectRepo) UpdateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	currentProject, err := repo.ProjectByID(ctx, project.ID)
	if err != nil {
		return nil, err
	}

	changes := currentProject.Changes(project)
	if len(changes) == 0 {
		return currentProject, nil
	}

	project.Sequence = currentProject.Sequence

	project.Sequence, err = repo.projectEvents.UpdateProject(ctx,
		&proj_model.ProjectChange{ID: project.ID, Payload: changes, Sequence: currentOrg.Sequence},
	)
	return repo.ProjectByID(ctx, project.ID)
}

func (repo *projectRepo) DeactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	project, err := repo.ProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := repo.projectEvents.IsProjectActive(ctx, project); err != nil {
		return nil, caos_errs.ThrowInvalidArgument(nil, "EVENT-r2fw1", "active")
	}

	project.Sequence, err = repo.projectEvents.DeactivateProject(ctx, project.ID, project.Sequence)
	project.State = proj_model.Inactive

	return project, err
}

func (repo *projectRepo) ReactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	project, err := repo.ProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := repo.projectEvents.IsProjectActive(ctx, project); err == nil {
		return nil, caos_errs.ThrowInvalidArgument(nil, "EVENT-r2fw1", "active")
	}

	project.Sequence, err = repo.projectEvents.ReactivateOrg(ctx, project.ID, project.Sequence)
	project.State = proj_model.Active

	return project, err
}
