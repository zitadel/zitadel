package projection

import (
	"context"

	"github.com/lib/pq"
	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type idpProjection struct {
	crdb.StatementHandler
}

const (
	IDPTable     = "zitadel.projections.idps"
	IDPOIDCTable = IDPTable + "_" + IDPOIDCSuffix
	IDPJWTTable  = IDPTable + "_" + IDPJWTSuffix
)

func newIDPProjection(ctx context.Context, config crdb.StatementHandlerConfig) *idpProjection {
	p := &idpProjection{}
	config.ProjectionName = IDPTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *idpProjection) reducers() []handler.AggregateReducer {
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
	IDPOIDCSuffix = "oidc_config"
	IDPJWTSuffix  = "jwt_config"

	IDPIDCol            = "id"
	IDPCreationDateCol  = "creation_date"
	IDPChangeDateCol    = "change_date"
	IDPSequenceCol      = "sequence"
	IDPResourceOwnerCol = "resource_owner"
	IDPStateCol         = "state"
	IDPNameCol          = "name"
	IDPStylingTypeCol   = "styling_type"
	IDPOwnerTypeCol     = "owner_type"
	IDPAutoRegisterCol  = "auto_register"
	IDPTypeCol          = "type"

	OIDCConfigIDPIDCol                 = "idp_id"
	OIDCConfigClientIDCol              = "client_id"
	OIDCConfigClientSecretCol          = "client_secret"
	OIDCConfigIssuerCol                = "issuer"
	OIDCConfigScopesCol                = "scopes"
	OIDCConfigDisplayNameMappingCol    = "display_name_mapping"
	OIDCConfigUsernameMappingCol       = "username_mapping"
	OIDCConfigAuthorizationEndpointCol = "authorization_endpoint"
	OIDCConfigTokenEndpointCol         = "token_endpoint"

	JWTConfigIDPIDCol        = "idp_id"
	JWTConfigIssuerCol       = "issuer"
	JWTConfigKeysEndpointCol = "keys_endpoint"
	JWTConfigHeaderNameCol   = "header_name"
	JWTConfigEndpointCol     = "endpoint"
)

func (p *idpProjection) reduceIDPAdded(event eventstore.Event) (*handler.Statement, error) {
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
		logging.LogWithFields("HANDL-hBriG", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPConfigAddedEventType, iam.IDPConfigAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-fcUdQ", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPIDCol, idpEvent.ConfigID),
			handler.NewCol(IDPCreationDateCol, idpEvent.CreationDate()),
			handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
			handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
			handler.NewCol(IDPResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
			handler.NewCol(IDPStateCol, domain.IDPConfigStateActive),
			handler.NewCol(IDPNameCol, idpEvent.Name),
			handler.NewCol(IDPStylingTypeCol, idpEvent.StylingType),
			handler.NewCol(IDPAutoRegisterCol, idpEvent.AutoRegister),
			handler.NewCol(IDPOwnerTypeCol, idpOwnerType),
		},
	), nil
}

func (p *idpProjection) reduceIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPConfigChangedEvent:
		idpEvent = e.IDPConfigChangedEvent
	case *iam.IDPConfigChangedEvent:
		idpEvent = e.IDPConfigChangedEvent
	default:
		logging.LogWithFields("HANDL-FFrph", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPConfigChangedEventType, iam.IDPConfigChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-NVvJD", "reduce.wrong.event.type")
	}

	cols := make([]handler.Column, 0, 5)
	if idpEvent.Name != nil {
		cols = append(cols, handler.NewCol(IDPNameCol, *idpEvent.Name))
	}
	if idpEvent.StylingType != nil {
		cols = append(cols, handler.NewCol(IDPStylingTypeCol, *idpEvent.StylingType))
	}
	if idpEvent.AutoRegister != nil {
		cols = append(cols, handler.NewCol(IDPAutoRegisterCol, *idpEvent.AutoRegister))
	}
	if len(cols) == 0 {
		return crdb.NewNoOpStatement(&idpEvent), nil
	}

	cols = append(cols,
		handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
		handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
	)

	return crdb.NewUpdateStatement(
		&idpEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
		},
	), nil
}

