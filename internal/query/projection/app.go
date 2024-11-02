package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	AppProjectionTable = "projections.apps7"
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
	AppOIDCConfigColumnBackChannelLogoutURI     = "back_channel_logout_uri"

	appSAMLTableSuffix             = "saml_configs"
	AppSAMLConfigColumnAppID       = "app_id"
	AppSAMLConfigColumnInstanceID  = "instance_id"
	AppSAMLConfigColumnEntityID    = "entity_id"
	AppSAMLConfigColumnMetadata    = "metadata"
	AppSAMLConfigColumnMetadataURL = "metadata_url"
)

type appProjection struct{}

func newAppProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(appProjection))
}

func (*appProjection) Name() string {
	return AppProjectionTable
}

func (*appProjection) Init() *old_handler.Check {
	return handler.NewMultiTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(AppColumnID, handler.ColumnTypeText),
			handler.NewColumn(AppColumnName, handler.ColumnTypeText),
			handler.NewColumn(AppColumnProjectID, handler.ColumnTypeText),
			handler.NewColumn(AppColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(AppColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(AppColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(AppColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(AppColumnState, handler.ColumnTypeEnum),
			handler.NewColumn(AppColumnSequence, handler.ColumnTypeInt64),
		},
			handler.NewPrimaryKey(AppColumnInstanceID, AppColumnID),
			handler.WithIndex(handler.NewIndex("project_id", []string{AppColumnProjectID})),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(AppAPIConfigColumnAppID, handler.ColumnTypeText),
			handler.NewColumn(AppAPIConfigColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(AppAPIConfigColumnClientID, handler.ColumnTypeText),
			handler.NewColumn(AppAPIConfigColumnClientSecret, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(AppAPIConfigColumnAuthMethod, handler.ColumnTypeEnum),
		},
			handler.NewPrimaryKey(AppAPIConfigColumnInstanceID, AppAPIConfigColumnAppID),
			appAPITableSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
			handler.WithIndex(handler.NewIndex("client_id", []string{AppAPIConfigColumnClientID})),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(AppOIDCConfigColumnAppID, handler.ColumnTypeText),
			handler.NewColumn(AppOIDCConfigColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(AppOIDCConfigColumnVersion, handler.ColumnTypeEnum),
			handler.NewColumn(AppOIDCConfigColumnClientID, handler.ColumnTypeText),
			handler.NewColumn(AppOIDCConfigColumnClientSecret, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(AppOIDCConfigColumnRedirectUris, handler.ColumnTypeTextArray, handler.Nullable()),
			handler.NewColumn(AppOIDCConfigColumnResponseTypes, handler.ColumnTypeEnumArray, handler.Nullable()),
			handler.NewColumn(AppOIDCConfigColumnGrantTypes, handler.ColumnTypeEnumArray, handler.Nullable()),
			handler.NewColumn(AppOIDCConfigColumnApplicationType, handler.ColumnTypeEnum),
			handler.NewColumn(AppOIDCConfigColumnAuthMethodType, handler.ColumnTypeEnum),
			handler.NewColumn(AppOIDCConfigColumnPostLogoutRedirectUris, handler.ColumnTypeTextArray, handler.Nullable()),
			handler.NewColumn(AppOIDCConfigColumnDevMode, handler.ColumnTypeBool),
			handler.NewColumn(AppOIDCConfigColumnAccessTokenType, handler.ColumnTypeEnum),
			handler.NewColumn(AppOIDCConfigColumnAccessTokenRoleAssertion, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(AppOIDCConfigColumnIDTokenRoleAssertion, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(AppOIDCConfigColumnIDTokenUserinfoAssertion, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(AppOIDCConfigColumnClockSkew, handler.ColumnTypeInt64, handler.Default(0)),
			handler.NewColumn(AppOIDCConfigColumnAdditionalOrigins, handler.ColumnTypeTextArray, handler.Nullable()),
			handler.NewColumn(AppOIDCConfigColumnSkipNativeAppSuccessPage, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(AppOIDCConfigColumnBackChannelLogoutURI, handler.ColumnTypeText, handler.Nullable()),
		},
			handler.NewPrimaryKey(AppOIDCConfigColumnInstanceID, AppOIDCConfigColumnAppID),
			appOIDCTableSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
			handler.WithIndex(handler.NewIndex("client_id", []string{AppOIDCConfigColumnClientID})),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(AppSAMLConfigColumnAppID, handler.ColumnTypeText),
			handler.NewColumn(AppSAMLConfigColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(AppSAMLConfigColumnEntityID, handler.ColumnTypeText),
			handler.NewColumn(AppSAMLConfigColumnMetadata, handler.ColumnTypeBytes),
			handler.NewColumn(AppSAMLConfigColumnMetadataURL, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(AppSAMLConfigColumnInstanceID, AppSAMLConfigColumnAppID),
			appSAMLTableSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
			handler.WithIndex(handler.NewIndex("entity_id", []string{AppSAMLConfigColumnEntityID})),
		),
	)
}

func (p *appProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
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
					Event:  project.APIConfigSecretHashUpdatedType,
					Reduce: p.reduceAPIConfigSecretHashUpdated,
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
					Event:  project.OIDCConfigSecretHashUpdatedType,
					Reduce: p.reduceOIDCConfigSecretHashUpdated,
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
			EventReducers: []handler.EventReducer{
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-1xYE6", "reduce.wrong.event.type %s", project.ApplicationAddedType)
	}
	return handler.NewCreateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ZJ8JA", "reduce.wrong.event.type %s", project.ApplicationChangedType)
	}
	if e.Name == "" {
		return handler.NewNoOpStatement(event), nil
	}
	return handler.NewUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-MVWxZ", "reduce.wrong.event.type %s", project.ApplicationDeactivatedType)
	}
	return handler.NewUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-D0HZO", "reduce.wrong.event.type %s", project.ApplicationReactivatedType)
	}
	return handler.NewUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y99aq", "reduce.wrong.event.type %s", project.ApplicationRemovedType)
	}
	return handler.NewDeleteStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-DlUlO", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}
	return handler.NewDeleteStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y99aq", "reduce.wrong.event.type %s", project.APIConfigAddedType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(AppAPIConfigColumnAppID, e.AppID),
				handler.NewCol(AppAPIConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(AppAPIConfigColumnClientID, e.ClientID),
				handler.NewCol(AppAPIConfigColumnClientSecret, crypto.SecretOrEncodedHash(e.ClientSecret, e.HashedSecret)),
				handler.NewCol(AppAPIConfigColumnAuthMethod, e.AuthMethodType),
			},
			handler.WithTableSuffix(appAPITableSuffix),
		),
		handler.AddUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-vnZKi", "reduce.wrong.event.type %s", project.APIConfigChangedType)
	}
	cols := make([]handler.Column, 0, 2)
	if e.AuthMethodType != nil {
		cols = append(cols, handler.NewCol(AppAPIConfigColumnAuthMethod, *e.AuthMethodType))
	}
	if len(cols) == 0 {
		return handler.NewNoOpStatement(e), nil
	}
	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			cols,
			[]handler.Condition{
				handler.NewCond(AppAPIConfigColumnAppID, e.AppID),
				handler.NewCond(AppAPIConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(appAPITableSuffix),
		),
		handler.AddUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ttb0I", "reduce.wrong.event.type %s", project.APIConfigSecretChangedType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(AppAPIConfigColumnClientSecret, crypto.SecretOrEncodedHash(e.ClientSecret, e.HashedSecret)),
			},
			[]handler.Condition{
				handler.NewCond(AppAPIConfigColumnAppID, e.AppID),
				handler.NewCond(AppAPIConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(appAPITableSuffix),
		),
		handler.AddUpdateStatement(
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

func (p *appProjection) reduceAPIConfigSecretHashUpdated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.APIConfigSecretHashUpdatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ttb0I", "reduce.wrong.event.type %s", project.APIConfigSecretHashUpdatedType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(AppAPIConfigColumnClientSecret, e.HashedSecret),
			},
			[]handler.Condition{
				handler.NewCond(AppAPIConfigColumnAppID, e.AppID),
				handler.NewCond(AppAPIConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(appAPITableSuffix),
		),
		handler.AddUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-GNHU1", "reduce.wrong.event.type %s", project.OIDCConfigAddedType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(AppOIDCConfigColumnAppID, e.AppID),
				handler.NewCol(AppOIDCConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(AppOIDCConfigColumnVersion, e.Version),
				handler.NewCol(AppOIDCConfigColumnClientID, e.ClientID),
				handler.NewCol(AppOIDCConfigColumnClientSecret, crypto.SecretOrEncodedHash(e.ClientSecret, e.HashedSecret)),
				handler.NewCol(AppOIDCConfigColumnRedirectUris, database.TextArray[string](e.RedirectUris)),
				handler.NewCol(AppOIDCConfigColumnResponseTypes, database.NumberArray[domain.OIDCResponseType](e.ResponseTypes)),
				handler.NewCol(AppOIDCConfigColumnGrantTypes, database.NumberArray[domain.OIDCGrantType](e.GrantTypes)),
				handler.NewCol(AppOIDCConfigColumnApplicationType, e.ApplicationType),
				handler.NewCol(AppOIDCConfigColumnAuthMethodType, e.AuthMethodType),
				handler.NewCol(AppOIDCConfigColumnPostLogoutRedirectUris, database.TextArray[string](e.PostLogoutRedirectUris)),
				handler.NewCol(AppOIDCConfigColumnDevMode, e.DevMode),
				handler.NewCol(AppOIDCConfigColumnAccessTokenType, e.AccessTokenType),
				handler.NewCol(AppOIDCConfigColumnAccessTokenRoleAssertion, e.AccessTokenRoleAssertion),
				handler.NewCol(AppOIDCConfigColumnIDTokenRoleAssertion, e.IDTokenRoleAssertion),
				handler.NewCol(AppOIDCConfigColumnIDTokenUserinfoAssertion, e.IDTokenUserinfoAssertion),
				handler.NewCol(AppOIDCConfigColumnClockSkew, e.ClockSkew),
				handler.NewCol(AppOIDCConfigColumnAdditionalOrigins, database.TextArray[string](e.AdditionalOrigins)),
				handler.NewCol(AppOIDCConfigColumnSkipNativeAppSuccessPage, e.SkipNativeAppSuccessPage),
				handler.NewCol(AppOIDCConfigColumnBackChannelLogoutURI, e.BackChannelLogoutURI),
			},
			handler.WithTableSuffix(appOIDCTableSuffix),
		),
		handler.AddUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-GNHU1", "reduce.wrong.event.type %s", project.OIDCConfigChangedType)
	}

	cols := make([]handler.Column, 0, 16)
	if e.Version != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnVersion, *e.Version))
	}
	if e.RedirectUris != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnRedirectUris, database.TextArray[string](*e.RedirectUris)))
	}
	if e.ResponseTypes != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnResponseTypes, database.NumberArray[domain.OIDCResponseType](*e.ResponseTypes)))
	}
	if e.GrantTypes != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnGrantTypes, database.NumberArray[domain.OIDCGrantType](*e.GrantTypes)))
	}
	if e.ApplicationType != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnApplicationType, *e.ApplicationType))
	}
	if e.AuthMethodType != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnAuthMethodType, *e.AuthMethodType))
	}
	if e.PostLogoutRedirectUris != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnPostLogoutRedirectUris, database.TextArray[string](*e.PostLogoutRedirectUris)))
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
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnAdditionalOrigins, database.TextArray[string](*e.AdditionalOrigins)))
	}
	if e.SkipNativeAppSuccessPage != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnSkipNativeAppSuccessPage, *e.SkipNativeAppSuccessPage))
	}
	if e.BackChannelLogoutURI != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnBackChannelLogoutURI, *e.BackChannelLogoutURI))
	}

	if len(cols) == 0 {
		return handler.NewNoOpStatement(e), nil
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			cols,
			[]handler.Condition{
				handler.NewCond(AppOIDCConfigColumnAppID, e.AppID),
				handler.NewCond(AppOIDCConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(appOIDCTableSuffix),
		),
		handler.AddUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-GNHU1", "reduce.wrong.event.type %s", project.OIDCConfigSecretChangedType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(AppOIDCConfigColumnClientSecret, crypto.SecretOrEncodedHash(e.ClientSecret, e.HashedSecret)),
			},
			[]handler.Condition{
				handler.NewCond(AppOIDCConfigColumnAppID, e.AppID),
				handler.NewCond(AppOIDCConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(appOIDCTableSuffix),
		),
		handler.AddUpdateStatement(
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

func (p *appProjection) reduceOIDCConfigSecretHashUpdated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.OIDCConfigSecretHashUpdatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-toSh1", "reduce.wrong.event.type %s", project.OIDCConfigSecretHashUpdatedType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(AppOIDCConfigColumnClientSecret, e.HashedSecret),
			},
			[]handler.Condition{
				handler.NewCond(AppOIDCConfigColumnAppID, e.AppID),
				handler.NewCond(AppOIDCConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(appOIDCTableSuffix),
		),
		handler.AddUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Hyd1f", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(AppColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(AppColumnResourceOwner, e.Aggregate().ID),
		},
	), nil
}

func (p *appProjection) reduceSAMLConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.SAMLConfigAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgument(nil, "HANDL-GMHU1", "reduce.wrong.event.type")
	}
	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(AppSAMLConfigColumnAppID, e.AppID),
				handler.NewCol(AppSAMLConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(AppSAMLConfigColumnEntityID, e.EntityID),
				handler.NewCol(AppSAMLConfigColumnMetadata, e.Metadata),
				handler.NewCol(AppSAMLConfigColumnMetadataURL, e.MetadataURL),
			},
			handler.WithTableSuffix(appSAMLTableSuffix),
		),
		handler.AddUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgument(nil, "HANDL-GMHU2", "reduce.wrong.event.type")
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
		return handler.NewNoOpStatement(e), nil
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			cols,
			[]handler.Condition{
				handler.NewCond(AppSAMLConfigColumnAppID, e.AppID),
				handler.NewCond(AppSAMLConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(appSAMLTableSuffix),
		),
		handler.AddUpdateStatement(
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
