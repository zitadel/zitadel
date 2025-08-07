package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/target"
)

const (
	TargetTable               = "projections.targets2"
	TargetIDCol               = "id"
	TargetCreationDateCol     = "creation_date"
	TargetChangeDateCol       = "change_date"
	TargetResourceOwnerCol    = "resource_owner"
	TargetInstanceIDCol       = "instance_id"
	TargetSequenceCol         = "sequence"
	TargetNameCol             = "name"
	TargetTargetType          = "target_type"
	TargetEndpointCol         = "endpoint"
	TargetTimeoutCol          = "timeout"
	TargetInterruptOnErrorCol = "interrupt_on_error"
	TargetSigningKey          = "signing_key"
)

type targetProjection struct{}

func newTargetProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(targetProjection))
}

func (*targetProjection) Name() string {
	return TargetTable
}

func (*targetProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(TargetIDCol, handler.ColumnTypeText),
			handler.NewColumn(TargetCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(TargetChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(TargetResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(TargetInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(TargetTargetType, handler.ColumnTypeEnum),
			handler.NewColumn(TargetSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(TargetNameCol, handler.ColumnTypeText),
			handler.NewColumn(TargetEndpointCol, handler.ColumnTypeText),
			handler.NewColumn(TargetTimeoutCol, handler.ColumnTypeInt64),
			handler.NewColumn(TargetInterruptOnErrorCol, handler.ColumnTypeBool),
			handler.NewColumn(TargetSigningKey, handler.ColumnTypeJSONB, handler.Nullable()),
		},
			handler.NewPrimaryKey(TargetInstanceIDCol, TargetIDCol),
		),
	)
}

func (p *targetProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: target.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  target.AddedEventType,
					Reduce: p.reduceTargetAdded,
				},
				{
					Event:  target.ChangedEventType,
					Reduce: p.reduceTargetChanged,
				},
				{
					Event:  target.RemovedEventType,
					Reduce: p.reduceTargetRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(TargetInstanceIDCol),
				},
			},
		},
	}
}

func (p *targetProjection) reduceTargetAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*target.AddedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(TargetInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCol(TargetResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(TargetIDCol, e.Aggregate().ID),
			handler.NewCol(TargetCreationDateCol, handler.OnlySetValueOnInsert(TargetTable, e.CreationDate())),
			handler.NewCol(TargetChangeDateCol, e.CreationDate()),
			handler.NewCol(TargetSequenceCol, e.Sequence()),
			handler.NewCol(TargetNameCol, e.Name),
			handler.NewCol(TargetEndpointCol, e.Endpoint),
			handler.NewCol(TargetTargetType, e.TargetType),
			handler.NewCol(TargetTimeoutCol, e.Timeout),
			handler.NewCol(TargetInterruptOnErrorCol, e.InterruptOnError),
			handler.NewCol(TargetSigningKey, e.SigningKey),
		},
	), nil
}

func (p *targetProjection) reduceTargetChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*target.ChangedEvent](event)
	if err != nil {
		return nil, err
	}
	values := []handler.Column{
		handler.NewCol(TargetChangeDateCol, e.CreationDate()),
		handler.NewCol(TargetSequenceCol, e.Sequence()),
		handler.NewCol(TargetResourceOwnerCol, e.Aggregate().ResourceOwner),
	}
	if e.Name != nil {
		values = append(values, handler.NewCol(TargetNameCol, *e.Name))
	}
	if e.TargetType != nil {
		values = append(values, handler.NewCol(TargetTargetType, *e.TargetType))
	}
	if e.Endpoint != nil {
		values = append(values, handler.NewCol(TargetEndpointCol, *e.Endpoint))
	}
	if e.Timeout != nil {
		values = append(values, handler.NewCol(TargetTimeoutCol, *e.Timeout))
	}
	if e.InterruptOnError != nil {
		values = append(values, handler.NewCol(TargetInterruptOnErrorCol, *e.InterruptOnError))
	}
	if e.SigningKey != nil {
		values = append(values, handler.NewCol(TargetSigningKey, e.SigningKey))
	}
	return handler.NewUpdateStatement(
		e,
		values,
		[]handler.Condition{
			handler.NewCond(TargetInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(TargetIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *targetProjection) reduceTargetRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*target.RemovedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(TargetInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(TargetIDCol, e.Aggregate().ID),
		},
	), nil
}
