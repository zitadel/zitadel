package eventsourcing

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	proj_model "github.com/caos/zitadel/internal/project/model"
	agent_model "github.com/caos/zitadel/internal/user_agent/model"
	"github.com/caos/zitadel/internal/user_agent/repository/eventsourcing/model"
)

type UserAgentEventstore struct {
	es_int.Eventstore
	agentCache *UserAgentCache
}

type UserAgentConfig struct {
	Eventstore es_int.Eventstore
	//Cache            *config.CacheConfig
}

func StartUserAgent(conf UserAgentConfig) (*UserAgentEventstore, error) {
	return &UserAgentEventstore{Eventstore: conf.Eventstore}, nil
}

func (es *UserAgentEventstore) UserAgentByID(ctx context.Context, id string) (*agent_model.UserAgent, error) {
	userAgent, sequence := es.agentCache.getUserAgent(id)

	query, err := UserAgentByIDQuery(userAgent.ID, userAgent.Sequence)
	if err != nil {
		return nil, err
	}

	agent := model.UserAgentFromModel(userAgent)
	err = es_sdk.Filter(ctx, es.FilterEvents, agent.AppendEvents, query)
	if err != nil {
		return nil, err
	}
	return model.UserAgentToModel(agent), nil
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
