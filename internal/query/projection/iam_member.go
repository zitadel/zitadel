package projection

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type iamMemberProjection struct {
	crdb.StatementHandler
}

const (
	IAMMemberProjectionTable = "zitadel.projections.iam_members"
)

func newIAMMemberProjection(ctx context.Context, config crdb.StatementHandlerConfig) *iamMemberProjection {
	p := &iamMemberProjection{}
	config.ProjectionName = IAMMemberProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *iamMemberProjection) reducers() []handler.AggregateReducer {
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

func (p *iamMemberProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-c8SBb", "seq", event.Sequence(), "expectedType", iam.MemberAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-pGNCu", "reduce.wrong.event.type")
	}
	return reduceMemberAdded(e.MemberAddedEvent, withMemberCol(IAMMemberIAMIDCol, e.Aggregate().ID))
}

func (p *iamMemberProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-QsjwO", "seq", event.Sequence(), "expected", iam.MemberChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-5WQcZ", "reduce.wrong.event.type")
	}
	return reduceMemberChanged(e.MemberChangedEvent)
}

func (p *iamMemberProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberCascadeRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-mOncs", "seq", event.Sequence(), "expected", iam.MemberCascadeRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Dmdf2", "reduce.wrong.event.type")
	}
	return reduceMemberCascadeRemoved(e.MemberCascadeRemovedEvent)
}

func (p *iamMemberProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-lW1Zv", "seq", event.Sequence(), "expected", iam.MemberRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-exVqy", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.UserID))
}

func (p *iamMemberProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-rBuvT", "seq", event.Sequence(), "expected", user.UserRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-mkDHF", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.Aggregate().ID))
}
