package view

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	stylingTyble = "adminapi.styling"
)

func (v *View) StylingByAggregateIDAndState(aggregateID string, state int32) (*model.LabelPolicyView, error) {
	return view.GetLabelPolicyByAggregateIDAndState(v.Db, stylingTyble, aggregateID, state)
}

func (v *View) PutStyling(policy *model.LabelPolicyView, event *models.Event) error {
	err := view.PutLabelPolicy(v.Db, stylingTyble, policy)
	if err != nil {
		return err
	}
	return v.ProcessedStylingSequence(event)
}

func (v *View) GetLatestStylingSequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(stylingTyble)
}

func (v *View) ProcessedStylingSequence(event *models.Event) error {
	return v.saveCurrentSequence(stylingTyble, event)
}

func (v *View) UpdateStylingSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(stylingTyble)
}

func (v *View) GetLatestStylingFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(stylingTyble, sequence)
}

func (v *View) ProcessedStylingFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
