package eventsourcing

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	proj_model "github.com/caos/zitadel/internal/project/model"
)

type UserAgentEventstore struct {
	es_int.Eventstore
}

type UserAgentConfig struct {
	Eventstore es_int.Eventstore
	//Cache            *config.CacheConfig
}

func StartUserAgent(conf UserAgentConfig) (*UserAgentEventstore, error) {
	return &UserAgentEventstore{Eventstore: conf.Eventstore}, nil
}

func (es *UserAgentEventstore) UserAgentByID(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
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

func (es *UserAgentEventstore) CreateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
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
