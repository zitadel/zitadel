package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
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
		case *org.PasswordAgePolicyRemovedEvent:
			wm.PasswordAgePolicyWriteModel.AppendEvents(&e.PasswordAgePolicyRemovedEvent)
		}
	}
}

func (wm *OrgPasswordAgePolicyWriteModel) Reduce() error {
	return wm.PasswordAgePolicyWriteModel.Reduce()
}

func (wm *OrgPasswordAgePolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.PasswordAgePolicyWriteModel.AggregateID).
		EventTypes(
			org.PasswordAgePolicyAddedEventType,
			org.PasswordAgePolicyChangedEventType,
			org.PasswordAgePolicyRemovedEventType).
		SearchQueryBuilder()
}

func (wm *OrgPasswordAgePolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	expireWarnDays,
	maxAgeDays uint64) (*org.PasswordAgePolicyChangedEvent, bool) {
	changes := make([]policy.PasswordAgePolicyChanges, 0)
	if wm.ExpireWarnDays != expireWarnDays {
		changes = append(changes, policy.ChangeExpireWarnDays(expireWarnDays))
	}
	if wm.MaxAgeDays != maxAgeDays {
		changes = append(changes, policy.ChangeMaxAgeDays(maxAgeDays))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := org.NewPasswordAgePolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
