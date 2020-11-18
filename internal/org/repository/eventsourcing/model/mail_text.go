package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddMailTextEvent(event *es_models.Event) error {
	o.MailText = new(iam_es_model.MailText)
	err := o.MailText.SetDataLabel(event)
	if err != nil {
		return err
	}
	o.MailText.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (o *Org) appendChangeMailTextEvent(event *es_models.Event) error {
	return o.MailText.SetDataLabel(event)
}

func (o *Org) appendRemoveMailTextEvent(event *es_models.Event) {
	o.MailText = nil
}
