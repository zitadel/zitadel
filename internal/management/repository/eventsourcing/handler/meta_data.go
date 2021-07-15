package handler

import (
	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	usr_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

type MetaData struct {
	handler
	subscription *v1.Subscription
}

func newMetaData(handler handler) *MetaData {
	h := &MetaData{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *MetaData) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

const (
	metaDataTable = "management.meta_data"
)

func (m *MetaData) ViewModel() string {
	return metaDataTable
}

func (m *MetaData) Subscription() *v1.Subscription {
	return m.subscription
}

func (_ *MetaData) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{usr_model.UserAggregate}
}

func (p *MetaData) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestMetaDataSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *MetaData) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := m.view.GetLatestMetaDataSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *MetaData) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case usr_model.UserAggregate:
		err = m.processMetaData(event)
	}
	return err
}

func (m *MetaData) processMetaData(event *es_models.Event) (err error) {
	customText := new(iam_model.MetaDataView)
	switch event.Type {
	case usr_model.UserMetaDataSet:
		data := new(iam_model.MetaDataView)
		err = data.SetData(event)
		if err != nil {
			return err
		}
		data, err = m.view.MetaDataByKey(event.AggregateID, data.Key)
		if err != nil && !caos_errs.IsNotFound(err) {
			return err
		}
		if caos_errs.IsNotFound(err) {
			err = nil
			data = new(iam_model.MetaDataView)
			data.CreationDate = event.CreationDate
		}
		err = data.AppendEvent(event)
	case usr_model.UserMetaDataRemoved:
		data := new(iam_model.MetaDataView)
		err = data.SetData(event)
		if err != nil {
			return err
		}
		return m.view.DeleteMetaData(event.AggregateID, data.Key, event)
	case usr_model.UserRemoved:
		return m.view.DeleteMetaDataByAggregateID(event.AggregateID, event)
	default:
		return m.view.ProcessedMetaDataSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutMetaData(customText, event)
}

func (m *MetaData) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-3m912", "id", event.AggregateID).WithError(err).Warn("something went wrong in custom text handler")
	return spooler.HandleError(event, err, m.view.GetLatestMetaDataFailedEvent, m.view.ProcessedMetaDataFailedEvent, m.view.ProcessedMetaDataSequence, m.errorCountUntilSkip)
}

func (o *MetaData) OnSuccess() error {
	return spooler.HandleSuccess(o.view.UpdateMetaDataSpoolerRunTimestamp)
}
