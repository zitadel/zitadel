package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddMailTextEvent(event *es_models.Event) error {
	mailText := &iam_es_model.MailText{}
	err := mailText.SetDataLabel(event)
	if err != nil {
		return err
	}
	mailText.ObjectRoot.CreationDate = event.CreationDate
	o.MailTexts = append(o.MailTexts, mailText)
	return nil
}

func (o *Org) appendChangeMailTextEvent(event *es_models.Event) error {
	mailText := &iam_es_model.MailText{}
	err := mailText.SetDataLabel(event)
	if err != nil {
		return err
	}
	mailText.ObjectRoot.ChangeDate = event.CreationDate
	if n, m := iam_es_model.GetMailText(o.MailTexts, mailText.MailTextType, mailText.Language); m != nil {
		o.MailTexts[n] = mailText
	}
	return nil
}

func (o *Org) appendRemoveMailTextEvent(event *es_models.Event) error {
	mailText := &iam_es_model.MailText{}
	err := mailText.SetDataLabel(event)
	if err != nil {
		return err
	}
	if n, m := iam_es_model.GetMailText(o.MailTexts, mailText.MailTextType, mailText.Language); m != nil {
		o.MailTexts[n] = o.MailTexts[len(o.MailTexts)-1]
		o.MailTexts[len(o.MailTexts)-1] = nil
		o.MailTexts = o.MailTexts[:len(o.MailTexts)-1]
	}
	return nil
}
