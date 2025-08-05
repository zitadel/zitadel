package debug_events

import (
	"fmt"

	resource_object "github.com/zitadel/zitadel/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
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

func eventsStateToPB(state *query.DebugEventState) *debug_events.State {
	return &debug_events.State{
		Details: resource_object.DomainToDetailsPb(&state.ObjectDetails, object.OwnerType_OWNER_TYPE_INSTANCE, state.ResourceOwner),
		Blob:    state.Blob,
	}
}

func eventStatesToPB(states []query.DebugEventState) []*debug_events.State {
	out := make([]*debug_events.State, len(states))
	for i, state := range states {
		out[i] = eventsStateToPB(&state)
	}
	return out
}
