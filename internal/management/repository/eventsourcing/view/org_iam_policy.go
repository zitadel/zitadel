package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	orgIAMPolicyTable = "management.org_iam_policies"
)

func (v *View) OrgIAMPolicyByAggregateID(aggregateID string) (*model.OrgIAMPolicyView, error) {
	return view.GetOrgIAMPolicyByAggregateID(v.Db, orgIAMPolicyTable, aggregateID)
}

func (v *View) PutOrgIAMPolicy(policy *model.OrgIAMPolicyView, event *models.Event) error {
	err := view.PutOrgIAMPolicy(v.Db, orgIAMPolicyTable, policy)
	if err != nil {
		return err
	}
	return v.ProcessedOrgIAMPolicySequence(event)
}

func (v *View) DeleteOrgIAMPolicy(aggregateID string, event *models.Event) error {
	err := view.DeleteOrgIAMPolicy(v.Db, orgIAMPolicyTable, aggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedOrgIAMPolicySequence(event)
}

func (v *View) GetLatestOrgIAMPolicySequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(orgIAMPolicyTable)
}

func (v *View) ProcessedOrgIAMPolicySequence(event *models.Event) error {
	return v.saveCurrentSequence(orgIAMPolicyTable, event)
}

func (v *View) UpdateOrgIAMPolicySpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(orgIAMPolicyTable)
}

func (v *View) GetLatestOrgIAMPolicyFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(orgIAMPolicyTable, sequence)
}

func (v *View) ProcessedOrgIAMPolicyFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
