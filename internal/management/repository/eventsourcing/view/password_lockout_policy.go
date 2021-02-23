package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	passwordLockoutPolicyTable = "management.password_lockout_policies"
)

func (v *View) PasswordLockoutPolicyByAggregateID(aggregateID string) (*model.PasswordLockoutPolicyView, error) {
	return view.GetPasswordLockoutPolicyByAggregateID(v.Db, passwordLockoutPolicyTable, aggregateID)
}

func (v *View) PutPasswordLockoutPolicy(policy *model.PasswordLockoutPolicyView, event *models.Event) error {
	err := view.PutPasswordLockoutPolicy(v.Db, passwordLockoutPolicyTable, policy)
	if err != nil {
		return err
	}
	return v.ProcessedPasswordLockoutPolicySequence(event)
}

func (v *View) DeletePasswordLockoutPolicy(aggregateID string, event *models.Event) error {
	err := view.DeletePasswordLockoutPolicy(v.Db, passwordLockoutPolicyTable, aggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedPasswordLockoutPolicySequence(event)
}

func (v *View) GetLatestPasswordLockoutPolicySequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(passwordLockoutPolicyTable)
}

func (v *View) ProcessedPasswordLockoutPolicySequence(event *models.Event) error {
	return v.saveCurrentSequence(passwordLockoutPolicyTable, event)
}

func (v *View) UpdatePasswordLockoutPolicySpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(passwordLockoutPolicyTable)
}

func (v *View) GetLatestPasswordLockoutPolicyFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(passwordLockoutPolicyTable, sequence)
}

func (v *View) ProcessedPasswordLockoutPolicyFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
