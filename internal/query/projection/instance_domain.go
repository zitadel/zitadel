package projection

import (
	"context"
	"fmt"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
)

type InstanceDomainProjection struct {
	crdb.StatementHandler
}

const (
	InstanceDomainTable = "zitadel.projections.instance_domains"
)

func NewInstanceDomainProjection(ctx context.Context, config crdb.StatementHandlerConfig) *InstanceDomainProjection {
	p := new(InstanceDomainProjection)
	config.ProjectionName = InstanceDomainTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *InstanceDomainProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.InstanceDomainAddedEventType,
					Reduce: p.reduceDomainAdded,
				},
				{
					Event:  iam.InstanceDomainRemovedEventType,
					Reduce: p.reduceDomainRemoved,
				},
			},
		},
	}
}

const (
	InstanceDomainCreationDateCol = "creation_date"
	InstanceDomainChangeDateCol   = "change_date"
	InstanceDomainSequenceCol     = "sequence"
	InstanceDomainDomainCol       = "domain"
	InstanceDomainInstanceIDCol   = "instance_id"
	InstanceDomainIsGeneratedCol  = "is_generated"
)

func (p *InstanceDomainProjection) reduceDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainAddedEvent)
	if !ok {
		logging.WithFields("seq", event.Sequence(), "expectedType", iam.InstanceDomainAddedEventType, "gottenType", fmt.Sprintf("%T", event)).Error("unexpected event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-DM2DI", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(InstanceDomainCreationDateCol, e.CreationDate()),
			handler.NewCol(InstanceDomainChangeDateCol, e.CreationDate()),
			handler.NewCol(InstanceDomainSequenceCol, e.Sequence()),
			handler.NewCol(InstanceDomainDomainCol, e.Domain),
			handler.NewCol(InstanceDomainInstanceIDCol, e.Aggregate().ID),
			handler.NewCol(InstanceDomainIsGeneratedCol, false),
		},
	), nil
}
func (p *InstanceDomainProjection) reduceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainRemovedEvent)
	if !ok {
		logging.WithFields("seq", event.Sequence(), "expectedType", iam.InstanceDomainRemovedEventType, "gottenType", fmt.Sprintf("%T", event)).Error("unexpected event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-gh1Mx", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(InstanceDomainDomainCol, e.Domain),
			handler.NewCond(InstanceDomainInstanceIDCol, e.Aggregate().ID),
		},
	), nil
}
