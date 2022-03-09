package projection

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/settings"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
)

const (
	DebugNotificationProviderTable = "zitadel.projections.notification_providers"

	DebugNotificationProviderAggIDCol         = "aggregate_id"
	DebugNotificationProviderCreationDateCol  = "creation_date"
	DebugNotificationProviderChangeDateCol    = "change_date"
	DebugNotificationProviderSequenceCol      = "sequence"
	DebugNotificationProviderResourceOwnerCol = "resource_owner"
	DebugNotificationProviderStateCol         = "state"
	DebugNotificationProviderTypeCol          = "provider_type"
	DebugNotificationProviderCompactCol       = "compact"
)

type DebugNotificationProviderProjection struct {
	crdb.StatementHandler
}

func NewDebugNotificationProviderProjection(ctx context.Context, config crdb.StatementHandlerConfig) *DebugNotificationProviderProjection {
	p := &DebugNotificationProviderProjection{}
	config.ProjectionName = DebugNotificationProviderTable
	config.Reducers = p.reducers()
	config.InitChecks = []*handler.Check{
		crdb.NewTableCheck(
			crdb.NewTable([]*crdb.Column{
				crdb.NewColumn(DebugNotificationProviderAggIDCol, crdb.ColumnTypeText),
				crdb.NewColumn(DebugNotificationProviderCreationDateCol, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(DebugNotificationProviderChangeDateCol, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(DebugNotificationProviderSequenceCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(DebugNotificationProviderResourceOwnerCol, crdb.ColumnTypeText),
				crdb.NewColumn(DebugNotificationProviderStateCol, crdb.ColumnTypeEnum),
				crdb.NewColumn(DebugNotificationProviderTypeCol, crdb.ColumnTypeEnum),
				crdb.NewColumn(DebugNotificationProviderCompactCol, crdb.ColumnTypeBool),
			},
				crdb.NewPrimaryKey(DebugNotificationProviderAggIDCol, DebugNotificationProviderTypeCol),
			),
		),
	}
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *DebugNotificationProviderProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.DebugNotificationProviderFileAddedEventType,
					Reduce: p.reduceDebugNotificationProviderAdded,
				},
				{
					Event:  iam.DebugNotificationProviderFileChangedEventType,
					Reduce: p.reduceDebugNotificationProviderChanged,
				},
				{
					Event:  iam.DebugNotificationProviderFileRemovedEventType,
					Reduce: p.reduceDebugNotificationProviderRemoved,
				},
				{
					Event:  iam.DebugNotificationProviderLogAddedEventType,
					Reduce: p.reduceDebugNotificationProviderAdded,
				},
				{
					Event:  iam.DebugNotificationProviderLogChangedEventType,
					Reduce: p.reduceDebugNotificationProviderChanged,
				},
				{
					Event:  iam.DebugNotificationProviderLogRemovedEventType,
					Reduce: p.reduceDebugNotificationProviderRemoved,
				},
			},
		},
	}
}

func (p *DebugNotificationProviderProjection) reduceDebugNotificationProviderAdded(event eventstore.Event) (*handler.Statement, error) {
	var providerEvent settings.DebugNotificationProviderAddedEvent
	var providerType domain.NotificationProviderType
	switch e := event.(type) {
	case *iam.DebugNotificationProviderFileAddedEvent:
		providerEvent = e.DebugNotificationProviderAddedEvent
		providerType = domain.NotificationProviderTypeFile
	case *iam.DebugNotificationProviderLogAddedEvent:
		providerEvent = e.DebugNotificationProviderAddedEvent
		providerType = domain.NotificationProviderTypeLog
	default:
		logging.WithFields("seq", event.Sequence(), "expectedTypes", []eventstore.EventType{iam.DebugNotificationProviderFileAddedEventType, iam.DebugNotificationProviderLogAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-pYPxS", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(&providerEvent, []handler.Column{
		handler.NewCol(DebugNotificationProviderAggIDCol, providerEvent.Aggregate().ID),
		handler.NewCol(DebugNotificationProviderCreationDateCol, providerEvent.CreationDate()),
		handler.NewCol(DebugNotificationProviderChangeDateCol, providerEvent.CreationDate()),
		handler.NewCol(DebugNotificationProviderSequenceCol, providerEvent.Sequence()),
		handler.NewCol(DebugNotificationProviderResourceOwnerCol, providerEvent.Aggregate().ResourceOwner),
		handler.NewCol(DebugNotificationProviderStateCol, domain.NotificationProviderStateActive),
		handler.NewCol(DebugNotificationProviderTypeCol, providerType),
		handler.NewCol(DebugNotificationProviderCompactCol, providerEvent.Compact),
	}), nil
}

func (p *DebugNotificationProviderProjection) reduceDebugNotificationProviderChanged(event eventstore.Event) (*handler.Statement, error) {
	var providerEvent settings.DebugNotificationProviderChangedEvent
	var providerType domain.NotificationProviderType
	switch e := event.(type) {
	case *iam.DebugNotificationProviderFileChangedEvent:
		providerEvent = e.DebugNotificationProviderChangedEvent
		providerType = domain.NotificationProviderTypeFile
	case *iam.DebugNotificationProviderLogChangedEvent:
		providerEvent = e.DebugNotificationProviderChangedEvent
		providerType = domain.NotificationProviderTypeLog
	default:
		logging.WithFields("seq", event.Sequence(), "expectedTypes", []eventstore.EventType{iam.DebugNotificationProviderFileChangedEventType, iam.DebugNotificationProviderLogChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-pYPxS", "reduce.wrong.event.type")
	}

	cols := []handler.Column{
		handler.NewCol(DebugNotificationProviderChangeDateCol, providerEvent.CreationDate()),
		handler.NewCol(DebugNotificationProviderSequenceCol, providerEvent.Sequence()),
	}
	if providerEvent.Compact != nil {
		cols = append(cols, handler.NewCol(DebugNotificationProviderCompactCol, *providerEvent.Compact))
	}

	return crdb.NewUpdateStatement(
		&providerEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(DebugNotificationProviderAggIDCol, providerEvent.Aggregate().ID),
			handler.NewCond(DebugNotificationProviderTypeCol, providerType),
		},
	), nil
}

func (p *DebugNotificationProviderProjection) reduceDebugNotificationProviderRemoved(event eventstore.Event) (*handler.Statement, error) {
	var providerEvent settings.DebugNotificationProviderRemovedEvent
	var providerType domain.NotificationProviderType
	switch e := event.(type) {
	case *iam.DebugNotificationProviderFileRemovedEvent:
		providerEvent = e.DebugNotificationProviderRemovedEvent
		providerType = domain.NotificationProviderTypeFile
	case *iam.DebugNotificationProviderLogRemovedEvent:
		providerEvent = e.DebugNotificationProviderRemovedEvent
		providerType = domain.NotificationProviderTypeLog
	default:
		logging.WithFields("seq", event.Sequence(), "expectedTypes", []eventstore.EventType{iam.DebugNotificationProviderFileRemovedEventType, iam.DebugNotificationProviderLogRemovedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-dow9f", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(
		&providerEvent,
		[]handler.Condition{
			handler.NewCond(DebugNotificationProviderAggIDCol, providerEvent.Aggregate().ID),
			handler.NewCond(DebugNotificationProviderTypeCol, providerType),
		},
	), nil
}
