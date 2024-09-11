package projection

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/debug_events"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	DebugEventsProjectionTable = "projections.debug_events"

	DebugEventsColumnID            = "id"
	DebugEventsColumnCreationDate  = "creation_date"
	DebugEventsColumnChangeDate    = "change_date"
	DebugEventsColumnResourceOwner = "resource_owner"
	DebugEventsColumnInstanceID    = "instance_id"
	DebugEventsColumnSequence      = "sequence"
	DebugEventsColumnBlob          = "blob"
)

type debugEventsProjection struct{}

func (*debugEventsProjection) Name() string {
	return DebugEventsProjectionTable
}

func newDebugEventsProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(debugEventsProjection))
}

// Init implements [handler.initializer]
func (p *debugEventsProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(DebugEventsColumnID, handler.ColumnTypeText),
			handler.NewColumn(DebugEventsColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(DebugEventsColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(DebugEventsColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(DebugEventsColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(DebugEventsColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(DebugEventsColumnBlob, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(DebugEventsColumnInstanceID, DebugEventsColumnID),
		),
	)
}

func (p *debugEventsProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: debug_events.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  debug_events.AddedEventType,
					Reduce: p.reduceDebugEventAdded,
				},
				{
					Event:  debug_events.ChangedEventType,
					Reduce: p.reduceDebugEventChanged,
				},
				{
					Event:  debug_events.RemovedEventType,
					Reduce: p.reduceDebugEventRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(DebugEventsColumnInstanceID),
				},
			},
		},
	}
}

func (p *debugEventsProjection) reduceDebugEventAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*debug_events.AddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-uYq4r", "reduce.wrong.event.type %s", debug_events.AddedEventType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddSleepStatement(e.ProjectionSleep),
		handler.AddCreateStatement([]handler.Column{
			handler.NewCol(DebugEventsColumnID, e.Aggregate().ID),
			handler.NewCol(DebugEventsColumnCreationDate, e.CreationDate()),
			handler.NewCol(DebugEventsColumnChangeDate, e.CreationDate()),
			handler.NewCol(DebugEventsColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(DebugEventsColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(DebugEventsColumnSequence, e.Sequence()),
			handler.NewCol(DebugEventsColumnBlob, gu.Value(e.Blob)),
		}),
	), nil
}

func (p *debugEventsProjection) reduceDebugEventChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*debug_events.ChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Bg8oM", "reduce.wrong.event.type %s", debug_events.ChangedEventType)
	}
	updateCols := []handler.Column{
		handler.NewCol(DebugEventsColumnChangeDate, e.CreationDate()),
		handler.NewCol(DebugEventsColumnSequence, e.Sequence()),
	}
	if e.Blob != nil {
		updateCols = append(updateCols,
			handler.NewCol(DebugEventsColumnBlob, *e.Blob),
		)
	}

	return handler.NewMultiStatement(
		e,
		handler.AddSleepStatement(e.ProjectionSleep),
		handler.AddUpdateStatement(updateCols,
			[]handler.Condition{
				handler.NewCond(DebugEventsColumnID, e.Aggregate().ID),
				handler.NewCond(DebugEventsColumnInstanceID, e.Aggregate().InstanceID),
			}),
	), nil
}

func (p *debugEventsProjection) reduceDebugEventRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*debug_events.RemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-DgMSg", "reduce.wrong.event.type %s", debug_events.RemovedEventType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddSleepStatement(e.ProjectionSleep),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(DebugEventsColumnID, e.Aggregate().ID),
				handler.NewCond(DebugEventsColumnInstanceID, e.Aggregate().InstanceID),
			}),
	), nil
}
