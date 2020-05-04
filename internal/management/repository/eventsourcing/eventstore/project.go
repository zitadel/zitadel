package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"

	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
)

type ProjectRepo struct {
	ProjectEvents *proj_event.ProjectEventstore
	View          *view.View
}

func (repo *ProjectRepo) ProjectByID(ctx context.Context, id string) (project *proj_model.Project, err error) {
	return repo.ProjectEvents.ProjectByID(ctx, id)
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

func (repo *ProjectRepo) SearchGrantedProjects(ctx context.Context, request *proj_model.GrantedProjectSearchRequest) (*proj_model.GrantedProjectSearchResponse, error) {
	projects, count, err := repo.View.SearchGrantedProjects(request)
	if err != nil {
		return nil, err
	}
	return &proj_model.GrantedProjectSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.GrantedProjectsToModel(projects),
	}, nil
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

func (repo *ProjectRepo) AddProjectRole(ctx context.Context, member *proj_model.ProjectRole) (*proj_model.ProjectRole, error) {
	return repo.ProjectEvents.AddProjectRole(ctx, member)
}

func (repo *ProjectRepo) ChangeProjectRole(ctx context.Context, member *proj_model.ProjectRole) (*proj_model.ProjectRole, error) {
	return repo.ProjectEvents.ChangeProjectRole(ctx, member)
}

func (repo *ProjectRepo) RemoveProjectRole(ctx context.Context, projectID, key string) error {
	member := proj_model.NewProjectRole(projectID, key)
	return repo.ProjectEvents.RemoveProjectRole(ctx, member)
}

func (repo *ProjectRepo) ApplicationByID(ctx context.Context, projectID, appID string) (app *proj_model.Application, err error) {
	return repo.ProjectEvents.ApplicationByIDs(ctx, projectID, appID)
}

func (repo *ProjectRepo) AddApplication(ctx context.Context, app *proj_model.Application) (*proj_model.Application, error) {
	return repo.ProjectEvents.AddApplication(ctx, app)
}

func (repo *ProjectRepo) ChangeApplication(ctx context.Context, app *proj_model.Application) (*proj_model.Application, error) {
	return repo.ProjectEvents.ChangeApplication(ctx, app)
}

func (repo *ProjectRepo) DeactivateApplication(ctx context.Context, projectID, appID string) (*proj_model.Application, error) {
	return repo.ProjectEvents.DeactivateApplication(ctx, projectID, appID)
}

func (repo *ProjectRepo) ReactivateApplication(ctx context.Context, projectID, appID string) (*proj_model.Application, error) {
	return repo.ProjectEvents.ReactivateApplication(ctx, projectID, appID)
}

func (repo *ProjectRepo) RemoveApplication(ctx context.Context, projectID, appID string) error {
	app := proj_model.NewApplication(projectID, appID)
	return repo.ProjectEvents.RemoveApplication(ctx, app)
}

func (repo *ProjectRepo) ChangeOIDCConfig(ctx context.Context, config *proj_model.OIDCConfig) (*proj_model.OIDCConfig, error) {
	return repo.ProjectEvents.ChangeOIDCConfig(ctx, config)
}

func (repo *ProjectRepo) ChangeOIDConfigSecret(ctx context.Context, projectID, appID string) (*proj_model.OIDCConfig, error) {
	return repo.ProjectEvents.ChangeOIDCConfigSecret(ctx, projectID, appID)
}

func (repo *ProjectRepo) ProjectGrantByID(ctx context.Context, projectID, appID string) (app *proj_model.ProjectGrant, err error) {
	return repo.ProjectEvents.ProjectGrantByIDs(ctx, projectID, appID)
}

func (repo *ProjectRepo) AddProjectGrant(ctx context.Context, app *proj_model.ProjectGrant) (*proj_model.ProjectGrant, error) {
	return repo.ProjectEvents.AddProjectGrant(ctx, app)
}

func (repo *ProjectRepo) ChangeProjectGrant(ctx context.Context, app *proj_model.ProjectGrant) (*proj_model.ProjectGrant, error) {
	return repo.ProjectEvents.ChangeProjectGrant(ctx, app)
}

func (repo *ProjectRepo) DeactivateProjectGrant(ctx context.Context, projectID, appID string) (*proj_model.ProjectGrant, error) {
	return repo.ProjectEvents.DeactivateProjectGrant(ctx, projectID, appID)
}

func (repo *ProjectRepo) ReactivateProjectGrant(ctx context.Context, projectID, appID string) (*proj_model.ProjectGrant, error) {
	return repo.ProjectEvents.ReactivateProjectGrant(ctx, projectID, appID)
}

func (repo *ProjectRepo) RemoveProjectGrant(ctx context.Context, projectID, appID string) error {
	app := proj_model.NewProjectGrant(projectID, appID)
	return repo.ProjectEvents.RemoveProjectGrant(ctx, app)
}

func (repo *ProjectRepo) ProjectGrantMemberByID(ctx context.Context, projectID, grantID, userID string) (member *proj_model.ProjectGrantMember, err error) {
	member = proj_model.NewProjectGrantMember(projectID, grantID, userID)
	return repo.ProjectEvents.ProjectGrantMemberByIDs(ctx, member)
}

func (repo *ProjectRepo) AddProjectGrantMember(ctx context.Context, member *proj_model.ProjectGrantMember) (*proj_model.ProjectGrantMember, error) {
	return repo.ProjectEvents.AddProjectGrantMember(ctx, member)
}

func (repo *ProjectRepo) ChangeProjectGrantMember(ctx context.Context, member *proj_model.ProjectGrantMember) (*proj_model.ProjectGrantMember, error) {
	return repo.ProjectEvents.ChangeProjectGrantMember(ctx, member)
}

func (repo *ProjectRepo) RemoveProjectGrantMember(ctx context.Context, projectID, grantID, userID string) error {
	member := proj_model.NewProjectGrantMember(projectID, grantID, userID)
	return repo.ProjectEvents.RemoveProjectGrantMember(ctx, member)
}
