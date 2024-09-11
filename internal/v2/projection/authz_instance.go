package projection

import (
	"slices"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/instance"
)

type AuthZInstance struct {
	Projection

	ID string

	DefaultOrgID    string
	ProjectID       string
	ConsoleClientID string
	ConsoleAppID    string
	DefaultLanguage language.Tag

	Domains []*InstanceDomain

	State *InstanceState
}

func NewAuthZInstanceFromEvent(event *v2_es.StorageEvent) *AuthZInstance {
	instance := &AuthZInstance{
		ID: event.Aggregate.ID,
	}
	err := instance.reduceAdded(event)
	logging.OnError(err).Error("could not reduce added event")

	return instance
}

func (i *AuthZInstance) Reducers() Reducers {
	if i.Projection.Reducers != nil {
		return i.Projection.Reducers
	}
	i.Projection.Reducers = Reducers{
		instance.AggregateType: {
			instance.AddedType:              i.reduceAdded,
			instance.DefaultOrgSetType:      i.reduceDefaultOrgSet,
			instance.ProjectSetType:         i.reduceProjectSet,
			instance.ConsoleSetType:         i.reduceConsoleSet,
			instance.DefaultLanguageSetType: i.reduceDefaultLanguageSet,
			instance.RemovedType:            i.reduceRemoved,

			instance.DomainAddedType:      i.reduceDomainAdded,
			instance.DomainVerifiedType:   i.reduceDomainVerified,
			instance.DomainPrimarySetType: i.reduceDomainPrimarySet,
			instance.DomainRemovedType:    i.reduceDomainRemoved,
		},
	}

	return i.Projection.Reducers
}

func (i *AuthZInstance) reduceAdded(event *v2_es.StorageEvent) error {
	if !i.ShouldReduce(event) {
		return nil
	}

	if i.State == nil {
		i.State = NewInstanceStateProjection(i.ID)
	}
	return i.Projection.Reduce(event, i.State.reduceAdded)
}

func (i *AuthZInstance) reduceDefaultOrgSet(event *v2_es.StorageEvent) error {
	if !i.ShouldReduce(event) {
		return nil
	}

	e, err := instance.DefaultOrgSetEventFromStorage(event)
	if err != nil {
		return err
	}
	i.DefaultOrgID = e.Payload.OrgID
	i.Projection.Set(event)
	return nil
}

func (i *AuthZInstance) reduceProjectSet(event *v2_es.StorageEvent) error {
	if !i.ShouldReduce(event) {
		return nil
	}

	e, err := instance.ProjectSetEventFromStorage(event)
	if err != nil {
		return err
	}
	i.ProjectID = e.Payload.ProjectID
	i.Projection.Set(event)
	return nil
}

func (i *AuthZInstance) reduceConsoleSet(event *v2_es.StorageEvent) error {
	if !i.ShouldReduce(event) {
		return nil
	}

	e, err := instance.ConsoleSetEventFromStorage(event)
	if err != nil {
		return err
	}
	i.ConsoleAppID = e.Payload.AppID
	i.ConsoleClientID = e.Payload.ClientID
	i.Projection.Set(event)
	return nil
}

func (i *AuthZInstance) reduceDefaultLanguageSet(event *v2_es.StorageEvent) error {
	if !i.ShouldReduce(event) {
		return nil
	}

	e, err := instance.DefaultLanguageSetEventFromStorage(event)
	if err != nil {
		return err
	}
	i.DefaultLanguage = e.Payload.Language
	i.Projection.Set(event)
	return nil
}

func (i *AuthZInstance) reduceRemoved(event *v2_es.StorageEvent) error {
	if !i.ShouldReduce(event) {
		return nil
	}

	return i.Projection.Reduce(event, i.State.reduceRemoved)
}

func (i *AuthZInstance) reduceDomainAdded(event *v2_es.StorageEvent) error {
	if !i.ShouldReduce(event) {
		return nil
	}

	i.Domains = append(i.Domains, NewInstanceDomainFromEvent(event))
	return nil
}

func (i *AuthZInstance) reduceDomainVerified(event *v2_es.StorageEvent) error {
	if !i.ShouldReduce(event) {
		return nil
	}

	domains := slices.Clone(i.Domains)
	for _, domain := range domains {
		err := domain.reduceVerified(event)
		if err != nil {
			return err
		}
	}
	i.Domains = domains
	i.Projection.Set(event)
	return nil
}

func (i *AuthZInstance) reduceDomainPrimarySet(event *v2_es.StorageEvent) error {
	if !i.ShouldReduce(event) {
		return nil
	}

	domains := slices.Clone(i.Domains)
	for _, domain := range domains {
		err := domain.reducePrimarySet(event)
		if err != nil {
			return err
		}
	}
	i.Domains = domains
	i.Projection.Set(event)
	return nil
}

func (i *AuthZInstance) reduceDomainRemoved(event *v2_es.StorageEvent) error {
	if !i.ShouldReduce(event) {
		return nil
	}

	e, err := instance.DomainRemovedEventFromStorage(event)
	if err != nil {
		return err
	}

	i.Domains = slices.DeleteFunc(i.Domains, func(domain *InstanceDomain) bool {
		return domain.Name == e.Payload.Name
	})
	i.Projection.Set(event)
	return nil
}

func (i *AuthZInstance) ShouldReduce(event *v2_es.StorageEvent) bool {
	return event.Aggregate.ID == i.ID && i.Projection.ShouldReduce(event)
}
