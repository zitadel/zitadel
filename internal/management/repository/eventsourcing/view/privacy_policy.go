package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	privacyPolicyTable = "management.privacy_policies"
)

func (v *View) PrivacyPolicyByAggregateID(aggregateID string) (*model.PrivacyPolicyView, error) {
	return view.GetPrivacyPolicyByAggregateID(v.Db, privacyPolicyTable, aggregateID)
}

func (v *View) PutPrivacyPolicy(policy *model.PrivacyPolicyView, event *models.Event) error {
	err := view.PutPrivacyPolicy(v.Db, privacyPolicyTable, policy)
	if err != nil {
		return err
	}
	return v.ProcessedPrivacyPolicySequence(event)
}

func (v *View) DeletePrivacyPolicy(aggregateID string, event *models.Event) error {
	err := view.DeletePrivacyPolicy(v.Db, privacyPolicyTable, aggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedPrivacyPolicySequence(event)
}

func (v *View) GetLatestPrivacyPolicySequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(privacyPolicyTable)
}

func (v *View) ProcessedPrivacyPolicySequence(event *models.Event) error {
	return v.saveCurrentSequence(privacyPolicyTable, event)
}

func (v *View) UpdatePrivacyPolicySpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(privacyPolicyTable)
}

func (v *View) GetLatestPrivacyPolicyFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(privacyPolicyTable, sequence)
}

func (v *View) ProcessedPrivacyPolicyFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
