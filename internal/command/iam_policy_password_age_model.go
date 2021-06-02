package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/policy"
)

type IAMPasswordAgePolicyWriteModel struct {
	PasswordAgePolicyWriteModel
}

func NewIAMPasswordAgePolicyWriteModel() *IAMPasswordAgePolicyWriteModel {
	return &IAMPasswordAgePolicyWriteModel{
		PasswordAgePolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *IAMPasswordAgePolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.PasswordAgePolicyAddedEvent:
			wm.PasswordAgePolicyWriteModel.AppendEvents(&e.PasswordAgePolicyAddedEvent)
		case *iam.PasswordAgePolicyChangedEvent:
			wm.PasswordAgePolicyWriteModel.AppendEvents(&e.PasswordAgePolicyChangedEvent)
		}
	}
}

func (wm *IAMPasswordAgePolicyWriteModel) Reduce() error {
	return wm.PasswordAgePolicyWriteModel.Reduce()
}

func (wm *IAMPasswordAgePolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.PasswordAgePolicyWriteModel.AggregateID).
		EventTypes(
			iam.PasswordAgePolicyAddedEventType,
			iam.PasswordAgePolicyChangedEventType).
		SearchQueryBuilder()
}

func (wm *IAMPasswordAgePolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	expireWarnDays,
	maxAgeDays uint64) (*iam.PasswordAgePolicyChangedEvent, bool) {
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
	changedEvent, err := iam.NewPasswordAgePolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
