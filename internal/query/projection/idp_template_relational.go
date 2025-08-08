package projection

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"

	db_domain "github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	IDPRelationalTable                = "zitadel.identity_providers"
	IDPRelationalOrgIdCol             = "org_id"
	IDPRelationalAutoRegisterCol      = "auto_register"
	IDPRelationalPayloadCol           = "payload"
	IDPRelationalOrgId                = "org_id"
	IDPRelationalAllowCreationCol     = "allow_creation"
	IDPRelationalAllowLinkingCol      = "allow_linking"
	IDPRelationalAllowAutoCreationCol = "allow_auto_creation"
	IDPRelationalAllowAutoUpdateCol   = "allow_auto_update"
	IDPRelationalAllowAutoLinkingCol  = "allow_auto_linking"
)

type idpTemplateRelationalProjection struct {
	idpRepo db_domain.IDProviderRepository
}

func newIDPTemplateRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	client := postgres.PGxPool(config.Client.Pool)
	idpRepo := repository.IDProviderRepository(client)
	return handler.NewHandler(ctx, &config, &idpTemplateRelationalProjection{
		idpRepo: idpRepo,
	})
}

func (*idpTemplateRelationalProjection) Name() string {
	return IDPRelationalTable
}

func (p *idpTemplateRelationalProjection) Reducers() []handler.AggregateReducer {
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
				{
					Event:  instance.OAuthIDPAddedEventType,
					Reduce: p.reduceOAuthIDPRelationalAdded,
				},
				{
					Event:  instance.OAuthIDPChangedEventType,
					Reduce: p.reduceOAuthIDPRelationalChanged,
				},
				{
					Event:  instance.OIDCIDPAddedEventType,
					Reduce: p.reduceOIDCIDPRelationalAdded,
				},
				{
					Event:  instance.OIDCIDPChangedEventType,
					Reduce: p.reduceOIDCIDPRelationalChanged,
				},
				{
					Event:  instance.OIDCIDPMigratedAzureADEventType,
					Reduce: p.reduceOIDCIDPRelationalMigratedAzureAD,
				},
				{
					Event:  instance.OIDCIDPMigratedGoogleEventType,
					Reduce: p.reduceOIDCIDPRelationalMigratedGoogle,
				},
				{
					Event:  instance.JWTIDPAddedEventType,
					Reduce: p.reduceJWTIDPRelationalAdded,
				},
				{
					Event:  instance.JWTIDPChangedEventType,
					Reduce: p.reduceJWTIDPRelationalChanged,
				},
				{
					Event:  instance.AzureADIDPAddedEventType,
					Reduce: p.reduceAzureADIDPRelationalAdded,
				},
				{
					Event:  instance.AzureADIDPChangedEventType,
					Reduce: p.reduceAzureADIDPRelationalChanged,
				},
				{
					Event:  instance.GitHubIDPAddedEventType,
					Reduce: p.reduceGitHubIDPRelationalAdded,
				},
				{
					Event:  instance.GitHubIDPChangedEventType,
					Reduce: p.reduceGitHubIDPRelationalChanged,
				},
				{
					Event:  instance.GitHubEnterpriseIDPAddedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPRelationalAdded,
				},
				{
					Event:  instance.GitHubEnterpriseIDPChangedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPRelationalChanged,
				},
				{
					Event:  instance.GitLabIDPAddedEventType,
					Reduce: p.reduceGitLabIDPRelationalAdded,
				},
				{
					Event:  instance.GitLabIDPChangedEventType,
					Reduce: p.reduceGitLabIDPRelationalChanged,
				},
				{
					Event:  instance.GitLabSelfHostedIDPAddedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPRelationalAdded,
				},
				{
					Event:  instance.GitLabSelfHostedIDPChangedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPRelationalChanged,
				},
				{
					Event:  instance.GoogleIDPAddedEventType,
					Reduce: p.reduceGoogleIDPRelationalAdded,
				},
				{
					Event:  instance.GoogleIDPChangedEventType,
					Reduce: p.reduceGoogleIDPRelationalChanged,
				},
				{
					Event:  instance.LDAPIDPAddedEventType,
					Reduce: p.reduceLDAPIDPAdded,
				},
				{
					Event:  instance.LDAPIDPChangedEventType,
					Reduce: p.reduceLDAPIDPChanged,
				},
				{
					Event:  instance.AppleIDPAddedEventType,
					Reduce: p.reduceAppleIDPAdded,
				},
				{
					Event:  instance.AppleIDPChangedEventType,
					Reduce: p.reduceAppleIDPChanged,
				},
				{
					Event:  instance.SAMLIDPAddedEventType,
					Reduce: p.reduceSAMLIDPAdded,
				},
				{
					Event:  instance.SAMLIDPChangedEventType,
					Reduce: p.reduceSAMLIDPChanged,
				},
				{
					Event:  instance.IDPRemovedEventType,
					Reduce: p.reduceIDPRemoved,
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
				{
					Event:  org.OAuthIDPAddedEventType,
					Reduce: p.reduceOAuthIDPRelationalAdded,
				},
				{
					Event:  org.OAuthIDPChangedEventType,
					Reduce: p.reduceOAuthIDPRelationalChanged,
				},
				{
					Event:  org.OIDCIDPAddedEventType,
					Reduce: p.reduceOIDCIDPRelationalAdded,
				},
				{
					Event:  org.OIDCIDPChangedEventType,
					Reduce: p.reduceOIDCIDPRelationalChanged,
				},
				{
					Event:  org.OIDCIDPMigratedAzureADEventType,
					Reduce: p.reduceOIDCIDPRelationalMigratedAzureAD,
				},
				{
					Event:  org.OIDCIDPMigratedGoogleEventType,
					Reduce: p.reduceOIDCIDPRelationalMigratedGoogle,
				},
				{
					Event:  org.JWTIDPAddedEventType,
					Reduce: p.reduceJWTIDPRelationalAdded,
				},
				{
					Event:  org.JWTIDPChangedEventType,
					Reduce: p.reduceJWTIDPRelationalChanged,
				},
				{
					Event:  org.AzureADIDPAddedEventType,
					Reduce: p.reduceAzureADIDPRelationalAdded,
				},
				{
					Event:  org.AzureADIDPChangedEventType,
					Reduce: p.reduceAzureADIDPRelationalChanged,
				},
				{
					Event:  org.GitHubIDPAddedEventType,
					Reduce: p.reduceGitHubIDPRelationalAdded,
				},
				{
					Event:  org.GitHubIDPChangedEventType,
					Reduce: p.reduceGitHubIDPRelationalChanged,
				},
				{
					Event:  org.GitHubEnterpriseIDPAddedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPRelationalAdded,
				},
				{
					Event:  org.GitHubEnterpriseIDPChangedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPRelationalChanged,
				},
				{
					Event:  org.GitLabIDPAddedEventType,
					Reduce: p.reduceGitLabIDPRelationalAdded,
				},
				{
					Event:  org.GitLabIDPChangedEventType,
					Reduce: p.reduceGitLabIDPRelationalChanged,
				},
				{
					Event:  org.GitLabSelfHostedIDPAddedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPRelationalAdded,
				},
				{
					Event:  org.GitLabSelfHostedIDPChangedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPRelationalChanged,
				},
				{
					Event:  org.GoogleIDPAddedEventType,
					Reduce: p.reduceGoogleIDPRelationalAdded,
				},
				{
					Event:  org.GoogleIDPChangedEventType,
					Reduce: p.reduceGoogleIDPRelationalChanged,
				},
				{
					Event:  org.LDAPIDPAddedEventType,
					Reduce: p.reduceLDAPIDPAdded,
				},
				{
					Event:  org.LDAPIDPChangedEventType,
					Reduce: p.reduceLDAPIDPChanged,
				},
				{
					Event:  org.AppleIDPAddedEventType,
					Reduce: p.reduceAppleIDPAdded,
				},
				{
					Event:  org.AppleIDPChangedEventType,
					Reduce: p.reduceAppleIDPChanged,
				},
				{
					Event:  org.SAMLIDPAddedEventType,
					Reduce: p.reduceSAMLIDPAdded,
				},
				{
					Event:  org.SAMLIDPChangedEventType,
					Reduce: p.reduceSAMLIDPChanged,
				},
				{
					Event:  org.IDPRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
			},
		},
	}
}

