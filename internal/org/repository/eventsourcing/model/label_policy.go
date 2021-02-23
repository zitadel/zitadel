package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddLabelPolicyEvent(event *es_models.Event) error {
	o.LabelPolicy = new(iam_es_model.LabelPolicy)
	err := o.LabelPolicy.SetDataLabel(event)
	if err != nil {
		return err
	}
	o.LabelPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (o *Org) appendChangeLabelPolicyEvent(event *es_models.Event) error {
	return o.LabelPolicy.SetDataLabel(event)
}

func (o *Org) appendRemoveLabelPolicyEvent(event *es_models.Event) {
	o.LabelPolicy = nil
}
