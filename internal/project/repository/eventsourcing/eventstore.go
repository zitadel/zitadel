package eventsourcing

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	proj_model "github.com/caos/zitadel/internal/project/model"
)

type ProjectEventstore struct {
	es_int.Eventstore
}

type ProjectConfig struct {
	es_int.Eventstore
}

func StartProject(conf ProjectConfig) (*ProjectEventstore, error) {
	return &ProjectEventstore{Eventstore: conf.Eventstore}, nil
}

func (es *ProjectEventstore) ProjectByID(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	query, err := ProjectByIDQuery(project.ID, project.Sequence)
	if err != nil {
		return nil, err
	}

	p := ProjectFromModel(project)
	err = es_sdk.Filter(ctx, es.FilterEvents, p.AppendEvents, query)
	if err != nil {
		return nil, err
	}
	return ProjectToModel(p), nil
}

func (es *ProjectEventstore) CreateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	if !project.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Name is required")
	}
	project.State = proj_model.Active
	repoProject := ProjectFromModel(project)

	createAggregate := ProjectCreateAggregate(es.AggregateCreator(), repoProject)
	err := es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	return ProjectToModel(repoProject), nil
}

func (es *ProjectEventstore) UpdateProject(ctx context.Context, existingProject *proj_model.Project, project *proj_model.Project) (*proj_model.Project, error) {
	if !project.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Name is required")
	}
	repoExisting := ProjectFromModel(existingProject)
	repoNew := ProjectFromModel(project)

	updateAggregate := ProjectUpdateAggregate(es.AggregateCreator(), repoExisting, repoNew)
	err := es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return nil, err
	}

	return ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) DeactivateProject(ctx context.Context, existing *proj_model.Project) (*proj_model.Project, error) {
	if !existing.IsActive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die45", "project must be active")
	}

	repoExisting := ProjectFromModel(existing)
	aggregate := ProjectDeactivateAggregate(es.AggregateCreator(), repoExisting)
	es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, aggregate)
	return ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) ReactivateProject(ctx context.Context, existing *proj_model.Project) (*proj_model.Project, error) {
	if existing.IsActive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die45", "project must be inactive")
	}

	repoExisting := ProjectFromModel(existing)
	aggregate := ProjectReactivateAggregate(es.AggregateCreator(), repoExisting)
	es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, aggregate)
	return ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) ProjectMemberByIDs(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	if member.UserID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-ld93d", "userID missing")
	}
	filter, err := ProjectByIDQuery(member.ID, member.Sequence)
	if err != nil {
		return nil, err
	}
	events, err := es.Eventstore.FilterEvents(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-8due3", "Could not find project events")
	}
	foundProject, err := ProjectFromEvents(nil, events...)
	if err != nil {
		return nil, err
	}
	for _, m := range foundProject.Members {
		if m.UserID == member.UserID {
			return ProjectMemberToModel(m), nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Could not find member in list")
}

func (es *ProjectEventstore) AddProjectMember(ctx context.Context, existing *proj_model.Project, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "UserID and Roles are required")
	}
	if existing.ContainsMember(member) {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-idke6", "User is already member of this Project")
	}
	repoProject := ProjectFromModel(existing)
	repoMember := ProjectMemberFromModel(member)
	projectAggregate, err := ProjectMemberAddedAggregate(ctx, es.Eventstore.AggregateCreator(), repoProject, repoMember)
	if err != nil {
		return nil, err
	}
	err = es.PushAggregates(ctx, projectAggregate)
	if err != nil {
		return nil, err
	}

	repoProject.AppendEvents(projectAggregate.Events...)
	for _, m := range repoProject.Members {
		if m.UserID == member.UserID {
			return ProjectMemberToModel(m), nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Could not find member in list")
}

func (es *ProjectEventstore) ChangeProjectMember(ctx context.Context, existing *proj_model.Project, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "UserID and Roles are required")
	}
	if !existing.ContainsMember(member) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-oe39f", "User is not member of this project")
	}
	repoProject := ProjectFromModel(existing)
	repoMember := ProjectMemberFromModel(member)
	projectAggregate, err := ProjectMemberChangedAggregate(ctx, es.Eventstore.AggregateCreator(), repoProject, repoMember)
	if err != nil {
		return nil, err
	}
	err = es.PushAggregates(ctx, projectAggregate)
	if err != nil {
		return nil, err
	}

	repoProject.AppendEvents(projectAggregate.Events...)
	for _, m := range repoProject.Members {
		if m.UserID == member.UserID {
			return ProjectMemberToModel(m), nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Could not find member in list")
}

func (es *ProjectEventstore) RemoveProjectMember(ctx context.Context, existing *proj_model.Project, member *proj_model.ProjectMember) error {
	if member.UserID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-d43fs", "UserID and Roles are required")
	}
	if !existing.ContainsMember(member) {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-swf34", "User is not member of this project")
	}
	repoProject := ProjectFromModel(existing)
	repoMember := ProjectMemberFromModel(member)
	projectAggregate, err := ProjectMemberRemovedAggregate(ctx, es.Eventstore.AggregateCreator(), repoProject, repoMember)
	if err != nil {
		return err
	}
	err = es.PushAggregates(ctx, projectAggregate)
	if err != nil {
		return err
	}

	repoProject.AppendEvents(projectAggregate.Events...)
	return nil
}
