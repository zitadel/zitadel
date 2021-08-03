package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	lockoutPolicyTable = "management.password_lockout_policies"
)

func (v *View) LockoutPolicyByAggregateID(aggregateID string) (*model.LockoutPolicyView, error) {
	return view.GetLockoutPolicyByAggregateID(v.Db, lockoutPolicyTable, aggregateID)
}

func (v *View) PutLockoutPolicy(policy *model.LockoutPolicyView, event *models.Event) error {
	err := view.PutLockoutPolicy(v.Db, lockoutPolicyTable, policy)
	if err != nil {
		return err
	}
	return v.ProcessedLockoutPolicySequence(event)
}

func (v *View) DeleteLockoutPolicy(aggregateID string, event *models.Event) error {
	err := view.DeleteLockoutPolicy(v.Db, lockoutPolicyTable, aggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedLockoutPolicySequence(event)
}

func (v *View) GetLatestLockoutPolicySequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(lockoutPolicyTable)
}

func (v *View) ProcessedLockoutPolicySequence(event *models.Event) error {
	return v.saveCurrentSequence(lockoutPolicyTable, event)
}

func (v *View) UpdateLockoutPolicySpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(lockoutPolicyTable)
}

func (v *View) GetLatestLockoutPolicyFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(lockoutPolicyTable, sequence)
}

func (v *View) ProcessedLockoutPolicyFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
