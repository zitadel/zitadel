package handler

import (
	"github.com/caos/zitadel/internal/eventstore/v1"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/key/repository/eventsourcing"
	es_model "github.com/caos/zitadel/internal/key/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/key/repository/view/model"
)

const (
	keyTable = "auth.keys"
)

type Key struct {
	handler
	subscription *v1.Subscription
	keyChan      chan<- *model.KeyView
}

func newKey(handler handler, keyChan chan<- *model.KeyView) *Key {
	h := &Key{
		handler: handler,
		keyChan: keyChan,
	}

	h.subscribe()

	return h
}

func (k *Key) subscribe() {
	k.subscription = k.es.Subscribe(k.AggregateTypes()...)
	go func() {
		for event := range k.subscription.Events {
			query.ReduceEvent(k, event)
		}
	}()
}

func (k *Key) ViewModel() string {
	return keyTable
}

func (k *Key) Subscription() *v1.Subscription {
	return k.subscription
}

func (_ *Key) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{es_model.KeyPairAggregate}
}

func (k *Key) CurrentSequence() (uint64, error) {
	sequence, err := k.view.GetLatestKeySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (k *Key) EventQuery() (*models.SearchQuery, error) {
	sequence, err := k.view.GetLatestKeySequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.KeyPairQuery(sequence.CurrentSequence), nil
}

func (k *Key) Reduce(event *models.Event) error {
	switch event.Type {
	case es_model.KeyPairAdded:
		privateKey, publicKey, err := view_model.KeysFromPairEvent(event)
		if err != nil {
			return err
		}
		if privateKey.Expiry.Before(time.Now()) && publicKey.Expiry.Before(time.Now()) {
			return k.view.ProcessedKeySequence(event)
		}
		err = k.view.PutKeys(privateKey, publicKey, event)
		if err != nil {
			return err
		}
		k.keyChan <- view_model.KeyViewToModel(privateKey)
		return nil
	default:
		return k.view.ProcessedKeySequence(event)
	}
}

func (k *Key) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-GHa3a", "id", event.AggregateID).WithError(err).Warn("something went wrong in key handler")
	return spooler.HandleError(event, err, k.view.GetLatestKeyFailedEvent, k.view.ProcessedKeyFailedEvent, k.view.ProcessedKeySequence, k.errorCountUntilSkip)
}

func (k *Key) OnSuccess() error {
	err := spooler.HandleSuccess(k.view.UpdateKeySpoolerRunTimestamp)
	logging.LogWithFields("SPOOL-vM9sd", "table", keyTable).OnError(err).Warn("could not process on success func")
	return err
}
