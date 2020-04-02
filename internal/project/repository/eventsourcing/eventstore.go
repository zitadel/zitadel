package eventsourcing

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
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
	filter, err := ProjectByIDQuery(project.ID, project.Sequence)
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
	return ProjectToModel(foundProject), nil
}

func (es *ProjectEventstore) CreateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	project.State = proj_model.Active
	repoProject := ProjectFromModel(project)
	projectAggregate, err := ProjectCreateAggregate(ctx, es.Eventstore.AggregateCreator(), repoProject)
	if err != nil {
		return nil, err
	}
	err = es.PushAggregates(ctx, projectAggregate)
	if err != nil {
		return nil, err
	}

	repoProject.AppendEvents(projectAggregate.Events...)
	return ProjectToModel(repoProject), nil
}

func (es *ProjectEventstore) UpdateProject(ctx context.Context, existing *proj_model.Project, new *proj_model.Project) (*proj_model.Project, error) {
	repoExisting := ProjectFromModel(existing)
	repoNew := ProjectFromModel(new)
	projectAggregate, err := ProjectUpdateAggregate(ctx, es.AggregateCreator(), repoExisting, repoNew)
	if err != nil {
		return nil, err
	}
	err = es.PushAggregates(ctx, projectAggregate)
	if err != nil {
		return nil, err
	}
	repoExisting.AppendEvents(projectAggregate.Events...)
	return ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) DeactivateProject(ctx context.Context, existing *proj_model.Project) (*proj_model.Project, error) {
	if !existing.IsActive() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "EVENT-die45", "project must be active")
	}
	repoExisting := ProjectFromModel(existing)
	projectAggregate, err := ProjectDeactivateAggregate(ctx, es.AggregateCreator(), repoExisting)
	if err != nil {
		return nil, err
	}
	err = es.PushAggregates(ctx, projectAggregate)
	if err != nil {
		return nil, err
	}
	repoExisting.AppendEvents(projectAggregate.Events...)
	return ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) ReactivateProject(ctx context.Context, existing *proj_model.Project) (*proj_model.Project, error) {
	if existing.IsActive() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "EVENT-die45", "project must be inactive")
	}
	repoExisting := ProjectFromModel(existing)
	projectAggregate, err := ProjectReactivateAggregate(ctx, es.AggregateCreator(), repoExisting)
	if err != nil {
		return nil, err
	}
	err = es.PushAggregates(ctx, projectAggregate)
	if err != nil {
		return nil, err
	}
	repoExisting.AppendEvents(projectAggregate.Events...)
	return ProjectToModel(repoExisting), nil
}
