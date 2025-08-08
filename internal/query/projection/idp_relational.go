package projection

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	IDPRelationalTable           = "zitadel.identity_providers"
	IDPRelationalOrgIdCol        = "org_id"
	IDPRelationalAutoRegisterCol = "auto_register"
	IDPRelationalPayloadCol      = "payload"
)

type idpRelationalProjection struct {
	idpRepo domain.IDProviderRepository
}

func newIDPRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	client := postgres.PGxPool(config.Client.Pool)
	idpRepo := repository.IDProviderRepository(client)

	return handler.NewHandler(ctx, &config, &idpRelationalProjection{
		idpRepo: idpRepo,
	})
}

func (*idpRelationalProjection) Name() string {
	return IDPRelationalTable
}

func (p *idpRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.IDPConfigAddedEventType,
					Reduce: p.reduceIDPRelationalAdded,
				},
				{
					Event:  instance.IDPConfigChangedEventType,
					Reduce: p.reduceIDPRelationalChanged,
				},
				{
					Event:  instance.IDPConfigDeactivatedEventType,
					Reduce: p.reduceIDRelationalPDeactivated,
				},
				{
					Event:  instance.IDPConfigReactivatedEventType,
					Reduce: p.reduceIDPRelationalReactivated,
				},
				{
					Event:  instance.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPRelationalRemoved,
				},
				{
					Event:  instance.IDPOIDCConfigAddedEventType,
					Reduce: p.reduceOIDCRelationalConfigAdded,
				},
				{
					Event:  instance.IDPOIDCConfigChangedEventType,
					Reduce: p.reduceOIDCRelationalConfigChanged,
				},
				{
					Event:  instance.IDPJWTConfigAddedEventType,
					Reduce: p.reduceJWTRelationalConfigAdded,
				},
				{
					Event:  instance.IDPJWTConfigChangedEventType,
					Reduce: p.reduceJWTRelationalConfigChanged,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.IDPConfigAddedEventType,
					Reduce: p.reduceIDPRelationalAdded,
				},
				{
					Event:  org.IDPConfigChangedEventType,
					Reduce: p.reduceIDPRelationalChanged,
				},
				{
					Event:  org.IDPConfigDeactivatedEventType,
					Reduce: p.reduceIDRelationalPDeactivated,
				},
				{
					Event:  org.IDPConfigReactivatedEventType,
					Reduce: p.reduceIDPRelationalReactivated,
				},
				{
					Event:  org.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPRelationalRemoved,
				},
				{
					Event:  org.IDPOIDCConfigAddedEventType,
					Reduce: p.reduceOIDCRelationalConfigAdded,
				},
				{
					Event:  org.IDPOIDCConfigChangedEventType,
					Reduce: p.reduceOIDCRelationalConfigChanged,
				},
				{
					Event:  org.IDPJWTConfigAddedEventType,
					Reduce: p.reduceJWTRelationalConfigAdded,
				},
				{
					Event:  org.IDPJWTConfigChangedEventType,
					Reduce: p.reduceJWTRelationalConfigChanged,
				},
			},
		},
	}
}

func (p *idpRelationalProjection) reduceIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPConfigAddedEvent:
		idpEvent = e.IDPConfigAddedEvent
	case *instance.IDPConfigAddedEvent:
		idpEvent = e.IDPConfigAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YcUdQ", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigAddedEventType, instance.IDPConfigAddedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPRelationalOrgIdCol, orgId),
			handler.NewCol(IDPIDCol, idpEvent.ConfigID),
			handler.NewCol(IDPStateCol, domain.IDPStateActive.String()),
			handler.NewCol(IDPNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeUnspecified.String()),
			handler.NewCol(IDPRelationalAutoRegisterCol, idpEvent.AutoRegister),
			handler.NewCol(IDPRelationalAllowCreationCol, true),
			handler.NewCol(IDPRelationalAllowAutoUpdateCol, false),
			handler.NewCol(IDPRelationalAllowLinkingCol, true),
			handler.NewCol(IDPRelationalAllowAutoLinkingCol, domain.IDPAutoLinkingOptionUnspecified.String()),
			handler.NewCol(IDPStylingTypeCol, idpEvent.StylingType),
			handler.NewCol(CreatedAt, idpEvent.CreationDate()),
		},
	), nil
}

