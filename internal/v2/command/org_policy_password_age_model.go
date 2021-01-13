package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

type OrgPasswordAgePolicyWriteModel struct {
	PasswordAgePolicyWriteModel
}

func NewOrgPasswordAgePolicyWriteModel(orgID string) *OrgPasswordAgePolicyWriteModel {
	return &OrgPasswordAgePolicyWriteModel{
		PasswordAgePolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgPasswordAgePolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.PasswordAgePolicyAddedEvent:
			wm.PasswordAgePolicyWriteModel.AppendEvents(&e.PasswordAgePolicyAddedEvent)
		case *org.PasswordAgePolicyChangedEvent:
			wm.PasswordAgePolicyWriteModel.AppendEvents(&e.PasswordAgePolicyChangedEvent)
		}
	}
}

func (wm *OrgPasswordAgePolicyWriteModel) Reduce() error {
	return wm.PasswordAgePolicyWriteModel.Reduce()
}

func (wm *OrgPasswordAgePolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.PasswordAgePolicyWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

func (wm *OrgPasswordAgePolicyWriteModel) NewChangedEvent(ctx context.Context, expireWarnDays, maxAgeDays uint64) (*org.PasswordAgePolicyChangedEvent, bool) {
	hasChanged := false
	changedEvent := org.NewPasswordAgePolicyChangedEvent(ctx)
	if wm.ExpireWarnDays != expireWarnDays {
		hasChanged = true
		changedEvent.ExpireWarnDays = &expireWarnDays
	}
	if wm.MaxAgeDays != maxAgeDays {
		hasChanged = true
		changedEvent.MaxAgeDays = &maxAgeDays
	}
	return changedEvent, hasChanged
}
