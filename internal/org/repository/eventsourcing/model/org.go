package model

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/zitadel/zitadel/internal/iam/repository/eventsourcing/model"
	org_model "github.com/zitadel/zitadel/internal/org/model"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Org struct {
	es_models.ObjectRoot `json:"-"`

	Name  string `json:"name,omitempty"`
	State int32  `json:"-"`

	Domains      []*OrgDomain               `json:"-"`
	DomainPolicy *iam_es_model.DomainPolicy `json:"-"`
}

func OrgToModel(org *Org) *org_model.Org {
	converted := &org_model.Org{
		ObjectRoot: org.ObjectRoot,
		Name:       org.Name,
		State:      org_model.OrgState(org.State),
		Domains:    OrgDomainsToModel(org.Domains),
	}
	if org.DomainPolicy != nil {
		converted.DomainPolicy = iam_es_model.DomainPolicyToModel(org.DomainPolicy)
	}
	return converted
}

func OrgFromEvents(org *Org, events ...eventstore.Event) (*Org, error) {
	if org == nil {
		org = new(Org)
	}

	return org, org.AppendEvents(events...)
}

func (o *Org) AppendEvents(events ...eventstore.Event) error {
	for _, event := range events {
		err := o.AppendEvent(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Org) AppendEvent(event eventstore.Event) (err error) {
	switch event.Type() {
	case org.OrgAddedEventType:
		err = o.SetData(event)
		if err != nil {
			return err
		}
	case org.OrgChangedEventType:
		err = o.SetData(event)
		if err != nil {
			return err
		}
	case org.OrgDeactivatedEventType:
		o.State = int32(org_model.OrgStateInactive)
	case org.OrgReactivatedEventType:
		o.State = int32(org_model.OrgStateActive)
	case org.OrgDomainAddedEventType:
		err = o.appendAddDomainEvent(event)
	case org.OrgDomainVerificationAddedEventType:
		err = o.appendVerificationDomainEvent(event)
	case org.OrgDomainVerifiedEventType:
		err = o.appendVerifyDomainEvent(event)
	case org.OrgDomainPrimarySetEventType:
		err = o.appendPrimaryDomainEvent(event)
	case org.OrgDomainRemovedEventType:
		err = o.appendRemoveDomainEvent(event)
	case org.DomainPolicyAddedEventType:
		err = o.appendAddDomainPolicyEvent(event)
	case org.DomainPolicyChangedEventType:
		err = o.appendChangeDomainPolicyEvent(event)
	case org.DomainPolicyRemovedEventType:
		o.appendRemoveDomainPolicyEvent()
	}
	if err != nil {
		return err
	}
	o.ObjectRoot.AppendEvent(event)
	return nil
}

func (o *Org) SetData(event eventstore.Event) error {
	err := event.Unmarshal(o)
	if err != nil {
		return zerrors.ThrowInternal(err, "EVENT-BpbQZ", "unable to unmarshal event")
	}
	return nil
}
