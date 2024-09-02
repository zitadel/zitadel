package debug_events

import (
	"fmt"

	"github.com/zitadel/zitadel/internal/command"
	debug_events "github.com/zitadel/zitadel/pkg/grpc/resources/debug_events/v3alpha"
)

func debugEventsFromRequest(req *debug_events.CreateDebugEventsRequest) *command.DebugEvents {
	reqEvents := req.GetEvents()
	events := make([]command.DebugEvent, len(reqEvents))
	for i, event := range reqEvents {
		events[i] = debugEventFromRequest(event)
	}

	return &command.DebugEvents{
		AggregateID: req.GetAggregateId(),
		Events:      events,
	}
}

func debugEventFromRequest(event *debug_events.Event) command.DebugEvent {
	switch e := event.Event.(type) {
	case *debug_events.Event_Add:
		return command.DebugEventAdded{
			ProjectionSleep: e.Add.GetProjectionSleep().AsDuration(),
			Blob:            e.Add.Blob,
		}

	case *debug_events.Event_Change:
		return command.DebugEventChanged{
			ProjectionSleep: e.Change.GetProjectionSleep().AsDuration(),
			Blob:            e.Change.Blob,
		}

	case *debug_events.Event_Remove:
		return command.DebugEventRemoved{
			ProjectionSleep: e.Remove.GetProjectionSleep().AsDuration(),
		}

	default:
		panic(fmt.Errorf("invalid debug event type %T", event.Event))
	}
}
