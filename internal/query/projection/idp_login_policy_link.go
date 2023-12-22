package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	IDPLoginPolicyLinkTable = "projections.idp_login_policy_links5"

	IDPLoginPolicyLinkIDPIDCol         = "idp_id"
	IDPLoginPolicyLinkAggregateIDCol   = "aggregate_id"
	IDPLoginPolicyLinkCreationDateCol  = "creation_date"
	IDPLoginPolicyLinkChangeDateCol    = "change_date"
	IDPLoginPolicyLinkSequenceCol      = "sequence"
	IDPLoginPolicyLinkResourceOwnerCol = "resource_owner"
	IDPLoginPolicyLinkInstanceIDCol    = "instance_id"
	IDPLoginPolicyLinkProviderTypeCol  = "provider_type"
	IDPLoginPolicyLinkOwnerRemovedCol  = "owner_removed"
)

type idpLoginPolicyLinkProjection struct{}

func newIDPLoginPolicyLinkProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(idpLoginPolicyLinkProjection))
}

func (*idpLoginPolicyLinkProjection) Name() string {
	return IDPLoginPolicyLinkTable
}

func (*idpLoginPolicyLinkProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(IDPLoginPolicyLinkIDPIDCol, handler.ColumnTypeText),
			handler.NewColumn(IDPLoginPolicyLinkAggregateIDCol, handler.ColumnTypeText),
			handler.NewColumn(IDPLoginPolicyLinkCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(IDPLoginPolicyLinkChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(IDPLoginPolicyLinkSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(IDPLoginPolicyLinkResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(IDPLoginPolicyLinkInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(IDPLoginPolicyLinkProviderTypeCol, handler.ColumnTypeEnum),
			handler.NewColumn(IDPLoginPolicyLinkOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(IDPLoginPolicyLinkInstanceIDCol, IDPLoginPolicyLinkAggregateIDCol, IDPLoginPolicyLinkIDPIDCol),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{IDPLoginPolicyLinkResourceOwnerCol})),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{IDPLoginPolicyLinkOwnerRemovedCol})),
		),
	)
}

func (p *idpLoginPolicyLinkProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
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
					Event:  org.LoginPolicyRemovedEventType,
					Reduce: p.reducePolicyRemoved,
				},
				{
					Event:  org.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPConfigRemoved,
				},
				{
					Event:  org.IDPRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.LoginPolicyIDPProviderAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  instance.LoginPolicyIDPProviderCascadeRemovedEventType,
					Reduce: p.reduceCascadeRemoved,
				},
				{
					Event:  instance.LoginPolicyIDPProviderRemovedEventType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  instance.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPConfigRemoved,
				},
				{
					Event:  instance.IDPRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(IDPUserLinkInstanceIDCol),
				},
			},
		},
	}
}

func (p *idpLoginPolicyLinkProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var (
		idp          policy.IdentityProviderAddedEvent
		providerType domain.IdentityProviderType
	)

	switch e := event.(type) {
	case *org.IdentityProviderAddedEvent:
		idp = e.IdentityProviderAddedEvent
		providerType = domain.IdentityProviderTypeOrg
	case *instance.IdentityProviderAddedEvent:
		idp = e.IdentityProviderAddedEvent
		providerType = domain.IdentityProviderTypeSystem
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Nlp55", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyIDPProviderAddedEventType, instance.LoginPolicyIDPProviderAddedEventType})
	}

	return handler.NewCreateStatement(&idp,
		[]handler.Column{
			handler.NewCol(IDPLoginPolicyLinkIDPIDCol, idp.IDPConfigID),
			handler.NewCol(IDPLoginPolicyLinkAggregateIDCol, idp.Aggregate().ID),
			handler.NewCol(IDPLoginPolicyLinkCreationDateCol, idp.CreationDate()),
			handler.NewCol(IDPLoginPolicyLinkChangeDateCol, idp.CreationDate()),
			handler.NewCol(IDPLoginPolicyLinkSequenceCol, idp.Sequence()),
			handler.NewCol(IDPLoginPolicyLinkResourceOwnerCol, idp.Aggregate().ResourceOwner),
			handler.NewCol(IDPLoginPolicyLinkInstanceIDCol, idp.Aggregate().InstanceID),
			handler.NewCol(IDPLoginPolicyLinkProviderTypeCol, providerType),
		},
	), nil
}

func (p *idpLoginPolicyLinkProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	var idp policy.IdentityProviderRemovedEvent

	switch e := event.(type) {
	case *org.IdentityProviderRemovedEvent:
		idp = e.IdentityProviderRemovedEvent
	case *instance.IdentityProviderRemovedEvent:
		idp = e.IdentityProviderRemovedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-tUMYY", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyIDPProviderRemovedEventType, instance.LoginPolicyIDPProviderRemovedEventType})
	}

	return handler.NewDeleteStatement(&idp,
		[]handler.Condition{
			handler.NewCond(IDPLoginPolicyLinkIDPIDCol, idp.IDPConfigID),
			handler.NewCond(IDPLoginPolicyLinkAggregateIDCol, idp.Aggregate().ID),
			handler.NewCond(IDPLoginPolicyLinkInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpLoginPolicyLinkProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	var idp policy.IdentityProviderCascadeRemovedEvent

	switch e := event.(type) {
	case *org.IdentityProviderCascadeRemovedEvent:
		idp = e.IdentityProviderCascadeRemovedEvent
	case *instance.IdentityProviderCascadeRemovedEvent:
		idp = e.IdentityProviderCascadeRemovedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-iCKSj", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyIDPProviderCascadeRemovedEventType, instance.LoginPolicyIDPProviderCascadeRemovedEventType})
	}

	return handler.NewDeleteStatement(&idp,
		[]handler.Condition{
			handler.NewCond(IDPLoginPolicyLinkIDPIDCol, idp.IDPConfigID),
			handler.NewCond(IDPLoginPolicyLinkAggregateIDCol, idp.Aggregate().ID),
			handler.NewCond(IDPLoginPolicyLinkInstanceIDCol, idp.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpLoginPolicyLinkProjection) reduceIDPConfigRemoved(event eventstore.Event) (*handler.Statement, error) {
	var conditions []handler.Condition

	switch e := event.(type) {
	case *org.IDPConfigRemovedEvent:
		conditions = []handler.Condition{
			handler.NewCond(IDPLoginPolicyLinkIDPIDCol, e.ConfigID),
			handler.NewCond(IDPLoginPolicyLinkResourceOwnerCol, event.Aggregate().ResourceOwner),
			handler.NewCond(IDPLoginPolicyLinkInstanceIDCol, event.Aggregate().InstanceID),
		}
	case *instance.IDPConfigRemovedEvent:
		conditions = []handler.Condition{
			handler.NewCond(IDPLoginPolicyLinkIDPIDCol, e.ConfigID),
			handler.NewCond(IDPLoginPolicyLinkInstanceIDCol, event.Aggregate().InstanceID),
		}
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-u6tze", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigRemovedEventType, instance.IDPConfigRemovedEventType})
	}

	return handler.NewDeleteStatement(event, conditions), nil
}

func (p *idpLoginPolicyLinkProjection) reduceIDPRemoved(event eventstore.Event) (*handler.Statement, error) {
	var conditions []handler.Condition

	switch e := event.(type) {
	case *org.IDPRemovedEvent:
		conditions = []handler.Condition{
			handler.NewCond(IDPLoginPolicyLinkIDPIDCol, e.ID),
			handler.NewCond(IDPLoginPolicyLinkResourceOwnerCol, event.Aggregate().ResourceOwner),
			handler.NewCond(IDPLoginPolicyLinkInstanceIDCol, event.Aggregate().InstanceID),
		}
	case *instance.IDPRemovedEvent:
		conditions = []handler.Condition{
			handler.NewCond(IDPLoginPolicyLinkIDPIDCol, e.ID),
			handler.NewCond(IDPLoginPolicyLinkInstanceIDCol, event.Aggregate().InstanceID),
		}
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SFED3", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPRemovedEventType, instance.IDPRemovedEventType})
	}

	return handler.NewDeleteStatement(event, conditions), nil
}

func (p *idpLoginPolicyLinkProjection) reducePolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.LoginPolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SF3dg", "reduce.wrong.event.type %s", org.LoginPolicyRemovedEventType)
	}
	return handler.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(IDPLoginPolicyLinkAggregateIDCol, e.Aggregate().ID),
			handler.NewCond(IDPLoginPolicyLinkInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpLoginPolicyLinkProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-YbhOv", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(IDPLoginPolicyLinkInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(IDPLoginPolicyLinkResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
