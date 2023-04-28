package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
)

const (
	AppProjectionTable = "projections.apps5"
	AppAPITable        = AppProjectionTable + "_" + appAPITableSuffix
	AppOIDCTable       = AppProjectionTable + "_" + appOIDCTableSuffix
	AppSAMLTable       = AppProjectionTable + "_" + appSAMLTableSuffix

	AppColumnID            = "id"
	AppColumnName          = "name"
	AppColumnProjectID     = "project_id"
	AppColumnCreationDate  = "creation_date"
	AppColumnChangeDate    = "change_date"
	AppColumnResourceOwner = "resource_owner"
	AppColumnInstanceID    = "instance_id"
	AppColumnState         = "state"
	AppColumnSequence      = "sequence"
	AppColumnOwnerRemoved  = "owner_removed"

	appAPITableSuffix              = "api_configs"
	AppAPIConfigColumnAppID        = "app_id"
	AppAPIConfigColumnInstanceID   = "instance_id"
	AppAPIConfigColumnClientID     = "client_id"
	AppAPIConfigColumnClientSecret = "client_secret"
	AppAPIConfigColumnAuthMethod   = "auth_method"

	appOIDCTableSuffix                          = "oidc_configs"
	AppOIDCConfigColumnAppID                    = "app_id"
	AppOIDCConfigColumnInstanceID               = "instance_id"
	AppOIDCConfigColumnVersion                  = "version"
	AppOIDCConfigColumnClientID                 = "client_id"
	AppOIDCConfigColumnClientSecret             = "client_secret"
	AppOIDCConfigColumnRedirectUris             = "redirect_uris"
	AppOIDCConfigColumnResponseTypes            = "response_types"
	AppOIDCConfigColumnGrantTypes               = "grant_types"
	AppOIDCConfigColumnApplicationType          = "application_type"
	AppOIDCConfigColumnAuthMethodType           = "auth_method_type"
	AppOIDCConfigColumnPostLogoutRedirectUris   = "post_logout_redirect_uris"
	AppOIDCConfigColumnDevMode                  = "is_dev_mode"
	AppOIDCConfigColumnAccessTokenType          = "access_token_type"
	AppOIDCConfigColumnAccessTokenRoleAssertion = "access_token_role_assertion"
	AppOIDCConfigColumnIDTokenRoleAssertion     = "id_token_role_assertion"
	AppOIDCConfigColumnIDTokenUserinfoAssertion = "id_token_userinfo_assertion"
	AppOIDCConfigColumnClockSkew                = "clock_skew"
	AppOIDCConfigColumnAdditionalOrigins        = "additional_origins"
	AppOIDCConfigColumnSkipNativeAppSuccessPage = "skip_native_app_success_page"

	appSAMLTableSuffix             = "saml_configs"
	AppSAMLConfigColumnAppID       = "app_id"
	AppSAMLConfigColumnInstanceID  = "instance_id"
	AppSAMLConfigColumnEntityID    = "entity_id"
	AppSAMLConfigColumnMetadata    = "metadata"
	AppSAMLConfigColumnMetadataURL = "metadata_url"
)

type appProjection struct {
	crdb.StatementHandler
}

