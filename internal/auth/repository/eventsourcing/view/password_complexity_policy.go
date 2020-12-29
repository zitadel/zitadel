package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	passwordComplexityPolicyTable = "auth.password_complexity_policies"
)

func (v *View) PasswordComplexityPolicyByAggregateID(aggregateID string) (*model.PasswordComplexityPolicyView, error) {
	return view.GetPasswordComplexityPolicyByAggregateID(v.Db, passwordComplexityPolicyTable, aggregateID)
}

func (v *View) PutPasswordComplexityPolicy(policy *model.PasswordComplexityPolicyView, event *models.Event) error {
	err := view.PutPasswordComplexityPolicy(v.Db, passwordComplexityPolicyTable, policy)
	if err != nil {
		return err
	}
	return v.ProcessedPasswordComplexityPolicySequence(event)
}

func (v *View) DeletePasswordComplexityPolicy(aggregateID string, event *models.Event) error {
	err := view.DeletePasswordComplexityPolicy(v.Db, passwordComplexityPolicyTable, aggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedPasswordComplexityPolicySequence(event)
}

func (v *View) GetLatestPasswordComplexityPolicySequence(aggregateType string) (*global_view.CurrentSequence, error) {
	return v.latestSequence(passwordComplexityPolicyTable, aggregateType)
}

func (v *View) ProcessedPasswordComplexityPolicySequence(event *models.Event) error {
	return v.saveCurrentSequence(passwordComplexityPolicyTable, event)
}

func (v *View) UpdatePasswordComplexityPolicySpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(passwordComplexityPolicyTable)
}

func (v *View) GetLatestPasswordComplexityPolicyFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(passwordComplexityPolicyTable, sequence)
}

func (v *View) ProcessedPasswordComplexityPolicyFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
