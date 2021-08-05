package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddLockoutPolicyEvent(event *es_models.Event) error {
	o.LockoutPolicy = new(iam_es_model.LockoutPolicy)
	err := o.LockoutPolicy.SetData(event)
	if err != nil {
		return err
	}
	o.LockoutPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (o *Org) appendChangeLockoutPolicyEvent(event *es_models.Event) error {
	return o.LockoutPolicy.SetData(event)
}

func (o *Org) appendRemoveLockoutPolicyEvent(event *es_models.Event) {
	o.LockoutPolicy = nil
}
