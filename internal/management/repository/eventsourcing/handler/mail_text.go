package handler

import (
	"github.com/caos/logging"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

type MailText struct {
	handler
}

const (
	mailTextTable = "management.mail_texts"
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
		AggregateTypeFilter(model.OrgAggregate, iam_es_model.IAMAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *MailText) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = m.processMailText(event)
	}
	return err
}

func (m *MailText) processMailText(event *models.Event) (err error) {
	text := new(iam_model.MailTextView)
	switch event.Type {
	case iam_es_model.MailTextAdded, model.MailTextAdded:
		err = text.AppendEvent(event)
	case iam_es_model.MailTextChanged, model.MailTextChanged:
		err = text.SetData(event)
		if err != nil {
			return err
		}
		text, err = m.view.MailTextByIDs(event.AggregateID, text.MailTextType, text.Language)
		if err != nil {
			return err
		}
		text.ChangeDate = event.CreationDate
		err = text.AppendEvent(event)
	case model.MailTextRemoved:
		err = text.SetData(event)
		return m.view.DeleteMailText(event.AggregateID, text.MailTextType, text.Language, event.Sequence)
	default:
		return m.view.ProcessedMailTextSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutMailText(text, text.Sequence)
}

func (m *MailText) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-4Djo9", "id", event.AggregateID).WithError(err).Warn("something went wrong in label text handler")
	return spooler.HandleError(event, err, m.view.GetLatestMailTextFailedEvent, m.view.ProcessedMailTextFailedEvent, m.view.ProcessedMailTextSequence, m.errorCountUntilSkip)
}
