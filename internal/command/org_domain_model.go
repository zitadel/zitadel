package command

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type OrgDomainWriteModel struct {
	eventstore.WriteModel

	Domain         string
	ValidationType domain.OrgDomainValidationType
	ValidationCode *crypto.CryptoValue
	Primary        bool
	Verified       bool

	State domain.OrgDomainState
}

func NewOrgDomainWriteModel(orgID string, domain string) *OrgDomainWriteModel {
	return &OrgDomainWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   orgID,
			ResourceOwner: orgID,
		},
		Domain: domain,
	}
}

func (wm *OrgDomainWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.DomainAddedEvent:
			if e.Domain != wm.Domain {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *org.DomainVerificationAddedEvent:
			if e.Domain != wm.Domain {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *org.DomainVerificationFailedEvent:
			if e.Domain != wm.Domain {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *org.DomainVerifiedEvent:
			if e.Domain != wm.Domain {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *org.DomainPrimarySetEvent:
			wm.WriteModel.AppendEvents(e)
		case *org.DomainRemovedEvent:
			if e.Domain != wm.Domain {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgDomainWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *org.DomainAddedEvent:
			wm.Domain = e.Domain
			wm.State = domain.OrgDomainStateActive
		case *org.DomainVerificationAddedEvent:
			wm.ValidationType = e.ValidationType
			wm.ValidationCode = e.ValidationCode
		case *org.DomainVerifiedEvent:
			wm.Verified = true
		case *org.DomainPrimarySetEvent:
			wm.Primary = e.Domain == wm.Domain
		case *org.DomainRemovedEvent:
			wm.State = domain.OrgDomainStateRemoved
			wm.Verified = false
			wm.Primary = false
			wm.ValidationType = domain.OrgDomainValidationTypeUnspecified
			wm.ValidationCode = nil
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OrgDomainWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.OrgDomainAddedEventType,
			org.OrgDomainVerifiedEventType,
			org.OrgDomainVerificationAddedEventType,
			org.OrgDomainVerifiedEventType,
			org.OrgDomainPrimarySetEventType,
			org.OrgDomainRemovedEventType).
		Builder()
}

type OrgDomainsWriteModel struct {
	eventstore.WriteModel

	Domains       []*Domain
	PrimaryDomain string
	OrgName       string
}

type Domain struct {
	Domain   string
	Verified bool
	State    domain.OrgDomainState
}

func NewOrgDomainsWriteModel(orgID string) *OrgDomainsWriteModel {
	return &OrgDomainsWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   orgID,
			ResourceOwner: orgID,
		},
		Domains: make([]*Domain, 0),
	}
}

func (wm *OrgDomainsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *org.OrgAddedEvent:
			wm.OrgName = e.Name
		case *org.OrgChangedEvent:
			wm.OrgName = e.Name
		case *org.DomainAddedEvent:
			wm.Domains = append(wm.Domains, &Domain{Domain: e.Domain, State: domain.OrgDomainStateActive})
		case *org.DomainVerifiedEvent:
			for _, d := range wm.Domains {
				if d.Domain == e.Domain {
					d.Verified = true
					continue
				}
			}
		case *org.DomainPrimarySetEvent:
			wm.PrimaryDomain = e.Domain
		case *org.DomainRemovedEvent:
			for _, d := range wm.Domains {
				if d.Domain == e.Domain {
					d.State = domain.OrgDomainStateRemoved
					continue
				}
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OrgDomainsWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.OrgAddedEventType,
			org.OrgChangedEventType,
			org.OrgDomainAddedEventType,
			org.OrgDomainVerifiedEventType,
			org.OrgDomainVerificationAddedEventType,
			org.OrgDomainVerifiedEventType,
			org.OrgDomainPrimarySetEventType,
			org.OrgDomainRemovedEventType).
		Builder()
}

type OrgDomainVerifiedWriteModel struct {
	eventstore.WriteModel

	Domain   string
	Verified bool
}

func NewOrgDomainVerifiedWriteModel(domain string) *OrgDomainVerifiedWriteModel {
	return &OrgDomainVerifiedWriteModel{
		WriteModel: eventstore.WriteModel{},
		Domain:     domain,
	}
}

func (wm *OrgDomainVerifiedWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.DomainVerifiedEvent:
			if e.Domain != wm.Domain {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *org.DomainRemovedEvent:
			if e.Domain != wm.Domain {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *org.OrgRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgDomainVerifiedWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *org.DomainVerifiedEvent:
			wm.Verified = true
			wm.ResourceOwner = e.Aggregate().ResourceOwner
		case *org.DomainRemovedEvent:
			wm.Verified = false
		case *org.OrgRemovedEvent:
			if wm.ResourceOwner != e.Aggregate().ID {
				continue
			}
			wm.Verified = false
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OrgDomainVerifiedWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(org.AggregateType).
		EventTypes(
			org.OrgDomainVerifiedEventType,
			org.OrgDomainRemovedEventType,
			org.OrgRemovedEventType).
		Builder()
}
