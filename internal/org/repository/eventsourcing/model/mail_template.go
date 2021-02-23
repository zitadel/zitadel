package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
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
	mailTemplate := &iam_es_model.MailTemplate{}
	err := mailTemplate.SetDataLabel(event)
	if err != nil {
		return err
	}
	mailTemplate.ObjectRoot.ChangeDate = event.CreationDate
	o.MailTemplate = mailTemplate
	return nil
}

func (o *Org) appendRemoveMailTemplateEvent(event *es_models.Event) {
	o.MailTemplate = nil
}
