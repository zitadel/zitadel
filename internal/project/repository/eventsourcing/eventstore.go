package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	proj_model "github.com/caos/zitadel/internal/project/model"
)

var eventstore es_int.Eventstore

type ProjectEventstore struct {
	es_int.Eventstore
}

type ProjectConfig struct {
	es_int.Eventstore
}

func StartProject(conf ProjectConfig) (*ProjectEventstore, error) {
	eventstore = conf.Eventstore
	if eventstore == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-ye89h", "config must set eventstore")
	}
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

	createAggregate := repoProject.ToCreateAggregate(es.AggregateCreator())
	err := es_sdk.Save(ctx, es.PushAggregates, createAggregate, repoProject.AppendEvents)
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

	updateAggregate := repoExisting.ToUpdateAggregate(es.AggregateCreator(), repoNew)
	err := es_sdk.Save(ctx, es.PushAggregates, updateAggregate, repoExisting.AppendEvents)
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
	es_sdk.Save(ctx, es.PushAggregates, repoExisting.ToDeactivateAggregate(es.AggregateCreator()), repoExisting.AppendEvents)
	return ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) ReactivateProject(ctx context.Context, existing *proj_model.Project) (*proj_model.Project, error) {
	if existing.IsActive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die45", "project must be inactive")
	}

	repoExisting := ProjectFromModel(existing)
	es_sdk.Save(ctx, es.PushAggregates, repoExisting.ToReactivateAggregate(es.AggregateCreator()), repoExisting.AppendEvents)
	return ProjectToModel(repoExisting), nil
}
