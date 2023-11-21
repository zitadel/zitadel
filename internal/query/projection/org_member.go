package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	OrgMemberProjectionTable = "projections.org_members4"
	OrgMemberOrgIDCol        = "org_id"
)

type orgMemberProjection struct {
	es handler.EventStore
}

func newOrgMemberProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, &orgMemberProjection{es: config.Eventstore})
}

func (*orgMemberProjection) Name() string {
	return OrgMemberProjectionTable
}

func (*orgMemberProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable(
			append(memberColumns, handler.NewColumn(OrgMemberOrgIDCol, handler.ColumnTypeText)),
			handler.NewPrimaryKey(MemberInstanceID, OrgMemberOrgIDCol, MemberUserIDCol),
			handler.WithIndex(handler.NewIndex("user_id", []string{MemberUserIDCol})),
			handler.WithIndex(
				handler.NewIndex("om_instance", []string{MemberInstanceID},
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

func (p *orgMemberProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.MemberAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.MemberChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  org.MemberCascadeRemovedEventType,
					Reduce: p.reduceCascadeRemoved,
				},
				{
					Event:  org.MemberRemovedEventType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOrgRemoved,
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
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(MemberInstanceID),
				},
			},
		},
	}
}

func (p *orgMemberProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.MemberAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-uYq4r", "reduce.wrong.event.type %s", org.MemberAddedEventType)
	}
	ctx := setMemberContext(e.Aggregate())
	userOwner, err := getResourceOwnerOfUser(ctx, p.es, e.Aggregate().InstanceID, e.UserID)
	if err != nil {
		return nil, err
	}
	return reduceMemberAdded(e.MemberAddedEvent, userOwner, withMemberCol(OrgMemberOrgIDCol, e.Aggregate().ID))
}

func (p *orgMemberProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.MemberChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Bg8oM", "reduce.wrong.event.type %s", org.MemberChangedEventType)
	}
	return reduceMemberChanged(e.MemberChangedEvent, withMemberCond(OrgMemberOrgIDCol, e.Aggregate().ID))
}

func (p *orgMemberProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.MemberCascadeRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-4twP2", "reduce.wrong.event.type %s", org.MemberCascadeRemovedEventType)
	}
	return reduceMemberCascadeRemoved(e.MemberCascadeRemovedEvent, withMemberCond(OrgMemberOrgIDCol, e.Aggregate().ID))
}

func (p *orgMemberProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.MemberRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-avatH", "reduce.wrong.event.type %s", org.MemberRemovedEventType)
	}
	return reduceMemberRemoved(e,
		withMemberCond(MemberUserIDCol, e.UserID),
		withMemberCond(OrgMemberOrgIDCol, e.Aggregate().ID),
	)
}

func (p *orgMemberProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-eBMqH", "reduce.wrong.event.type %s", user.UserRemovedType)
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.Aggregate().ID))
}

func (p *orgMemberProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-jnGAV", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}
	return handler.NewMultiStatement(
		e,
		multiReduceMemberOwnerRemoved(e),
		multiReduceMemberUserOwnerRemoved(e),
	), nil
}
