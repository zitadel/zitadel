package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/idpconfig"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/lib/pq"
)

type IDPProjection struct {
	crdb.StatementHandler
}

const (
	idpProjection = "zitadel.projections.idps"
)

func NewIDPProjection(ctx context.Context, config crdb.StatementHandlerConfig) *IDPProjection {
	p := &IDPProjection{}
	config.ProjectionName = idpProjection
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *IDPProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.IDPConfigAddedEventType,
					Reduce: p.reduceIDPAdded,
				},
				{
					Event:  iam.IDPConfigChangedEventType,
					Reduce: p.reduceIDPChanged,
				},
				{
					Event:  iam.IDPConfigDeactivatedEventType,
					Reduce: p.reduceIDPDeactivated,
				},
				{
					Event:  iam.IDPConfigReactivatedEventType,
					Reduce: p.reduceIDPReactivated,
				},
				{
					Event:  iam.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
				{
					Event:  iam.IDPOIDCConfigAddedEventType,
					Reduce: p.reduceOIDCConfigAdded,
				},
				{
					Event:  iam.IDPOIDCConfigChangedEventType,
					Reduce: p.reduceOIDCConfigChanged,
				},
				{
					Event:  iam.IDPJWTConfigAddedEventType,
					Reduce: p.reduceJWTConfigAdded,
				},
				{
					Event:  iam.IDPJWTConfigChangedEventType,
					Reduce: p.reduceJWTConfigChanged,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.IDPConfigAddedEventType,
					Reduce: p.reduceIDPAdded,
				},
				{
					Event:  org.IDPConfigChangedEventType,
					Reduce: p.reduceIDPChanged,
				},
				{
					Event:  org.IDPConfigDeactivatedEventType,
					Reduce: p.reduceIDPDeactivated,
				},
				{
					Event:  org.IDPConfigReactivatedEventType,
					Reduce: p.reduceIDPReactivated,
				},
				{
					Event:  org.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
				{
					Event:  org.IDPOIDCConfigAddedEventType,
					Reduce: p.reduceOIDCConfigAdded,
				},
				{
					Event:  org.IDPOIDCConfigChangedEventType,
					Reduce: p.reduceOIDCConfigChanged,
				},
				{
					Event:  org.IDPJWTConfigAddedEventType,
					Reduce: p.reduceJWTConfigAdded,
				},
				{
					Event:  org.IDPJWTConfigChangedEventType,
					Reduce: p.reduceJWTConfigChanged,
				},
			},
		},
	}
}

const (
	idpOIDCSuffix = "oidc_config"
	idpJWTSuffix  = "jwt_config"

	idpIDCol           = "id"
	idpStateCol        = "state"
	idpNameCol         = "name"
	idpStylingTypeCol  = "styling_type"
	idpOwnerCol        = "owner"
	idpAutoRegisterCol = "auto_register"

	oidcConfigIDPIDCol                 = "idp_id"
	oidcConfigClientIDCol              = "client_id"
	oidcConfigClientSecretCol          = "client_secret"
	oidcConfigIssuerCol                = "issuer"
	oidcConfigScopesCol                = "scopes"
	oidcConfigDisplayNameMappingCol    = "display_name_mapping"
	oidcConfigUsernameMappingCol       = "username_mapping"
	oidcConfigAuthorizationEndpointCol = "authorization_endpoint"
	oidcConfigTokenEndpointCol         = "token_endpoint"

	jwtConfigIDPIDCol        = "idp_id"
	jwtConfigIssuerCol       = "issuer"
	jwtConfigKeysEndpointCol = "keys_endpoint"
	jwtConfigHeaderNameCol   = "header_name"
	jwtConfigEndpointCol     = "endpoint"
)

func (p *IDPProjection) reduceIDPAdded(event eventstore.EventReader) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.IDPConfigAddedEvent:
		idpEvent = e.IDPConfigAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *iam.IDPConfigAddedEvent:
		idpEvent = e.IDPConfigAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		logging.LogWithFields("HANDL-hBriG", "seq", e.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPConfigAddedEventType, iam.IDPConfigAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-fcUdQ", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(idpIDCol, idpEvent.ConfigID),
			handler.NewCol(idpStateCol, domain.IDPConfigStateActive),
			handler.NewCol(idpNameCol, idpEvent.Name),
			handler.NewCol(idpStylingTypeCol, idpEvent.StylingType),
			handler.NewCol(idpAutoRegisterCol, idpEvent.AutoRegister),
			handler.NewCol(idpOwnerCol, idpOwnerType),
		},
	), nil
}

