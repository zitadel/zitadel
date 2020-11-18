package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddMailTemplateEvent(event *es_models.Event) error {
	o.MailTemplate = new(iam_es_model.MailTemplate)
	err := o.MailTemplate.SetDataLabel(event)
	if err != nil {
		return err
	}
	o.MailTemplate.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (o *Org) appendChangeMailTemplateEvent(event *es_models.Event) error {
	return o.MailTemplate.SetDataLabel(event)
}

func (o *Org) appendRemoveMailTemplateEvent(event *es_models.Event) {
	o.MailTemplate = nil
}
