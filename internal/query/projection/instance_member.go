package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	InstanceMemberProjectionTable = "projections.instance_members4"

	InstanceMemberIAMIDCol = "id"
)

type instanceMemberProjection struct {
	es handler.EventStore
}

func newInstanceMemberProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, &instanceMemberProjection{es: config.Eventstore})
}

func (*instanceMemberProjection) Name() string {
	return InstanceMemberProjectionTable
}

func (*instanceMemberProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable(
			append(memberColumns, handler.NewColumn(InstanceColumnID, handler.ColumnTypeText)),
			handler.NewPrimaryKey(MemberInstanceID, InstanceColumnID, MemberUserIDCol),
			handler.WithIndex(handler.NewIndex("user_id", []string{MemberUserIDCol})),
			handler.WithIndex(
				handler.NewIndex("im_instance", []string{MemberInstanceID},
					handler.WithInclude(
						MemberCreationDate,
						MemberChangeDate,
						MemberRolesCol,
						MemberSequence,
						MemberResourceOwner,
					),
				),
			),
		),
	)
}

func (p *instanceMemberProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
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
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(AppColumnInstanceID),
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceUserOwnerRemoved,
				},
			},
		},
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceUserRemoved,
				},
			},
		},
	}
}

func (p *instanceMemberProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.MemberAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-pGNCu", "reduce.wrong.event.type %s", instance.MemberAddedEventType)
	}
	ctx := setMemberContext(e.Aggregate())
	userOwner, err := getResourceOwnerOfUser(ctx, p.es, e.Aggregate().InstanceID, e.UserID)
	if err != nil {
		return nil, err
	}
	return reduceMemberAdded(e.MemberAddedEvent, userOwner, withMemberCol(InstanceMemberIAMIDCol, e.Aggregate().ID))
}

func (p *instanceMemberProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.MemberChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-5WQcZ", "reduce.wrong.event.type %s", instance.MemberChangedEventType)
	}
	return reduceMemberChanged(e.MemberChangedEvent)
}

func (p *instanceMemberProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.MemberCascadeRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Dmdf2", "reduce.wrong.event.type %s", instance.MemberCascadeRemovedEventType)
	}
	return reduceMemberCascadeRemoved(e.MemberCascadeRemovedEvent)
}

func (p *instanceMemberProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.MemberRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-exVqy", "reduce.wrong.event.type %s", instance.MemberRemovedEventType)
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.UserID))
}

func (p *instanceMemberProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-mkDHF", "reduce.wrong.event.type %s", user.UserRemovedType)
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.Aggregate().ID))
}

func (p *instanceMemberProjection) reduceUserOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-mkDHa", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}
	return reduceMemberUserOwnerRemoved(e)
}
