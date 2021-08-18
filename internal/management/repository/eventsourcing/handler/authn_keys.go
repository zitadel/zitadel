package handler

import (
	"github.com/caos/zitadel/internal/eventstore/v1"
	"time"

	"github.com/caos/logging"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	key_model "github.com/caos/zitadel/internal/key/repository/view/model"
	proj_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	user_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	authnKeysTable = "management.authn_keys"
)

type AuthNKeys struct {
	handler
	subscription *v1.Subscription
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

func (k *AuthNKeys) Subscription() *v1.Subscription {
	return k.subscription
}

func (_ *AuthNKeys) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{user_model.UserAggregate, proj_model.ProjectAggregate}
}

func (k *AuthNKeys) CurrentSequence() (uint64, error) {
	sequence, err := k.view.GetLatestAuthNKeySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (k *AuthNKeys) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := k.view.GetLatestAuthNKeySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(k.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (k *AuthNKeys) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case user_model.UserAggregate,
		proj_model.ProjectAggregate:
		err = k.processAuthNKeys(event)
	}
	return err
}

func (k *AuthNKeys) processAuthNKeys(event *es_models.Event) (err error) {
	key := new(key_model.AuthNKeyView)
	switch event.Type {
	case user_model.MachineKeyAdded,
		proj_model.ClientKeyAdded:
		err = key.AppendEvent(event)
		if key.ExpirationDate.Before(time.Now()) {
			return k.view.ProcessedAuthNKeySequence(event)
		}
	case user_model.MachineKeyRemoved:
		err = key.SetUserData(event)
		if err != nil {
			return err
		}
		return k.view.DeleteAuthNKey(key.ID, event)
	case proj_model.ClientKeyRemoved:
		err = key.SetClientData(event)
		if err != nil {
			return err
		}
		return k.view.DeleteAuthNKey(key.ID, event)
	case user_model.UserRemoved,
		proj_model.ApplicationRemoved:
		return k.view.DeleteAuthNKeysByObjectID(event.AggregateID, event)
	default:
		return k.view.ProcessedAuthNKeySequence(event)
	}
	if err != nil {
		return err
	}
	return k.view.PutAuthNKey(key, event)
}

func (d *AuthNKeys) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-S9fe", "id", event.AggregateID).WithError(err).Warn("something went wrong in machine key handler")
	return spooler.HandleError(event, err, d.view.GetLatestAuthNKeyFailedEvent, d.view.ProcessedAuthNKeyFailedEvent, d.view.ProcessedAuthNKeySequence, d.errorCountUntilSkip)
}

func (d *AuthNKeys) OnSuccess() error {
	return spooler.HandleSuccess(d.view.UpdateAuthNKeySpoolerRunTimestamp)
}
