package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/debug_events"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type DebugEvents struct {
	AggregateID string
	Events      []DebugEvent
}

type DebugEvent interface {
	isADebugEvent()
}

type DebugEventAdded struct {
	ProjectionSleep time.Duration
	Blob            *string
}

type DebugEventChanged struct {
	ProjectionSleep time.Duration
	Blob            *string
}

type DebugEventRemoved struct {
	ProjectionSleep time.Duration
}

func (DebugEventAdded) isADebugEvent()   {}
func (DebugEventChanged) isADebugEvent() {}
func (DebugEventRemoved) isADebugEvent() {}

func (c *Commands) CreateDebugEvents(ctx context.Context, dbe *DebugEvents) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	model := NewDebugEventsWriteModel(dbe.AggregateID, authz.GetInstance(ctx).InstanceID())
	if err = c.eventstore.FilterToQueryReducer(ctx, model); err != nil {
		return nil, err
	}
	aggr := debug_events.AggregateFromWriteModel(ctx, &model.WriteModel)

	cmds := make([]eventstore.Command, len(dbe.Events))
	for i, event := range dbe.Events {
		var cmd eventstore.Command
		switch e := event.(type) {
		case DebugEventAdded:
			if model.State.Exists() {
				return nil, zerrors.ThrowAlreadyExists(nil, "COMMAND-Aex6j", "debug aggregate already exists")
			}
			cmd = debug_events.NewAddedEvent(ctx, aggr, e.ProjectionSleep, e.Blob)
		case DebugEventChanged:
			if !model.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "COMMAND-Thie6", "debug aggregate not found")
			}
			cmd = debug_events.NewChangedEvent(ctx, aggr, e.ProjectionSleep, e.Blob)
		case DebugEventRemoved:
			if !model.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "COMMAND-Ohna9", "debug aggregate not found")
			}
			cmd = debug_events.NewRemovedEvent(ctx, aggr, e.ProjectionSleep)
		}

		cmds[i] = cmd
		// be sure the state of the last event is reduced before handling the next one.
		model.reduceEvent(cmd.(eventstore.Event))
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}
