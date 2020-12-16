package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type IAMPasswordAgePolicyWriteModel struct {
	PasswordAgePolicyWriteModel
}

func NewIAMPasswordAgePolicyWriteModel(iamID string) *IAMPasswordAgePolicyWriteModel {
	return &IAMPasswordAgePolicyWriteModel{
		PasswordAgePolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: iamID,
			},
		},
	}
}

func (wm *IAMPasswordAgePolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.PasswordAgePolicyAddedEvent:
			wm.PasswordAgePolicyWriteModel.AppendEvents(&e.PassowordAgePolicyAddedEvent)
		case *iam.PasswordAgePolicyChangedEvent:
			wm.PasswordAgePolicyWriteModel.AppendEvents(&e.PasswordAgePolicyChangedEvent)
		}
	}
}

func (wm *IAMPasswordAgePolicyWriteModel) Reduce() error {
	return wm.PasswordAgePolicyWriteModel.Reduce()
}

func (wm *IAMPasswordAgePolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.PasswordAgePolicyWriteModel.AggregateID)
}

func (wm *IAMPasswordAgePolicyWriteModel) NewChangedEvent(expireWarnDays, maxAgeDays uint64) (*iam.PasswordAgePolicyChangedEvent, bool) {
	hasChanged := false
	changedEvent := &iam.PasswordAgePolicyChangedEvent{}
	if wm.ExpireWarnDays == expireWarnDays {
		hasChanged = true
		changedEvent.ExpireWarnDays = expireWarnDays
	}
	if wm.MaxAgeDays == maxAgeDays {
		hasChanged = true
		changedEvent.MaxAgeDays = maxAgeDays
	}
	return changedEvent, hasChanged
}
