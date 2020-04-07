package eventsourcing

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
)

const (
	projectVersion = "v1"
)

type Project struct {
	es_models.ObjectRoot
	Name  string `json:"name,omitempty"`
	State int32  `json:"-"`
}

func ProjectFromModel(project *model.Project) *Project {
	return &Project{
		ObjectRoot: es_models.ObjectRoot{
			ID:           project.ObjectRoot.ID,
			Sequence:     project.Sequence,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
		},
		Name:  project.Name,
		State: model.ProjectStateToInt(project.State),
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

func ProjectFromEvents(project *Project, events ...*es_models.Event) (*Project, error) {
	if project == nil {
		project = &Project{}
	}

	return project, project.AppendEvents(events...)
}

func (p *Project) Changes(changed *Project) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.Name != "" && p.Name != changed.Name {
		changes["name"] = changed.Name
	}
	return changes
}

func (p *Project) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := p.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) AppendEvent(event *es_models.Event) error {
	p.ObjectRoot.AppendEvent(event)

	switch event.Type {
	case model.ProjectAdded, model.ProjectChanged:
		if err := json.Unmarshal(event.Data, p); err != nil {
			logging.Log("EVEN-idl93").WithError(err).Error("could not unmarshal event data")
			return err
		}
		p.State = model.ProjectStateToInt(model.Active)
		return nil
	case model.ProjectDeactivated:
		return p.appendDeactivatedEvent()
	case model.ProjectReactivated:
		return p.appendReactivatedEvent()
	}

	return nil
}

func (p *Project) appendDeactivatedEvent() error {
	p.State = model.ProjectStateToInt(model.Inactive)
	return nil
}

func (p *Project) appendReactivatedEvent() error {
	p.State = model.ProjectStateToInt(model.Active)
	return nil
}

func (p *Project) ToCreateAggregate(aggCreator *es_models.AggregateCreator) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if p == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kdie6", "project should not be nil")
		}
		var err error
		id, err := idGenerator.NextID()
		if err != nil {
			return nil, err
		}
		p.ID = strconv.FormatUint(id, 10)

		agg, err := ProjectAggregate(ctx, aggCreator, p.ID, p.Sequence)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.ProjectAdded, p)
	}
}

func (p *Project) ToUpdateAggregate(aggCreator *es_models.AggregateCreator, updatedProject *Project) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if p == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dk93d", "existing project should not be nil")
		}
		if updatedProject == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "new project should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, p.ID, p.Sequence)
		if err != nil {
			return nil, err
		}
		changes := p.Changes(updatedProject)
		return agg.AppendEvent(model.ProjectChanged, changes)
	}
}

func (p *Project) ToDeactivateAggregate(aggCreator *es_models.AggregateCreator) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if p == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dk93d", "existing project should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, p.ID, p.Sequence)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.ProjectDeactivated, nil)
	}
}

func (p *Project) ToReactivateAggregate(aggCreator *es_models.AggregateCreator) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if p == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dk93d", "existing project should not be nil")
		}
		agg, err := ProjectAggregate(ctx, aggCreator, p.ID, p.Sequence)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.ProjectReactivated, nil)
	}
}
