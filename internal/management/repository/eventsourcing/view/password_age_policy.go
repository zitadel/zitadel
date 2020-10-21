package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	passwordAgePolicyTable = "management.password_age_policies"
)

func (v *View) PasswordAgePolicyByAggregateID(aggregateID string) (*model.PasswordAgePolicyView, error) {
	return view.GetPasswordAgePolicyByAggregateID(v.Db, passwordAgePolicyTable, aggregateID)
}

func (v *View) PutPasswordAgePolicy(policy *model.PasswordAgePolicyView, sequence uint64) error {
	err := view.PutPasswordAgePolicy(v.Db, passwordAgePolicyTable, policy)
	if err != nil {
		return err
	}
	return v.ProcessedPasswordAgePolicySequence(sequence)
}

func (v *View) DeletePasswordAgePolicy(aggregateID string, eventSequence uint64) error {
	err := view.DeletePasswordAgePolicy(v.Db, passwordAgePolicyTable, aggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedPasswordAgePolicySequence(eventSequence)
}

func (v *View) GetLatestPasswordAgePolicySequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(passwordAgePolicyTable)
}

func (v *View) ProcessedPasswordAgePolicySequence(eventSequence uint64) error {
	return v.saveCurrentSequence(passwordAgePolicyTable, eventSequence)
}

func (v *View) GetLatestPasswordAgePolicyFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(passwordAgePolicyTable, sequence)
}

func (v *View) ProcessedPasswordAgePolicyFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
