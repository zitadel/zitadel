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

func (es *ProjectEventstore) ProjectByID(ctx context.Context, project *proj_model.Project) error {
	filter := ProjectByIDQuery(project.ID, project.Sequence)
	events, err := es.Eventstore.FilterEvents(ctx, filter)
	if err != nil {
		return err
	}
	foundProject, err := FromEvents(nil, events...)

	*project = *ProjectToModel(foundProject)
	return err
}

func (es *ProjectEventstore) CreateProject(ctx context.Context, project *proj_model.Project) (err error) {
	projectAggregate, err := ProjectCreateEvents(ctx, es.Eventstore.AggregateCreator(), project)
	if err != nil {
		return err
	}
	err = es.PushAggregates(ctx, projectAggregate)
	if err != nil {
		return err
	}
	project.ID = projectAggregate.ID
	project.Sequence = projectAggregate.Events[len(projectAggregate.Events)-1].Sequence
	return es.ProjectByID(ctx, project)
}

func (es *ProjectEventstore) UpdateProject(ctx context.Context, existing *proj_model.Project, new *proj_model.Project) (err error) {
	projectAggregate, err := ProjectUpdateEvents(ctx, es.AggregateCreator(), existing, new)
	if err != nil {
		return err
	}
	err = es.PushAggregates(ctx, projectAggregate)
	if err != nil {
		return err
	}
	new.ID = existing.ID
	new.Sequence = existing.Sequence
	return es.ProjectByID(ctx, new)
}

func (es *ProjectEventstore) DeactivateProject(ctx context.Context, existing *proj_model.Project) error {
	if !existing.IsActive() {
		return caos_errs.ThrowInvalidArgument(nil, "EVENT-die45", "project must be active")
	}
	projectAggregate, err := ProjectDeactivateEvents(ctx, es.AggregateCreator(), existing)
	if err != nil {
		return err
	}
	err = es.PushAggregates(ctx, projectAggregate)
	if err != nil {
		return err
	}
	existing.Sequence = projectAggregate.Events[len(projectAggregate.Events)-1].Sequence
	existing.State = proj_model.Inactive
	return nil
}

func (es *ProjectEventstore) ReactivateProject(ctx context.Context, existing *proj_model.Project) error {
	if existing.IsActive() {
		return caos_errs.ThrowInvalidArgument(nil, "EVENT-die45", "project must be inactive")
	}
	projectAggregate, err := ProjectReactivateEvents(ctx, es.AggregateCreator(), existing)
	if err != nil {
		return err
	}
	err = es.PushAggregates(ctx, projectAggregate)
	if err != nil {
		return err
	}
	existing.Sequence = projectAggregate.Events[len(projectAggregate.Events)-1].Sequence
	existing.State = proj_model.Active
	return nil
}
