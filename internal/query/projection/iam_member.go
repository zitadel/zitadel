package projection

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMMemberProjection struct {
	crdb.StatementHandler
}

const (
	IAMMemberProjectionTable = "zitadel.projections.iam_members"
)

func NewIAMMemberProjection(ctx context.Context, config crdb.StatementHandlerConfig) *IAMMemberProjection {
	p := &IAMMemberProjection{}
	config.ProjectionName = IAMMemberProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *IAMMemberProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.MemberAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  iam.MemberChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  iam.MemberCascadeRemovedEventType,
					Reduce: p.reduceCascadeRemoved,
				},
				{
					Event:  iam.MemberRemovedEventType,
					Reduce: p.reduceRemoved,
				},
			},
		},
	}
}

type IAMMemberColumn string

const (
	IAMMemberIAMIDCol = "iam_id"
)

func (p *IAMMemberProjection) reduceAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-c8SBb", "seq", event.Sequence(), "expectedType", iam.MemberAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-pGNCu", "reduce.wrong.event.type")
	}
	return reduceMemberAdded(e.MemberAddedEvent, withMemberCol(IAMMemberIAMIDCol, e.Aggregate().ID))
}

func (p *IAMMemberProjection) reduceChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-QsjwO", "seq", event.Sequence(), "expected", iam.MemberChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-5WQcZ", "reduce.wrong.event.type")
	}
	return reduceMemberChanged(e.MemberChangedEvent, withMemberCond(IAMMemberIAMIDCol, e.Aggregate().ID))
}

func (p *IAMMemberProjection) reduceCascadeRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberCascadeRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-mOncs", "seq", event.Sequence(), "expected", iam.MemberCascadeRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Dmdf2", "reduce.wrong.event.type")
	}
	return reduceMemberCascadeRemoved(e.MemberCascadeRemovedEvent, withMemberCond(IAMMemberIAMIDCol, e.Aggregate().ID))
}

func (p *IAMMemberProjection) reduceRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-lW1Zv", "seq", event.Sequence(), "expected", iam.MemberRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-exVqy", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e.MemberRemovedEvent, withMemberCond(IAMMemberIAMIDCol, e.Aggregate().ID))
}
