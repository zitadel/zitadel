package projection

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

var _ Projection = (*Instance)(nil)

type Instance struct {
	ID string

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
	Name         string
	Removed      bool

	DefaultOrgID string
	IAMProjectID string
	ConsoleID    string
	ConsoleAppID string
	DefaultLang  language.Tag
	Domains      []*InstanceDomain
}

func NewInstance(id string) *Instance {
	return &Instance{
		ID: id,
	}
}

type InstanceDomain struct {
	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
	Domain       string
	InstanceID   string
	IsGenerated  bool
	IsPrimary    bool
}

func (i *Instance) Reduce(events []eventstore.Event) {
	for _, event := range events {
		i.ChangeDate = event.CreationDate()
		i.Sequence = event.Sequence()

		switch e := event.(type) {
		case *instance.InstanceAddedEvent:
			i.reduceInstanceAddedEvent(e)
		case *instance.InstanceChangedEvent:
			i.reduceInstanceChangedEvent(e)
		case *instance.InstanceRemovedEvent:
			i.reduceInstanceRemovedEvent(e)
		case *instance.DefaultOrgSetEvent:
			i.reduceDefaultOrgSetEvent(e)
		case *instance.ProjectSetEvent:
			i.reduceProjectSetEvent(e)
		case *instance.ConsoleSetEvent:
			i.reduceConsoleSetEvent(e)
		case *instance.DefaultLanguageSetEvent:
			i.reduceDefaultLanguageSetEvent(e)
		case *instance.DomainAddedEvent:
			i.reduceDomainAddedEvent(e)
		case *instance.DomainPrimarySetEvent:
			i.reduceDomainPrimarySetEvent(e)
		case *instance.DomainRemovedEvent:
			i.reduceDomainRemovedEvent(e)
		}
	}
}

func (i *Instance) SearchQuery(context.Context) *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(i.ID).
		OrderAsc().
		AddQuery().
		AggregateIDs(i.ID).
		AggregateTypes(instance.AggregateType).
		EventTypes(
			instance.InstanceAddedEventType,
			instance.InstanceChangedEventType,
			instance.InstanceRemovedEventType,

			instance.DefaultOrgSetEventType,
			instance.ProjectSetEventType,
			instance.ConsoleSetEventType,
			instance.DefaultLanguageSetEventType,

			instance.InstanceDomainAddedEventType,
			instance.InstanceDomainPrimarySetEventType,
			instance.InstanceDomainRemovedEventType,
		).
		Builder()
}

func (i *Instance) reduceInstanceAddedEvent(event *instance.InstanceAddedEvent) {
	i.ID = event.Aggregate().ID
	i.CreationDate = event.CreationDate()
	i.Name = event.Name
}

func (i *Instance) reduceInstanceChangedEvent(event *instance.InstanceChangedEvent) {
	i.Name = event.Name
}

func (i *Instance) reduceInstanceRemovedEvent(event *instance.InstanceRemovedEvent) {
	i.Removed = true
}

func (i *Instance) reduceDefaultOrgSetEvent(event *instance.DefaultOrgSetEvent) {
	i.DefaultOrgID = event.OrgID
}

func (i *Instance) reduceProjectSetEvent(event *instance.ProjectSetEvent) {
	i.IAMProjectID = event.ProjectID
}

func (i *Instance) reduceConsoleSetEvent(event *instance.ConsoleSetEvent) {
	i.ConsoleAppID = event.AppID
	i.ConsoleID = event.ClientID
}

func (i *Instance) reduceDefaultLanguageSetEvent(event *instance.DefaultLanguageSetEvent) {
	i.DefaultLang = event.Language
}

func (i *Instance) reduceDomainAddedEvent(event *instance.DomainAddedEvent) {
	i.Domains = append(i.Domains, &InstanceDomain{
		CreationDate: event.CreationDate(),
		ChangeDate:   event.CreationDate(),
		Sequence:     event.Sequence(),
		Domain:       event.Domain,
		InstanceID:   i.ID,
		IsGenerated:  event.Generated,
	})
}

func (i *Instance) reduceDomainPrimarySetEvent(event *instance.DomainPrimarySetEvent) {
	for _, domain := range i.Domains {
		domain.IsPrimary = domain.Domain == event.Domain
	}
}

func (i *Instance) reduceDomainRemovedEvent(event *instance.DomainRemovedEvent) {
	for idx, domain := range i.Domains {
		if domain.Domain != event.Domain {
			continue
		}
		i.Domains[idx] = i.Domains[len(i.Domains)-1]
		i.Domains[len(i.Domains)-1] = nil
		i.Domains = i.Domains[:len(i.Domains)-1]
		return
	}
}
