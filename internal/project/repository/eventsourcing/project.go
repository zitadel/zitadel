package eventsourcing

import (
	"context"
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
	changes := make(map[string]interface{}, 2)
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

func ProjectByIDQuery(id string, latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.ProjectAggregate).
		LatestSequenceFilter(latestSequence).
		AggregateIDFilter(id)
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

	return agg.AppendEvent(model.AddedProject, project)
}

func ProjectUpdateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *Project, new *Project) (*es_models.Aggregate, error) {
	agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	changes := existing.Changes(new)
	return agg.AppendEvent(model.ChangedProject, changes)
}

func ProjectDeactivateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *Project) (*es_models.Aggregate, error) {
	agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	return agg.AppendEvent(model.DeactivatedProject, nil)
}

func ProjectReactivateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *Project) (*es_models.Aggregate, error) {
	agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	return agg.AppendEvent(model.ReactivatedProject, nil)
}
