package projection

import (
	"context"

	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/idpconfig"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
)

const (
	IDPTable     = "projections.idps"
	IDPOIDCTable = IDPTable + "_" + IDPOIDCSuffix
	IDPJWTTable  = IDPTable + "_" + IDPJWTSuffix

	IDPOIDCSuffix = "oidc_config"
	IDPJWTSuffix  = "jwt_config"

	IDPIDCol            = "id"
	IDPCreationDateCol  = "creation_date"
	IDPChangeDateCol    = "change_date"
	IDPSequenceCol      = "sequence"
	IDPResourceOwnerCol = "resource_owner"
	IDPInstanceIDCol    = "instance_id"
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

type IDPProjection struct {
	crdb.StatementHandler
}

func NewIDPProjection(ctx context.Context, config crdb.StatementHandlerConfig) *IDPProjection {
	p := new(IDPProjection)
	config.ProjectionName = IDPTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(IDPIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(IDPChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(IDPSequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(IDPResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(IDPNameCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPStylingTypeCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(IDPOwnerTypeCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(IDPAutoRegisterCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(IDPTypeCol, crdb.ColumnTypeEnum),
		},
			crdb.NewPrimaryKey(IDPInstanceIDCol, IDPIDCol),
			crdb.NewIndex("ro_idx", []string{IDPResourceOwnerCol}),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(OIDCConfigIDPIDCol, crdb.ColumnTypeText, crdb.DeleteCascade(IDPIDCol)),
			crdb.NewColumn(OIDCConfigClientIDCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(OIDCConfigClientSecretCol, crdb.ColumnTypeJSONB, crdb.Nullable()),
			crdb.NewColumn(OIDCConfigIssuerCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(OIDCConfigScopesCol, crdb.ColumnTypeTextArray, crdb.Nullable()),
			crdb.NewColumn(OIDCConfigDisplayNameMappingCol, crdb.ColumnTypeEnum, crdb.Nullable()),
			crdb.NewColumn(OIDCConfigUsernameMappingCol, crdb.ColumnTypeEnum, crdb.Nullable()),
			crdb.NewColumn(OIDCConfigAuthorizationEndpointCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(OIDCConfigTokenEndpointCol, crdb.ColumnTypeEnum, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(OIDCConfigIDPIDCol),
			IDPOIDCSuffix,
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(JWTConfigIDPIDCol, crdb.ColumnTypeText, crdb.DeleteCascade(IDPIDCol)),
			crdb.NewColumn(JWTConfigIssuerCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(JWTConfigKeysEndpointCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(JWTConfigHeaderNameCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(JWTConfigEndpointCol, crdb.ColumnTypeText, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(JWTConfigIDPIDCol),
			IDPJWTSuffix,
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *IDPProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.IDPConfigAddedEventType,
					Reduce: p.reduceIDPAdded,
				},
				{
					Event:  instance.IDPConfigChangedEventType,
					Reduce: p.reduceIDPChanged,
				},
				{
					Event:  instance.IDPConfigDeactivatedEventType,
					Reduce: p.reduceIDPDeactivated,
				},
				{
					Event:  instance.IDPConfigReactivatedEventType,
					Reduce: p.reduceIDPReactivated,
				},
				{
					Event:  instance.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
				{
					Event:  instance.IDPOIDCConfigAddedEventType,
					Reduce: p.reduceOIDCConfigAdded,
				},
				{
					Event:  instance.IDPOIDCConfigChangedEventType,
					Reduce: p.reduceOIDCConfigChanged,
				},
				{
					Event:  instance.IDPJWTConfigAddedEventType,
					Reduce: p.reduceJWTConfigAdded,
				},
				{
					Event:  instance.IDPJWTConfigChangedEventType,
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

func (p *IDPProjection) reduceIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.IDPConfigAddedEvent:
		idpEvent = e.IDPConfigAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.IDPConfigAddedEvent:
		idpEvent = e.IDPConfigAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-fcUdQ", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigAddedEventType, instance.IDPConfigAddedEventType})
	}

	return crdb.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPIDCol, idpEvent.ConfigID),
			handler.NewCol(IDPCreationDateCol, idpEvent.CreationDate()),
			handler.NewCol(IDPChangeDateCol, idpEvent.CreationDate()),
			handler.NewCol(IDPSequenceCol, idpEvent.Sequence()),
			handler.NewCol(IDPResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
			handler.NewCol(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPStateCol, domain.IDPConfigStateActive),
			handler.NewCol(IDPNameCol, idpEvent.Name),
			handler.NewCol(IDPStylingTypeCol, idpEvent.StylingType),
			handler.NewCol(IDPAutoRegisterCol, idpEvent.AutoRegister),
			handler.NewCol(IDPOwnerTypeCol, idpOwnerType),
		},
	), nil
}

func (p *IDPProjection) reduceIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPConfigChangedEvent:
		idpEvent = e.IDPConfigChangedEvent
	case *instance.IDPConfigChangedEvent:
		idpEvent = e.IDPConfigChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-NVvJD", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigChangedEventType, instance.IDPConfigChangedEventType})
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

func (p *IDPProjection) reduceIDPDeactivated(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigDeactivatedEvent
	switch e := event.(type) {
	case *org.IDPConfigDeactivatedEvent:
		idpEvent = e.IDPConfigDeactivatedEvent
	case *instance.IDPConfigDeactivatedEvent:
		idpEvent = e.IDPConfigDeactivatedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-94O5l", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigDeactivatedEventType, instance.IDPConfigDeactivatedEventType})
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

func (p *IDPProjection) reduceIDPReactivated(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigReactivatedEvent
	switch e := event.(type) {
	case *org.IDPConfigReactivatedEvent:
		idpEvent = e.IDPConfigReactivatedEvent
	case *instance.IDPConfigReactivatedEvent:
		idpEvent = e.IDPConfigReactivatedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-I8QyS", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigReactivatedEventType, instance.IDPConfigReactivatedEventType})
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

func (p *IDPProjection) reduceIDPRemoved(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigRemovedEvent
	switch e := event.(type) {
	case *org.IDPConfigRemovedEvent:
		idpEvent = e.IDPConfigRemovedEvent
	case *instance.IDPConfigRemovedEvent:
		idpEvent = e.IDPConfigRemovedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-B4zy8", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigRemovedEventType, instance.IDPConfigRemovedEventType})
	}

	return crdb.NewDeleteStatement(
		&idpEvent,
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
		},
	), nil
}

func (p *IDPProjection) reduceOIDCConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.OIDCConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPOIDCConfigAddedEvent:
		idpEvent = e.OIDCConfigAddedEvent
	case *instance.IDPOIDCConfigAddedEvent:
		idpEvent = e.OIDCConfigAddedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-2FuAA", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigAddedEventType, instance.IDPOIDCConfigAddedEventType})
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

func (p *IDPProjection) reduceOIDCConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.OIDCConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPOIDCConfigChangedEvent:
		idpEvent = e.OIDCConfigChangedEvent
	case *instance.IDPOIDCConfigChangedEvent:
		idpEvent = e.OIDCConfigChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-x2IVI", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigChangedEventType, instance.IDPOIDCConfigChangedEventType})
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

func (p *IDPProjection) reduceJWTConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.JWTConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPJWTConfigAddedEvent:
		idpEvent = e.JWTConfigAddedEvent
	case *instance.IDPJWTConfigAddedEvent:
		idpEvent = e.JWTConfigAddedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-qvPdb", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigAddedEventType, instance.IDPJWTConfigAddedEventType})
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

func (p *IDPProjection) reduceJWTConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.JWTConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPJWTConfigChangedEvent:
		idpEvent = e.JWTConfigChangedEvent
	case *instance.IDPJWTConfigChangedEvent:
		idpEvent = e.JWTConfigChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-x2IVI", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigChangedEventType, instance.IDPJWTConfigChangedEventType})
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
