package model

import (
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/project/model"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type Project struct {
	es_models.ObjectRoot
	Name                 string `json:"name,omitempty"`
	ProjectRoleAssertion bool   `json:"projectRoleAssertion,omitempty"`
	ProjectRoleCheck     bool   `json:"projectRoleCheck,omitempty"`
	HasProjectCheck      bool   `json:"hasProjectCheck,omitempty"`
	State                int32  `json:"-"`
	OIDCApplications     []*oidcApp
}

type oidcApp struct {
	AppID    string `json:"appId"`
	ClientID string `json:"clientId,omitempty"`
}

func ProjectToModel(project *Project) *model.Project {
	apps := make([]*model.Application, len(project.OIDCApplications))
	for i, application := range project.OIDCApplications {
		apps[i] = &model.Application{OIDCConfig: &model.OIDCConfig{ClientID: application.ClientID}}
	}
	return &model.Project{
		ObjectRoot:           project.ObjectRoot,
		Name:                 project.Name,
		ProjectRoleAssertion: project.ProjectRoleAssertion,
		ProjectRoleCheck:     project.ProjectRoleCheck,
		State:                model.ProjectState(project.State),
		Applications:         apps,
	}
}

func ProjectFromEvents(project *Project, events ...eventstore.Event) (*Project, error) {
	if project == nil {
		project = &Project{}
	}

	return project, project.AppendEvents(events...)
}

func (p *Project) AppendEvents(events ...eventstore.Event) error {
	for _, event := range events {
		if err := p.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) AppendEvent(event eventstore.Event) error {
	p.ObjectRoot.AppendEvent(event)

	switch event.Type() {
	case project.ProjectAddedType, project.ProjectChangedType:
		return p.AppendAddProjectEvent(event)
	case project.ProjectDeactivatedType:
		return p.appendDeactivatedEvent()
	case project.ProjectReactivatedType:
		return p.appendReactivatedEvent()
	case project.ProjectRemovedType:
		return p.appendRemovedEvent()
	case project.OIDCConfigAddedType:
		return p.appendOIDCConfig(event)
	case project.ApplicationRemovedType:
		return p.appendApplicationRemoved(event)
	}
	return nil
}

func (p *Project) AppendAddProjectEvent(event eventstore.Event) error {
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

func (p *Project) appendOIDCConfig(event eventstore.Event) error {
	appEvent := new(oidcApp)
	if err := event.Unmarshal(appEvent); err != nil {
		return err
	}
	p.OIDCApplications = append(p.OIDCApplications, appEvent)
	return nil
}

func (p *Project) appendApplicationRemoved(event eventstore.Event) error {
	appEvent := new(oidcApp)
	if err := event.Unmarshal(appEvent); err != nil {
		return err
	}
	for i := len(p.OIDCApplications) - 1; i >= 0; i-- {
		if p.OIDCApplications[i].AppID == appEvent.AppID {
			p.OIDCApplications[i] = p.OIDCApplications[len(p.OIDCApplications)-1]
			p.OIDCApplications[len(p.OIDCApplications)-1] = nil
			p.OIDCApplications = p.OIDCApplications[:len(p.OIDCApplications)-1]
			return nil
		}
	}
	return nil
}

func (p *Project) SetData(event eventstore.Event) error {
	if err := event.Unmarshal(p); err != nil {
		logging.Log("EVEN-lo9sr").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
