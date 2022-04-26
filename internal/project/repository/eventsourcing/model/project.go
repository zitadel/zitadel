package model

import (
	"encoding/json"

	"github.com/caos/logging"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/project/model"
)

const (
	ProjectVersion = "v1"
)

type Project struct {
	es_models.ObjectRoot
	Name                 string `json:"name,omitempty"`
	ProjectRoleAssertion bool   `json:"projectRoleAssertion,omitempty"`
	ProjectRoleCheck     bool   `json:"projectRoleCheck,omitempty"`
	HasProjectCheck      bool   `json:"hasProjectCheck,omitempty"`
	State                int32  `json:"-"`
}

func ProjectToModel(project *Project) *model.Project {
	return &model.Project{
		ObjectRoot:           project.ObjectRoot,
		Name:                 project.Name,
		ProjectRoleAssertion: project.ProjectRoleAssertion,
		ProjectRoleCheck:     project.ProjectRoleCheck,
		State:                model.ProjectState(project.State),
	}
}

func ProjectFromEvents(project *Project, events ...*es_models.Event) (*Project, error) {
	if project == nil {
		project = &Project{}
	}

	return project, project.AppendEvents(events...)
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
	case ProjectAdded, ProjectChanged:
		return p.AppendAddProjectEvent(event)
	case ProjectDeactivated:
		return p.appendDeactivatedEvent()
	case ProjectReactivated:
		return p.appendReactivatedEvent()
	case ProjectRemoved:
		return p.appendRemovedEvent()
	}
	return nil
}

func (p *Project) AppendAddProjectEvent(event *es_models.Event) error {
	p.SetData(event)
	p.State = int32(model.ProjectStateActive)
	return nil
}

func (p *Project) appendDeactivatedEvent() error {
	p.State = int32(model.ProjectStateInactive)
	return nil
}

func (p *Project) appendReactivatedEvent() error {
	p.State = int32(model.ProjectStateActive)
	return nil
}

func (p *Project) appendRemovedEvent() error {
	p.State = int32(model.ProjectStateRemoved)
	return nil
}

func (p *Project) SetData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, p); err != nil {
		logging.Log("EVEN-lo9sr").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
