package projection

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	IDPRelationalTable                = "zitadel.identity_providers"
	IDPRelationalOrgIdCol             = "org_id"
	IDPRelationalAllowAutoCreationCol = "allow_auto_creation"
)

type idpRelationalProjection struct{}

func newIDPRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(idpRelationalProjection))
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
				// {
				// 	Event:  instance.IDPConfigChangedEventType,
				// 	Reduce: p.reduceIDPChanged,
				// },
				// {
				// 	Event:  instance.IDPConfigDeactivatedEventType,
				// 	Reduce: p.reduceIDPDeactivated,
				// },
				// {
				// 	Event:  instance.IDPConfigReactivatedEventType,
				// 	Reduce: p.reduceIDPReactivated,
				// },
				// {
				// 	Event:  instance.IDPConfigRemovedEventType,
				// 	Reduce: p.reduceIDPRemoved,
				// },
				// {
				// 	Event:  instance.IDPOIDCConfigAddedEventType,
				// 	Reduce: p.reduceOIDCConfigAdded,
				// },
				// {
				// 	Event:  instance.IDPOIDCConfigChangedEventType,
				// 	Reduce: p.reduceOIDCConfigChanged,
				// },
				// {
				// 	Event:  instance.IDPJWTConfigAddedEventType,
				// 	Reduce: p.reduceJWTConfigAdded,
				// },
				// {
				// 	Event:  instance.IDPJWTConfigChangedEventType,
				// 	Reduce: p.reduceJWTConfigChanged,
				// },
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-uYR5R", "reduce.wrong.event.type %s", org.OrgAddedEventType)
	}

	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(IDPIDCol, e.ConfigID),
			handler.NewCol(IDPRelationalOrgIdCol, nil),
			handler.NewCol(IDPInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCol(IDPStateCol, domain.IDPStateActive.String()),
			handler.NewCol(IDPNameCol, e.Name),
			handler.NewCol(IDPStylingTypeCol, e.StylingType),
			handler.NewCol(IDPRelationalAllowAutoCreationCol, e.AutoRegister),
			handler.NewCol(IDPTypeCol, domain.IDPTypeOIDC.String()),
			handler.NewCol(CreatedAt, e.CreationDate()),
			handler.NewCol(UpdatedAt, e.CreationDate()),
		},
	), nil
}

// func (p *idpRelationalProjection) reduceIDPChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.IDPConfigChangedEvent
// 	switch e := event.(type) {
// 	case *org.IDPConfigChangedEvent:
// 		idpEvent = e.IDPConfigChangedEvent
// 	case *instance.IDPConfigChangedEvent:
// 		idpEvent = e.IDPConfigChangedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-NVvJD", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigChangedEventType, instance.IDPConfigChangedEventType})
// 	}

// 	cols := make([]handler.Column, 0, 5)
// 	if idpEvent.Name != nil {
// 		cols = append(cols, handler.NewCol(IDPNameCol, *idpEvent.Name))
// 	}
// 	if idpEvent.StylingType != nil {
// 		cols = append(cols, handler.NewCol(IDPStylingTypeCol, *idpEvent.StylingType))
// 	}
// 	if idpEvent.AutoRegister != nil {
// 		cols = append(cols, handler.NewCol(IDPAutoRegisterCol, *idpEvent.AutoRegister))
// 	}
// 	if len(cols) == 0 {
// 		return handler.NewNoOpStatement(&idpEvent), nil
// 	}

// 	cols = append(cols,
// 		handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
// 		handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
// 	)

// 	return handler.NewUpdateStatement(
// 		&idpEvent,
// 		cols,
// 		[]handler.Condition{
// 			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
// 			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 		},
// 	), nil
// }

// func (p *idpRelationalProjection) reduceIDPDeactivated(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.IDPConfigDeactivatedEvent
// 	switch e := event.(type) {
// 	case *org.IDPConfigDeactivatedEvent:
// 		idpEvent = e.IDPConfigDeactivatedEvent
// 	case *instance.IDPConfigDeactivatedEvent:
// 		idpEvent = e.IDPConfigDeactivatedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-94O5l", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigDeactivatedEventType, instance.IDPConfigDeactivatedEventType})
// 	}

// 	return handler.NewUpdateStatement(
// 		&idpEvent,
// 		[]handler.Column{
// 			handler.NewCol(IDPStateCol, domain.IDPConfigStateInactive),
// 			handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
// 			handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
// 			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 		},
// 	), nil
// }

