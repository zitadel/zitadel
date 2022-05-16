package projection

import (
	"context"

	"github.com/lib/pq"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/project"
)

const (
	AppProjectionTable = "projections.apps"
	AppAPITable        = AppProjectionTable + "_" + appAPITableSuffix
	AppOIDCTable       = AppProjectionTable + "_" + appOIDCTableSuffix

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
	AppAPIConfigColumnClientID     = "client_id"
	AppAPIConfigColumnClientSecret = "client_secret"
	AppAPIConfigColumnAuthMethod   = "auth_method"

	appOIDCTableSuffix                          = "oidc_configs"
	AppOIDCConfigColumnAppID                    = "app_id"
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
)

type AppProjection struct {
	crdb.StatementHandler
}

func NewAppProjection(ctx context.Context, config crdb.StatementHandlerConfig) *AppProjection {
	p := new(AppProjection)
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
		},
			crdb.NewPrimaryKey(AppColumnInstanceID, ActionIDCol),
			crdb.WithIndex(crdb.NewIndex("project_id_idx", []string{AppColumnProjectID})),
			crdb.WithConstraint(crdb.NewConstraint("id_unique", []string{AppColumnID})),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(AppAPIConfigColumnAppID, crdb.ColumnTypeText, crdb.DeleteCascade(AppColumnID)),
			crdb.NewColumn(AppAPIConfigColumnClientID, crdb.ColumnTypeText),
			crdb.NewColumn(AppAPIConfigColumnClientSecret, crdb.ColumnTypeJSONB, crdb.Nullable()),
			crdb.NewColumn(AppAPIConfigColumnAuthMethod, crdb.ColumnTypeEnum),
		},
			crdb.NewPrimaryKey(AppAPIConfigColumnAppID),
			appAPITableSuffix,
			crdb.WithIndex(crdb.NewIndex("client_id_idx", []string{AppAPIConfigColumnClientID})),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(AppOIDCConfigColumnAppID, crdb.ColumnTypeText, crdb.DeleteCascade(AppColumnID)),
			crdb.NewColumn(AppOIDCConfigColumnVersion, crdb.ColumnTypeText),
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
		},
			crdb.NewPrimaryKey(AppOIDCConfigColumnAppID),
			appOIDCTableSuffix,
			crdb.WithIndex(crdb.NewIndex("client_id_idx", []string{AppOIDCConfigColumnClientID})),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *AppProjection) reducers() []handler.AggregateReducer {
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
			},
		},
	}
}

func (p *AppProjection) reduceAppAdded(event eventstore.Event) (*handler.Statement, error) {
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

func (p *AppProjection) reduceAppChanged(event eventstore.Event) (*handler.Statement, error) {
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
		},
	), nil
}

func (p *AppProjection) reduceAppDeactivated(event eventstore.Event) (*handler.Statement, error) {
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
		},
	), nil
}

func (p *AppProjection) reduceAppReactivated(event eventstore.Event) (*handler.Statement, error) {
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
		},
	), nil
}

func (p *AppProjection) reduceAppRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ApplicationRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Y99aq", "reduce.wrong.event.type %s", project.ApplicationRemovedType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(AppColumnID, e.AppID),
		},
	), nil
}

func (p *AppProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-DlUlO", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(AppColumnProjectID, e.Aggregate().ID),
		},
	), nil
}

func (p *AppProjection) reduceAPIConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.APIConfigAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Y99aq", "reduce.wrong.event.type %s", project.APIConfigAddedType)
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(AppAPIConfigColumnAppID, e.AppID),
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
			},
		),
	), nil
}

func (p *AppProjection) reduceAPIConfigChanged(event eventstore.Event) (*handler.Statement, error) {
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
			},
		),
	), nil
}

func (p *AppProjection) reduceAPIConfigSecretChanged(event eventstore.Event) (*handler.Statement, error) {
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
			},
		),
	), nil
}

func (p *AppProjection) reduceOIDCConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.OIDCConfigAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-GNHU1", "reduce.wrong.event.type %s", project.OIDCConfigAddedType)
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(AppOIDCConfigColumnAppID, e.AppID),
				handler.NewCol(AppOIDCConfigColumnVersion, e.Version),
				handler.NewCol(AppOIDCConfigColumnClientID, e.ClientID),
				handler.NewCol(AppOIDCConfigColumnClientSecret, e.ClientSecret),
				handler.NewCol(AppOIDCConfigColumnRedirectUris, pq.StringArray(e.RedirectUris)),
				handler.NewCol(AppOIDCConfigColumnResponseTypes, pq.Array(e.ResponseTypes)),
				handler.NewCol(AppOIDCConfigColumnGrantTypes, pq.Array(e.GrantTypes)),
				handler.NewCol(AppOIDCConfigColumnApplicationType, e.ApplicationType),
				handler.NewCol(AppOIDCConfigColumnAuthMethodType, e.AuthMethodType),
				handler.NewCol(AppOIDCConfigColumnPostLogoutRedirectUris, pq.StringArray(e.PostLogoutRedirectUris)),
				handler.NewCol(AppOIDCConfigColumnDevMode, e.DevMode),
				handler.NewCol(AppOIDCConfigColumnAccessTokenType, e.AccessTokenType),
				handler.NewCol(AppOIDCConfigColumnAccessTokenRoleAssertion, e.AccessTokenRoleAssertion),
				handler.NewCol(AppOIDCConfigColumnIDTokenRoleAssertion, e.IDTokenRoleAssertion),
				handler.NewCol(AppOIDCConfigColumnIDTokenUserinfoAssertion, e.IDTokenUserinfoAssertion),
				handler.NewCol(AppOIDCConfigColumnClockSkew, e.ClockSkew),
				handler.NewCol(AppOIDCConfigColumnAdditionalOrigins, pq.StringArray(e.AdditionalOrigins)),
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
			},
		),
	), nil
}

func (p *AppProjection) reduceOIDCConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.OIDCConfigChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-GNHU1", "reduce.wrong.event.type %s", project.OIDCConfigChangedType)
	}

	cols := make([]handler.Column, 0, 15)
	if e.Version != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnVersion, *e.Version))
	}
	if e.RedirectUris != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnRedirectUris, pq.StringArray(*e.RedirectUris)))
	}
	if e.ResponseTypes != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnResponseTypes, pq.Array(*e.ResponseTypes)))
	}
	if e.GrantTypes != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnGrantTypes, pq.Array(*e.GrantTypes)))
	}
	if e.ApplicationType != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnApplicationType, *e.ApplicationType))
	}
	if e.AuthMethodType != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnAuthMethodType, *e.AuthMethodType))
	}
	if e.PostLogoutRedirectUris != nil {
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnPostLogoutRedirectUris, pq.StringArray(*e.PostLogoutRedirectUris)))
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
		cols = append(cols, handler.NewCol(AppOIDCConfigColumnAdditionalOrigins, pq.StringArray(*e.AdditionalOrigins)))
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
			},
		),
	), nil
}

func (p *AppProjection) reduceOIDCConfigSecretChanged(event eventstore.Event) (*handler.Statement, error) {
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
			},
		),
	), nil
}
