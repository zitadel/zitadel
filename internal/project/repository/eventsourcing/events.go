package eventsourcing

import (
	"github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/proto"
)

func FromEvents(project *Project, events ...*es_api.EventResponse) (*Project, error) {
	if project == nil {
		project = &Project{}
	}

	return project, project.AppendEvents(events...)
}

func (p *Project) AppendEvents(events ...*es_api.EventResponse) error {
	for _, event := range events {
		if err := p.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) AppendEvent(event *es_api.EventResponse) error {
	p.ObjectRoot.AppendEvent(event)

	switch event.Type {
	case model.AddedProject, model.ChangedProject:
		p.State = model.ProjectStateToInt(model.Active)
		return proto.FromPBStruct(p, event.Data)
	case model.DeactivatedProject:
		return p.appendDeactivatedEvent()
	case model.ReactivatedProject:
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