func newAppProjection(ctx context.Context, config crdb.StatementHandlerConfig) *appProjection {
	p := new(appProjection)
	config.ProjectionName = AppProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(AppColumnID, crdb.ColumnTypeText),
			crdb.NewColumn(AppColumnName, crdb.ColumnTypeText),
			crdb.NewColumn(AppColumnProjectID, crdb.ColumnTypeText),
			crdb.NewColumn(AppColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(AppColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(AppColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(AppColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(AppColumnState, crdb.ColumnTypeEnum),
			crdb.NewColumn(AppColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(AppColumnOwnerRemoved, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(AppColumnInstanceID, AppColumnID),
			crdb.WithIndex(crdb.NewIndex("project_id", []string{AppColumnProjectID})),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{AppColumnOwnerRemoved})),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(AppAPIConfigColumnAppID, crdb.ColumnTypeText),
			crdb.NewColumn(AppAPIConfigColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(AppAPIConfigColumnClientID, crdb.ColumnTypeText),
			crdb.NewColumn(AppAPIConfigColumnClientSecret, crdb.ColumnTypeJSONB, crdb.Nullable()),
			crdb.NewColumn(AppAPIConfigColumnAuthMethod, crdb.ColumnTypeEnum),
		},
			crdb.NewPrimaryKey(AppAPIConfigColumnInstanceID, AppAPIConfigColumnAppID),
			appAPITableSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
			crdb.WithIndex(crdb.NewIndex("client_id", []string{AppAPIConfigColumnClientID})),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(AppOIDCConfigColumnAppID, crdb.ColumnTypeText),
			crdb.NewColumn(AppOIDCConfigColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(AppOIDCConfigColumnVersion, crdb.ColumnTypeEnum),
			crdb.NewColumn(AppOIDCConfigColumnClientID, crdb.ColumnTypeText),
			crdb.NewColumn(AppOIDCConfigColumnClientSecret, crdb.ColumnTypeJSONB, crdb.Nullable()),
			crdb.NewColumn(AppOIDCConfigColumnRedirectUris, crdb.ColumnTypeTextArray, crdb.Nullable()),
			crdb.NewColumn(AppOIDCConfigColumnResponseTypes, crdb.ColumnTypeEnumArray, crdb.Nullable()),
			crdb.NewColumn(AppOIDCConfigColumnGrantTypes, crdb.ColumnTypeEnumArray, crdb.Nullable()),
			crdb.NewColumn(AppOIDCConfigColumnApplicationType, crdb.ColumnTypeEnum),
			crdb.NewColumn(AppOIDCConfigColumnAuthMethodType, crdb.ColumnTypeEnum),
			crdb.NewColumn(AppOIDCConfigColumnPostLogoutRedirectUris, crdb.ColumnTypeTextArray, crdb.Nullable()),
			crdb.NewColumn(AppOIDCConfigColumnDevMode, crdb.ColumnTypeBool),
			crdb.NewColumn(AppOIDCConfigColumnAccessTokenType, crdb.ColumnTypeEnum),
			crdb.NewColumn(AppOIDCConfigColumnAccessTokenRoleAssertion, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(AppOIDCConfigColumnIDTokenRoleAssertion, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(AppOIDCConfigColumnIDTokenUserinfoAssertion, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(AppOIDCConfigColumnClockSkew, crdb.ColumnTypeInt64, crdb.Default(0)),
			crdb.NewColumn(AppOIDCConfigColumnAdditionalOrigins, crdb.ColumnTypeTextArray, crdb.Nullable()),
			crdb.NewColumn(AppOIDCConfigColumnSkipNativeAppSuccessPage, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(AppOIDCConfigColumnInstanceID, AppOIDCConfigColumnAppID),
			appOIDCTableSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
			crdb.WithIndex(crdb.NewIndex("client_id", []string{AppOIDCConfigColumnClientID})),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(AppSAMLConfigColumnAppID, crdb.ColumnTypeText),
			crdb.NewColumn(AppSAMLConfigColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(AppSAMLConfigColumnEntityID, crdb.ColumnTypeText),
			crdb.NewColumn(AppSAMLConfigColumnMetadata, crdb.ColumnTypeBytes),
			crdb.NewColumn(AppSAMLConfigColumnMetadataURL, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(AppSAMLConfigColumnInstanceID, AppSAMLConfigColumnAppID),
			appSAMLTableSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
			crdb.WithIndex(crdb.NewIndex("entity_id", []string{AppSAMLConfigColumnEntityID})),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *appProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  project.ApplicationAddedType,
					Reduce: p.reduceAppAdded,
				},
				{
					Event:  project.ApplicationChangedType,
					Reduce: p.reduceAppChanged,
				},
				{
					Event:  project.ApplicationDeactivatedType,
					Reduce: p.reduceAppDeactivated,
				},
				{
					Event:  project.ApplicationReactivatedType,
					Reduce: p.reduceAppReactivated,
				},
				{
					Event:  project.ApplicationRemovedType,
					Reduce: p.reduceAppRemoved,
				},
				{
					Event:  project.ProjectRemovedType,
					Reduce: p.reduceProjectRemoved,
				},
				{
					Event:  project.APIConfigAddedType,
					Reduce: p.reduceAPIConfigAdded,
				},
				{
					Event:  project.APIConfigChangedType,
					Reduce: p.reduceAPIConfigChanged,
				},
				{
					Event:  project.APIConfigSecretChangedType,
					Reduce: p.reduceAPIConfigSecretChanged,
				},
				{
					Event:  project.OIDCConfigAddedType,
					Reduce: p.reduceOIDCConfigAdded,
				},
				{
					Event:  project.OIDCConfigChangedType,
					Reduce: p.reduceOIDCConfigChanged,
				},
				{
					Event:  project.OIDCConfigSecretChangedType,
					Reduce: p.reduceOIDCConfigSecretChanged,
				},
				{
					Event:  project.SAMLConfigAddedType,
					Reduce: p.reduceSAMLConfigAdded,
				},
				{
					Event:  project.SAMLConfigChangedType,
					Reduce: p.reduceSAMLConfigChanged,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(AppColumnInstanceID),
				},
			},
		},
	}
}

func (p *appProjection) reduceAppAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ApplicationAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-1xYE6", "reduce.wrong.event.type %s", project.ApplicationAddedType)
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(AppColumnID, e.AppID),
			handler.NewCol(AppColumnName, e.Name),
			handler.NewCol(AppColumnProjectID, e.Aggregate().ID),
			handler.NewCol(AppColumnCreationDate, e.CreationDate()),
			handler.NewCol(AppColumnChangeDate, e.CreationDate()),
			handler.NewCol(AppColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(AppColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(AppColumnState, domain.AppStateActive),
			handler.NewCol(AppColumnSequence, e.Sequence()),
		},
	), nil
}

func (p *appProjection) reduceAppChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ApplicationChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ZJ8JA", "reduce.wrong.event.type %s", project.ApplicationChangedType)
	}
	if e.Name == "" {
		return crdb.NewNoOpStatement(event), nil
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(AppColumnName, e.Name),
			handler.NewCol(AppColumnChangeDate, e.CreationDate()),
			handler.NewCol(AppColumnSequence, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(AppColumnID, e.AppID),
			handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *appProjection) reduceAppDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ApplicationDeactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-MVWxZ", "reduce.wrong.event.type %s", project.ApplicationDeactivatedType)
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(AppColumnState, domain.AppStateInactive),
			handler.NewCol(AppColumnChangeDate, e.CreationDate()),
			handler.NewCol(AppColumnSequence, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(AppColumnID, e.AppID),
			handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *appProjection) reduceAppReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ApplicationReactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-D0HZO", "reduce.wrong.event.type %s", project.ApplicationReactivatedType)
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(AppColumnState, domain.AppStateActive),
			handler.NewCol(AppColumnChangeDate, e.CreationDate()),
			handler.NewCol(AppColumnSequence, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(AppColumnID, e.AppID),
			handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *appProjection) reduceAppRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ApplicationRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Y99aq", "reduce.wrong.event.type %s", project.ApplicationRemovedType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(AppColumnID, e.AppID),
			handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *appProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-DlUlO", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(AppColumnProjectID, e.Aggregate().ID),
			handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *appProjection) reduceAPIConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.APIConfigAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Y99aq", "reduce.wrong.event.type %s", project.APIConfigAddedType)
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(AppAPIConfigColumnAppID, e.AppID),
				handler.NewCol(AppAPIConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(AppAPIConfigColumnClientID, e.ClientID),
				handler.NewCol(AppAPIConfigColumnClientSecret, e.ClientSecret),
				handler.NewCol(AppAPIConfigColumnAuthMethod, e.AuthMethodType),
			},
			crdb.WithTableSuffix(appAPITableSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(AppColumnChangeDate, e.CreationDate()),
				handler.NewCol(AppColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(AppColumnID, e.AppID),
				handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *appProjection) reduceAPIConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.APIConfigChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-vnZKi", "reduce.wrong.event.type %s", project.APIConfigChangedType)
	}
	cols := make([]handler.Column, 0, 2)
	if e.ClientSecret != nil {
		cols = append(cols, handler.NewCol(AppAPIConfigColumnClientSecret, e.ClientSecret))
	}
	if e.AuthMethodType != nil {
		cols = append(cols, handler.NewCol(AppAPIConfigColumnAuthMethod, *e.AuthMethodType))
	}
	if len(cols) == 0 {
		return crdb.NewNoOpStatement(e), nil
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			cols,
			[]handler.Condition{
				handler.NewCond(AppAPIConfigColumnAppID, e.AppID),
				handler.NewCond(AppAPIConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(appAPITableSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(AppColumnChangeDate, e.CreationDate()),
				handler.NewCol(AppColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(AppColumnID, e.AppID),
				handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *appProjection) reduceAPIConfigSecretChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.APIConfigSecretChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ttb0I", "reduce.wrong.event.type %s", project.APIConfigSecretChangedType)
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(AppAPIConfigColumnClientSecret, e.ClientSecret),
			},
			[]handler.Condition{
				handler.NewCond(AppAPIConfigColumnAppID, e.AppID),
				handler.NewCond(AppAPIConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(appAPITableSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(AppColumnChangeDate, e.CreationDate()),
				handler.NewCol(AppColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(AppColumnID, e.AppID),
				handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *appProjection) reduceOIDCConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.OIDCConfigAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-GNHU1", "reduce.wrong.event.type %s", project.OIDCConfigAddedType)
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(AppOIDCConfigColumnAppID, e.AppID),
				handler.NewCol(AppOIDCConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(AppOIDCConfigColumnVersion, e.Version),
				handler.NewCol(AppOIDCConfigColumnClientID, e.ClientID),
				handler.NewCol(AppOIDCConfigColumnClientSecret, e.ClientSecret),
				handler.NewCol(AppOIDCConfigColumnRedirectUris, database.StringArray(e.RedirectUris)),
				handler.NewCol(AppOIDCConfigColumnResponseTypes, database.EnumArray[domain.OIDCResponseType](e.ResponseTypes)),
				handler.NewCol(AppOIDCConfigColumnGrantTypes, database.EnumArray[domain.OIDCGrantType](e.GrantTypes)),
				handler.NewCol(AppOIDCConfigColumnApplicationType, e.ApplicationType),
				handler.NewCol(AppOIDCConfigColumnAuthMethodType, e.AuthMethodType),
				handler.NewCol(AppOIDCConfigColumnPostLogoutRedirectUris, database.StringArray(e.PostLogoutRedirectUris)),
				handler.NewCol(AppOIDCConfigColumnDevMode, e.DevMode),
				handler.NewCol(AppOIDCConfigColumnAccessTokenType, e.AccessTokenType),
				handler.NewCol(AppOIDCConfigColumnAccessTokenRoleAssertion, e.AccessTokenRoleAssertion),
				handler.NewCol(AppOIDCConfigColumnIDTokenRoleAssertion, e.IDTokenRoleAssertion),
				handler.NewCol(AppOIDCConfigColumnIDTokenUserinfoAssertion, e.IDTokenUserinfoAssertion),
				handler.NewCol(AppOIDCConfigColumnClockSkew, e.ClockSkew),
				handler.NewCol(AppOIDCConfigColumnAdditionalOrigins, database.StringArray(e.AdditionalOrigins)),
				handler.NewCol(AppOIDCConfigColumnSkipNativeAppSuccessPage, e.SkipNativeAppSuccessPage),
			},
			crdb.WithTableSuffix(appOIDCTableSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(AppColumnChangeDate, e.CreationDate()),
				handler.NewCol(AppColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(AppColumnID, e.AppID),
				handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *appProjection) reduceOIDCConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.OIDCConfigChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-GNHU1", "reduce.wrong.event.type %s", project.OIDCConfigChangedType)
	}

	cols := make([]handler.Column, 0, 15)
	if e.Version != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnVersion, *e.Version))
	}
	if e.RedirectUris != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnRedirectUris, database.StringArray(*e.RedirectUris)))
	}
	if e.ResponseTypes != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnResponseTypes, database.EnumArray[domain.OIDCResponseType](*e.ResponseTypes)))
	}
	if e.GrantTypes != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnGrantTypes, database.EnumArray[domain.OIDCGrantType](*e.GrantTypes)))
	}
	if e.ApplicationType != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnApplicationType, *e.ApplicationType))
	}
	if e.AuthMethodType != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnAuthMethodType, *e.AuthMethodType))
	}
	if e.PostLogoutRedirectUris != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnPostLogoutRedirectUris, database.StringArray(*e.PostLogoutRedirectUris)))
	}
	if e.DevMode != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnDevMode, *e.DevMode))
	}
	if e.AccessTokenType != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnAccessTokenType, *e.AccessTokenType))
	}
	if e.AccessTokenRoleAssertion != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnAccessTokenRoleAssertion, *e.AccessTokenRoleAssertion))
	}
	if e.IDTokenRoleAssertion != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnIDTokenRoleAssertion, *e.IDTokenRoleAssertion))
	}
	if e.IDTokenUserinfoAssertion != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnIDTokenUserinfoAssertion, *e.IDTokenUserinfoAssertion))
	}
	if e.ClockSkew != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnClockSkew, *e.ClockSkew))
	}
	if e.AdditionalOrigins != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnAdditionalOrigins, database.StringArray(*e.AdditionalOrigins)))
	}
	if e.SkipNativeAppSuccessPage != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnSkipNativeAppSuccessPage, *e.SkipNativeAppSuccessPage))
	}

	if len(cols) == 0 {
		return crdb.NewNoOpStatement(e), nil
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			cols,
			[]handler.Condition{
				handler.NewCond(AppOIDCConfigColumnAppID, e.AppID),
				handler.NewCond(AppOIDCConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(appOIDCTableSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(AppColumnChangeDate, e.CreationDate()),
				handler.NewCol(AppColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(AppColumnID, e.AppID),
				handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *appProjection) reduceOIDCConfigSecretChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.OIDCConfigSecretChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-GNHU1", "reduce.wrong.event.type %s", project.OIDCConfigSecretChangedType)
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(AppOIDCConfigColumnClientSecret, e.ClientSecret),
			},
			[]handler.Condition{
				handler.NewCond(AppOIDCConfigColumnAppID, e.AppID),
				handler.NewCond(AppOIDCConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(appOIDCTableSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(AppColumnChangeDate, e.CreationDate()),
				handler.NewCol(AppColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(AppColumnID, e.AppID),
				handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *appProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-Hyd1f", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(AppColumnChangeDate, e.CreationDate()),
			handler.NewCol(AppColumnSequence, e.Sequence()),
			handler.NewCol(AppColumnOwnerRemoved, true),
		},
		[]handler.Condition{
			handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(AppColumnResourceOwner, e.Aggregate().ID),
		},
	), nil
}

func (p *appProjection) reduceSAMLConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.SAMLConfigAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-GMHU1", "reduce.wrong.event.type")
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(AppSAMLConfigColumnAppID, e.AppID),
				handler.NewCol(AppSAMLConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(AppSAMLConfigColumnEntityID, e.EntityID),
				handler.NewCol(AppSAMLConfigColumnMetadata, e.Metadata),
				handler.NewCol(AppSAMLConfigColumnMetadataURL, e.MetadataURL),
			},
			crdb.WithTableSuffix(appSAMLTableSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(AppColumnChangeDate, e.CreationDate()),
				handler.NewCol(AppColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(AppColumnID, e.AppID),
				handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *appProjection) reduceSAMLConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.SAMLConfigChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-GMHU2", "reduce.wrong.event.type")
	}

	cols := make([]handler.Column, 0, 3)
	if e.Metadata != nil {
		cols = append(cols, handler.NewCol(AppSAMLConfigColumnMetadata, e.Metadata))
	}
	if e.MetadataURL != nil {
		cols = append(cols, handler.NewCol(AppSAMLConfigColumnMetadataURL, *e.MetadataURL))
	}
	if e.EntityID != "" {
		cols = append(cols, handler.NewCol(AppSAMLConfigColumnEntityID, e.EntityID))
	}

	if len(cols) == 0 {
		return crdb.NewNoOpStatement(e), nil
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			cols,
			[]handler.Condition{
				handler.NewCond(AppSAMLConfigColumnAppID, e.AppID),
				handler.NewCond(AppSAMLConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(appSAMLTableSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(AppColumnChangeDate, e.CreationDate()),
				handler.NewCol(AppColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(AppColumnID, e.AppID),
				handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}