func (p *idpTemplateRelationalProjection) reduceIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
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

func (p *idpTemplateRelationalProjection) reduceIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
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

func (p *idpTemplateRelationalProjection) reduceIDRelationalPDeactivated(event eventstore.Event) (*handler.Statement, error) {
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

func (p *idpTemplateRelationalProjection) reduceIDPRelationalReactivated(event eventstore.Event) (*handler.Statement, error) {
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

func (p *idpTemplateRelationalProjection) reduceIDPRelationalRemoved(event eventstore.Event) (*handler.Statement, error) {
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

func (p *idpTemplateRelationalProjection) reduceOIDCRelationalConfigAdded(event eventstore.Event) (*handler.Statement, error) {
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

func (p *idpTemplateRelationalProjection) reduceOIDCRelationalConfigChanged(event eventstore.Event) (*handler.Statement, error) {
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

func (p *idpTemplateRelationalProjection) reduceJWTRelationalConfigAdded(event eventstore.Event) (*handler.Statement, error) {
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

func (p *idpTemplateRelationalProjection) reduceJWTRelationalConfigChanged(event eventstore.Event) (*handler.Statement, error) {
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

func (p *idpTemplateRelationalProjection) reduceOAuthIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.OAuthIDPAddedEvent
	switch e := event.(type) {
	case *org.OAuthIDPAddedEvent:
		idpEvent = e.OAuthIDPAddedEvent
	case *instance.OAuthIDPAddedEvent:
		idpEvent = e.OAuthIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Yap9ihb", "reduce.wrong.event.type %v", []eventstore.EventType{org.OAuthIDPAddedEventType, instance.OAuthIDPAddedEventType})
	}

	oauth := db_domain.OAuth{
		ClientID:              idpEvent.ClientID,
		ClientSecret:          idpEvent.ClientSecret,
		AuthorizationEndpoint: idpEvent.AuthorizationEndpoint,
		TokenEndpoint:         idpEvent.TokenEndpoint,
		UserEndpoint:          idpEvent.UserEndpoint,
		Scopes:                idpEvent.Scopes,
		IDAttribute:           idpEvent.IDAttribute,
		UsePKCE:               idpEvent.UsePKCE,
	}

	payload, err := json.Marshal(oauth)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPRelationalOrgId, orgId),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, db_domain.IDPStateActive.String()),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, db_domain.IDPTypeOAuth.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(idpEvent.AutoLinkingOption).String()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
				handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceOAuthIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.OAuthIDPChangedEvent
	switch e := event.(type) {
	case *org.OAuthIDPChangedEvent:
		idpEvent = e.OAuthIDPChangedEvent
	case *instance.OAuthIDPChangedEvent:
		idpEvent = e.OAuthIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OAuthIDPChangedEventType, instance.OAuthIDPChangedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	oauth, err := p.idpRepo.GetOAuth(context.Background(), p.idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges, &columns)

	payload := &oauth.OAuth
	payloadChanged := reduceOAuthIDPRelationalChangedColumns(payload, &idpEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCond(IDPRelationalOrgId, orgId),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceOIDCIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.OIDCIDPAddedEvent
	switch e := event.(type) {
	case *org.OIDCIDPAddedEvent:
		idpEvent = e.OIDCIDPAddedEvent
	case *instance.OIDCIDPAddedEvent:
		idpEvent = e.OIDCIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Ys02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPAddedEventType, instance.OIDCIDPAddedEventType})
	}

	payload, err := json.Marshal(idpEvent)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPRelationalOrgId, orgId),
				handler.NewCol(IDPTemplateStateCol, db_domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, db_domain.IDPTypeOIDC.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(idpEvent.AutoLinkingOption).String()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
				handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceOIDCIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.OIDCIDPChangedEvent
	switch e := event.(type) {
	case *org.OIDCIDPChangedEvent:
		idpEvent = e.OIDCIDPChangedEvent
	case *instance.OIDCIDPChangedEvent:
		idpEvent = e.OIDCIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPChangedEventType, instance.OIDCIDPChangedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	oidc, err := p.idpRepo.GetOIDC(context.Background(), p.idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges, &columns)

	payload := &oidc.OIDC
	payloadChanged := reduceOIDCIDPRelationalChangedColumns(payload, &idpEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCond(IDPRelationalOrgId, orgId),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceOIDCIDPRelationalMigratedAzureAD(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.OIDCIDPMigratedAzureADEvent
	switch e := event.(type) {
	case *org.OIDCIDPMigratedAzureADEvent:
		idpEvent = e.OIDCIDPMigratedAzureADEvent
	case *instance.OIDCIDPMigratedAzureADEvent:
		idpEvent = e.OIDCIDPMigratedAzureADEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPMigratedAzureADEventType, instance.OIDCIDPMigratedAzureADEventType})
	}

	azure := db_domain.Azure{
		ClientID:        idpEvent.ClientID,
		ClientSecret:    idpEvent.ClientSecret,
		Scopes:          idpEvent.Scopes,
		Tenant:          idpEvent.Tenant,
		IsEmailVerified: idpEvent.IsEmailVerified,
	}

	payload, err := json.Marshal(azure)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, db_domain.IDPTypeAzure.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(idpEvent.AutoLinkingOption).String()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCond(IDPRelationalOrgId, orgId),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceOIDCIDPRelationalMigratedGoogle(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.OIDCIDPMigratedGoogleEvent
	switch e := event.(type) {
	case *org.OIDCIDPMigratedGoogleEvent:
		idpEvent = e.OIDCIDPMigratedGoogleEvent
	case *instance.OIDCIDPMigratedGoogleEvent:
		idpEvent = e.OIDCIDPMigratedGoogleEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPMigratedGoogleEventType, instance.OIDCIDPMigratedGoogleEventType})
	}

	google := db_domain.Google{
		ClientID:     idpEvent.ClientID,
		ClientSecret: idpEvent.ClientSecret,
		Scopes:       idpEvent.Scopes,
	}

	payload, err := json.Marshal(google)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, db_domain.IDPTypeGoogle.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(idpEvent.AutoLinkingOption).String()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCond(IDPRelationalOrgId, orgId),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceJWTIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.JWTIDPAddedEvent
	switch e := event.(type) {
	case *org.JWTIDPAddedEvent:
		idpEvent = e.JWTIDPAddedEvent
	case *instance.JWTIDPAddedEvent:
		idpEvent = e.JWTIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Yopi2s", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPAddedEventType, instance.JWTIDPAddedEventType})
	}

	jwt := db_domain.JWT{
		JWTEndpoint:  idpEvent.JWTEndpoint,
		Issuer:       idpEvent.Issuer,
		KeysEndpoint: idpEvent.KeysEndpoint,
		HeaderName:   idpEvent.HeaderName,
	}

	payload, err := json.Marshal(jwt)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPRelationalOrgId, orgId),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateStateCol, db_domain.IDPStateActive.String()),
				handler.NewCol(IDPTemplateTypeCol, db_domain.IDPTypeJWT.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(idpEvent.AutoLinkingOption).String()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
				handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceJWTIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.JWTIDPChangedEvent
	switch e := event.(type) {
	case *org.JWTIDPChangedEvent:
		idpEvent = e.JWTIDPChangedEvent
	case *instance.JWTIDPChangedEvent:
		idpEvent = e.JWTIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPChangedEventType, instance.JWTIDPChangedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	jwt, err := p.idpRepo.GetJWT(context.Background(), p.idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges, &columns)

	payload := &jwt.JWT
	payloadChanged := reduceJWTIDPRelationalChangedColumns(payload, &idpEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCond(IDPRelationalOrgId, orgId),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceAzureADIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.AzureADIDPAddedEvent
	switch e := event.(type) {
	case *org.AzureADIDPAddedEvent:
		idpEvent = e.AzureADIDPAddedEvent
	case *instance.AzureADIDPAddedEvent:
		idpEvent = e.AzureADIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.AzureADIDPAddedEventType, instance.AzureADIDPAddedEventType})
	}

	azure := db_domain.Azure{
		ClientID:        idpEvent.ClientID,
		ClientSecret:    idpEvent.ClientSecret,
		Scopes:          idpEvent.Scopes,
		Tenant:          idpEvent.Tenant,
		IsEmailVerified: idpEvent.IsEmailVerified,
	}

	payload, err := json.Marshal(azure)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPRelationalOrgId, orgId),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, db_domain.IDPTypeAzure.String()),
				handler.NewCol(IDPTemplateStateCol, db_domain.IDPStateActive.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(idpEvent.AutoLinkingOption).String()),
				handler.NewCol(CreatedAt, idpEvent.CreationDate()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceAzureADIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.AzureADIDPChangedEvent
	switch e := event.(type) {
	case *org.AzureADIDPChangedEvent:
		idpEvent = e.AzureADIDPChangedEvent
	case *instance.AzureADIDPChangedEvent:
		idpEvent = e.AzureADIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.AzureADIDPChangedEventType, instance.AzureADIDPChangedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	oauth, err := p.idpRepo.GetOAzureAD(context.Background(), p.idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges, &columns)

	payload := &oauth.Azure
	payloadChanged := reduceAzureADIDPRelationalChangedColumns(payload, &idpEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCond(IDPRelationalOrgId, orgId),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceGitHubIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitHubIDPAddedEvent
	switch e := event.(type) {
	case *org.GitHubIDPAddedEvent:
		idpEvent = e.GitHubIDPAddedEvent
	case *instance.GitHubIDPAddedEvent:
		idpEvent = e.GitHubIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPAddedEventType, instance.GitHubIDPAddedEventType})
	}

	github := db_domain.Github{
		ClientID:     idpEvent.ClientID,
		ClientSecret: idpEvent.ClientSecret,
		Scopes:       idpEvent.Scopes,
	}

	payload, err := json.Marshal(github)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPRelationalOrgId, orgId),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, db_domain.IDPTypeGitHub.String()),
				handler.NewCol(IDPTemplateStateCol, db_domain.IDPStateActive.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(idpEvent.AutoLinkingOption).String()),
				handler.NewCol(CreatedAt, idpEvent.CreationDate()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceGitHubIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitHubIDPChangedEvent
	switch e := event.(type) {
	case *org.GitHubIDPChangedEvent:
		idpEvent = e.GitHubIDPChangedEvent
	case *instance.GitHubIDPChangedEvent:
		idpEvent = e.GitHubIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPChangedEventType, instance.GitHubIDPChangedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	github, err := p.idpRepo.GetGithub(context.Background(), p.idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges, &columns)

	payload := &github.Github
	payloadChanged := reduceGitHubIDPRelationalChangedColumns(payload, &idpEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCond(IDPRelationalOrgId, orgId),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceGitHubEnterpriseIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitHubEnterpriseIDPAddedEvent
	switch e := event.(type) {
	case *org.GitHubEnterpriseIDPAddedEvent:
		idpEvent = e.GitHubEnterpriseIDPAddedEvent
	case *instance.GitHubEnterpriseIDPAddedEvent:
		idpEvent = e.GitHubEnterpriseIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Yf3g2a", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubEnterpriseIDPAddedEventType, instance.GitHubEnterpriseIDPAddedEventType})
	}

	githubEnterprise := db_domain.GithubEnterprise{
		ClientID:              idpEvent.ClientID,
		ClientSecret:          idpEvent.ClientSecret,
		AuthorizationEndpoint: idpEvent.AuthorizationEndpoint,
		TokenEndpoint:         idpEvent.TokenEndpoint,
		UserEndpoint:          idpEvent.UserEndpoint,
		Scopes:                idpEvent.Scopes,
	}

	payload, err := json.Marshal(githubEnterprise)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPRelationalOrgId, orgId),
				handler.NewCol(IDPTemplateStateCol, db_domain.IDPStateActive.String()),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, db_domain.IDPTypeGitHubEnterprise.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(idpEvent.AutoLinkingOption).String()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
				handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceGitHubEnterpriseIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitHubEnterpriseIDPChangedEvent
	switch e := event.(type) {
	case *org.GitHubEnterpriseIDPChangedEvent:
		idpEvent = e.GitHubEnterpriseIDPChangedEvent
	case *instance.GitHubEnterpriseIDPChangedEvent:
		idpEvent = e.GitHubEnterpriseIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YDg3g", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubEnterpriseIDPChangedEventType, instance.GitHubEnterpriseIDPChangedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	githubEnterprise, err := p.idpRepo.GetGithubEnterprise(context.Background(), p.idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges, &columns)

	payload := &githubEnterprise.GithubEnterprise
	payloadChanged := reduceGitHubEnterpriseIDPRelationalChangedColumns(payload, &idpEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCond(IDPRelationalOrgId, orgId),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceGitLabIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitLabIDPAddedEvent
	switch e := event.(type) {
	case *org.GitLabIDPAddedEvent:
		idpEvent = e.GitLabIDPAddedEvent
	case *instance.GitLabIDPAddedEvent:
		idpEvent = e.GitLabIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabIDPAddedEventType, instance.GitLabIDPAddedEventType})
	}

	gitlab := db_domain.Gitlab{
		ClientID:     idpEvent.ClientID,
		ClientSecret: idpEvent.ClientSecret,
		Scopes:       idpEvent.Scopes,
	}

	payload, err := json.Marshal(gitlab)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPRelationalOrgId, orgId),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, db_domain.IDPTypeGitLab.String()),
				handler.NewCol(IDPTemplateStateCol, db_domain.IDPStateActive.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(idpEvent.AutoLinkingOption).String()),
				handler.NewCol(CreatedAt, idpEvent.CreationDate()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceGitLabIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitLabIDPChangedEvent
	switch e := event.(type) {
	case *org.GitLabIDPChangedEvent:
		idpEvent = e.GitLabIDPChangedEvent
	case *instance.GitLabIDPChangedEvent:
		idpEvent = e.GitLabIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabIDPChangedEventType, instance.GitLabIDPChangedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	oauth, err := p.idpRepo.GetGitlab(context.Background(), p.idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges, &columns)

	payload := &oauth.Gitlab
	payloadChanged := reduceGitLabIDPRelationalChangedColumns(payload, &idpEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCond(IDPRelationalOrgId, orgId),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceGitLabSelfHostedIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitLabSelfHostedIDPAddedEvent
	switch e := event.(type) {
	case *org.GitLabSelfHostedIDPAddedEvent:
		idpEvent = e.GitLabSelfHostedIDPAddedEvent
	case *instance.GitLabSelfHostedIDPAddedEvent:
		idpEvent = e.GitLabSelfHostedIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YAF3gw", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabSelfHostedIDPAddedEventType, instance.GitLabSelfHostedIDPAddedEventType})
	}

	gitlabSelfHosting := db_domain.GitlabSelfHosting{
		Issuer:       idpEvent.Issuer,
		ClientID:     idpEvent.ClientID,
		ClientSecret: idpEvent.ClientSecret,
		Scopes:       idpEvent.Scopes,
	}

	payload, err := json.Marshal(gitlabSelfHosting)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPRelationalOrgId, orgId),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, db_domain.IDPTypeGitLabSelfHosted.String()),
				handler.NewCol(IDPTemplateStateCol, db_domain.IDPStateActive.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(idpEvent.AutoLinkingOption).String()),
				handler.NewCol(CreatedAt, idpEvent.CreationDate()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceGitLabSelfHostedIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitLabSelfHostedIDPChangedEvent
	switch e := event.(type) {
	case *org.GitLabSelfHostedIDPChangedEvent:
		idpEvent = e.GitLabSelfHostedIDPChangedEvent
	case *instance.GitLabSelfHostedIDPChangedEvent:
		idpEvent = e.GitLabSelfHostedIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YAf3g2", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabSelfHostedIDPChangedEventType, instance.GitLabSelfHostedIDPChangedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	gitlabSelfHosted, err := p.idpRepo.GetGitlabSelfHosting(context.Background(), p.idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges, &columns)

	payload := &gitlabSelfHosted.GitlabSelfHosting
	payloadChanged := reduceGitLabSelfHostedIDPRelationalChangedColumns(payload, &idpEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCond(IDPRelationalOrgId, orgId),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceGoogleIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GoogleIDPAddedEvent
	switch e := event.(type) {
	case *org.GoogleIDPAddedEvent:
		idpEvent = e.GoogleIDPAddedEvent
	case *instance.GoogleIDPAddedEvent:
		idpEvent = e.GoogleIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Yp9ihb", "reduce.wrong.event.type %v", []eventstore.EventType{org.GoogleIDPAddedEventType, instance.GoogleIDPAddedEventType})
	}

	google := db_domain.Google{
		ClientID:     idpEvent.ClientID,
		ClientSecret: idpEvent.ClientSecret,
		Scopes:       idpEvent.Scopes,
	}

	payload, err := json.Marshal(google)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}
	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPRelationalOrgId, orgId),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, db_domain.IDPTypeGoogle.String()),
				handler.NewCol(IDPTemplateStateCol, db_domain.IDPStateActive.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(idpEvent.AutoLinkingOption).String()),
				handler.NewCol(CreatedAt, idpEvent.CreationDate()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceGoogleIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GoogleIDPChangedEvent
	switch e := event.(type) {
	case *org.GoogleIDPChangedEvent:
		idpEvent = e.GoogleIDPChangedEvent
	case *instance.GoogleIDPChangedEvent:
		idpEvent = e.GoogleIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GoogleIDPChangedEventType, instance.GoogleIDPChangedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	oauth, err := p.idpRepo.GetGoogle(context.Background(), p.idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges, &columns)

	payload := &oauth.Google
	payloadChanged := reduceGoogleIDPRelationalChangedColumns(payload, &idpEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCond(IDPRelationalOrgId, orgId),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceLDAPIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.LDAPIDPAddedEvent
	switch e := event.(type) {
	case *org.LDAPIDPAddedEvent:
		idpEvent = e.LDAPIDPAddedEvent
	case *instance.LDAPIDPAddedEvent:
		idpEvent = e.LDAPIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-9s02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPAddedEventType, instance.LDAPIDPAddedEventType})
	}

	ldap := db_domain.LDAP{
		Servers:           idpEvent.Servers,
		StartTLS:          idpEvent.StartTLS,
		BaseDN:            idpEvent.BaseDN,
		BindDN:            idpEvent.BindDN,
		BindPassword:      idpEvent.BindPassword,
		UserBase:          idpEvent.UserBase,
		UserObjectClasses: idpEvent.UserObjectClasses,
		UserFilters:       idpEvent.UserFilters,
		Timeout:           idpEvent.Timeout,
		LDAPAttributes: db_domain.LDAPAttributes{
			IDAttribute:                idpEvent.IDAttribute,
			FirstNameAttribute:         idpEvent.FirstNameAttribute,
			LastNameAttribute:          idpEvent.LastNameAttribute,
			DisplayNameAttribute:       idpEvent.DisplayNameAttribute,
			NickNameAttribute:          idpEvent.NickNameAttribute,
			PreferredUsernameAttribute: idpEvent.PreferredUsernameAttribute,
			EmailAttribute:             idpEvent.EmailAttribute,
			EmailVerifiedAttribute:     idpEvent.EmailVerifiedAttribute,
			PhoneAttribute:             idpEvent.PhoneAttribute,
			PhoneVerifiedAttribute:     idpEvent.PhoneVerifiedAttribute,
			PreferredLanguageAttribute: idpEvent.PreferredLanguageAttribute,
			AvatarURLAttribute:         idpEvent.AvatarURLAttribute,
			ProfileAttribute:           idpEvent.ProfileAttribute,
		},
	}

	payload, err := json.Marshal(ldap)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPRelationalOrgId, orgId),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, db_domain.IDPTypeLDAP.String()),
				handler.NewCol(IDPTemplateStateCol, db_domain.IDPStateActive.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(idpEvent.AutoLinkingOption).String()),
				handler.NewCol(CreatedAt, idpEvent.CreationDate()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceLDAPIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.LDAPIDPChangedEvent
	switch e := event.(type) {
	case *org.LDAPIDPChangedEvent:
		idpEvent = e.LDAPIDPChangedEvent
	case *instance.LDAPIDPChangedEvent:
		idpEvent = e.LDAPIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPChangedEventType, instance.LDAPIDPChangedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	oauth, err := p.idpRepo.GetLDAP(context.Background(), p.idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges, &columns)

	payload := &oauth.LDAP
	payloadChanged := reduceLDAPIDPRelationalChangedColumns(payload, &idpEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCond(IDPRelationalOrgId, orgId),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceAppleIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.AppleIDPAddedEvent
	switch e := event.(type) {
	case *org.AppleIDPAddedEvent:
		idpEvent = e.AppleIDPAddedEvent
	case *instance.AppleIDPAddedEvent:
		idpEvent = e.AppleIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YFvg3", "reduce.wrong.event.type %v", []eventstore.EventType{org.AppleIDPAddedEventType /*, instance.AppleIDPAddedEventType*/})
	}

	apple := db_domain.Apple{
		ClientID:   idpEvent.ClientID,
		TeamID:     idpEvent.TeamID,
		KeyID:      idpEvent.KeyID,
		PrivateKey: idpEvent.PrivateKey,
		Scopes:     idpEvent.Scopes,
	}

	payload, err := json.Marshal(apple)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPRelationalOrgId, orgId),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, db_domain.IDPTypeApple.String()),
				handler.NewCol(IDPTemplateStateCol, db_domain.IDPStateActive.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(idpEvent.AutoLinkingOption).String()),
				handler.NewCol(CreatedAt, idpEvent.CreationDate()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceAppleIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.AppleIDPChangedEvent
	switch e := event.(type) {
	case *org.AppleIDPChangedEvent:
		idpEvent = e.AppleIDPChangedEvent
	case *instance.AppleIDPChangedEvent:
		idpEvent = e.AppleIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YBez3", "reduce.wrong.event.type %v", []eventstore.EventType{org.AppleIDPChangedEventType /*, instance.AppleIDPChangedEventType*/})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	apple, err := p.idpRepo.GetApple(context.Background(), p.idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges, &columns)

	payload := &apple.Apple
	payloadChanged := reduceAppleIDPRelationalChangedColumns(payload, &idpEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCond(IDPRelationalOrgId, orgId),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceSAMLIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.SAMLIDPAddedEvent
	switch e := event.(type) {
	case *org.SAMLIDPAddedEvent:
		idpEvent = e.SAMLIDPAddedEvent
	case *instance.SAMLIDPAddedEvent:
		idpEvent = e.SAMLIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Ys02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.SAMLIDPAddedEventType, instance.SAMLIDPAddedEventType})
	}

	saml := db_domain.SAML{
		Metadata:                      idpEvent.Metadata,
		Key:                           idpEvent.Key,
		Certificate:                   idpEvent.Certificate,
		Binding:                       idpEvent.Binding,
		WithSignedRequest:             idpEvent.WithSignedRequest,
		NameIDFormat:                  idpEvent.NameIDFormat,
		TransientMappingAttributeName: idpEvent.TransientMappingAttributeName,
		FederatedLogoutEnabled:        idpEvent.FederatedLogoutEnabled,
	}

	payload, err := json.Marshal(saml)
	if err != nil {
		return nil, err
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPRelationalOrgId, orgId),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, db_domain.IDPTypeSAML.String()),
				handler.NewCol(IDPTemplateStateCol, db_domain.IDPStateActive.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(idpEvent.AutoLinkingOption).String()),
				handler.NewCol(CreatedAt, idpEvent.CreationDate()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceSAMLIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.SAMLIDPChangedEvent
	switch e := event.(type) {
	case *org.SAMLIDPChangedEvent:
		idpEvent = e.SAMLIDPChangedEvent
	case *instance.SAMLIDPChangedEvent:
		idpEvent = e.SAMLIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y7c0fii4ad", "reduce.wrong.event.type %v", []eventstore.EventType{org.SAMLIDPChangedEventType, instance.SAMLIDPChangedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	saml, err := p.idpRepo.GetSAML(context.Background(), p.idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges, &columns)

	payload := &saml.SAML
	payloadChanged := reduceSAMLIDPRelationalChangedColumns(payload, &idpEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCond(IDPRelationalOrgId, orgId),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceIDPRemoved(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.RemovedEvent
	switch e := event.(type) {
	case *org.IDPRemovedEvent:
		idpEvent = e.RemovedEvent
	case *instance.IDPRemovedEvent:
		idpEvent = e.RemovedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Ybcvwin2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPRemovedEventType, instance.IDPRemovedEventType})
	}

	var orgId *string
	if idpEvent.Aggregate().ResourceOwner != idpEvent.Agg.InstanceID {
		orgId = &idpEvent.Aggregate().ResourceOwner
	}

	return handler.NewDeleteStatement(
		&idpEvent,
		[]handler.Condition{
			handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCond(IDPRelationalOrgId, orgId),
		},
	), nil
}

func reduceIDPRelationalChangedTemplateColumns(name *string, optionChanges idp.OptionChanges, cols *[]handler.Column) {
	if name != nil {
		*cols = append(*cols, handler.NewCol(IDPTemplateNameCol, *name))
	}
	if optionChanges.IsCreationAllowed != nil {
		*cols = append(*cols, handler.NewCol(IDPRelationalAllowCreationCol, *optionChanges.IsCreationAllowed))
	}
	if optionChanges.IsLinkingAllowed != nil {
		*cols = append(*cols, handler.NewCol(IDPRelationalAllowLinkingCol, *optionChanges.IsLinkingAllowed))
	}
	if optionChanges.IsAutoCreation != nil {
		*cols = append(*cols, handler.NewCol(IDPRelationalAllowAutoCreationCol, *optionChanges.IsAutoCreation))
	}
	if optionChanges.IsAutoUpdate != nil {
		*cols = append(*cols, handler.NewCol(IDPRelationalAllowAutoUpdateCol, *optionChanges.IsAutoUpdate))
	}
	if optionChanges.AutoLinkingOption != nil {
		*cols = append(*cols, handler.NewCol(IDPRelationalAllowAutoLinkingCol, db_domain.IDPAutoLinkingOption(*optionChanges.AutoLinkingOption).String()))
	}
}

func reduceOAuthIDPRelationalChangedColumns(payload *db_domain.OAuth, idpEvent *idp.OAuthIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.AuthorizationEndpoint != nil {
		payloadChange = true
		payload.AuthorizationEndpoint = *idpEvent.AuthorizationEndpoint
	}
	if idpEvent.TokenEndpoint != nil {
		payloadChange = true
		payload.TokenEndpoint = *idpEvent.TokenEndpoint
	}
	if idpEvent.UserEndpoint != nil {
		payloadChange = true
		payload.UserEndpoint = *idpEvent.UserEndpoint
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	if idpEvent.IDAttribute != nil {
		payloadChange = true
		payload.IDAttribute = *idpEvent.IDAttribute
	}
	if idpEvent.UsePKCE != nil {
		payloadChange = true
		payload.UsePKCE = *idpEvent.UsePKCE
	}
	return payloadChange
}

func reduceOIDCIDPRelationalChangedColumns(payload *db_domain.OIDC, idpEvent *idp.OIDCIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = *idpEvent.ClientSecret
	}
	if idpEvent.Issuer != nil {
		payloadChange = true
		payload.Issuer = *idpEvent.Issuer
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	if idpEvent.IsIDTokenMapping != nil {
		payloadChange = true
		payload.IsIDTokenMapping = *idpEvent.IsIDTokenMapping
	}
	if idpEvent.UsePKCE != nil {
		payloadChange = true
		payload.UsePKCE = *idpEvent.UsePKCE
	}
	return payloadChange
}

func reduceJWTIDPRelationalChangedColumns(payload *db_domain.JWT, idpEvent *idp.JWTIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.JWTEndpoint != nil {
		payloadChange = true
		payload.JWTEndpoint = *idpEvent.JWTEndpoint
	}
	if idpEvent.KeysEndpoint != nil {
		payloadChange = true
		payload.KeysEndpoint = *idpEvent.KeysEndpoint
	}
	if idpEvent.HeaderName != nil {
		payloadChange = true
		payload.HeaderName = *idpEvent.HeaderName
	}
	if idpEvent.Issuer != nil {
		payloadChange = true
		payload.Issuer = *idpEvent.Issuer
	}
	return payloadChange
}

func reduceAzureADIDPRelationalChangedColumns(payload *db_domain.Azure, idpEvent *idp.AzureADIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	if idpEvent.Tenant != nil {
		payloadChange = true
		payload.Tenant = *idpEvent.Tenant
	}
	if idpEvent.IsEmailVerified != nil {
		payloadChange = true
		payload.IsEmailVerified = *idpEvent.IsEmailVerified
	}
	return payloadChange
}

func reduceGitHubIDPRelationalChangedColumns(payload *db_domain.Github, idpEvent *idp.GitHubIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChange
}

func reduceGitHubEnterpriseIDPRelationalChangedColumns(payload *db_domain.GithubEnterprise, idpEvent *idp.GitHubEnterpriseIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.AuthorizationEndpoint != nil {
		payloadChange = true
		payload.AuthorizationEndpoint = *idpEvent.AuthorizationEndpoint
	}
	if idpEvent.TokenEndpoint != nil {
		payloadChange = true
		payload.TokenEndpoint = *idpEvent.TokenEndpoint
	}
	if idpEvent.UserEndpoint != nil {
		payloadChange = true
		payload.UserEndpoint = *idpEvent.UserEndpoint
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChange
}

func reduceGitLabIDPRelationalChangedColumns(payload *db_domain.Gitlab, idpEvent *idp.GitLabIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChange
}

func reduceGitLabSelfHostedIDPRelationalChangedColumns(payload *db_domain.GitlabSelfHosting, idpEvent *idp.GitLabSelfHostedIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.Issuer != nil {
		payloadChange = true
		payload.Issuer = *idpEvent.Issuer
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChange
}

func reduceGoogleIDPRelationalChangedColumns(payload *db_domain.Google, idpEvent *idp.GoogleIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChange
}

func reduceLDAPIDPRelationalChangedColumns(payload *db_domain.LDAP, idpEvent *idp.LDAPIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.Servers != nil {
		payloadChange = true
		payload.Servers = idpEvent.Servers
	}
	if idpEvent.StartTLS != nil {
		payloadChange = true
		payload.StartTLS = *idpEvent.StartTLS
	}
	if idpEvent.BaseDN != nil {
		payloadChange = true
		payload.BaseDN = *idpEvent.BaseDN
	}
	if idpEvent.BindDN != nil {
		payloadChange = true
		payload.BindDN = *idpEvent.BindDN
	}
	if idpEvent.BindPassword != nil {
		payloadChange = true
		payload.BindPassword = idpEvent.BindPassword
	}
	if idpEvent.UserBase != nil {
		payloadChange = true
		payload.UserBase = *idpEvent.UserBase
	}
	if idpEvent.UserObjectClasses != nil {
		payloadChange = true
		payload.UserObjectClasses = idpEvent.UserObjectClasses
	}
	if idpEvent.UserFilters != nil {
		payloadChange = true
		payload.UserFilters = idpEvent.UserFilters
	}
	if idpEvent.Timeout != nil {
		payloadChange = true
		payload.Timeout = *idpEvent.Timeout
	}
	if idpEvent.RootCA != nil {
		payloadChange = true
		payload.RootCA = idpEvent.RootCA
	}
	if idpEvent.IDAttribute != nil {
		payloadChange = true
		payload.IDAttribute = *idpEvent.IDAttribute
	}
	if idpEvent.FirstNameAttribute != nil {
		payloadChange = true
		payload.FirstNameAttribute = *idpEvent.FirstNameAttribute
	}
	if idpEvent.LastNameAttribute != nil {
		payloadChange = true
		payload.LastNameAttribute = *idpEvent.LastNameAttribute
	}
	if idpEvent.DisplayNameAttribute != nil {
		payloadChange = true
		payload.DisplayNameAttribute = *idpEvent.DisplayNameAttribute
	}
	if idpEvent.NickNameAttribute != nil {
		payloadChange = true
		payload.NickNameAttribute = *idpEvent.NickNameAttribute
	}
	if idpEvent.PreferredUsernameAttribute != nil {
		payloadChange = true
		payload.PreferredUsernameAttribute = *idpEvent.PreferredUsernameAttribute
	}
	if idpEvent.EmailAttribute != nil {
		payloadChange = true
		payload.EmailAttribute = *idpEvent.EmailAttribute
	}
	if idpEvent.EmailVerifiedAttribute != nil {
		payloadChange = true
		payload.EmailVerifiedAttribute = *idpEvent.EmailVerifiedAttribute
	}
	if idpEvent.PhoneAttribute != nil {
		payloadChange = true
		payload.PhoneAttribute = *idpEvent.PhoneAttribute
	}
	if idpEvent.PhoneVerifiedAttribute != nil {
		payloadChange = true
		payload.PhoneVerifiedAttribute = *idpEvent.PhoneVerifiedAttribute
	}
	if idpEvent.PreferredLanguageAttribute != nil {
		payloadChange = true
		payload.PreferredLanguageAttribute = *idpEvent.PreferredLanguageAttribute
	}
	if idpEvent.AvatarURLAttribute != nil {
		payloadChange = true
		payload.AvatarURLAttribute = *idpEvent.AvatarURLAttribute
	}
	if idpEvent.ProfileAttribute != nil {
		payloadChange = true
		payload.ProfileAttribute = *idpEvent.ProfileAttribute
	}
	return payloadChange
}

func reduceAppleIDPRelationalChangedColumns(payload *domain.Apple, idpEvent *idp.AppleIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.TeamID != nil {
		payloadChange = true
		payload.TeamID = *idpEvent.TeamID
	}
	if idpEvent.KeyID != nil {
		payloadChange = true
		payload.KeyID = *idpEvent.KeyID
	}
	if idpEvent.PrivateKey != nil {
		payloadChange = true
		payload.PrivateKey = idpEvent.PrivateKey
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChange
}

func reduceSAMLIDPRelationalChangedColumns(payload *domain.SAML, idpEvent *idp.SAMLIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.Metadata != nil {
		payloadChange = true
		payload.Metadata = idpEvent.Metadata
	}
	if idpEvent.Key != nil {
		payloadChange = true
		payload.Key = idpEvent.Key
	}
	if idpEvent.Certificate != nil {
		payloadChange = true
		payload.Certificate = idpEvent.Certificate
	}
	if idpEvent.Binding != nil {
		payloadChange = true
		payload.Binding = *idpEvent.Binding
	}
	if idpEvent.WithSignedRequest != nil {
		payloadChange = true
		payload.WithSignedRequest = *idpEvent.WithSignedRequest
	}
	if idpEvent.NameIDFormat != nil {
		payloadChange = true
		payload.NameIDFormat = idpEvent.NameIDFormat
	}
	if idpEvent.TransientMappingAttributeName != nil {
		payloadChange = true
		payload.TransientMappingAttributeName = *idpEvent.TransientMappingAttributeName
	}
	if idpEvent.FederatedLogoutEnabled != nil {
		payloadChange = true
		payload.FederatedLogoutEnabled = *idpEvent.FederatedLogoutEnabled
	}
	return payloadChange
}
