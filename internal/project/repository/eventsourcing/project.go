package eventsourcing

import (
	"context"
	"strconv"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	"github.com/sony/sonyflake"
)

var idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})

func ProjectByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dke74", "id should be filled")
	}
	return ProjectQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func ProjectQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.ProjectAggregate).
		LatestSequenceFilter(latestSequence)
}

func ProjectAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, id string, sequence uint64) (*es_models.Aggregate, error) {
	return aggCreator.NewAggregate(ctx, id, model.ProjectAggregate, projectVersion, sequence)
}

func ProjectCreateAggregate(aggCreator *es_models.AggregateCreator, project *Project) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if project == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kdie6", "project should not be nil")
		}
		var err error
		id, err := idGenerator.NextID()
		if err != nil {
			return nil, err
		}
		project.ID = strconv.FormatUint(id, 10)

		agg, err := ProjectAggregate(ctx, aggCreator, project.ID, project.Sequence)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.ProjectAdded, project)
	}
}

func ProjectUpdateAggregate(aggCreator *es_models.AggregateCreator, existing *Project, new *Project) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dk93d", "existing project should not be nil")
		}
		if new == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "new project should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		changes := existing.Changes(new)
		return agg.AppendEvent(model.ProjectChanged, changes)
	}
}

func ProjectDeactivateAggregate(aggCreator *es_models.AggregateCreator, project *Project) func(ctx context.Context) (*es_models.Aggregate, error) {
	return projectStateAggregate(aggCreator, project, model.ProjectDeactivated)
}

func ProjectReactivateAggregate(aggCreator *es_models.AggregateCreator, project *Project) func(ctx context.Context) (*es_models.Aggregate, error) {
	return projectStateAggregate(aggCreator, project, model.ProjectReactivated)
}

func projectStateAggregate(aggCreator *es_models.AggregateCreator, project *Project, state models.EventType) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if project == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-37dur", "existing project should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, project.ID, project.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(state, nil)
	}
}

func ProjectMemberAddedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, member *ProjectMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-di38f", "existing project should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectMemberAdded, member)
	}
}

func ProjectMemberChangedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, member *ProjectMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sle3d", "existing project should not be nil")
		}

		agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectMemberChanged, member)
	}
}

func ProjectMemberRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *Project, member *ProjectMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slo9e", "existing project should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.ProjectMemberRemoved, member)
	}
}
