package model

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	iam_es_model "github.com/zitadel/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddDomainPolicyEvent(event eventstore.Event) error {
	o.DomainPolicy = new(iam_es_model.DomainPolicy)
	err := o.DomainPolicy.SetData(event)
	if err != nil {
		return err
	}
	o.DomainPolicy.ObjectRoot.CreationDate = event.CreatedAt()
	return nil
}

func (o *Org) appendChangeDomainPolicyEvent(event eventstore.Event) error {
	return o.DomainPolicy.SetData(event)
}

func (o *Org) appendRemoveDomainPolicyEvent() {
	o.DomainPolicy = nil
}
