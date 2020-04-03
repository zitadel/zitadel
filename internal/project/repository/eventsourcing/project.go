package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	"github.com/sony/sonyflake"
	"strconv"
)

var idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})

const (
	projectVersion = "v1"
)

type Project struct {
	es_models.ObjectRoot
	Name  string `json:"name,omitempty"`
	State int32  `json:"-"`
}

func (p *Project) Changes(changed *Project) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.Name != "" && p.Name != changed.Name {
		changes["name"] = changed.Name
	}
	return changes
}

func ProjectFromModel(project *model.Project) *Project {
	return &Project{
		Name:  project.Name,
		State: model.ProjectStateToInt(project.State),
		ObjectRoot: es_models.ObjectRoot{
			ID:           project.ObjectRoot.ID,
			Sequence:     project.Sequence,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
		},
	}
}

func ProjectToModel(project *Project) *model.Project {
	return &model.Project{
		ObjectRoot: es_models.ObjectRoot{
			ID:           project.ID,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
			Sequence:     project.Sequence,
		},
		Name:  project.Name,
		State: model.ProjectStateFromInt(project.State),
	}
}

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

func ProjectCreateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, project *Project) (*es_models.Aggregate, error) {
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

func ProjectUpdateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *Project, new *Project) (*es_models.Aggregate, error) {
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

func ProjectDeactivateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *Project) (*es_models.Aggregate, error) {
	if existing == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ueh45", "existing project should not be nil")
	}
	agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	return agg.AppendEvent(model.ProjectDeactivated, nil)
}

func ProjectReactivateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *Project) (*es_models.Aggregate, error) {
	if existing == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-37dur", "existing project should not be nil")
	}
	agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	return agg.AppendEvent(model.ProjectReactivated, nil)
}
