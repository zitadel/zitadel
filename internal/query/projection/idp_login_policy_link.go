package projection

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type idpLoginPolicyLinkProjection struct {
	crdb.StatementHandler
}

const (
	IDPLoginPolicyLinkTable = "zitadel.projections.idp_login_policy_links"
)

func newIDPLoginPolicyLinkProjection(ctx context.Context, config crdb.StatementHandlerConfig) *idpLoginPolicyLinkProjection {
	p := &idpLoginPolicyLinkProjection{}
	config.ProjectionName = IDPLoginPolicyLinkTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *idpLoginPolicyLinkProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.LoginPolicyIDPProviderAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.LoginPolicyIDPProviderCascadeRemovedEventType,
					Reduce: p.reduceCascadeRemoved,
				},
				{
					Event:  org.LoginPolicyIDPProviderRemovedEventType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOrgRemoved,
				},
				{
					Event:  org.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPConfigRemoved,
				},
			},
		},
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.LoginPolicyIDPProviderAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  iam.LoginPolicyIDPProviderCascadeRemovedEventType,
					Reduce: p.reduceCascadeRemoved,
				},
				{
					Event:  iam.LoginPolicyIDPProviderRemovedEventType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  iam.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPConfigRemoved,
				},
			},
		},
	}
}

const (
	IDPLoginPolicyLinkIDPIDCol         = "idp_id"
	IDPLoginPolicyLinkAggregateIDCol   = "aggregate_id"
	IDPLoginPolicyLinkCreationDateCol  = "creation_date"
	IDPLoginPolicyLinkChangeDateCol    = "change_date"
	IDPLoginPolicyLinkSequenceCol      = "sequence"
	IDPLoginPolicyLinkResourceOwnerCol = "resource_owner"
	IDPLoginPolicyLinkProviderTypeCol  = "provider_type"
)

func (p *idpLoginPolicyLinkProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var (
		idp          policy.IdentityProviderAddedEvent
		providerType domain.IdentityProviderType
	)

	switch e := event.(type) {
	case *org.IdentityProviderAddedEvent:
		idp = e.IdentityProviderAddedEvent
		providerType = domain.IdentityProviderTypeOrg
	case *iam.IdentityProviderAddedEvent:
		idp = e.IdentityProviderAddedEvent
		providerType = domain.IdentityProviderTypeSystem
	default:
		logging.LogWithFields("HANDL-oce92", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicyIDPProviderAddedEventType, iam.LoginPolicyIDPProviderAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Nlp55", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(&idp,
		[]handler.Column{
			handler.NewCol(IDPLoginPolicyLinkIDPIDCol, idp.IDPConfigID),
			handler.NewCol(IDPLoginPolicyLinkAggregateIDCol, idp.Aggregate().ID),
			handler.NewCol(IDPLoginPolicyLinkCreationDateCol, idp.CreationDate()),
			handler.NewCol(IDPLoginPolicyLinkChangeDateCol, idp.CreationDate()),
			handler.NewCol(IDPLoginPolicyLinkSequenceCol, idp.Sequence()),
			handler.NewCol(IDPLoginPolicyLinkResourceOwnerCol, idp.Aggregate().ResourceOwner),
			handler.NewCol(IDPLoginPolicyLinkProviderTypeCol, providerType),
		},
	), nil
}

func (p *idpLoginPolicyLinkProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	var idp policy.IdentityProviderRemovedEvent

	switch e := event.(type) {
	case *org.IdentityProviderRemovedEvent:
		idp = e.IdentityProviderRemovedEvent
	case *iam.IdentityProviderRemovedEvent:
		idp = e.IdentityProviderRemovedEvent
	default:
		logging.LogWithFields("HANDL-vAH3I", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicyIDPProviderRemovedEventType, iam.LoginPolicyIDPProviderRemovedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-tUMYY", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(&idp,
		[]handler.Condition{
			handler.NewCond(IDPLoginPolicyLinkIDPIDCol, idp.IDPConfigID),
			handler.NewCond(IDPLoginPolicyLinkAggregateIDCol, idp.Aggregate().ID),
		},
	), nil
}

func (p *idpLoginPolicyLinkProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	var idp policy.IdentityProviderCascadeRemovedEvent

	switch e := event.(type) {
	case *org.IdentityProviderCascadeRemovedEvent:
		idp = e.IdentityProviderCascadeRemovedEvent
	case *iam.IdentityProviderCascadeRemovedEvent:
		idp = e.IdentityProviderCascadeRemovedEvent
	default:
		logging.LogWithFields("HANDL-7lZaf", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicyIDPProviderCascadeRemovedEventType, iam.LoginPolicyIDPProviderCascadeRemovedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-iCKSj", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(&idp,
		[]handler.Condition{
			handler.NewCond(IDPLoginPolicyLinkIDPIDCol, idp.IDPConfigID),
			handler.NewCond(IDPLoginPolicyLinkAggregateIDCol, idp.Aggregate().ID),
		},
	), nil
}

func (p *idpLoginPolicyLinkProjection) reduceIDPConfigRemoved(event eventstore.Event) (*handler.Statement, error) {
	var idpID string

	switch e := event.(type) {
	case *org.IDPConfigRemovedEvent:
		idpID = e.ConfigID
	case *iam.IDPConfigRemovedEvent:
		idpID = e.ConfigID
	default:
		logging.LogWithFields("HANDL-aJvob", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPConfigRemovedEventType, iam.IDPConfigRemovedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-u6tze", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(event,
		[]handler.Condition{
			handler.NewCond(IDPLoginPolicyLinkIDPIDCol, idpID),
			handler.NewCond(IDPLoginPolicyLinkResourceOwnerCol, event.Aggregate().ResourceOwner),
		},
	), nil
}

func (p *idpLoginPolicyLinkProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-WTYC1", "seq", event.Sequence(), "expectedType", org.OrgRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-QSoSe", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(IDPLoginPolicyLinkResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
