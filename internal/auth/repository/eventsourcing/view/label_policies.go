package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	labelPolicyTable = "auth.label_policies"
)

func (v *View) LabelPolicyByAggregateID(aggregateID string) (*model.LabelPolicyView, error) {
	return view.GetLabelPolicyByAggregateID(v.Db, labelPolicyTable, aggregateID)
}

func (v *View) PutLabelPolicy(policy *model.LabelPolicyView, event *models.Event) error {
	err := view.PutLabelPolicy(v.Db, labelPolicyTable, policy)
	if err != nil {
		return err
	}
	return v.ProcessedLabelPolicySequence(event)
}

func (v *View) DeleteLabelPolicy(aggregateID string, event *models.Event) error {
	err := view.DeleteLabelPolicy(v.Db, labelPolicyTable, aggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedLabelPolicySequence(event)
}

func (v *View) GetLatestLabelPolicySequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(labelPolicyTable)
}

func (v *View) ProcessedLabelPolicySequence(event *models.Event) error {
	return v.saveCurrentSequence(labelPolicyTable, event)
}

func (v *View) UpdateLabelPolicySpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(labelPolicyTable)
}

func (v *View) GetLatestLabelPolicyFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(labelPolicyTable, sequence)
}

func (v *View) ProcessedLabelPolicyFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