// func (p *idpRelationalProjection) reduceIDPReactivated(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.IDPConfigReactivatedEvent
// 	switch e := event.(type) {
// 	case *org.IDPConfigReactivatedEvent:
// 		idpEvent = e.IDPConfigReactivatedEvent
// 	case *instance.IDPConfigReactivatedEvent:
// 		idpEvent = e.IDPConfigReactivatedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-I8QyS", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigReactivatedEventType, instance.IDPConfigReactivatedEventType})
// 	}

// 	return handler.NewUpdateStatement(
// 		&idpEvent,
// 		[]handler.Column{
// 			handler.NewCol(IDPStateCol, domain.IDPConfigStateActive),
// 			handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
// 			handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
// 			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 		},
// 	), nil
// }

// func (p *idpRelationalProjection) reduceIDPRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.IDPConfigRemovedEvent
// 	switch e := event.(type) {
// 	case *org.IDPConfigRemovedEvent:
// 		idpEvent = e.IDPConfigRemovedEvent
// 	case *instance.IDPConfigRemovedEvent:
// 		idpEvent = e.IDPConfigRemovedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-B4zy8", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigRemovedEventType, instance.IDPConfigRemovedEventType})
// 	}

// 	return handler.NewDeleteStatement(
// 		&idpEvent,
// 		[]handler.Condition{
// 			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
// 			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 		},
// 	), nil
// }

// func (p *idpRelationalProjection) reduceOIDCConfigAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.OIDCConfigAddedEvent
// 	switch e := event.(type) {
// 	case *org.IDPOIDCConfigAddedEvent:
// 		idpEvent = e.OIDCConfigAddedEvent
// 	case *instance.IDPOIDCConfigAddedEvent:
// 		idpEvent = e.OIDCConfigAddedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-2FuAA", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigAddedEventType, instance.IDPOIDCConfigAddedEventType})
// 	}

// 	return handler.NewMultiStatement(&idpEvent,
// 		handler.AddUpdateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
// 				handler.NewCol(IDPTypeCol, domain.IDPConfigTypeOIDC),
// 			},
// 			[]handler.Condition{
// 				handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
// 				handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(OIDCConfigIDPIDCol, idpEvent.IDPConfigID),
// 				handler.NewCol(OIDCConfigInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(OIDCConfigClientIDCol, idpEvent.ClientID),
// 				handler.NewCol(OIDCConfigClientSecretCol, idpEvent.ClientSecret),
// 				handler.NewCol(OIDCConfigIssuerCol, idpEvent.Issuer),
// 				handler.NewCol(OIDCConfigScopesCol, database.TextArray[string](idpEvent.Scopes)),
// 				handler.NewCol(OIDCConfigDisplayNameMappingCol, idpEvent.IDPDisplayNameMapping),
// 				handler.NewCol(OIDCConfigUsernameMappingCol, idpEvent.UserNameMapping),
// 				handler.NewCol(OIDCConfigAuthorizationEndpointCol, idpEvent.AuthorizationEndpoint),
// 				handler.NewCol(OIDCConfigTokenEndpointCol, idpEvent.TokenEndpoint),
// 			},
// 			handler.WithTableSuffix(IDPOIDCSuffix),
// 		),
// 	), nil
// }

// func (p *idpRelationalProjection) reduceOIDCConfigChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.OIDCConfigChangedEvent
// 	switch e := event.(type) {
// 	case *org.IDPOIDCConfigChangedEvent:
// 		idpEvent = e.OIDCConfigChangedEvent
// 	case *instance.IDPOIDCConfigChangedEvent:
// 		idpEvent = e.OIDCConfigChangedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x2IVI", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigChangedEventType, instance.IDPOIDCConfigChangedEventType})
// 	}

// 	cols := make([]handler.Column, 0, 8)

// 	if idpEvent.ClientID != nil {
// 		cols = append(cols, handler.NewCol(OIDCConfigClientIDCol, *idpEvent.ClientID))
// 	}
// 	if idpEvent.ClientSecret != nil {
// 		cols = append(cols, handler.NewCol(OIDCConfigClientSecretCol, idpEvent.ClientSecret))
// 	}
// 	if idpEvent.Issuer != nil {
// 		cols = append(cols, handler.NewCol(OIDCConfigIssuerCol, *idpEvent.Issuer))
// 	}
// 	if idpEvent.AuthorizationEndpoint != nil {
// 		cols = append(cols, handler.NewCol(OIDCConfigAuthorizationEndpointCol, *idpEvent.AuthorizationEndpoint))
// 	}
// 	if idpEvent.TokenEndpoint != nil {
// 		cols = append(cols, handler.NewCol(OIDCConfigTokenEndpointCol, *idpEvent.TokenEndpoint))
// 	}
// 	if idpEvent.Scopes != nil {
// 		cols = append(cols, handler.NewCol(OIDCConfigScopesCol, database.TextArray[string](idpEvent.Scopes)))
// 	}
// 	if idpEvent.IDPDisplayNameMapping != nil {
// 		cols = append(cols, handler.NewCol(OIDCConfigDisplayNameMappingCol, *idpEvent.IDPDisplayNameMapping))
// 	}
// 	if idpEvent.UserNameMapping != nil {
// 		cols = append(cols, handler.NewCol(OIDCConfigUsernameMappingCol, *idpEvent.UserNameMapping))
// 	}

