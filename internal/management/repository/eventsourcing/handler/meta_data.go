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

type Metadata struct {
	handler
	subscription *v1.Subscription
}

func newMetadata(handler handler) *Metadata {
	h := &Metadata{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *Metadata) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

const (
	metadataTable = "management.meta_data"
)

func (m *Metadata) ViewModel() string {
	return metadataTable
}

func (m *Metadata) Subscription() *v1.Subscription {
	return m.subscription
}

func (_ *Metadata) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{usr_model.UserAggregate}
}

func (p *Metadata) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestMetadataSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *Metadata) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := m.view.GetLatestMetadataSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *Metadata) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case usr_model.UserAggregate:
		err = m.processMetadata(event)
	}
	return err
}

func (m *Metadata) processMetadata(event *es_models.Event) (err error) {
	metaData := new(iam_model.MetadataView)
	switch event.Type {
	case usr_model.UserMetadataSet:
		err = metaData.SetData(event)
		if err != nil {
			return err
		}
		metaData, err = m.view.MetadataByKey(event.AggregateID, metaData.Key)
		if err != nil && !caos_errs.IsNotFound(err) {
			return err
		}
		if caos_errs.IsNotFound(err) {
			err = nil
			metaData = new(iam_model.MetadataView)
			metaData.CreationDate = event.CreationDate
		}
		err = metaData.AppendEvent(event)
	case usr_model.UserMetadataRemoved:
		data := new(iam_model.MetadataView)
		err = data.SetData(event)
		if err != nil {
			return err
		}
		return m.view.DeleteMetadata(event.AggregateID, data.Key, event)
	case usr_model.UserRemoved:
		return m.view.DeleteMetadataByAggregateID(event.AggregateID, event)
	default:
		return m.view.ProcessedMetadataSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutMetadata(metaData, event)
}

func (m *Metadata) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-3m912", "id", event.AggregateID).WithError(err).Warn("something went wrong in custom text handler")
	return spooler.HandleError(event, err, m.view.GetLatestMetadataFailedEvent, m.view.ProcessedMetadataFailedEvent, m.view.ProcessedMetadataSequence, m.errorCountUntilSkip)
}

func (o *Metadata) OnSuccess() error {
	return spooler.HandleSuccess(o.view.UpdateMetadataSpoolerRunTimestamp)
}
