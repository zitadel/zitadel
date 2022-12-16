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
	CSP          InstanceCSP

	DefaultOrgID string
	IAMProjectID string
	ConsoleID    string
	ConsoleAppID string
	DefaultLang  language.Tag
	Domains      []*InstanceDomain

	Host string
}

type InstanceCSP struct {
	Enabled        bool
	AllowedOrigins []string
}

func NewInstance(id, host string) *Instance {
	return &Instance{
		ID:   id,
		Host: host,
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
			i.reduceInstanceAdded(e)
		case *instance.InstanceChangedEvent:
			i.reduceInstanceChanged(e)
		case *instance.InstanceRemovedEvent:
			i.reduceInstanceRemoved(e)
		case *instance.DefaultOrgSetEvent:
			i.reduceDefaultOrgSet(e)
		case *instance.ProjectSetEvent:
			i.reduceProjectSet(e)
		case *instance.ConsoleSetEvent:
			i.reduceConsoleSet(e)
		case *instance.DefaultLanguageSetEvent:
			i.reduceDefaultLanguageSet(e)
		case *instance.DomainAddedEvent:
			i.reduceDomainAdded(e)
		case *instance.DomainPrimarySetEvent:
			i.reduceDomainPrimarySet(e)
		case *instance.DomainRemovedEvent:
			i.reduceDomainRemoved(e)
		case *instance.SecurityPolicySetEvent:
			i.reduceSecurityPolicySet(e)
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

			instance.SecurityPolicySetEventType,
		).
		Builder()
}

func (i *Instance) reduceInstanceAdded(event *instance.InstanceAddedEvent) {
	i.ID = event.Aggregate().ID
	i.CreationDate = event.CreationDate()
	i.Name = event.Name
}

func (i *Instance) reduceInstanceChanged(event *instance.InstanceChangedEvent) {
	i.Name = event.Name
}

func (i *Instance) reduceInstanceRemoved(event *instance.InstanceRemovedEvent) {
	i.Removed = true
}

func (i *Instance) reduceDefaultOrgSet(event *instance.DefaultOrgSetEvent) {
	i.DefaultOrgID = event.OrgID
}

func (i *Instance) reduceProjectSet(event *instance.ProjectSetEvent) {
	i.IAMProjectID = event.ProjectID
}

func (i *Instance) reduceConsoleSet(event *instance.ConsoleSetEvent) {
	i.ConsoleAppID = event.AppID
	i.ConsoleID = event.ClientID
}

func (i *Instance) reduceDefaultLanguageSet(event *instance.DefaultLanguageSetEvent) {
	i.DefaultLang = event.Language
}

func (i *Instance) reduceDomainAdded(event *instance.DomainAddedEvent) {
	i.Domains = append(i.Domains, &InstanceDomain{
		CreationDate: event.CreationDate(),
		ChangeDate:   event.CreationDate(),
		Sequence:     event.Sequence(),
		Domain:       event.Domain,
		InstanceID:   i.ID,
		IsGenerated:  event.Generated,
	})
}

func (i *Instance) reduceDomainPrimarySet(event *instance.DomainPrimarySetEvent) {
	for _, domain := range i.Domains {
		domain.IsPrimary = domain.Domain == event.Domain
	}
}

func (i *Instance) reduceDomainRemoved(event *instance.DomainRemovedEvent) {
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

func (i *Instance) reduceSecurityPolicySet(event *instance.SecurityPolicySetEvent) {
	if event.Enabled != nil {
		i.CSP.Enabled = *event.Enabled
	}
	if event.AllowedOrigins != nil {
		i.CSP.AllowedOrigins = *event.AllowedOrigins
	}
}
