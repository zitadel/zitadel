package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

const (
	InstanceAllowedDomainTable = "projections.instance_allowed_domains"

	InstanceAllowedDomainInstanceIDCol   = "instance_id"
	InstanceAllowedDomainCreationDateCol = "creation_date"
	InstanceAllowedDomainChangeDateCol   = "change_date"
	InstanceAllowedDomainSequenceCol     = "sequence"
	InstanceAllowedDomainDomainCol       = "domain"
)

type instanceAllowedDomainProjection struct{}

func newInstanceAllowedDomainProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(instanceAllowedDomainProjection))
}

func (*instanceAllowedDomainProjection) Name() string {
	return InstanceAllowedDomainTable
}

func (*instanceAllowedDomainProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(InstanceAllowedDomainInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(InstanceAllowedDomainCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(InstanceAllowedDomainChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(InstanceAllowedDomainSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(InstanceAllowedDomainDomainCol, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(InstanceAllowedDomainInstanceIDCol, InstanceAllowedDomainDomainCol),
			handler.WithIndex(
				handler.NewIndex("instance_allowed_domain", []string{InstanceAllowedDomainDomainCol},
					handler.WithInclude(InstanceAllowedDomainCreationDateCol, InstanceAllowedDomainChangeDateCol, InstanceAllowedDomainSequenceCol),
				),
			),
		),
	)
}

func (p *instanceAllowedDomainProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.AllowedDomainAddedEventType,
					Reduce: p.reduceDomainAdded,
				},
				{
					Event:  instance.AllowedDomainRemovedEventType,
					Reduce: p.reduceDomainRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(InstanceAllowedDomainInstanceIDCol),
				},
			},
		},
	}
}

func (p *instanceAllowedDomainProjection) reduceDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.AllowedDomainAddedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(InstanceAllowedDomainCreationDateCol, e.CreatedAt()),
			handler.NewCol(InstanceAllowedDomainChangeDateCol, e.CreatedAt()),
			handler.NewCol(InstanceAllowedDomainSequenceCol, e.Sequence()),
			handler.NewCol(InstanceAllowedDomainDomainCol, e.Domain),
			handler.NewCol(InstanceAllowedDomainInstanceIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *instanceAllowedDomainProjection) reduceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.AllowedDomainRemovedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(InstanceAllowedDomainDomainCol, e.Domain),
			handler.NewCond(InstanceAllowedDomainInstanceIDCol, e.Aggregate().ID),
		},
	), nil
}
