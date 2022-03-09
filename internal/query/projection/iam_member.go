package projection

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/user"
)

const (
	IAMMemberProjectionTable = "zitadel.projections.iam_members"

	IAMMemberIAMIDCol = "iam_id"
)

type IAMMemberProjection struct {
	crdb.StatementHandler
}

func NewIAMMemberProjection(ctx context.Context, config crdb.StatementHandlerConfig) *IAMMemberProjection {
	p := new(IAMMemberProjection)
	config.ProjectionName = IAMMemberProjectionTable
	config.Reducers = p.reducers()
	config.InitChecks = []*handler.Check{
		crdb.NewTableCheck(
			crdb.NewTable(
				append(memberColumns, crdb.NewColumn(IAMColumnID, crdb.ColumnTypeEnum)),
				crdb.NewPrimaryKey(IAMColumnID, MemberUserIDCol),
			),
		),
		crdb.NewIndexCheck(crdb.NewIndex("user_idx", []string{MemberUserIDCol})),
	}
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

func (p *IAMMemberProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-c8SBb", "seq", event.Sequence(), "expectedType", iam.MemberAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-pGNCu", "reduce.wrong.event.type")
	}
	return reduceMemberAdded(e.MemberAddedEvent, withMemberCol(IAMMemberIAMIDCol, e.Aggregate().ID))
}

func (p *IAMMemberProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-QsjwO", "seq", event.Sequence(), "expected", iam.MemberChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-5WQcZ", "reduce.wrong.event.type")
	}
	return reduceMemberChanged(e.MemberChangedEvent)
}

func (p *IAMMemberProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberCascadeRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-mOncs", "seq", event.Sequence(), "expected", iam.MemberCascadeRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Dmdf2", "reduce.wrong.event.type")
	}
	return reduceMemberCascadeRemoved(e.MemberCascadeRemovedEvent)
}

func (p *IAMMemberProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-lW1Zv", "seq", event.Sequence(), "expected", iam.MemberRemovedEventType).Error("wrong event type")
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
