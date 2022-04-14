package projection

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/user"
)

const (
	InstanceMemberProjectionTable = "projections.instance_members"

	InstanceMemberIAMIDCol = "id"
)

type InstanceMemberProjection struct {
	crdb.StatementHandler
}

func NewInstanceMemberProjection(ctx context.Context, config crdb.StatementHandlerConfig) *InstanceMemberProjection {
	p := new(InstanceMemberProjection)
	config.ProjectionName = InstanceMemberProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable(
			append(memberColumns, crdb.NewColumn(InstanceColumnID, crdb.ColumnTypeText)),
			crdb.NewPrimaryKey(MemberInstanceID, InstanceColumnID, MemberUserIDCol),
			crdb.WithIndex(crdb.NewIndex("user_idx", []string{MemberUserIDCol})),
		),
	)

	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *InstanceMemberProjection) reducers() []handler.AggregateReducer {
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

func (p *InstanceMemberProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.MemberAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-pGNCu", "reduce.wrong.event.type %s", instance.MemberAddedEventType)
	}
	return reduceMemberAdded(e.MemberAddedEvent, withMemberCol(InstanceMemberIAMIDCol, e.Aggregate().ID))
}

func (p *InstanceMemberProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.MemberChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-5WQcZ", "reduce.wrong.event.type %s", instance.MemberChangedEventType)
	}
	return reduceMemberChanged(e.MemberChangedEvent)
}

func (p *InstanceMemberProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.MemberCascadeRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Dmdf2", "reduce.wrong.event.type %s", instance.MemberCascadeRemovedEventType)
	}
	return reduceMemberCascadeRemoved(e.MemberCascadeRemovedEvent)
}

func (p *InstanceMemberProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.MemberRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-exVqy", "reduce.wrong.event.type %s", instance.MemberRemovedEventType)
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.UserID))
}

func (p *InstanceMemberProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-mkDHF", "reduce.wrong.event.type %s", user.UserRemovedType)
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.Aggregate().ID))
}
