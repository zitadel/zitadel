package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	mailTemplateTable = "management.mail_templates"
)

func (v *View) MailTemplateByAggregateID(aggregateID string) (*model.MailTemplateView, error) {
	return view.GetMailTemplateByAggregateID(v.Db, mailTemplateTable, aggregateID)
}

func (v *View) PutMailTemplate(template *model.MailTemplateView, sequence uint64) error {
	err := view.PutMailTemplate(v.Db, mailTemplateTable, template)
	if err != nil {
		return err
	}
	return v.ProcessedMailTemplateSequence(sequence)
}

func (v *View) DeleteMailTemplate(aggregateID string, eventSequence uint64) error {
	err := view.DeleteMailTemplate(v.Db, mailTemplateTable, aggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedMailTemplateSequence(eventSequence)
}

func (v *View) GetLatestMailTemplateSequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(mailTemplateTable)
}

func (v *View) ProcessedMailTemplateSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(mailTemplateTable, eventSequence)
}

func (v *View) GetLatestMailTemplateFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(mailTemplateTable, sequence)
}

func (v *View) ProcessedMailTemplateFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
