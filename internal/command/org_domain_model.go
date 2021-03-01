package command

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
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

func (wm *OrgDomainWriteModel) AppendEvents(events ...eventstore.EventReader) {
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
		}
	}
	return nil
}

func (wm *OrgDomainWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			org.OrgDomainAddedEventType,
			org.OrgDomainVerifiedEventType,
			org.OrgDomainVerificationAddedEventType,
			org.OrgDomainVerifiedEventType,
			org.OrgDomainPrimarySetEventType,
			org.OrgDomainRemovedEventType)
}
