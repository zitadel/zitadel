package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

const (
	InstanceTrustedDomainTable = "projections.instance_trusted_domains"

	InstanceTrustedDomainInstanceIDCol   = "instance_id"
	InstanceTrustedDomainCreationDateCol = "creation_date"
	InstanceTrustedDomainChangeDateCol   = "change_date"
	InstanceTrustedDomainSequenceCol     = "sequence"
	InstanceTrustedDomainDomainCol       = "domain"
)

type instanceTrustedDomainProjection struct{}

func newInstanceTrustedDomainProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(instanceTrustedDomainProjection))
}

func (*instanceTrustedDomainProjection) Name() string {
	return InstanceTrustedDomainTable
}

func (*instanceTrustedDomainProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(InstanceTrustedDomainInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(InstanceTrustedDomainCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(InstanceTrustedDomainChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(InstanceTrustedDomainSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(InstanceTrustedDomainDomainCol, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(InstanceTrustedDomainInstanceIDCol, InstanceTrustedDomainDomainCol),
			handler.WithIndex(
				handler.NewIndex("instance_trusted_domain", []string{InstanceTrustedDomainDomainCol},
					handler.WithInclude(InstanceTrustedDomainCreationDateCol, InstanceTrustedDomainChangeDateCol, InstanceTrustedDomainSequenceCol),
				),
			),
		),
	)
}

func (p *instanceTrustedDomainProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.TrustedDomainAddedEventType,
					Reduce: p.reduceDomainAdded,
				},
				{
					Event:  instance.TrustedDomainRemovedEventType,
					Reduce: p.reduceDomainRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(InstanceTrustedDomainInstanceIDCol),
				},
			},
		},
	}
}

func (p *instanceTrustedDomainProjection) reduceDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.TrustedDomainAddedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(InstanceTrustedDomainCreationDateCol, e.CreatedAt()),
			handler.NewCol(InstanceTrustedDomainChangeDateCol, e.CreatedAt()),
			handler.NewCol(InstanceTrustedDomainSequenceCol, e.Sequence()),
			handler.NewCol(InstanceTrustedDomainDomainCol, e.Domain),
			handler.NewCol(InstanceTrustedDomainInstanceIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *instanceTrustedDomainProjection) reduceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.TrustedDomainRemovedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(InstanceTrustedDomainDomainCol, e.Domain),
			handler.NewCond(InstanceTrustedDomainInstanceIDCol, e.Aggregate().ID),
		},
	), nil
}
