package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
)

const (
	idpConfigTable = "adminapi.idp_configs"
)

type IDPConfig struct {
	handler
	subscription *v1.Subscription
}

func newIDPConfig(handler handler) *IDPConfig {
	h := &IDPConfig{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (i *IDPConfig) subscribe() {
	i.subscription = i.es.Subscribe(i.AggregateTypes()...)
	go func() {
		for event := range i.subscription.Events {
			query.ReduceEvent(i, event)
		}
	}()
}

func (i *IDPConfig) ViewModel() string {
	return idpConfigTable
}

func (i *IDPConfig) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.IAMAggregate}
}

func (i *IDPConfig) CurrentSequence() (uint64, error) {
	sequence, err := i.view.GetLatestIDPConfigSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (i *IDPConfig) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := i.view.GetLatestIDPConfigSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(i.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (i *IDPConfig) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate:
		err = i.processIDPConfig(event)
	}
	return err
}

func (i *IDPConfig) processIDPConfig(event *es_models.Event) (err error) {
	idp := new(iam_view_model.IDPConfigView)
	switch event.Type {
	case model.IDPConfigAdded:
		err = idp.AppendEvent(iam_model.IDPProviderTypeSystem, event)
	case model.IDPConfigChanged,
		model.OIDCIDPConfigAdded,
		model.OIDCIDPConfigChanged:
		err = idp.SetData(event)
		if err != nil {
			return err
		}
		idp, err = i.view.IDPConfigByID(idp.IDPConfigID)
		if err != nil {
			return err
		}
		err = idp.AppendEvent(iam_model.IDPProviderTypeSystem, event)
	case model.IDPConfigRemoved:
		err = idp.SetData(event)
		if err != nil {
			return err
		}
		return i.view.DeleteIDPConfig(idp.IDPConfigID, event)
	default:
		return i.view.ProcessedIDPConfigSequence(event)
	}
	if err != nil {
		return err
	}
	return i.view.PutIDPConfig(idp, event)
}

func (i *IDPConfig) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-Mslo9", "id", event.AggregateID).WithError(err).Warn("something went wrong in idp config handler")
	return spooler.HandleError(event, err, i.view.GetLatestIDPConfigFailedEvent, i.view.ProcessedIDPConfigFailedEvent, i.view.ProcessedIDPConfigSequence, i.errorCountUntilSkip)
}

func (i *IDPConfig) OnSuccess() error {
	return spooler.HandleSuccess(i.view.UpdateIDPConfigSpoolerRunTimestamp)
}
