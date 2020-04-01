package eventsourcing

import (
	"context"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
)

type ProjectRepo struct {
	ProjectEvents *proj_event.ProjectEventstore
	//view      *view.View
}

func (repo *ProjectRepo) ProjectByID(ctx context.Context, id string) (project *proj_model.Project, err error) {
	//viewProject, err := repo.view.OrgByID(id)
	//if err != nil && !caos_errs.IsNotFound(err) {
	//	return nil, err
	//}
	//if viewProject != nil {
	//	project = org_view.OrgToModel(viewProject)
	//} else {
	project = proj_model.NewProject(id)
	//}
	return project, repo.ProjectEvents.ProjectByID(ctx, project)
}

func (repo *ProjectRepo) CreateProject(ctx context.Context, name string) (*proj_model.Project, error) {
	project := &proj_model.Project{Name: name}
	err := repo.ProjectEvents.CreateProject(ctx, project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (repo *ProjectRepo) UpdateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	existingProject, err := repo.ProjectByID(ctx, project.ID)
	if err != nil {
		return nil, err
	}

	err = repo.ProjectEvents.UpdateProject(ctx, existingProject, project)
	if err != nil {
		return nil, err
	}
	return project, err
}

func (repo *ProjectRepo) DeactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	project, err := repo.ProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = repo.ProjectEvents.DeactivateProject(ctx, project)
	if err != nil {
		return nil, err
	}
	return project, err
}

func (repo *ProjectRepo) ReactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	project, err := repo.ProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = repo.ProjectEvents.ReactivateProject(ctx, project)
	if err != nil {
		return nil, err
	}
	return project, err
}
