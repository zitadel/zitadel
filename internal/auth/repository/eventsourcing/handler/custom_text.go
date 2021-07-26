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

type CustomText struct {
	handler
	subscription *v1.Subscription
}

func newCustomText(handler handler) *CustomText {
	h := &CustomText{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *CustomText) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

const (
	customTextTable = "auth.custom_texts"
)

func (m *CustomText) ViewModel() string {
	return customTextTable
}

func (m *CustomText) Subscription() *v1.Subscription {
	return m.subscription
}

func (_ *CustomText) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *CustomText) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestCustomTextSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *CustomText) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := m.view.GetLatestCustomTextSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *CustomText) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = m.processCustomText(event)
	}
	return err
}

func (m *CustomText) processCustomText(event *es_models.Event) (err error) {
	customText := new(iam_model.CustomTextView)
	switch event.Type {
	case iam_es_model.CustomTextSet, model.CustomTextSet:
		text := new(iam_model.CustomTextView)
		err = text.SetData(event)
		if err != nil {
			return err
		}
		customText, err = m.view.CustomTextByIDs(event.AggregateID, text.Template, text.Key, text.Language)
		if err != nil && !caos_errs.IsNotFound(err) {
			return err
		}
		if caos_errs.IsNotFound(err) {
			err = nil
			customText = new(iam_model.CustomTextView)
			customText.Language = text.Language
			customText.Template = text.Template
			customText.CreationDate = event.CreationDate
		}
		err = customText.AppendEvent(event)
	case iam_es_model.CustomTextRemoved, model.CustomTextRemoved:
		text := new(iam_model.CustomTextView)
		err = text.SetData(event)
		if err != nil {
			return err
		}
		return m.view.DeleteCustomText(event.AggregateID, text.Template, text.Language, text.Key, event)
	case iam_es_model.CustomTextMessageRemoved, model.CustomTextMessageRemoved:
		text := new(iam_model.CustomTextView)
		err = text.SetData(event)
		if err != nil {
			return err
		}
		return m.view.DeleteCustomTextTemplate(event.AggregateID, text.Template, text.Language, event)
	default:
		return m.view.ProcessedCustomTextSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutCustomText(customText, event)
}

func (m *CustomText) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-3m0fs", "id", event.AggregateID).WithError(err).Warn("something went wrong in custom text handler")
	return spooler.HandleError(event, err, m.view.GetLatestCustomTextFailedEvent, m.view.ProcessedCustomTextFailedEvent, m.view.ProcessedCustomTextSequence, m.errorCountUntilSkip)
}

func (o *CustomText) OnSuccess() error {
	return spooler.HandleSuccess(o.view.UpdateCustomTextSpoolerRunTimestamp)
}
