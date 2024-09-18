package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	debug "github.com/zitadel/zitadel/internal/repository/debug_events"
)

type DebugEventsWriteModel struct {
	eventstore.WriteModel
	State domain.DebugEventsState
	Blob  string
}

func NewDebugEventsWriteModel(aggregateID, resourceOwner string) *DebugEventsWriteModel {
	return &DebugEventsWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   aggregateID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *DebugEventsWriteModel) AppendEvents(events ...eventstore.Event) {
	wm.WriteModel.AppendEvents(events...)
}

func (wm *DebugEventsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		wm.reduceEvent(event)
	}
	return wm.WriteModel.Reduce()
}

func (wm *DebugEventsWriteModel) reduceEvent(event eventstore.Event) {
	if event.Aggregate().ID != wm.AggregateID {
		return
	}
	switch e := event.(type) {
	case *debug.AddedEvent:
		wm.State = domain.DebugEventsStateInitial
		if e.Blob != nil {
			wm.Blob = *e.Blob
		}
	case *debug.ChangedEvent:
		wm.State = domain.DebugEventsStateChanged
		if e.Blob != nil {
			wm.Blob = *e.Blob
		}
	case *debug.RemovedEvent:
		wm.State = domain.DebugEventsStateRemoved
		wm.Blob = ""
	}
}

func (wm *DebugEventsWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(debug.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			debug.AddedEventType,
			debug.ChangedEventType,
			debug.RemovedEventType,
		).
		Builder()
}
