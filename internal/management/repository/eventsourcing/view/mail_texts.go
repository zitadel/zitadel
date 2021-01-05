package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	mailTextTable = "management.mail_texts"
)

func (v *View) MailTextsByAggregateID(aggregateID string) ([]*model.MailTextView, error) {
	return view.GetMailTexts(v.Db, mailTextTable, aggregateID)
}

func (v *View) MailTextByIDs(aggregateID string, textType string, language string) (*model.MailTextView, error) {
	return view.GetMailTextByIDs(v.Db, mailTextTable, aggregateID, textType, language)
}

func (v *View) PutMailText(template *model.MailTextView, event *models.Event) error {
	err := view.PutMailText(v.Db, mailTextTable, template)
	if err != nil {
		return err
	}
	return v.ProcessedMailTextSequence(event)
}

func (v *View) DeleteMailText(aggregateID string, textType string, language string, event *models.Event) error {
	err := view.DeleteMailText(v.Db, mailTextTable, aggregateID, textType, language)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedMailTextSequence(event)
}

func (v *View) GetLatestMailTextSequence(aggregateType string) (*global_view.CurrentSequence, error) {
	return v.latestSequence(mailTextTable, aggregateType)
}

func (v *View) ProcessedMailTextSequence(event *models.Event) error {
	return v.saveCurrentSequence(mailTextTable, event)
}

func (v *View) UpdateMailTextSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(mailTextTable)
}

func (v *View) GetLatestMailTextFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(mailTextTable, sequence)
}

func (v *View) ProcessedMailTextFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