func (p *IDPProjection) reduceIDPChanged(event eventstore.EventReader) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPConfigChangedEvent:
		idpEvent = e.IDPConfigChangedEvent
	case *iam.IDPConfigChangedEvent:
		idpEvent = e.IDPConfigChangedEvent
	default:
		logging.LogWithFields("HANDL-FFrph", "seq", e.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPConfigChangedEventType, iam.IDPConfigChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-NVvJD", "reduce.wrong.event.type")
	}

	cols := make([]handler.Column, 0, 3)
	if idpEvent.Name != nil {
		cols = append(cols, handler.NewCol(idpNameCol, *idpEvent.Name))
	}
	if idpEvent.StylingType != nil {
		cols = append(cols, handler.NewCol(idpStylingTypeCol, *idpEvent.StylingType))
	}
	if idpEvent.AutoRegister != nil {
		cols = append(cols, handler.NewCol(idpAutoRegisterCol, *idpEvent.AutoRegister))
	}
	if len(cols) == 0 {
		return crdb.NewNoOpStatement(&idpEvent), nil
	}

	return crdb.NewUpdateStatement(
		&idpEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(idpIDCol, idpEvent.ConfigID),
		},
	), nil
}

func (p *IDPProjection) reduceIDPDeactivated(event eventstore.EventReader) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigDeactivatedEvent
	switch e := event.(type) {
	case *org.IDPConfigDeactivatedEvent:
		idpEvent = e.IDPConfigDeactivatedEvent
	case *iam.IDPConfigDeactivatedEvent:
		idpEvent = e.IDPConfigDeactivatedEvent
	default:
		logging.LogWithFields("HANDL-1s33a", "seq", e.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPConfigDeactivatedEventType, iam.IDPConfigDeactivatedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-94O5l", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(idpStateCol, domain.IDPConfigStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(idpIDCol, idpEvent.ConfigID),
		},
	), nil
}

func (p *IDPProjection) reduceIDPReactivated(event eventstore.EventReader) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigReactivatedEvent
	switch e := event.(type) {
	case *org.IDPConfigReactivatedEvent:
		idpEvent = e.IDPConfigReactivatedEvent
	case *iam.IDPConfigReactivatedEvent:
		idpEvent = e.IDPConfigReactivatedEvent
	default:
		logging.LogWithFields("HANDL-Zgzpt", "seq", e.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPConfigReactivatedEventType, iam.IDPConfigReactivatedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-I8QyS", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(idpStateCol, domain.IDPConfigStateActive),
		},
		[]handler.Condition{
			handler.NewCond(idpIDCol, idpEvent.ConfigID),
		},
	), nil
}

func (p *IDPProjection) reduceIDPRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigRemovedEvent
	switch e := event.(type) {
	case *org.IDPConfigRemovedEvent:
		idpEvent = e.IDPConfigRemovedEvent
	case *iam.IDPConfigRemovedEvent:
		idpEvent = e.IDPConfigRemovedEvent
	default:
		logging.LogWithFields("HANDL-JJasT", "seq", e.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPConfigRemovedEventType, iam.IDPConfigRemovedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-B4zy8", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(
		&idpEvent,
		[]handler.Condition{
			handler.NewCond(idpIDCol, idpEvent.ConfigID),
		},
	), nil
}

func (p *IDPProjection) reduceOIDCConfigAdded(event eventstore.EventReader) (*handler.Statement, error) {
	var idpEvent idpconfig.OIDCConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPOIDCConfigAddedEvent:
		idpEvent = e.OIDCConfigAddedEvent
	case *iam.IDPOIDCConfigAddedEvent:
		idpEvent = e.OIDCConfigAddedEvent
	default:
		logging.LogWithFields("HANDL-DCmeB", "seq", e.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPOIDCConfigAddedEventType, iam.IDPOIDCConfigAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-2FuAA", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(oidcConfigIDPIDCol, idpEvent.IDPConfigID),
			handler.NewCol(oidcConfigClientIDCol, idpEvent.ClientID),
			handler.NewCol(oidcConfigClientSecretCol, idpEvent.ClientSecret),
			handler.NewCol(oidcConfigIssuerCol, idpEvent.Issuer),
			handler.NewCol(oidcConfigScopesCol, pq.StringArray(idpEvent.Scopes)),
			handler.NewCol(oidcConfigDisplayNameMappingCol, idpEvent.IDPDisplayNameMapping),
			handler.NewCol(oidcConfigUsernameMappingCol, idpEvent.UserNameMapping),
			handler.NewCol(oidcConfigAuthorizationEndpointCol, idpEvent.AuthorizationEndpoint),
			handler.NewCol(oidcConfigTokenEndpointCol, idpEvent.TokenEndpoint),
		},
		crdb.WithTableSuffix(idpOIDCSuffix),
	), nil
}

