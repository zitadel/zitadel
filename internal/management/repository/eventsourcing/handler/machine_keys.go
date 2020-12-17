package handler

import (
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	usr_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type MachineKeys struct {
	handler
}

const (
	machineKeysTable = "management.machine_keys"
)

func (d *MachineKeys) ViewModel() string {
	return machineKeysTable
}

func (_ *MachineKeys) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.UserAggregate}
}

func (k *MachineKeys) CurrentSequence(event *models.Event) (uint64, error) {
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
			return d.view.ProcessedMachineKeySequence(event.Sequence, event.CreationDate)
		}
	case model.MachineKeyRemoved:
		err = key.SetData(event)
		if err != nil {
			return err
		}
		return d.view.DeleteMachineKey(key.ID, event.Sequence, event.CreationDate)
	case model.UserRemoved:
		return d.view.DeleteMachineKeysByUserID(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return d.view.ProcessedMachineKeySequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return d.view.PutMachineKey(key, key.Sequence, event.CreationDate)
}

func (d *MachineKeys) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-S9fe", "id", event.AggregateID).WithError(err).Warn("something went wrong in machine key handler")
	return spooler.HandleError(event, err, d.view.GetLatestMachineKeyFailedEvent, d.view.ProcessedMachineKeyFailedEvent, d.view.ProcessedMachineKeySequence, d.errorCountUntilSkip)
}

func (d *MachineKeys) OnSuccess() error {
	return spooler.HandleSuccess(d.view.UpdateMachineKeySpoolerRunTimestamp)
}
