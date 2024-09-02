package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/debug_events"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	DebugProjectionTable = "projections.debug"

	DebugColumnID            = "id"
	DebugColumnCreationDate  = "creation_date"
	DebugColumnChangeDate    = "change_date"
	DebugColumnResourceOwner = "resource_owner"
	DebugColumnInstanceID    = "instance_id"
	DebugColumnSequence      = "sequence"
)

type debugProjection struct{}

func (*debugProjection) Name() string {
	return DebugProjectionTable
}

func newDebugProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(debugProjection))
}

// Init implements [handler.initializer]
func (p *debugProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(DebugColumnID, handler.ColumnTypeText),
			handler.NewColumn(DebugColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(DebugColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(DebugColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(DebugColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(DebugColumnSequence, handler.ColumnTypeInt64),
		},
			handler.NewPrimaryKey(DebugColumnInstanceID, DebugColumnID),
		),
	)
}

func (p *debugProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: debug_events.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  debug_events.AddedEventType,
					Reduce: p.reduceDebugAdded,
				},
				{
					Event:  debug_events.ChangedEventType,
					Reduce: p.reduceDebugChanged,
				},
				{
					Event:  debug_events.RemovedEventType,
					Reduce: p.reduceDebugRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(DebugColumnInstanceID),
				},
			},
		},
	}
}

func (p *debugProjection) reduceDebugAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*debug_events.AddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-uYq4r", "reduce.wrong.event.type %s", debug_events.AddedEventType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddSleepStatement(e.ProjectionSleep),
		handler.AddCreateStatement([]handler.Column{
			handler.NewCol(DebugColumnID, e.Aggregate().ID),
			handler.NewCol(DebugColumnCreationDate, e.CreationDate()),
			handler.NewCol(DebugColumnChangeDate, e.CreationDate()),
			handler.NewCol(DebugColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(DebugColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(DebugColumnSequence, e.Sequence()),
		}),
	), nil
}

func (p *debugProjection) reduceDebugChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*debug_events.ChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Bg8oM", "reduce.wrong.event.type %s", debug_events.ChangedEventType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddSleepStatement(e.ProjectionSleep),
		handler.AddUpdateStatement([]handler.Column{
			handler.NewCol(DebugColumnChangeDate, e.CreationDate()),
			handler.NewCol(DebugColumnSequence, e.Sequence()),
		},
			[]handler.Condition{
				handler.NewCond(DebugColumnID, e.Aggregate().ID),
				handler.NewCond(DebugColumnInstanceID, e.Aggregate().InstanceID),
			}),
	), nil
}

func (p *debugProjection) reduceDebugRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*debug_events.RemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-DgMSg", "reduce.wrong.event.type %s", debug_events.RemovedEventType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddSleepStatement(e.ProjectionSleep),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(DebugColumnID, e.Aggregate().ID),
				handler.NewCond(DebugColumnInstanceID, e.Aggregate().InstanceID),
			}),
	), nil
}
