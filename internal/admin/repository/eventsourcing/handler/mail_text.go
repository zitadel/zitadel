package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
)

type MailText struct {
	handler
	subscription *eventstore.Subscription
}

func newMailText(handler handler) *MailText {
	h := &MailText{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *MailText) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

const (
	mailTextTable = "adminapi.mail_texts"
)

func (m *MailText) ViewModel() string {
	return mailTextTable
}

func (_ *MailText) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{iam_es_model.IAMAggregate}
}

func (p *MailText) CurrentSequence(event *models.Event) (uint64, error) {
	sequence, err := p.view.GetLatestMailTextSequence(string(event.AggregateType))
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *MailText) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestMailTextSequence("")
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
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
		return m.view.ProcessedMailTextSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutMailText(mailText, event)
}

func (m *MailText) OnError(event *models.Event, err error) error {
	logging.LogWithFields("HANDL-5jk84", "id", event.AggregateID).WithError(err).Warn("something went wrong in label mailText handler")
	return spooler.HandleError(event, err, m.view.GetLatestMailTextFailedEvent, m.view.ProcessedMailTextFailedEvent, m.view.ProcessedMailTextSequence, m.errorCountUntilSkip)
}

func (o *MailText) OnSuccess() error {
	return spooler.HandleSuccess(o.view.UpdateMailTextSpoolerRunTimestamp)
}
