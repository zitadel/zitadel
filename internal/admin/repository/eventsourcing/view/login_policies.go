package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	loginPolicyTable = "adminapi.login_policies"
)

func (v *View) AllDefaultLoginPolicies() ([]*model.LoginPolicyView, error) {
	return view.GetDefaultLoginPolicies(v.Db, loginPolicyTable)
}

func (v *View) LoginPolicyByAggregateID(aggregateID string) (*model.LoginPolicyView, error) {
	return view.GetLoginPolicyByAggregateID(v.Db, loginPolicyTable, aggregateID)
}

func (v *View) PutLoginPolicy(policy *model.LoginPolicyView, event *models.Event) error {
	err := view.PutLoginPolicy(v.Db, loginPolicyTable, policy)
	if err != nil {
		return err
	}
	return v.ProcessedLoginPolicySequence(event)
}

func (v *View) PutLoginPolicies(policies []*model.LoginPolicyView, event *models.Event) error {
	err := view.PutLoginPolicies(v.Db, loginPolicyTable, policies...)
	if err != nil {
		return err
	}
	return v.ProcessedLoginPolicySequence(event)
}

func (v *View) DeleteLoginPolicy(aggregateID string, event *models.Event) error {
	err := view.DeleteLoginPolicy(v.Db, loginPolicyTable, aggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedLoginPolicySequence(event)
}

func (v *View) GetLatestLoginPolicySequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(loginPolicyTable)
}

func (v *View) ProcessedLoginPolicySequence(event *models.Event) error {
	return v.saveCurrentSequence(loginPolicyTable, event)
}

func (v *View) UpdateLoginPolicySpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(loginPolicyTable)
}

func (v *View) GetLatestLoginPolicyFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(loginPolicyTable, sequence)
}

func (v *View) ProcessedLoginPolicyFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
