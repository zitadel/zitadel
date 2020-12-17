package handler

import (
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/key/repository/eventsourcing"
	es_model "github.com/caos/zitadel/internal/key/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/key/repository/view/model"
)

type Key struct {
	handler
}

const (
	keyTable = "auth.keys"
)

func (k *Key) ViewModel() string {
	return keyTable
}

func (_ *Key) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{es_model.KeyPairAggregate}
}

func (k *Key) CurrentSequence(event *models.Event) (uint64, error) {
	sequence, err := k.view.GetLatestKeySequence(string(event.AggregateType))
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (k *Key) EventQuery() (*models.SearchQuery, error) {
	sequence, err := k.view.GetLatestKeySequence("")
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
		return k.view.PutKeys(privateKey, publicKey, event)
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
