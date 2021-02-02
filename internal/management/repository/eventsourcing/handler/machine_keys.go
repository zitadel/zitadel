package handler

import (
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	usr_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	machineKeysTable = "management.machine_keys"
)

type MachineKeys struct {
	handler
	subscription *eventstore.Subscription
}

func newMachineKeys(handler handler) *MachineKeys {
	h := &MachineKeys{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *MachineKeys) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (d *MachineKeys) ViewModel() string {
	return machineKeysTable
}

func (_ *MachineKeys) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.UserAggregate}
}

func (k *MachineKeys) CurrentSequence() (uint64, error) {
	sequence, err := k.view.GetLatestMachineKeySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (d *MachineKeys) EventQuery() (*models.SearchQuery, error) {
	sequence, err := d.view.GetLatestMachineKeySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(d.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (d *MachineKeys) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.UserAggregate:
		err = d.processMachineKeys(event)
	}
	return err
}

func (d *MachineKeys) processMachineKeys(event *models.Event) (err error) {
	key := new(usr_model.MachineKeyView)
	switch event.Type {
	case model.MachineKeyAdded:
		err = key.AppendEvent(event)
		if key.ExpirationDate.Before(time.Now()) {
			return d.view.ProcessedMachineKeySequence(event)
		}
	case model.MachineKeyRemoved:
		err = key.SetData(event)
		if err != nil {
			return err
		}
		return d.view.DeleteMachineKey(key.ID, event)
	case model.UserRemoved:
		return d.view.DeleteMachineKeysByUserID(event.AggregateID, event)
	default:
		return d.view.ProcessedMachineKeySequence(event)
	}
	if err != nil {
		return err
	}
	return d.view.PutMachineKey(key, event)
}

func (d *MachineKeys) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-S9fe", "id", event.AggregateID).WithError(err).Warn("something went wrong in machine key handler")
	return spooler.HandleError(event, err, d.view.GetLatestMachineKeyFailedEvent, d.view.ProcessedMachineKeyFailedEvent, d.view.ProcessedMachineKeySequence, d.errorCountUntilSkip)
}

func (d *MachineKeys) OnSuccess() error {
	return spooler.HandleSuccess(d.view.UpdateMachineKeySpoolerRunTimestamp)
}
