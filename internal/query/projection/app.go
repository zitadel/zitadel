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
	"github.com/zitadel/zitadel/internal/repository/project"
)

type appProjection struct {
	crdb.StatementHandler
}

const (
	AppProjectionTable = "zitadel.projections.apps"
	AppAPITable        = AppProjectionTable + "_" + appAPITableSuffix
	AppOIDCTable       = AppProjectionTable + "_" + appOIDCTableSuffix
)

func newAppProjection(ctx context.Context, config crdb.StatementHandlerConfig) *appProjection {
	p := &appProjection{}
	config.ProjectionName = AppProjectionTable
	config.Reducers = p.reducers()
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
			},
		},
	}
}

const (
	AppColumnID            = "id"
	AppColumnName          = "name"
	AppColumnProjectID     = "project_id"
	AppColumnCreationDate  = "creation_date"
	AppColumnChangeDate    = "change_date"
	AppColumnResourceOwner = "resource_owner"
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

func (p *appProjection) reduceAppAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ApplicationAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-OzK4m", "seq", event.Sequence(), "expectedType", project.ApplicationAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-1xYE6", "reduce.wrong.event.type")
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
			handler.NewCol(AppColumnState, domain.AppStateActive),
			handler.NewCol(AppColumnSequence, e.Sequence()),
		},
	), nil
}

func (p *appProjection) reduceAppChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ApplicationChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-4Fjh2", "seq", event.Sequence(), "expectedType", project.ApplicationChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-ZJ8JA", "reduce.wrong.event.type")
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

func (p *appProjection) reduceAppDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ApplicationDeactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-hZ9to", "seq", event.Sequence(), "expectedType", project.ApplicationDeactivatedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-MVWxZ", "reduce.wrong.event.type")
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

func (p *appProjection) reduceAppReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ApplicationReactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-AbK3B", "seq", event.Sequence(), "expectedType", project.ApplicationReactivatedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-D0HZO", "reduce.wrong.event.type")
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

func (p *appProjection) reduceAppRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ApplicationRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-tdRId", "seq", event.Sequence(), "expectedType", project.ApplicationRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Y99aq", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(AppColumnID, e.AppID),
		},
	), nil
}

func (p *appProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-ZxQnj", "seq", event.Sequence(), "expectedType", project.ProjectRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-DlUlO", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(AppColumnProjectID, e.Aggregate().ID),
		},
	), nil
}

func (p *appProjection) reduceAPIConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.APIConfigAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-tdRId", "seq", event.Sequence(), "expectedType", project.APIConfigAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Y99aq", "reduce.wrong.event.type")
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

func (p *appProjection) reduceAPIConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.APIConfigChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-C6b4f", "seq", event.Sequence(), "expectedType", project.APIConfigChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-vnZKi", "reduce.wrong.event.type")
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

func (p *appProjection) reduceAPIConfigSecretChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.APIConfigSecretChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-dssSI", "seq", event.Sequence(), "expectedType", project.APIConfigSecretChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-ttb0I", "reduce.wrong.event.type")
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

func (p *appProjection) reduceOIDCConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.OIDCConfigAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-nlDQv", "seq", event.Sequence(), "expectedType", project.OIDCConfigAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-GNHU1", "reduce.wrong.event.type")
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

func (p *appProjection) reduceOIDCConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.OIDCConfigChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-nlDQv", "seq", event.Sequence(), "expectedType", project.OIDCConfigChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-GNHU1", "reduce.wrong.event.type")
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

func (p *appProjection) reduceOIDCConfigSecretChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.OIDCConfigSecretChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-nlDQv", "seq", event.Sequence(), "expectedType", project.OIDCConfigSecretChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-GNHU1", "reduce.wrong.event.type")
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