// 	if len(cols) == 0 {
// 		return handler.NewNoOpStatement(&idpEvent), nil
// 	}

// 	return handler.NewMultiStatement(&idpEvent,
// 		handler.AddUpdateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
// 			},
// 			[]handler.Condition{
// 				handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
// 				handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),
// 		handler.AddUpdateStatement(
// 			cols,
// 			[]handler.Condition{
// 				handler.NewCond(OIDCConfigIDPIDCol, idpEvent.IDPConfigID),
// 				handler.NewCond(OIDCConfigInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 			handler.WithTableSuffix(IDPOIDCSuffix),
// 		),
// 	), nil
// }

// func (p *idpRelationalProjection) reduceJWTConfigAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.JWTConfigAddedEvent
// 	switch e := event.(type) {
// 	case *org.IDPJWTConfigAddedEvent:
// 		idpEvent = e.JWTConfigAddedEvent
// 	case *instance.IDPJWTConfigAddedEvent:
// 		idpEvent = e.JWTConfigAddedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qvPdb", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigAddedEventType, instance.IDPJWTConfigAddedEventType})
// 	}

// 	return handler.NewMultiStatement(&idpEvent,
// 		handler.AddUpdateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
// 				handler.NewCol(IDPTypeCol, domain.IDPConfigTypeJWT),
// 			},
// 			[]handler.Condition{
// 				handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
// 				handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),

// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(JWTConfigIDPIDCol, idpEvent.IDPConfigID),
// 				handler.NewCol(JWTConfigInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(JWTConfigEndpointCol, idpEvent.JWTEndpoint),
// 				handler.NewCol(JWTConfigIssuerCol, idpEvent.Issuer),
// 				handler.NewCol(JWTConfigKeysEndpointCol, idpEvent.KeysEndpoint),
// 				handler.NewCol(JWTConfigHeaderNameCol, idpEvent.HeaderName),
// 			},
// 			handler.WithTableSuffix(IDPJWTSuffix),
// 		),
// 	), nil
// }

// func (p *idpRelationalProjection) reduceJWTConfigChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.JWTConfigChangedEvent
// 	switch e := event.(type) {
// 	case *org.IDPJWTConfigChangedEvent:
// 		idpEvent = e.JWTConfigChangedEvent
// 	case *instance.IDPJWTConfigChangedEvent:
// 		idpEvent = e.JWTConfigChangedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x2IVI", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigChangedEventType, instance.IDPJWTConfigChangedEventType})
// 	}

// 	cols := make([]handler.Column, 0, 4)

// 	if idpEvent.JWTEndpoint != nil {
// 		cols = append(cols, handler.NewCol(JWTConfigEndpointCol, *idpEvent.JWTEndpoint))
// 	}
// 	if idpEvent.Issuer != nil {
// 		cols = append(cols, handler.NewCol(JWTConfigIssuerCol, *idpEvent.Issuer))
// 	}
// 	if idpEvent.KeysEndpoint != nil {
// 		cols = append(cols, handler.NewCol(JWTConfigKeysEndpointCol, *idpEvent.KeysEndpoint))
// 	}
// 	if idpEvent.HeaderName != nil {
// 		cols = append(cols, handler.NewCol(JWTConfigHeaderNameCol, *idpEvent.HeaderName))
// 	}

// 	if len(cols) == 0 {
// 		return handler.NewNoOpStatement(&idpEvent), nil
// 	}

// 	return handler.NewMultiStatement(&idpEvent,
// 		handler.AddUpdateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
// 			},
// 			[]handler.Condition{
// 				handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
// 				handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),
// 		handler.AddUpdateStatement(
// 			cols,
// 			[]handler.Condition{
// 				handler.NewCond(JWTConfigIDPIDCol, idpEvent.IDPConfigID),
// 				handler.NewCond(JWTConfigInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 			handler.WithTableSuffix(IDPJWTSuffix),
// 		),
// 	), nil
// }

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
