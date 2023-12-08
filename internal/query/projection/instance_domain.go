package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	InstanceDomainTable = "projections.instance_domains"

	InstanceDomainInstanceIDCol   = "instance_id"
	InstanceDomainCreationDateCol = "creation_date"
	InstanceDomainChangeDateCol   = "change_date"
	InstanceDomainSequenceCol     = "sequence"
	InstanceDomainDomainCol       = "domain"
	InstanceDomainIsGeneratedCol  = "is_generated"
	InstanceDomainIsPrimaryCol    = "is_primary"
)

type instanceDomainProjection struct{}

func newInstanceDomainProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(instanceDomainProjection))
}

func (*instanceDomainProjection) Name() string {
	return InstanceDomainTable
}

func (*instanceDomainProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(InstanceDomainInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(InstanceDomainCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(InstanceDomainChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(InstanceDomainSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(InstanceDomainDomainCol, handler.ColumnTypeText),
			handler.NewColumn(InstanceDomainIsGeneratedCol, handler.ColumnTypeBool),
			handler.NewColumn(InstanceDomainIsPrimaryCol, handler.ColumnTypeBool),
		},
			handler.NewPrimaryKey(InstanceDomainInstanceIDCol, InstanceDomainDomainCol),
			handler.WithIndex(
				handler.NewIndex("instance_domain", []string{InstanceDomainDomainCol},
					handler.WithInclude(InstanceDomainCreationDateCol, InstanceDomainChangeDateCol, InstanceDomainSequenceCol, InstanceDomainIsGeneratedCol, InstanceDomainIsPrimaryCol),
				),
			),
		),
	)
}

func (p *instanceDomainProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceDomainAddedEventType,
					Reduce: p.reduceDomainAdded,
				},
				{
					Event:  instance.InstanceDomainPrimarySetEventType,
					Reduce: p.reduceDomainPrimarySet,
				},
				{
					Event:  instance.InstanceDomainRemovedEventType,
					Reduce: p.reduceDomainRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(InstanceDomainInstanceIDCol),
				},
			},
		},
	}
}

func (p *instanceDomainProjection) reduceDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-38nNf", "reduce.wrong.event.type %s", instance.InstanceDomainAddedEventType)
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(InstanceDomainCreationDateCol, e.CreationDate()),
			handler.NewCol(InstanceDomainChangeDateCol, e.CreationDate()),
			handler.NewCol(InstanceDomainSequenceCol, e.Sequence()),
			handler.NewCol(InstanceDomainDomainCol, e.Domain),
			handler.NewCol(InstanceDomainInstanceIDCol, e.Aggregate().ID),
			handler.NewCol(InstanceDomainIsGeneratedCol, e.Generated),
			handler.NewCol(InstanceDomainIsPrimaryCol, false),
		},
	), nil
}

func (p *instanceDomainProjection) reduceDomainPrimarySet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainPrimarySetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-f8nlw", "reduce.wrong.event.type %s", instance.InstanceDomainPrimarySetEventType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(InstanceDomainChangeDateCol, e.CreationDate()),
				handler.NewCol(InstanceDomainSequenceCol, e.Sequence()),
				handler.NewCol(InstanceDomainIsPrimaryCol, false),
			},
			[]handler.Condition{
				handler.NewCond(InstanceDomainInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCond(InstanceDomainIsPrimaryCol, true),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(InstanceDomainChangeDateCol, e.CreationDate()),
				handler.NewCol(InstanceDomainSequenceCol, e.Sequence()),
				handler.NewCol(InstanceDomainIsPrimaryCol, true),
			},
			[]handler.Condition{
				handler.NewCond(InstanceDomainDomainCol, e.Domain),
				handler.NewCond(InstanceDomainInstanceIDCol, e.Aggregate().ID),
			},
		),
	), nil
}

func (p *instanceDomainProjection) reduceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-388Nk", "reduce.wrong.event.type %s", instance.InstanceDomainRemovedEventType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(InstanceDomainDomainCol, e.Domain),
			handler.NewCond(InstanceDomainInstanceIDCol, e.Aggregate().ID),
		},
	), nil
}
