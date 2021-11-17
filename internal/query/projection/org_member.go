package projection

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/org"
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
					Reduce: p.reduceOrgAdded,
				},
				{
					Event:  org.MemberChangedEventType,
					Reduce: p.reduceOrgChanged,
				},
				{
					Event:  org.MemberCascadeRemovedEventType,
					Reduce: p.reduceOrgMemberCascadeRemoved,
				},
				{
					Event:  org.MemberRemovedEventType,
					Reduce: p.reduceOrgMemberRemoved,
				},
			},
		},
	}
}

type OrgMemberColumn string

const (
	OrgMemberOrgIDCol  = "org_id"
	OrgMemberUserIDCol = "user_id"
	OrgMemberRolesCol  = "roles"
)

func (p *OrgMemberProjection) reduceOrgAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.MemberAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-BoKBr", "seq", event.Sequence(), "expectedType", org.MemberAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-uYq4r", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgMemberOrgIDCol, e.Aggregate().ResourceOwner),
			handler.NewCol(OrgMemberUserIDCol, e.UserID),
			handler.NewCol(OrgMemberRolesCol, e.Roles),
		},
	), nil
}

func (p *OrgMemberProjection) reduceOrgChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.MemberChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-bfqNl", "seq", event.Sequence(), "expected", org.MemberChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Bg8oM", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgMemberRolesCol, e.Roles),
		},
		[]handler.Condition{
			handler.NewCond(OrgMemberOrgIDCol, e.Aggregate().ResourceOwner),
			handler.NewCond(OrgMemberUserIDCol, e.UserID),
		},
	), nil
}

func (p *OrgMemberProjection) reduceOrgMemberCascadeRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.MemberCascadeRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zgb6w", "seq", event.Sequence(), "expected", org.MemberCascadeRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-4twP2", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(OrgMemberOrgIDCol, e.Aggregate().ResourceOwner),
			handler.NewCond(OrgMemberUserIDCol, e.UserID),
		},
	), nil
}

func (p *OrgMemberProjection) reduceOrgMemberRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.MemberRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-KPyxE", "seq", event.Sequence(), "expected", org.MemberRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-avatH", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(OrgMemberOrgIDCol, e.Aggregate().ResourceOwner),
			handler.NewCond(OrgMemberUserIDCol, e.UserID),
		},
	), nil
}
