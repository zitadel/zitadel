package handler

import (
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	key_model "github.com/caos/zitadel/internal/key/repository/view/model"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	authnKeysTable = "auth.authn_keys"
)

type AuthNKeys struct {
	handler
	subscription *eventstore.Subscription
}

func newAuthNKeys(handler handler) *AuthNKeys {
	h := &AuthNKeys{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (k *AuthNKeys) subscribe() {
	k.subscription = k.es.Subscribe(k.AggregateTypes()...)
	go func() {
		for event := range k.subscription.Events {
			query.ReduceEvent(k, event)
		}
	}()
}

func (k *AuthNKeys) ViewModel() string {
	return authnKeysTable
}

func (_ *AuthNKeys) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.UserAggregate}
}

func (k *AuthNKeys) CurrentSequence(event *models.Event) (uint64, error) {
	sequence, err := k.view.GetLatestAuthNKeySequence(string(event.AggregateType))
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (k *AuthNKeys) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := k.view.GetLatestAuthNKeySequence("")
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(k.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (k *AuthNKeys) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.UserAggregate:
		err = k.processAuthNKeys(event)
	}
	return err
}

func (k *AuthNKeys) processAuthNKeys(event *es_models.Event) (err error) {
	key := new(key_model.AuthNKeyView)
	switch event.Type {
	case model.MachineKeyAdded:
		err = key.AppendEvent(event)
		if key.ExpirationDate.Before(time.Now()) {
			return k.view.ProcessedAuthNKeySequence(event)
		}
	case model.MachineKeyRemoved:
		err = key.SetData(event)
		if err != nil {
			return err
		}
		return k.view.DeleteAuthNKey(key.ID, event)
	case model.UserRemoved:
		return k.view.DeleteAuthNKeysByObjectID(event.AggregateID, event)
	default:
		return k.view.ProcessedAuthNKeySequence(event)
	}
	if err != nil {
		return err
	}
	return k.view.PutAuthNKey(key, event)
}

func (k *AuthNKeys) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-S9fe", "id", event.AggregateID).WithError(err).Warn("something went wrong in authn key handler")
	return spooler.HandleError(event, err, k.view.GetLatestAuthNKeyFailedEvent, k.view.ProcessedAuthNKeyFailedEvent, k.view.ProcessedAuthNKeySequence, k.errorCountUntilSkip)
}

func (k *AuthNKeys) OnSuccess() error {
	return spooler.HandleSuccess(k.view.UpdateAuthNKeySpoolerRunTimestamp)
}
