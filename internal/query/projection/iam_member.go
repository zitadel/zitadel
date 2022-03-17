package projection

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/user"
)

type IAMMemberProjection struct {
	crdb.StatementHandler
}

const (
	IAMMemberProjectionTable = "zitadel.projections.iam_members"
)

func NewIAMMemberProjection(ctx context.Context, config crdb.StatementHandlerConfig) *IAMMemberProjection {
	p := new(IAMMemberProjection)
	config.ProjectionName = IAMMemberProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *IAMMemberProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.MemberAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  instance.MemberChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  instance.MemberCascadeRemovedEventType,
					Reduce: p.reduceCascadeRemoved,
				},
				{
					Event:  instance.MemberRemovedEventType,
					Reduce: p.reduceRemoved,
				},
			},
		},
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceUserRemoved,
				},
			},
		},
	}
}

type IAMMemberColumn string

const (
	IAMMemberIAMIDCol = "iam_id"
)

func (p *IAMMemberProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.MemberAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-c8SBb", "seq", event.Sequence(), "expectedType", instance.MemberAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-pGNCu", "reduce.wrong.event.type")
	}
	return reduceMemberAdded(e.MemberAddedEvent, withMemberCol(IAMMemberIAMIDCol, e.Aggregate().ID))
}

func (p *IAMMemberProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.MemberChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-QsjwO", "seq", event.Sequence(), "expected", instance.MemberChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-5WQcZ", "reduce.wrong.event.type")
	}
	return reduceMemberChanged(e.MemberChangedEvent)
}

func (p *IAMMemberProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.MemberCascadeRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-mOncs", "seq", event.Sequence(), "expected", instance.MemberCascadeRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Dmdf2", "reduce.wrong.event.type")
	}
	return reduceMemberCascadeRemoved(e.MemberCascadeRemovedEvent)
}

func (p *IAMMemberProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.MemberRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-lW1Zv", "seq", event.Sequence(), "expected", instance.MemberRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-exVqy", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.UserID))
}

func (p *IAMMemberProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-rBuvT", "seq", event.Sequence(), "expected", user.UserRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-mkDHF", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.Aggregate().ID))
}
