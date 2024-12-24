package model

import (
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/groups/model"
	"github.com/zitadel/zitadel/internal/repository/group"
)

type Group struct {
	es_models.ObjectRoot
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	State       int32  `json:"-"`
}

func GroupToModel(group *Group) *model.Group {
	return &model.Group{
		ObjectRoot:  group.ObjectRoot,
		Name:        group.Name,
		Description: group.Description,
		State:       model.GroupState(group.State),
	}
}

func GroupFromEvents(group *Group, events ...eventstore.Event) (*Group, error) {
	if group == nil {
		group = &Group{}
	}

	return group, group.AppendEvents(events...)
}

func (g *Group) AppendEvents(events ...eventstore.Event) error {
	for _, event := range events {
		if err := g.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (g *Group) AppendEvent(event eventstore.Event) error {
	g.ObjectRoot.AppendEvent(event)

	switch event.Type() {
	case group.GroupAddedType, group.GroupChangedType:
		return g.AppendAddGroupEvent(event)
	case group.GroupDeactivatedType:
		return g.appendDeactivatedEvent()
	case group.GroupReactivatedType:
		return g.appendReactivatedEvent()
	case group.GroupRemovedType:
		return g.appendRemovedEvent()
	}
	return nil
}

func (g *Group) AppendAddGroupEvent(event eventstore.Event) error {
	if err := g.SetData(event); err != nil {
		return err
	}
	g.State = int32(model.GroupStateActive)
	return nil
}

func (g *Group) appendDeactivatedEvent() error {
	g.State = int32(model.GroupStateInactive)
	return nil
}

func (g *Group) appendReactivatedEvent() error {
	g.State = int32(model.GroupStateActive)
	return nil
}

func (g *Group) appendRemovedEvent() error {
	g.State = int32(model.GroupStateRemoved)
	return nil
}

func (g *Group) SetData(event eventstore.Event) error {
	if err := event.Unmarshal(g); err != nil {
		logging.Log("EVEN-lo8tr").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
