package projection

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type idpUserLinkProjection struct {
	crdb.StatementHandler
}

func newIDPUserLinkProjection(ctx context.Context, config crdb.StatementHandlerConfig) *idpUserLinkProjection {
	p := &idpUserLinkProjection{}
	config.ProjectionName = IDPUserLinkTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *idpUserLinkProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.UserIDPLinkAddedType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  user.UserIDPLinkCascadeRemovedType,
					Reduce: p.reduceCascadeRemoved,
				},
				{
					Event:  user.UserIDPLinkRemovedType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceUserRemoved,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPConfigRemoved,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOrgRemoved,
				},
			},
		},
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPConfigRemoved,
				},
			},
		},
	}
}

const (
	IDPUserLinkTable             = "zitadel.projections.idp_user_links"
	IDPUserLinkIDPIDCol          = "idp_id"
	IDPUserLinkUserIDCol         = "user_id"
	IDPUserLinkExternalUserIDCol = "external_user_id"
	IDPUserLinkCreationDateCol   = "creation_date"
	IDPUserLinkChangeDateCol     = "change_date"
	IDPUserLinkSequenceCol       = "sequence"
	IDPUserLinkResourceOwnerCol  = "resource_owner"
	IDPUserLinkDisplayNameCol    = "display_name"
)

func (p *idpUserLinkProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserIDPLinkAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-v2qC3", "seq", event.Sequence(), "expectedType", user.UserIDPLinkAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-DpmXq", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(e,
		[]handler.Column{
			handler.NewCol(IDPUserLinkIDPIDCol, e.IDPConfigID),
			handler.NewCol(IDPUserLinkUserIDCol, e.Aggregate().ID),
			handler.NewCol(IDPUserLinkExternalUserIDCol, e.ExternalUserID),
			handler.NewCol(IDPUserLinkCreationDateCol, e.CreationDate()),
			handler.NewCol(IDPUserLinkChangeDateCol, e.CreationDate()),
			handler.NewCol(IDPUserLinkSequenceCol, e.Sequence()),
			handler.NewCol(IDPUserLinkResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(IDPUserLinkDisplayNameCol, e.DisplayName),
		},
	), nil
}

func (p *idpUserLinkProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserIDPLinkRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zX5m9", "seq", event.Sequence(), "expectedType", user.UserIDPLinkRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-AZmfJ", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(IDPUserLinkIDPIDCol, e.IDPConfigID),
			handler.NewCond(IDPUserLinkUserIDCol, e.Aggregate().ID),
			handler.NewCond(IDPUserLinkExternalUserIDCol, e.ExternalUserID),
		},
	), nil
}

func (p *idpUserLinkProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserIDPLinkCascadeRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-I0s2H", "seq", event.Sequence(), "expectedType", user.UserIDPLinkCascadeRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-jQpv9", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(IDPUserLinkIDPIDCol, e.IDPConfigID),
			handler.NewCond(IDPUserLinkUserIDCol, e.Aggregate().ID),
			handler.NewCond(IDPUserLinkExternalUserIDCol, e.ExternalUserID),
		},
	), nil
}

func (p *idpUserLinkProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zX5m9", "seq", event.Sequence(), "expectedType", org.OrgRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-AZmfJ", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(IDPUserLinkResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}

func (p *idpUserLinkProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-yM6u6", "seq", event.Sequence(), "expectedType", user.UserRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-uwlWE", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(IDPUserLinkUserIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *idpUserLinkProjection) reduceIDPConfigRemoved(event eventstore.Event) (*handler.Statement, error) {
	var idpID string

	switch e := event.(type) {
	case *org.IDPConfigRemovedEvent:
		idpID = e.ConfigID
	case *iam.IDPConfigRemovedEvent:
		idpID = e.ConfigID
	default:
		logging.LogWithFields("HANDL-7lZaf", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPConfigRemovedEventType, iam.IDPConfigRemovedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-iCKSj", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(event,
		[]handler.Condition{
			handler.NewCond(IDPUserLinkIDPIDCol, idpID),
			handler.NewCond(IDPUserLinkResourceOwnerCol, event.Aggregate().ResourceOwner),
		},
	), nil
}
