package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddPasswordComplexityPolicyEvent(event *es_models.Event) error {
	o.PasswordComplexityPolicy = new(iam_es_model.PasswordComplexityPolicy)
	err := o.PasswordComplexityPolicy.SetData(event)
	if err != nil {
		return err
	}
	o.PasswordComplexityPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (o *Org) appendChangePasswordComplexityPolicyEvent(event *es_models.Event) error {
	return o.PasswordComplexityPolicy.SetData(event)
}

func (o *Org) appendRemovePasswordComplexityPolicyEvent(event *es_models.Event) {
	o.PasswordComplexityPolicy = nil
}