func (p *idpProjection) reduceIDPDeactivated(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigDeactivatedEvent
	switch e := event.(type) {
	case *org.IDPConfigDeactivatedEvent:
		idpEvent = e.IDPConfigDeactivatedEvent
	case *iam.IDPConfigDeactivatedEvent:
		idpEvent = e.IDPConfigDeactivatedEvent
	default:
		logging.LogWithFields("HANDL-1s33a", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPConfigDeactivatedEventType, iam.IDPConfigDeactivatedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-94O5l", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPStateCol, domain.IDPConfigStateInactive),
			handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
			handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
		},
	), nil
}

func (p *idpProjection) reduceIDPReactivated(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigReactivatedEvent
	switch e := event.(type) {
	case *org.IDPConfigReactivatedEvent:
		idpEvent = e.IDPConfigReactivatedEvent
	case *iam.IDPConfigReactivatedEvent:
		idpEvent = e.IDPConfigReactivatedEvent
	default:
		logging.LogWithFields("HANDL-Zgzpt", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPConfigReactivatedEventType, iam.IDPConfigReactivatedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-I8QyS", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPStateCol, domain.IDPConfigStateActive),
			handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
			handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
		},
	), nil
}

func (p *idpProjection) reduceIDPRemoved(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigRemovedEvent
	switch e := event.(type) {
	case *org.IDPConfigRemovedEvent:
		idpEvent = e.IDPConfigRemovedEvent
	case *iam.IDPConfigRemovedEvent:
		idpEvent = e.IDPConfigRemovedEvent
	default:
		logging.LogWithFields("HANDL-JJasT", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPConfigRemovedEventType, iam.IDPConfigRemovedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-B4zy8", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(
		&idpEvent,
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
		},
	), nil
}

func (p *idpProjection) reduceOIDCConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.OIDCConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPOIDCConfigAddedEvent:
		idpEvent = e.OIDCConfigAddedEvent
	case *iam.IDPOIDCConfigAddedEvent:
		idpEvent = e.OIDCConfigAddedEvent
	default:
		logging.LogWithFields("HANDL-DCmeB", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPOIDCConfigAddedEventType, iam.IDPOIDCConfigAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-2FuAA", "reduce.wrong.event.type")
	}

	return crdb.NewMultiStatement(&idpEvent,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTypeCol, domain.IDPConfigTypeOIDC),
			},
			[]handler.Condition{
				handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(OIDCConfigIDPIDCol, idpEvent.IDPConfigID),
				handler.NewCol(OIDCConfigClientIDCol, idpEvent.ClientID),
				handler.NewCol(OIDCConfigClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(OIDCConfigIssuerCol, idpEvent.Issuer),
				handler.NewCol(OIDCConfigScopesCol, pq.StringArray(idpEvent.Scopes)),
				handler.NewCol(OIDCConfigDisplayNameMappingCol, idpEvent.IDPDisplayNameMapping),
				handler.NewCol(OIDCConfigUsernameMappingCol, idpEvent.UserNameMapping),
				handler.NewCol(OIDCConfigAuthorizationEndpointCol, idpEvent.AuthorizationEndpoint),
				handler.NewCol(OIDCConfigTokenEndpointCol, idpEvent.TokenEndpoint),
			},
			crdb.WithTableSuffix(IDPOIDCSuffix),
		),
	), nil
}

