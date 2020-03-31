package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
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
	return es_models.NewSearchQuery(0, false).
		AggregateTypeFilter(model.ProjectAggregate).
		LatestSequenceFilter(latestSequence).
		AggregateIDFilter(id)
}

func ProjectQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery(0, false).
		AggregateTypeFilter(model.ProjectAggregate).
		LatestSequenceFilter(latestSequence)
}

func ProjectCreateEvents(ctx context.Context, aggCreator eventstore.AggregateCreator, project *model.Project) (*eventstore.Aggregate, error) {
	var err error
	id, err := idGenerator.NextID()
	if err != nil {
		return nil, err
	}
	project.ID = strconv.FormatUint(id, 10)

	return createdProject(ctx, aggCreator, project)
}

func ProjectUpdateEvents(project *model.ProjectChange) *pkg.SaveAggregate {
	return updatedProject(project)
}

func projectAggregate(ctx context.Context, aggCreator eventstore.AggregateCreator, p *model.Project) (*eventstore.Aggregate, error) {
	return aggCreator.NewAggregate(ctx, p.ID, model.ProjectAggregate, "v1", p.Sequence)
}

func createdProject(ctx context.Context, aggCreator eventstore.AggregateCreator, p *model.Project) (*eventstore.Aggregate, error) {
	agg, err := projectAggregate(ctx, aggCreator, p)
	if err != nil {
		return nil, err
	}
	return agg.AppendEvent(model.AddedProject, p)
}

func updatedProject(p *model.ProjectChange) *pkg.SaveAggregate {
	return pkg.NewSaveAggregate(p.ID, model.ProjectAggregate, p.Sequence,
		pkg.SaveEvent{Type: model.ChangedProject, Payload: p.Payload},
	)
}

func ProjectDeactivateEvents(id string, sequence uint64) *pkg.SaveAggregate {
	return pkg.NewSaveAggregate(id, model.ProjectAggregate, sequence,
		pkg.SaveEvent{model.DeactivatedProject, nil},
	)
}

func ProjectReactivateEvents(id string, sequence uint64) *pkg.SaveAggregate {
	return pkg.NewSaveAggregate(id, model.ProjectAggregate, sequence,
		pkg.SaveEvent{model.ReactivatedProject, nil},
	)
}
