package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddOrgIAMPolicyEvent(event *es_models.Event) error {
	o.OrgIamPolicy = new(iam_es_model.OrgIAMPolicy)
	err := o.OrgIamPolicy.SetData(event)
	if err != nil {
		return err
	}
	o.OrgIamPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (o *Org) appendChangeOrgIAMPolicyEvent(event *es_models.Event) error {
	return o.OrgIamPolicy.SetData(event)
}

func (o *Org) appendRemoveOrgIAMPolicyEvent() {
	o.OrgIamPolicy = nil
}