func (p *idpProjection) reduceOIDCConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.OIDCConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPOIDCConfigChangedEvent:
		idpEvent = e.OIDCConfigChangedEvent
	case *iam.IDPOIDCConfigChangedEvent:
		idpEvent = e.OIDCConfigChangedEvent
	default:
		logging.LogWithFields("HANDL-VyBm2", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPOIDCConfigChangedEventType, iam.IDPOIDCConfigChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-x2IVI", "reduce.wrong.event.type")
	}

	cols := make([]handler.Column, 0, 8)

	if idpEvent.ClientID != nil {
		cols = append(cols, handler.NewCol(OIDCConfigClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		cols = append(cols, handler.NewCol(OIDCConfigClientSecretCol, idpEvent.ClientSecret))
	}
	if idpEvent.Issuer != nil {
		cols = append(cols, handler.NewCol(OIDCConfigIssuerCol, *idpEvent.Issuer))
	}
	if idpEvent.AuthorizationEndpoint != nil {
		cols = append(cols, handler.NewCol(OIDCConfigAuthorizationEndpointCol, *idpEvent.AuthorizationEndpoint))
	}
	if idpEvent.TokenEndpoint != nil {
		cols = append(cols, handler.NewCol(OIDCConfigTokenEndpointCol, *idpEvent.TokenEndpoint))
	}
	if idpEvent.Scopes != nil {
		cols = append(cols, handler.NewCol(OIDCConfigScopesCol, pq.StringArray(idpEvent.Scopes)))
	}
	if idpEvent.IDPDisplayNameMapping != nil {
		cols = append(cols, handler.NewCol(OIDCConfigDisplayNameMappingCol, *idpEvent.IDPDisplayNameMapping))
	}
	if idpEvent.UserNameMapping != nil {
		cols = append(cols, handler.NewCol(OIDCConfigUsernameMappingCol, *idpEvent.UserNameMapping))
	}

	if len(cols) == 0 {
		return crdb.NewNoOpStatement(&idpEvent), nil
	}

	return crdb.NewMultiStatement(&idpEvent,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
			},
		),
		crdb.AddUpdateStatement(
			cols,
			[]handler.Condition{
				handler.NewCond(OIDCConfigIDPIDCol, idpEvent.IDPConfigID),
			},
			crdb.WithTableSuffix(IDPOIDCSuffix),
		),
	), nil
}

func (p *idpProjection) reduceJWTConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.JWTConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPJWTConfigAddedEvent:
		idpEvent = e.JWTConfigAddedEvent
	case *iam.IDPJWTConfigAddedEvent:
		idpEvent = e.JWTConfigAddedEvent
	default:
		logging.LogWithFields("HANDL-228q7", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPJWTConfigAddedEventType, iam.IDPJWTConfigAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-qvPdb", "reduce.wrong.event.type")
	}

	return crdb.NewMultiStatement(&idpEvent,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTypeCol, domain.IDPConfigTypeJWT),
			},
			[]handler.Condition{
				handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
			},
		),

		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(OIDCConfigIDPIDCol, idpEvent.IDPConfigID),
				handler.NewCol(JWTConfigEndpointCol, idpEvent.JWTEndpoint),
				handler.NewCol(JWTConfigIssuerCol, idpEvent.Issuer),
				handler.NewCol(JWTConfigKeysEndpointCol, idpEvent.KeysEndpoint),
				handler.NewCol(JWTConfigHeaderNameCol, idpEvent.HeaderName),
			},
			crdb.WithTableSuffix(IDPJWTSuffix),
		),
	), nil
}

func (p *idpProjection) reduceJWTConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.JWTConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPJWTConfigChangedEvent:
		idpEvent = e.JWTConfigChangedEvent
	case *iam.IDPJWTConfigChangedEvent:
		idpEvent = e.JWTConfigChangedEvent
	default:
		logging.LogWithFields("HANDL-VyBm2", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.IDPJWTConfigChangedEventType, iam.IDPJWTConfigChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-x2IVI", "reduce.wrong.event.type")
	}

	cols := make([]handler.Column, 0, 4)

	if idpEvent.JWTEndpoint != nil {
		cols = append(cols, handler.NewCol(JWTConfigEndpointCol, *idpEvent.JWTEndpoint))
	}
	if idpEvent.Issuer != nil {
		cols = append(cols, handler.NewCol(JWTConfigIssuerCol, *idpEvent.Issuer))
	}
	if idpEvent.KeysEndpoint != nil {
		cols = append(cols, handler.NewCol(JWTConfigKeysEndpointCol, *idpEvent.KeysEndpoint))
	}
	if idpEvent.HeaderName != nil {
		cols = append(cols, handler.NewCol(JWTConfigHeaderNameCol, *idpEvent.HeaderName))
	}

	if len(cols) == 0 {
		return crdb.NewNoOpStatement(&idpEvent), nil
	}

	return crdb.NewMultiStatement(&idpEvent,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
			},
		),
		crdb.AddUpdateStatement(
			cols,
			[]handler.Condition{
				handler.NewCond(OIDCConfigIDPIDCol, idpEvent.IDPConfigID),
			},
			crdb.WithTableSuffix(IDPJWTSuffix),
		),
	), nil
}
