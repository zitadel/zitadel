package projection

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/user"
)

type OrgMemberProjection struct {
	crdb.StatementHandler
}

const (
	OrgMemberProjectionTable = "zitadel.projections.org_members"
)

func NewOrgMemberProjection(ctx context.Context, config crdb.StatementHandlerConfig) *OrgMemberProjection {
	p := &OrgMemberProjection{}
	config.ProjectionName = OrgMemberProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *OrgMemberProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
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
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceUserRemoved,
				},
			},
		},
	}
}

type OrgMemberColumn string

const (
	OrgMemberOrgIDCol = "org_id"
)

func (p *OrgMemberProjection) reduceAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.MemberAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-BoKBr", "seq", event.Sequence(), "expectedType", org.MemberAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-uYq4r", "reduce.wrong.event.type")
	}
	return reduceMemberAdded(e.MemberAddedEvent, withMemberCol(OrgMemberOrgIDCol, e.Aggregate().ID))
}

func (p *OrgMemberProjection) reduceChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.MemberChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-bfqNl", "seq", event.Sequence(), "expected", org.MemberChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Bg8oM", "reduce.wrong.event.type")
	}
	return reduceMemberChanged(e.MemberChangedEvent, withMemberCond(OrgMemberOrgIDCol, e.Aggregate().ID))
}

func (p *OrgMemberProjection) reduceCascadeRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.MemberCascadeRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zgb6w", "seq", event.Sequence(), "expected", org.MemberCascadeRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-4twP2", "reduce.wrong.event.type")
	}
	return reduceMemberCascadeRemoved(e.MemberCascadeRemovedEvent, withMemberCond(OrgMemberOrgIDCol, e.Aggregate().ID))
}

func (p *OrgMemberProjection) reduceRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.MemberRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-KPyxE", "seq", event.Sequence(), "expected", org.MemberRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-avatH", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e,
		withMemberCond(MemberUserIDCol, e.UserID),
		withMemberCond(OrgMemberOrgIDCol, e.Aggregate().ID),
	)
}

func (p *OrgMemberProjection) reduceUserRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-f5pgn", "seq", event.Sequence(), "expected", user.UserRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-eBMqH", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.Aggregate().ID))
}

func (p *OrgMemberProjection) reduceOrgRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	//TODO: as soon as org deletion is implemented:
	// Case: The user has resource owner A and an org has resource owner B
	// if org B deleted it works
	// if org A is deleted, the membership wouldn't be deleted
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-E5lDs", "seq", event.Sequence(), "expected", org.OrgRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-jnGAV", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e, withMemberCond(OrgMemberOrgIDCol, e.Aggregate().ID))
}
