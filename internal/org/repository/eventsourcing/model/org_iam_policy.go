package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddOrgIAMPolicyEvent(event *es_models.Event) error {
	o.OrgIAMPolicy = new(iam_es_model.OrgIAMPolicy)
	err := o.OrgIAMPolicy.SetData(event)
	if err != nil {
		return err
	}
	o.OrgIAMPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (o *Org) appendChangeOrgIAMPolicyEvent(event *es_models.Event) error {
	return o.OrgIAMPolicy.SetData(event)
}

func (o *Org) appendRemoveOrgIAMPolicyEvent() {
	o.OrgIAMPolicy = nil
}