func (p *IDPProjection) reduceOIDCConfigChanged(event eventstore.EventReader) (*handler.Statement, error) {
	var idpEvent idpconfig.OIDCConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPOIDCConfigChangedEvent:
		idpEvent = e.OIDCConfigChangedEvent
	case *iam.IDPOIDCConfigChangedEvent:
		idpEvent = e.OIDCConfigChangedEvent
	default:
		logging.LogWithFields("HANDL-VyBm2", "seq", e.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPOIDCConfigChangedEventType, iam.IDPOIDCConfigChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-x2IVI", "reduce.wrong.event.type")
	}

	cols := make([]handler.Column, 0, 8)

	if idpEvent.ClientID != nil {
		cols = append(cols, handler.NewCol(oidcConfigClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		cols = append(cols, handler.NewCol(oidcConfigClientSecretCol, *idpEvent.ClientSecret))
	}
	if idpEvent.Issuer != nil {
		cols = append(cols, handler.NewCol(oidcConfigIssuerCol, *idpEvent.Issuer))
	}
	if idpEvent.AuthorizationEndpoint != nil {
		cols = append(cols, handler.NewCol(oidcConfigAuthorizationEndpointCol, *idpEvent.AuthorizationEndpoint))
	}
	if idpEvent.TokenEndpoint != nil {
		cols = append(cols, handler.NewCol(oidcConfigTokenEndpointCol, *idpEvent.TokenEndpoint))
	}
	if idpEvent.Scopes != nil {
		cols = append(cols, handler.NewCol(oidcConfigScopesCol, idpEvent.Scopes))
	}
	if idpEvent.IDPDisplayNameMapping != nil {
		cols = append(cols, handler.NewCol(oidcConfigDisplayNameMappingCol, *idpEvent.IDPDisplayNameMapping))
	}
	if idpEvent.UserNameMapping != nil {
		cols = append(cols, handler.NewCol(oidcConfigUsernameMappingCol, *idpEvent.UserNameMapping))
	}

	if len(cols) == 0 {
		return crdb.NewNoOpStatement(&idpEvent), nil
	}

	return crdb.NewUpdateStatement(
		&idpEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(oidcConfigIDPIDCol, idpEvent.IDPConfigID),
		},
		crdb.WithTableSuffix(idpOIDCSuffix),
	), nil
}

func (p *IDPProjection) reduceJWTConfigAdded(event eventstore.EventReader) (*handler.Statement, error) {
	var idpEvent idpconfig.JWTConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPJWTConfigAddedEvent:
		idpEvent = e.JWTConfigAddedEvent
	case *iam.IDPJWTConfigAddedEvent:
		idpEvent = e.JWTConfigAddedEvent
	default:
		logging.LogWithFields("HANDL-228q7", "seq", e.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPJWTConfigAddedEventType, iam.IDPJWTConfigAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-qvPdb", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(oidcConfigIDPIDCol, idpEvent.IDPConfigID),
			handler.NewCol(jwtConfigEndpointCol, idpEvent.JWTEndpoint),
			handler.NewCol(jwtConfigIssuerCol, idpEvent.Issuer),
			handler.NewCol(jwtConfigKeysEndpointCol, idpEvent.KeysEndpoint),
			handler.NewCol(jwtConfigHeaderNameCol, idpEvent.HeaderName),
		},
		crdb.WithTableSuffix(idpJWTSuffix),
	), nil
}

func (p *IDPProjection) reduceJWTConfigChanged(event eventstore.EventReader) (*handler.Statement, error) {
	var idpEvent idpconfig.JWTConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPJWTConfigChangedEvent:
		idpEvent = e.JWTConfigChangedEvent
	case *iam.IDPJWTConfigChangedEvent:
		idpEvent = e.JWTConfigChangedEvent
	default:
		logging.LogWithFields("HANDL-VyBm2", "seq", e.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPJWTConfigChangedEventType, iam.IDPJWTConfigChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-x2IVI", "reduce.wrong.event.type")
	}

	cols := make([]handler.Column, 0, 4)

	if idpEvent.JWTEndpoint != nil {
		cols = append(cols, handler.NewCol(jwtConfigEndpointCol, *idpEvent.JWTEndpoint))
	}
	if idpEvent.Issuer != nil {
		cols = append(cols, handler.NewCol(jwtConfigIssuerCol, *idpEvent.Issuer))
	}
	if idpEvent.KeysEndpoint != nil {
		cols = append(cols, handler.NewCol(jwtConfigKeysEndpointCol, *idpEvent.Issuer))
	}
	if idpEvent.HeaderName != nil {
		cols = append(cols, handler.NewCol(jwtConfigHeaderNameCol, *idpEvent.HeaderName))
	}

	if len(cols) == 0 {
		return crdb.NewNoOpStatement(&idpEvent), nil
	}

	return crdb.NewUpdateStatement(
		&idpEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(oidcConfigIDPIDCol, idpEvent.IDPConfigID),
		},
		crdb.WithTableSuffix(idpJWTSuffix),
	), nil
}
