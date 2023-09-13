package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/settings"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

const (
	DebugNotificationProviderTable = "projections.notification_providers"

	DebugNotificationProviderAggIDCol         = "aggregate_id"
	DebugNotificationProviderCreationDateCol  = "creation_date"
	DebugNotificationProviderChangeDateCol    = "change_date"
	DebugNotificationProviderSequenceCol      = "sequence"
	DebugNotificationProviderResourceOwnerCol = "resource_owner"
	DebugNotificationProviderInstanceIDCol    = "instance_id"
	DebugNotificationProviderStateCol         = "state"
	DebugNotificationProviderTypeCol          = "provider_type"
	DebugNotificationProviderCompactCol       = "compact"
)

type debugNotificationProviderProjection struct{}

func newDebugNotificationProviderProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(debugNotificationProviderProjection))
}

func (*debugNotificationProviderProjection) Name() string {
	return DebugNotificationProviderTable
}

func (*debugNotificationProviderProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(DebugNotificationProviderAggIDCol, handler.ColumnTypeText),
			handler.NewColumn(DebugNotificationProviderCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(DebugNotificationProviderChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(DebugNotificationProviderSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(DebugNotificationProviderResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(DebugNotificationProviderInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(DebugNotificationProviderStateCol, handler.ColumnTypeEnum),
			handler.NewColumn(DebugNotificationProviderTypeCol, handler.ColumnTypeEnum),
			handler.NewColumn(DebugNotificationProviderCompactCol, handler.ColumnTypeBool),
		},
			handler.NewPrimaryKey(DebugNotificationProviderInstanceIDCol, DebugNotificationProviderAggIDCol, DebugNotificationProviderTypeCol),
		),
	)
}

func (p *debugNotificationProviderProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.DebugNotificationProviderFileAddedEventType,
					Reduce: p.reduceDebugNotificationProviderAdded,
				},
				{
					Event:  instance.DebugNotificationProviderFileChangedEventType,
					Reduce: p.reduceDebugNotificationProviderChanged,
				},
				{
					Event:  instance.DebugNotificationProviderFileRemovedEventType,
					Reduce: p.reduceDebugNotificationProviderRemoved,
				},
				{
					Event:  instance.DebugNotificationProviderLogAddedEventType,
					Reduce: p.reduceDebugNotificationProviderAdded,
				},
				{
					Event:  instance.DebugNotificationProviderLogChangedEventType,
					Reduce: p.reduceDebugNotificationProviderChanged,
				},
				{
					Event:  instance.DebugNotificationProviderLogRemovedEventType,
					Reduce: p.reduceDebugNotificationProviderRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(DebugNotificationProviderInstanceIDCol),
				},
			},
		},
	}
}

func (p *debugNotificationProviderProjection) reduceDebugNotificationProviderAdded(event eventstore.Event) (*handler.Statement, error) {
	var providerEvent settings.DebugNotificationProviderAddedEvent
	var providerType domain.NotificationProviderType
	switch e := event.(type) {
	case *instance.DebugNotificationProviderFileAddedEvent:
		providerEvent = e.DebugNotificationProviderAddedEvent
		providerType = domain.NotificationProviderTypeFile
	case *instance.DebugNotificationProviderLogAddedEvent:
		providerEvent = e.DebugNotificationProviderAddedEvent
		providerType = domain.NotificationProviderTypeLog
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-pYPxS", "reduce.wrong.event.type %v", []eventstore.EventType{instance.DebugNotificationProviderFileAddedEventType, instance.DebugNotificationProviderLogAddedEventType})
	}

	return handler.NewCreateStatement(&providerEvent, []handler.Column{
		handler.NewCol(DebugNotificationProviderAggIDCol, providerEvent.Aggregate().ID),
		handler.NewCol(DebugNotificationProviderCreationDateCol, providerEvent.CreationDate()),
		handler.NewCol(DebugNotificationProviderChangeDateCol, providerEvent.CreationDate()),
		handler.NewCol(DebugNotificationProviderSequenceCol, providerEvent.Sequence()),
		handler.NewCol(DebugNotificationProviderResourceOwnerCol, providerEvent.Aggregate().ResourceOwner),
		handler.NewCol(DebugNotificationProviderInstanceIDCol, providerEvent.Aggregate().InstanceID),
		handler.NewCol(DebugNotificationProviderStateCol, domain.NotificationProviderStateActive),
		handler.NewCol(DebugNotificationProviderTypeCol, providerType),
		handler.NewCol(DebugNotificationProviderCompactCol, providerEvent.Compact),
	}), nil
}

func (p *debugNotificationProviderProjection) reduceDebugNotificationProviderChanged(event eventstore.Event) (*handler.Statement, error) {
	var providerEvent settings.DebugNotificationProviderChangedEvent
	var providerType domain.NotificationProviderType
	switch e := event.(type) {
	case *instance.DebugNotificationProviderFileChangedEvent:
		providerEvent = e.DebugNotificationProviderChangedEvent
		providerType = domain.NotificationProviderTypeFile
	case *instance.DebugNotificationProviderLogChangedEvent:
		providerEvent = e.DebugNotificationProviderChangedEvent
		providerType = domain.NotificationProviderTypeLog
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-pYPxS", "reduce.wrong.event.type %v", []eventstore.EventType{instance.DebugNotificationProviderFileChangedEventType, instance.DebugNotificationProviderLogChangedEventType})
	}

	cols := []handler.Column{
		handler.NewCol(DebugNotificationProviderChangeDateCol, providerEvent.CreationDate()),
		handler.NewCol(DebugNotificationProviderSequenceCol, providerEvent.Sequence()),
	}
	if providerEvent.Compact != nil {
		cols = append(cols, handler.NewCol(DebugNotificationProviderCompactCol, *providerEvent.Compact))
	}

	return handler.NewUpdateStatement(
		&providerEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(DebugNotificationProviderAggIDCol, providerEvent.Aggregate().ID),
			handler.NewCond(DebugNotificationProviderTypeCol, providerType),
			handler.NewCond(DebugNotificationProviderInstanceIDCol, providerEvent.Aggregate().InstanceID),
		},
	), nil
}

func (p *debugNotificationProviderProjection) reduceDebugNotificationProviderRemoved(event eventstore.Event) (*handler.Statement, error) {
	var providerEvent settings.DebugNotificationProviderRemovedEvent
	var providerType domain.NotificationProviderType
	switch e := event.(type) {
	case *instance.DebugNotificationProviderFileRemovedEvent:
		providerEvent = e.DebugNotificationProviderRemovedEvent
		providerType = domain.NotificationProviderTypeFile
	case *instance.DebugNotificationProviderLogRemovedEvent:
		providerEvent = e.DebugNotificationProviderRemovedEvent
		providerType = domain.NotificationProviderTypeLog
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-dow9f", "reduce.wrong.event.type %v", []eventstore.EventType{instance.DebugNotificationProviderFileRemovedEventType, instance.DebugNotificationProviderLogRemovedEventType})
	}

	return handler.NewDeleteStatement(
		&providerEvent,
		[]handler.Condition{
			handler.NewCond(DebugNotificationProviderAggIDCol, providerEvent.Aggregate().ID),
			handler.NewCond(DebugNotificationProviderTypeCol, providerType),
			handler.NewCond(DebugNotificationProviderInstanceIDCol, providerEvent.Aggregate().InstanceID),
		},
	), nil
}
