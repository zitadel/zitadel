package view

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	messageTextTable = "adminapi.message_texts"
)

func (v *View) MessageTexts(aggregateID string) ([]*model.MessageTextView, error) {
	return view.GetMessageTexts(v.Db, messageTextTable, aggregateID)
}

func (v *View) MessageTextByIDs(aggregateID, textType, lang string) (*model.MessageTextView, error) {
	return view.GetMessageTextByIDs(v.Db, messageTextTable, aggregateID, textType, lang)
}

func (v *View) PutMessageText(template *model.MessageTextView, event *models.Event) error {
	err := view.PutMessageText(v.Db, messageTextTable, template)
	if err != nil {
		return err
	}
	return v.ProcessedMessageTextSequence(event)
}

func (v *View) GetLatestMessageTextSequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(messageTextTable)
}

func (v *View) ProcessedMessageTextSequence(event *models.Event) error {
	return v.saveCurrentSequence(messageTextTable, event)
}

func (v *View) UpdateMessageTextSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(messageTextTable)
}

func (v *View) GetLatestMessageTextFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(messageTextTable, sequence)
}

func (v *View) ProcessedMessageTextFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
