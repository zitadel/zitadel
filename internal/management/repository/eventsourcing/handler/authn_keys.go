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
	authNKeyTable = "management.authn_keys"
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

func (m *AuthNKeys) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (d *AuthNKeys) ViewModel() string {
	return authNKeyTable
}

func (_ *AuthNKeys) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.UserAggregate}
}

func (k *AuthNKeys) CurrentSequence() (uint64, error) {
	sequence, err := k.view.GetLatestAuthNKeySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (d *AuthNKeys) EventQuery() (*models.SearchQuery, error) {
	sequence, err := d.view.GetLatestAuthNKeySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(d.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (d *AuthNKeys) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.UserAggregate:
		err = d.processAuthNKeys(event)
	}
	return err
}

func (d *AuthNKeys) processAuthNKeys(event *models.Event) (err error) {
	key := new(key_model.AuthNKeyView)
	switch event.Type {
	case model.MachineKeyAdded:
		err = key.AppendEvent(event)
		if key.ExpirationDate.Before(time.Now()) {
			return d.view.ProcessedAuthNKeySequence(event)
		}
	case model.MachineKeyRemoved:
		err = key.SetData(event)
		if err != nil {
			return err
		}
		return d.view.DeleteAuthNKey(key.ID, event)
	case model.UserRemoved:
		return d.view.DeleteAuthNKeysByObjectID(event.AggregateID, event)
	default:
		return d.view.ProcessedAuthNKeySequence(event)
	}
	if err != nil {
		return err
	}
	return d.view.PutAuthNKey(key, event)
}

func (d *AuthNKeys) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-S9fe", "id", event.AggregateID).WithError(err).Warn("something went wrong in machine key handler")
	return spooler.HandleError(event, err, d.view.GetLatestAuthNKeyFailedEvent, d.view.ProcessedAuthNKeyFailedEvent, d.view.ProcessedAuthNKeySequence, d.errorCountUntilSkip)
}

func (d *AuthNKeys) OnSuccess() error {
	return spooler.HandleSuccess(d.view.UpdateAuthNKeySpoolerRunTimestamp)
}
