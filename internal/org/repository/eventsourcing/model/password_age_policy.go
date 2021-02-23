package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddPasswordAgePolicyEvent(event *es_models.Event) error {
	o.PasswordAgePolicy = new(iam_es_model.PasswordAgePolicy)
	err := o.PasswordAgePolicy.SetData(event)
	if err != nil {
		return err
	}
	o.PasswordAgePolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (o *Org) appendChangePasswordAgePolicyEvent(event *es_models.Event) error {
	return o.PasswordAgePolicy.SetData(event)
}

func (o *Org) appendRemovePasswordAgePolicyEvent(event *es_models.Event) {
	o.PasswordAgePolicy = nil
}
