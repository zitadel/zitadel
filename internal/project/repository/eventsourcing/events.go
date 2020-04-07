package eventsourcing

import (
	"encoding/json"
	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
)

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
	case model.ProjectMemberAdded:
		return p.appendAddMemberEvent(event)
	case model.ProjectMemberChanged:
		return p.appendChangeMemberEvent(event)
	case model.ProjectMemberRemoved:
		return p.appendRemoveMemberEvent(event)
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

func (p *Project) appendAddMemberEvent(event *es_models.Event) error {
	member, err := getMemberData(event)
	if err != nil {
		return nil
	}
	member.ObjectRoot.CreationDate = event.CreationDate
	p.Members = append(p.Members, member)
	return nil
}

func (p *Project) appendChangeMemberEvent(event *es_models.Event) error {
	member, err := getMemberData(event)
	if err != nil {
		return nil
	}
	for i, m := range p.Members {
		if m.UserID == member.UserID {
			p.Members[i] = member
		}
	}
	return nil
}

func (p *Project) appendRemoveMemberEvent(event *es_models.Event) error {
	member, err := getMemberData(event)
	if err != nil {
		return nil
	}
	for i, m := range p.Members {
		if m.UserID == member.UserID {
			p.Members[i] = p.Members[len(p.Members)-1]
			p.Members[len(p.Members)-1] = nil
			p.Members = p.Members[:len(p.Members)-1]
		}
	}
	return nil
}

func getMemberData(event *es_models.Event) (*ProjectMember, error) {
	member := &ProjectMember{}
	member.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, member); err != nil {
		logging.Log("EVEN-e4dkp").WithError(err).Error("could not unmarshal event data")
		return nil, err
	}
	return member, nil
}
