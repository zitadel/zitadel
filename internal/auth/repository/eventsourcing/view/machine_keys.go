package view

import (
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"time"
)

const (
	machineKeyTable = "auth.machine_keys"
)

func (v *View) MachineKeyByIDs(userID, keyID string) (*model.MachineKeyView, error) {
	return view.MachineKeyByIDs(v.Db, machineKeyTable, userID, keyID)
}

func (v *View) MachineKeysByUserID(userID string) ([]*model.MachineKeyView, error) {
	return view.MachineKeysByUserID(v.Db, machineKeyTable, userID)
}

func (v *View) MachineKeyByID(keyID string) (*model.MachineKeyView, error) {
	return view.MachineKeyByID(v.Db, machineKeyTable, keyID)
}

func (v *View) SearchMachineKeys(request *usr_model.MachineKeySearchRequest) ([]*model.MachineKeyView, uint64, error) {
	return view.SearchMachineKeys(v.Db, machineKeyTable, request)
}

func (v *View) PutMachineKey(key *model.MachineKeyView, sequence uint64, eventTimestamp time.Time) error {
	err := view.PutMachineKey(v.Db, machineKeyTable, key)
	if err != nil {
		return err
	}
	if sequence != 0 {
		return v.ProcessedMachineKeySequence(sequence, eventTimestamp)
	}
	return nil
}

func (v *View) DeleteMachineKey(keyID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteMachineKey(v.Db, machineKeyTable, keyID)
	if err != nil {
		return nil
	}
	return v.ProcessedMachineKeySequence(eventSequence, eventTimestamp)
}

func (v *View) DeleteMachineKeysByUserID(userID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteMachineKey(v.Db, machineKeyTable, userID)
	if err != nil {
		return nil
	}
	return v.ProcessedMachineKeySequence(eventSequence, eventTimestamp)
}

func (v *View) GetLatestMachineKeySequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(machineKeyTable)
}

func (v *View) ProcessedMachineKeySequence(eventSequence uint64, eventTimestamp time.Time) error {
	return v.saveCurrentSequence(machineKeyTable, eventSequence, eventTimestamp)
}

func (v *View) UpdateMachineKeySpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(machineKeyTable)
}

func (v *View) GetLatestMachineKeyFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(machineKeyTable, sequence)
}

func (v *View) ProcessedMachineKeyFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
