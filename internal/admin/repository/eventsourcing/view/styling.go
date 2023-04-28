package view

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/iam/repository/view"
	"github.com/zitadel/zitadel/internal/iam/repository/view/model"
	global_view "github.com/zitadel/zitadel/internal/view/repository"
)

const (
	stylingTyble = "adminapi.styling2"
)

func (v *View) StylingByAggregateIDAndState(aggregateID, instanceID string, state int32) (*model.LabelPolicyView, error) {
	return view.GetStylingByAggregateIDAndState(v.Db, stylingTyble, aggregateID, instanceID, state)
}

func (v *View) PutStyling(policy *model.LabelPolicyView, event *models.Event) error {
	err := view.PutStyling(v.Db, stylingTyble, policy)
	if err != nil {
		return err
	}
	return v.ProcessedStylingSequence(event)
}

func (v *View) DeleteInstanceStyling(event *models.Event) error {
	err := view.DeleteInstanceStyling(v.Db, stylingTyble, event.InstanceID)
	if err != nil {
		return err
	}
	return v.ProcessedStylingSequence(event)
}

func (v *View) UpdateOrgOwnerRemovedStyling(event *models.Event) error {
	err := view.UpdateOrgOwnerRemovedStyling(v.Db, stylingTyble, event.InstanceID, event.AggregateID)
	if err != nil {
		return err
	}
	return v.ProcessedStylingSequence(event)
}

func (v *View) GetLatestStylingSequence(ctx context.Context, instanceID string) (*global_view.CurrentSequence, error) {
	return v.latestSequence(ctx, stylingTyble, instanceID)
}

func (v *View) GetLatestStylingSequences(ctx context.Context, instanceIDs []string) ([]*global_view.CurrentSequence, error) {
	return v.latestSequences(ctx, stylingTyble, instanceIDs)
}

func (v *View) ProcessedStylingSequence(event *models.Event) error {
	return v.saveCurrentSequence(stylingTyble, event)
}

func (v *View) UpdateStylingSpoolerRunTimestamp(instanceIDs []string) error {
	return v.updateSpoolerRunSequence(stylingTyble, instanceIDs)
}

func (v *View) GetLatestStylingFailedEvent(sequence uint64, instanceID string) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(stylingTyble, instanceID, sequence)
}

func (v *View) ProcessedStylingFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
