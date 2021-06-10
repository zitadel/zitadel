package handler

import (
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

type MessageText struct {
	handler
	subscription *v1.Subscription
}

func newMessageText(handler handler) *MessageText {
	h := &MessageText{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *MessageText) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

const (
	messageTextTable = "management.message_texts"
)

func (m *MessageText) ViewModel() string {
	return messageTextTable
}

func (_ *MessageText) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *MessageText) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestMessageTextSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *MessageText) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := m.view.GetLatestMessageTextSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *MessageText) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = m.processMessageText(event)
	}
	return err
}

func (m *MessageText) processMessageText(event *es_models.Event) (err error) {
	message := new(iam_model.MessageTextView)
	switch event.Type {
	case iam_es_model.CustomTextSet, model.CustomTextSet,
		iam_es_model.CustomTextRemoved, model.CustomTextRemoved:
		text := new(iam_model.CustomText)
		err = text.SetData(event)
		if err != nil {
			return err
		}
		message, err = m.view.MessageTextByIDs(event.AggregateID, text.Template, text.Language.String())
		if err != nil && !caos_errs.IsNotFound(err) {
			return err
		}
		if caos_errs.IsNotFound(err) {
			err = nil
			message = new(iam_model.MessageTextView)
			message.Language = text.Language.String()
			message.MessageTextType = text.Template
			message.CreationDate = event.CreationDate
		}
		err = message.AppendEvent(event)
	case model.CustomTextMessageRemoved:
		text := new(iam_model.CustomText)
		err = text.SetData(event)
		if err != nil {
			return err
		}
		return m.view.DeleteMessageText(event.AggregateID, text.Template, text.Language.String(), event)
	default:
		return m.view.ProcessedMessageTextSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutMessageText(message, event)
}

func (m *MessageText) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-4Djo9", "id", event.AggregateID).WithError(err).Warn("something went wrong in label text handler")
	return spooler.HandleError(event, err, m.view.GetLatestMessageTextFailedEvent, m.view.ProcessedMessageTextFailedEvent, m.view.ProcessedMessageTextSequence, m.errorCountUntilSkip)
}

func (o *MessageText) OnSuccess() error {
	return spooler.HandleSuccess(o.view.UpdateMessageTextSpoolerRunTimestamp)
}
