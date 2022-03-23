package projection

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/user"
)

const (
	IAMMemberProjectionTable = "projections.iam_members"

	IAMMemberIAMIDCol = "iam_id"
)

type IAMMemberProjection struct {
	crdb.StatementHandler
}

func NewIAMMemberProjection(ctx context.Context, config crdb.StatementHandlerConfig) *IAMMemberProjection {
	p := new(IAMMemberProjection)
	config.ProjectionName = IAMMemberProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable(
			append(memberColumns, crdb.NewColumn(IAMColumnID, crdb.ColumnTypeText)),
			crdb.NewPrimaryKey(MemberInstanceID, IAMColumnID, MemberUserIDCol),
			crdb.NewIndex("user_idx", []string{MemberUserIDCol}),
		),
	)

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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-pGNCu", "reduce.wrong.event.type %s", iam.MemberAddedEventType)
	}
	return reduceMemberAdded(e.MemberAddedEvent, withMemberCol(IAMMemberIAMIDCol, e.Aggregate().ID))
}

func (p *IAMMemberProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-5WQcZ", "reduce.wrong.event.type %s", iam.MemberChangedEventType)
	}
	return reduceMemberChanged(e.MemberChangedEvent)
}

func (p *IAMMemberProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberCascadeRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Dmdf2", "reduce.wrong.event.type %s", iam.MemberCascadeRemovedEventType)
	}
	return reduceMemberCascadeRemoved(e.MemberCascadeRemovedEvent)
}

func (p *IAMMemberProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-exVqy", "reduce.wrong.event.type %s", iam.MemberRemovedEventType)
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.UserID))
}

func (p *IAMMemberProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-mkDHF", "reduce.wrong.event.type %s", user.UserRemovedType)
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.Aggregate().ID))
}
