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

func (repo *projectRepo) UpdateOrg(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
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
		&org_model.OrgChange{ID: org.ID, Payload: changes, Sequence: currentOrg.Sequence},
		uniqueDomain,
		uniqueName,
	)
	return repo.OrgByID(ctx, org.ID)
}

func unique(ctx context.Context, id interface{}, isAvailable func(ctx context.Context, name string, sequence uint64) (uint64, error)) (uint64, error, bool) {
	sequence := uint64(0)
	if id == nil || id.(string) == "" {
		return 0, nil, false
	}
	sequence, err := isAvailable(ctx, id.(string), sequence)
	return sequence, err, true

}

func (repo *projectRepo) DeactivateOrg(ctx context.Context, id string) (*proj_model.Project, error) {
	org, err := repo.OrgByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := repo.orgEvents.IsOrgActive(ctx, org); err != nil {
		return nil, caos_errs.ThrowInvalidArgument(nil, "EVENT-r2fw1", "active")
	}

	org.Sequence, err = repo.orgEvents.DeactivateOrg(ctx, org.ID, org.Sequence)
	org.State = org_model.Inactive

	return org, err
}

func (repo *projectRepo) ReactivateOrg(ctx context.Context, id string) (*proj_model.Project, error) {
	org, err := repo.OrgByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := repo.orgEvents.IsOrgActive(ctx, org); err == nil {
		return nil, caos_errs.ThrowInvalidArgument(nil, "EVENT-r2fw1", "active")
	}

	org.Sequence, err = repo.orgEvents.ReactivateOrg(ctx, org.ID, org.Sequence)
	org.State = org_model.Active

	return org, err
}
