package projection

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
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
					Reduce: p.reduceJWTConfigChanged,
				},
				// {
				// 	Event:  instance.InstanceRemovedEventType,
				// 	Reduce: reduceInstanceRemovedHelper(IDPInstanceIDCol),
				// },
			},
		},
	}
}

func (p *idpRelationalProjection) reduceIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.IDPConfigAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-fcUdQ", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigAddedEventType, instance.IDPConfigAddedEventType})
	}

	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(IDPInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCol(IDPRelationalOrgIdCol, nil),
			handler.NewCol(IDPIDCol, e.ConfigID),
			handler.NewCol(IDPStateCol, domain.IDPStateActive.String()),
			handler.NewCol(IDPNameCol, e.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeUnspecified.String()),
			handler.NewCol(IDPRelationalAutoRegisterCol, e.AutoRegister),
			handler.NewCol(IDPRelationalAllowCreationCol, true),
			handler.NewCol(IDPRelationalAllowAutoUpdateCol, false),
			handler.NewCol(IDPRelationalAllowLinkingCol, true),
			handler.NewCol(IDPRelationalAllowAutoLinkingCol, domain.IDPAutoLinkingOptionUnspecified.String()),
			handler.NewCol(IDPStylingTypeCol, e.StylingType),
			handler.NewCol(CreatedAt, e.CreationDate()),
		},
	), nil
}

func (p *idpRelationalProjection) reduceIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.IDPConfigChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-NVvJD", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigChangedEventType, instance.IDPConfigChangedEventType})
	}

	cols := make([]handler.Column, 0, 5)
	if e.Name != nil {
		cols = append(cols, handler.NewCol(IDPNameCol, *e.Name))
	}
	if e.StylingType != nil {
		cols = append(cols, handler.NewCol(IDPStylingTypeCol, *e.StylingType))
	}
	if e.AutoRegister != nil {
		cols = append(cols, handler.NewCol(IDPRelationalAutoRegisterCol, *e.AutoRegister))
	}
	if len(cols) == 0 {
		return handler.NewNoOpStatement(e), nil
	}

	return handler.NewUpdateStatement(
		e,
		cols,
		[]handler.Condition{
			handler.NewCond(IDPIDCol, e.ConfigID),
			handler.NewCond(IDPInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpRelationalProjection) reduceIDRelationalPDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.IDPConfigDeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-94O5l", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigDeactivatedEventType, instance.IDPConfigDeactivatedEventType})
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(IDPStateCol, domain.IDPStateInactive.String()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, e.ConfigID),
			handler.NewCond(IDPInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpRelationalProjection) reduceIDPRelationalReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.IDPConfigReactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-I8QyS", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigReactivatedEventType, instance.IDPConfigReactivatedEventType})
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(IDPStateCol, domain.IDPStateActive.String()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, e.ConfigID),
			handler.NewCond(IDPInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpRelationalProjection) reduceIDPRelationalRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.IDPConfigRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-B4zy8", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigRemovedEventType, instance.IDPConfigRemovedEventType})
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(IDPIDCol, e.ConfigID),
			handler.NewCond(IDPInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpRelationalProjection) reduceOIDCRelationalConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.IDPOIDCConfigAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-2FuAA", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigAddedEventType, instance.IDPOIDCConfigAddedEventType})
	}

	payload, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(IDPRelationalPayloadCol, payload),
			handler.NewCol(IDPTypeCol, domain.IDPTypeOIDC.String()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, e.IDPConfigID),
		},
	), nil
}

func (p *idpRelationalProjection) reduceOIDCRelationalConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.IDPOIDCConfigChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x2IBI", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigChangedEventType, instance.IDPOIDCConfigChangedEventType})
	}

	oidc, err := p.idpRepo.GetOIDC(context.Background(), p.idpRepo.IDCondition(e.IDPConfigID), e.Agg.InstanceID, nil)
	if err != nil {
		return nil, err
	}

	if e.ClientID != nil {
		oidc.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		oidc.ClientSecret = *e.ClientSecret
	}
	if e.Issuer != nil {
		oidc.Issuer = *e.Issuer
	}
	if e.AuthorizationEndpoint != nil {
		oidc.AuthorizationEndpoint = *e.AuthorizationEndpoint
	}
	if e.TokenEndpoint != nil {
		oidc.TokenEndpoint = *e.TokenEndpoint
	}
	if e.Scopes != nil {
		oidc.Scopes = e.Scopes
	}
	if e.IDPDisplayNameMapping != nil {
		oidc.IDPDisplayNameMapping = domain.OIDCMappingField(*e.IDPDisplayNameMapping)
	}
	if e.UserNameMapping != nil {
		oidc.UserNameMapping = domain.OIDCMappingField(*e.UserNameMapping)
	}

	payload, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(IDPRelationalPayloadCol, payload),
			handler.NewCol(IDPTypeCol, domain.IDPTypeOIDC.String()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, e.IDPConfigID),
		},
	), nil
}

func (p *idpRelationalProjection) reduceJWTRelationalConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.IDPJWTConfigAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qvPdb", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigAddedEventType, instance.IDPJWTConfigAddedEventType})
	}
	payload, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(IDPRelationalPayloadCol, payload),
			handler.NewCol(IDPTypeCol, domain.IDPTypeJWT.String()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, e.IDPConfigID),
		},
	), nil
}

func (p *idpRelationalProjection) reduceJWTConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.IDPJWTConfigChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-P2I9I", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigChangedEventType, instance.IDPJWTConfigChangedEventType})
	}

	jwt, err := p.idpRepo.GetJWT(context.Background(), p.idpRepo.IDCondition(e.IDPConfigID), e.Agg.InstanceID, nil)
	if err != nil {
		return nil, err
	}

	if e.JWTEndpoint != nil {
		jwt.JWTEndpoint = *e.JWTEndpoint
	}
	if e.Issuer != nil {
		jwt.Issuer = *e.Issuer
	}
	if e.KeysEndpoint != nil {
		jwt.KeysEndpoint = *e.KeysEndpoint
	}
	if e.HeaderName != nil {
		jwt.HeaderName = *e.HeaderName
	}

	payload, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(IDPRelationalPayloadCol, payload),
			handler.NewCol(IDPTypeCol, domain.IDPTypeJWT.String()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, e.IDPConfigID),
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
