package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddPasswordLockoutPolicyEvent(event *es_models.Event) error {
	o.PasswordLockoutPolicy = new(iam_es_model.LockoutPolicy)
	err := o.PasswordLockoutPolicy.SetData(event)
	if err != nil {
		return err
	}
	o.PasswordLockoutPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (o *Org) appendChangePasswordLockoutPolicyEvent(event *es_models.Event) error {
	return o.PasswordLockoutPolicy.SetData(event)
}

func (o *Org) appendRemovePasswordLockoutPolicyEvent(event *es_models.Event) {
	o.PasswordLockoutPolicy = nil
}
