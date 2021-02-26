package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	passwordAgePolicyTable = "adminapi.password_age_policies"
)

func (v *View) PasswordAgePolicyByAggregateID(aggregateID string) (*model.PasswordAgePolicyView, error) {
	return view.GetPasswordAgePolicyByAggregateID(v.Db, passwordAgePolicyTable, aggregateID)
}

func (v *View) PutPasswordAgePolicy(policy *model.PasswordAgePolicyView, event *models.Event) error {
	err := view.PutPasswordAgePolicy(v.Db, passwordAgePolicyTable, policy)
	if err != nil {
		return err
	}
	return v.ProcessedPasswordAgePolicySequence(event)
}

func (v *View) DeletePasswordAgePolicy(aggregateID string, event *models.Event) error {
	err := view.DeletePasswordAgePolicy(v.Db, passwordAgePolicyTable, aggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedPasswordAgePolicySequence(event)
}

func (v *View) GetLatestPasswordAgePolicySequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(passwordAgePolicyTable)
}

func (v *View) ProcessedPasswordAgePolicySequence(event *models.Event) error {
	return v.saveCurrentSequence(passwordAgePolicyTable, event)
}

func (v *View) UpdateProcessedPasswordAgePolicySpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(passwordAgePolicyTable)
}

func (v *View) GetLatestPasswordAgePolicyFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(passwordAgePolicyTable, sequence)
}

func (v *View) ProcessedPasswordAgePolicyFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
