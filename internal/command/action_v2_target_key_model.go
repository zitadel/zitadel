package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/target"
)

type TargetKeyWriteModel struct {
	eventstore.WriteModel
	KeyID string

	KeyExists    bool
	TargetExists bool
}

func NewTargetKeyWriteModel(targetID, keyID string, resourceOwner string) *TargetKeyWriteModel {
	return &TargetKeyWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   targetID,
			ResourceOwner: resourceOwner,
		},
		KeyID: keyID,
	}
}

func (wm *TargetKeyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *target.KeyAddedEvent:
			if wm.KeyID != e.KeyID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *target.KeyRemovedEvent:
			if wm.KeyID != e.KeyID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *target.AddedEvent:
			wm.WriteModel.AppendEvents(e)
		case *target.RemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *TargetKeyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *target.KeyAddedEvent:
			wm.KeyID = e.KeyID
			wm.KeyExists = true
		case *target.KeyRemovedEvent:
			wm.KeyExists = false
		case *target.AddedEvent:
			wm.TargetExists = true
		case *target.RemovedEvent:
			wm.TargetExists = false
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *TargetKeyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(target.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(target.KeyAddedEventType,
			target.KeyRemovedEventType,
			target.AddedEventType,
			target.RemovedEventType).
		Builder()
}