func (p *idpRelationalProjection) reduceIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPConfigChangedEvent:
		idpEvent = e.IDPConfigChangedEvent
	case *instance.IDPConfigChangedEvent:
		idpEvent = e.IDPConfigChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YVvJD", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigChangedEventType, instance.IDPConfigChangedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	cols := make([]handler.Column, 0, 5)
	if idpEvent.Name != nil {
		cols = append(cols, handler.NewCol(IDPNameCol, *idpEvent.Name))
	}
	if idpEvent.StylingType != nil {
		cols = append(cols, handler.NewCol(IDPStylingTypeCol, *idpEvent.StylingType))
	}
	if idpEvent.AutoRegister != nil {
		cols = append(cols, handler.NewCol(IDPRelationalAutoRegisterCol, *idpEvent.AutoRegister))
	}
	if len(cols) == 0 {
		return handler.NewNoOpStatement(&idpEvent), nil
	}

	return handler.NewUpdateStatement(
		&idpEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCond(IDPRelationalOrgId, orgId),
		},
	), nil
}

func (p *idpRelationalProjection) reduceIDRelationalPDeactivated(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigDeactivatedEvent
	switch e := event.(type) {
	case *org.IDPConfigDeactivatedEvent:
		idpEvent = e.IDPConfigDeactivatedEvent
	case *instance.IDPConfigDeactivatedEvent:
		idpEvent = e.IDPConfigDeactivatedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y4O5l", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigDeactivatedEventType, instance.IDPConfigDeactivatedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPStateCol, domain.IDPStateInactive.String()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCond(IDPRelationalOrgId, orgId),
		},
	), nil
}

func (p *idpRelationalProjection) reduceIDPRelationalReactivated(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigReactivatedEvent
	switch e := event.(type) {
	case *org.IDPConfigReactivatedEvent:
		idpEvent = e.IDPConfigReactivatedEvent
	case *instance.IDPConfigReactivatedEvent:
		idpEvent = e.IDPConfigReactivatedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y8QyS", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigReactivatedEventType, instance.IDPConfigReactivatedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPStateCol, domain.IDPStateActive.String()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCond(IDPRelationalOrgId, orgId),
		},
	), nil
}

func (p *idpRelationalProjection) reduceIDPRelationalRemoved(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigRemovedEvent
	switch e := event.(type) {
	case *org.IDPConfigRemovedEvent:
		idpEvent = e.IDPConfigRemovedEvent
	case *instance.IDPConfigRemovedEvent:
		idpEvent = e.IDPConfigRemovedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y4zy8", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigRemovedEventType, instance.IDPConfigRemovedEventType})
	}
	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewDeleteStatement(
		&idpEvent,
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCond(IDPRelationalOrgId, orgId),
		},
	), nil
}

func (p *idpRelationalProjection) reduceOIDCRelationalConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.OIDCConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPOIDCConfigAddedEvent:
		idpEvent = e.OIDCConfigAddedEvent
	case *instance.IDPOIDCConfigAddedEvent:
		idpEvent = e.OIDCConfigAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YFuAA", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigAddedEventType, instance.IDPOIDCConfigAddedEventType})
	}

	payload, err := json.Marshal(idpEvent)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPRelationalPayloadCol, payload),
			handler.NewCol(IDPTypeCol, domain.IDPTypeOIDC.String()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCond(IDPRelationalOrgId, orgId),
		},
	), nil
}

