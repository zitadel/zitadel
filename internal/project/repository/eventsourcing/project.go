package eventsourcing

import (
	"context"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	"github.com/sony/sonyflake"
	"strconv"
)

var idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})

type Project struct {
	es_models.ObjectRoot
	Name  string `json:"name,omitempty"`
	State int32  `json:"-"`
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

func ProjectCreateEvents(ctx context.Context, aggCreator *es_models.AggregateCreator, project *model.Project) (*es_models.Aggregate, error) {
	var err error
	id, err := idGenerator.NextID()
	if err != nil {
		return nil, err
	}
	project.ID = strconv.FormatUint(id, 10)

	return createdProject(ctx, aggCreator, project)
}

func ProjectUpdateEvents(agg *es_models.Aggregate, project *model.ProjectChange) (*es_models.Aggregate, error) {
	return updatedProject(agg, project)
}

func projectAggregate(ctx context.Context, aggCreator es_models.AggregateCreator, p *model.Project) (*es_models.Aggregate, error) {
	return aggCreator.NewAggregate(ctx, p.ID, model.ProjectAggregate, "v1", p.Sequence)
}

func createdProject(ctx context.Context, aggCreator *es_models.AggregateCreator, p *model.Project) (*es_models.Aggregate, error) {
	agg, err := projectAggregate(ctx, aggCreator, p)
	if err != nil {
		return nil, err
	}
	return agg.AppendEvent(model.AddedProject, p)
}

func updatedProject(agg *es_models.Aggregate, p *model.ProjectChange) (*es_models.Aggregate, error) {
	return agg.AppendEvent(model.ChangedProject, p)
}

func ProjectDeactivateEvents(agg *es_models.Aggregate) (*es_models.Aggregate, error) {
	return agg.AppendEvent(model.DeactivatedProject, nil)
}

func ProjectReactivateEvents(agg *es_models.Aggregate) (*es_models.Aggregate, error) {
	return agg.AppendEvent(model.ReactivatedProject, nil)
}
