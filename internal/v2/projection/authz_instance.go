package projection

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/eventstore"
	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/instance"
)

type AuthZInstance struct {
	projection

	ID string

	DefaultOrgID    string
	ProjectID       string
	ConsoleClientID string
	ConsoleAppID    string
	DefaultLanguage language.Tag

	State *InstanceState
}

// InterestedIn implements model.
func (i *AuthZInstance) InterestedIn() map[eventstore.AggregateType][]eventstore.EventType {
	return map[eventstore.AggregateType][]eventstore.EventType{
		eventstore.AggregateType(instance.AggregateType): {
			eventstore.EventType(instance.AddedType),
			eventstore.EventType(instance.DefaultOrgSetType),
			eventstore.EventType(instance.ProjectSetType),
			eventstore.EventType(instance.ConsoleSetType),
			eventstore.EventType(instance.DefaultLanguageSetType),
			eventstore.EventType(instance.RemovedType),
		},
	}
}

func (i *AuthZInstance) Reduce(events ...*v2_es.StorageEvent) (err error) {
	for _, event := range events {
		if event.Aggregate.ID != i.ID {
			continue
		}
		i.projection.reduce(event)
		switch event.Type {
		case instance.AddedType:
			err = i.reduceAdded(event)
		case instance.DefaultOrgSetType:
			err = i.reduceDefaultOrgSet(event)
		case instance.ProjectSetType:
			err = i.reduceProjectSet(event)
		case instance.ConsoleSetType:
			err = i.reduceConsoleSet(event)
		case instance.DefaultLanguageSetType:
			err = i.reduceDefaultLanguageSet(event)
		}
		if err != nil {
			return err
		}
		if i.State == nil {
			i.State = NewInstanceStateProjection(i.ID)
		}
		if err = i.State.Reduce(event); err != nil {
			return err
		}

		i.position = event.Position
	}
	return nil
}

func (i *AuthZInstance) reduceAdded(event *v2_es.StorageEvent) error {
	e, err := instance.AddedEventFromStorage(event)
	if err != nil {
		return err
	}
	i.ID = e.Aggregate.ID
	return nil
}

func (i *AuthZInstance) reduceDefaultOrgSet(event *v2_es.StorageEvent) error {
	e, err := instance.DefaultOrgSetEventFromStorage(event)
	if err != nil {
		return err
	}
	i.DefaultOrgID = e.Payload.OrgID
	return nil
}

func (i *AuthZInstance) reduceProjectSet(event *v2_es.StorageEvent) error {
	e, err := instance.ProjectSetEventFromStorage(event)
	if err != nil {
		return err
	}
	i.ProjectID = e.Payload.ProjectID
	return nil
}

func (i *AuthZInstance) reduceConsoleSet(event *v2_es.StorageEvent) error {
	e, err := instance.ConsoleSetEventFromStorage(event)
	if err != nil {
		return err
	}
	i.ConsoleAppID = e.Payload.AppID
	i.ConsoleClientID = e.Payload.ClientID
	return nil
}

func (i *AuthZInstance) reduceDefaultLanguageSet(event *v2_es.StorageEvent) error {
	e, err := instance.DefaultLanguageSetEventFromStorage(event)
	if err != nil {
		return err
	}
	i.DefaultLanguage = e.Payload.Language
	return nil
}
