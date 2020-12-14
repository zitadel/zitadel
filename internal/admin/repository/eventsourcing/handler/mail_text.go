package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
)

type MailText struct {
	handler
}

const (
	mailTextTable = "adminapi.mail_texts"
)

func (m *MailText) ViewModel() string {
	return mailTextTable
}

func (m *MailText) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestMailTextSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.IAMAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *MailText) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate:
		err = m.processMailText(event)
	}
	return err
}

func (m *MailText) processMailText(event *models.Event) (err error) {
	mailText := new(iam_model.MailTextView)
	switch event.Type {
	case model.MailTextAdded:
		err = mailText.AppendEvent(event)
	case model.MailTextChanged:
		err = mailText.SetData(event)
		if err != nil {
			return err
		}
		mailText, err = m.view.MailTextByIDs(event.AggregateID, mailText.MailTextType, mailText.Language)
		if err != nil {
			return err
		}
		err = mailText.AppendEvent(event)
	default:
		return m.view.ProcessedMailTextSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutMailText(mailText, mailText.Sequence)
}

func (m *MailText) OnError(event *models.Event, err error) error {
	logging.LogWithFields("HANDL-5jk84", "id", event.AggregateID).WithError(err).Warn("something went wrong in label mailText handler")
	return spooler.HandleError(event, err, m.view.GetLatestMailTextFailedEvent, m.view.ProcessedMailTextFailedEvent, m.view.ProcessedMailTextSequence, m.errorCountUntilSkip)
}
