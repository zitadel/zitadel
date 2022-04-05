package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddDomainPolicyEvent(event *es_models.Event) error {
	o.DomainPolicy = new(iam_es_model.DomainPolicy)
	err := o.DomainPolicy.SetData(event)
	if err != nil {
		return err
	}
	o.DomainPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (o *Org) appendChangeDomainPolicyEvent(event *es_models.Event) error {
	return o.DomainPolicy.SetData(event)
}

func (o *Org) appendRemoveDomainPolicyEvent() {
	o.DomainPolicy = nil
}
