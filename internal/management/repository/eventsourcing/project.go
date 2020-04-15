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
	//viewProject, err := repo.view.ProjectByID(id)
	//if err != nil && !caos_errs.IsNotFound(err) {
	//	return nil, err
	//}
	//if viewProject != nil {
	//	project = org_view.ProjectToModel(viewProject)
	//} else {
	project = proj_model.NewProject(id)
	//}
	return repo.ProjectEvents.ProjectByID(ctx, project)
}

func (repo *ProjectRepo) CreateProject(ctx context.Context, name string) (*proj_model.Project, error) {
	project := &proj_model.Project{Name: name}
	return repo.ProjectEvents.CreateProject(ctx, project)
}

func (repo *ProjectRepo) UpdateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	return repo.ProjectEvents.UpdateProject(ctx, project)
}

func (repo *ProjectRepo) DeactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	return repo.ProjectEvents.DeactivateProject(ctx, id)
}

func (repo *ProjectRepo) ReactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	return repo.ProjectEvents.ReactivateProject(ctx, id)
}

func (repo *ProjectRepo) ProjectMemberByID(ctx context.Context, projectID, userID string) (member *proj_model.ProjectMember, err error) {
	member = proj_model.NewProjectMember(projectID, userID)
	return repo.ProjectEvents.ProjectMemberByIDs(ctx, member)
}

func (repo *ProjectRepo) AddProjectMember(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	return repo.ProjectEvents.AddProjectMember(ctx, member)
}

func (repo *ProjectRepo) ChangeProjectMember(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	return repo.ProjectEvents.ChangeProjectMember(ctx, member)
}

func (repo *ProjectRepo) RemoveProjectMember(ctx context.Context, projectID, userID string) error {
	member := proj_model.NewProjectMember(projectID, userID)
	return repo.ProjectEvents.RemoveProjectMember(ctx, member)
}
