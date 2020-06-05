package handler

import (
	"time"

	es_model "github.com/caos/zitadel/internal/key/repository/eventsourcing/model"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/key/repository/eventsourcing"
	view_model "github.com/caos/zitadel/internal/key/repository/view/model"
)

type Key struct {
	handler
}

const (
	keyTable = "auth.keys"
)

func (k *Key) MinimumCycleDuration() time.Duration { return k.cycleDuration }

func (k *Key) ViewModel() string {
	return keyTable
}

func (k *Key) EventQuery() (*models.SearchQuery, error) {
	sequence, err := k.view.GetLatestKeySequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.KeyPairQuery(sequence), nil
}

func (k *Key) Process(event *models.Event) error {
	switch event.Type {
	case es_model.KeyPairAdded:
		privateKey, publicKey, err := view_model.KeysFromPairEvent(event)
		if err != nil {
			return err
		}
		return k.view.PutKeys(privateKey, publicKey, event.Sequence)
	default:
		return k.view.ProcessedKeySequence(event.Sequence)
	}
	return nil
}

func (k *Key) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-GHa3a", "id", event.AggregateID).WithError(err).Warn("something went wrong in key handler")
	return spooler.HandleError(event, err, k.view.GetLatestKeyFailedEvent, k.view.ProcessedKeyFailedEvent, k.view.ProcessedKeySequence, k.errorCountUntilSkip)
}
