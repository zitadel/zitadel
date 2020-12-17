package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

type HumanEmailWriteModel struct {
	eventstore.WriteModel

	Email string
}

func NewHumanEmailWriteModel(userID string) *HumanEmailWriteModel {
	return &HumanEmailWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: userID,
		},
	}
}

func (wm *HumanEmailWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanEmailChangedEvent:
			wm.AppendEvents(e)
		case *user.HumanEmailVerifiedEvent:
			wm.AppendEvents(e)
			//TODO: Handle relevant User Events (remove, etc)
		}
	}
}

func (wm *HumanEmailWriteModel) Reduce() error {
	//TODO: implement
	return nil
}

func (wm *HumanEmailWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID)
}

func (wm *HumanEmailWriteModel) NewChangedEvent(
	ctx context.Context,
	email string,
) (*user.HumanEmailChangedEvent, bool) {
	hasChanged := false
	changedEvent := user.NewHumanEmailChangedEvent(ctx)
	if wm.Email != email {
		hasChanged = true
		changedEvent.EmailAddress = email
	}
	return changedEvent, hasChanged
}