func (p *idpRelationalProjection) reduceOIDCRelationalConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.OIDCConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPOIDCConfigChangedEvent:
		idpEvent = e.OIDCConfigChangedEvent
	case *instance.IDPOIDCConfigChangedEvent:
		idpEvent = e.OIDCConfigChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y2IVI", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigChangedEventType, instance.IDPOIDCConfigChangedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	oidc, err := p.idpRepo.GetOIDC(context.Background(), p.idpRepo.IDCondition(idpEvent.IDPConfigID), idpEvent.Agg.InstanceID, orgId)
	if err != nil {
		return nil, err
	}

	if idpEvent.ClientID != nil {
		oidc.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		oidc.ClientSecret = *idpEvent.ClientSecret
	}
	if idpEvent.Issuer != nil {
		oidc.Issuer = *idpEvent.Issuer
	}
	if idpEvent.AuthorizationEndpoint != nil {
		oidc.AuthorizationEndpoint = *idpEvent.AuthorizationEndpoint
	}
	if idpEvent.TokenEndpoint != nil {
		oidc.TokenEndpoint = *idpEvent.TokenEndpoint
	}
	if idpEvent.Scopes != nil {
		oidc.Scopes = idpEvent.Scopes
	}
	if idpEvent.IDPDisplayNameMapping != nil {
		oidc.IDPDisplayNameMapping = domain.OIDCMappingField(*idpEvent.IDPDisplayNameMapping)
	}
	if idpEvent.UserNameMapping != nil {
		oidc.UserNameMapping = domain.OIDCMappingField(*idpEvent.UserNameMapping)
	}

	payload, err := json.Marshal(idpEvent)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPRelationalPayloadCol, payload),
			handler.NewCol(IDPTypeCol, domain.IDPTypeOIDC.String()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCond(IDPRelationalOrgId, orgId),
		},
	), nil
}

func (p *idpRelationalProjection) reduceJWTRelationalConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.JWTConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPJWTConfigAddedEvent:
		idpEvent = e.JWTConfigAddedEvent
	case *instance.IDPJWTConfigAddedEvent:
		idpEvent = e.JWTConfigAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YvPdb", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigAddedEventType, instance.IDPJWTConfigAddedEventType})
	}

	payload, err := json.Marshal(idpEvent)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPRelationalPayloadCol, payload),
			handler.NewCol(IDPTypeCol, domain.IDPTypeJWT.String()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCond(IDPRelationalOrgId, orgId),
		},
	), nil
}

func (p *idpRelationalProjection) reduceJWTRelationalConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.JWTConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPJWTConfigChangedEvent:
		idpEvent = e.JWTConfigChangedEvent
	case *instance.IDPJWTConfigChangedEvent:
		idpEvent = e.JWTConfigChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y2IVI", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigChangedEventType, instance.IDPJWTConfigChangedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	jwt, err := p.idpRepo.GetJWT(context.Background(), p.idpRepo.IDCondition(idpEvent.IDPConfigID), idpEvent.Agg.InstanceID, orgId)
	if err != nil {
		return nil, err
	}

	if idpEvent.JWTEndpoint != nil {
		jwt.JWTEndpoint = *idpEvent.JWTEndpoint
	}
	if idpEvent.Issuer != nil {
		jwt.Issuer = *idpEvent.Issuer
	}
	if idpEvent.KeysEndpoint != nil {
		jwt.KeysEndpoint = *idpEvent.KeysEndpoint
	}
	if idpEvent.HeaderName != nil {
		jwt.HeaderName = *idpEvent.HeaderName
	}

	payload, err := json.Marshal(idpEvent)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPRelationalPayloadCol, payload),
			handler.NewCol(IDPTypeCol, domain.IDPTypeJWT.String()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCond(IDPRelationalOrgId, orgId),
		},
	), nil
}

// func (p *idpProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	e, ok := event.(*org.OrgRemovedEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-YsbQC", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
// 	}

// 	return handler.NewDeleteStatement(
// 		e,
// 		[]handler.Condition{
// 			handler.NewCond(IDPInstanceIDCol, e.Aggregate().InstanceID),
// 			handler.NewCond(IDPResourceOwnerCol, e.Aggregate().ID),
// 		},
// 	), nil
// }
