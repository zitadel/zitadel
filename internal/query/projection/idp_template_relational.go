package projection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	IDPRelationalAllowCreationCol     = "allow_creation"
	IDPRelationalAllowLinkingCol      = "allow_linking"
	IDPRelationalAllowAutoCreationCol = "allow_auto_creation"
	IDPRelationalAllowAutoUpdateCol   = "allow_auto_update"
	IDPRelationalAllowAutoLinkingCol  = "allow_auto_linking"
)

type idpTemplateRelationalProjection struct {
	idpRepo domain.IDProviderRepository
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
					Reduce: p.reduceJWTIDPReducedAdded,
				},
				{
					Event:  instance.JWTIDPChangedEventType,
					Reduce: p.reduceJWTIDPRelationalChanged,
				},
				// TODO
				// {
				// 	Event:  instance.IDPConfigAddedEventType,
				// 	Reduce: p.reduceOldConfigAdded,
				// },
				// TODO
				// 		{
				// 			Event:  instance.IDPConfigChangedEventType,
				// 			Reduce: p.reduceOldConfigChanged,
				// 		},
				// TODO
				// 		{
				// 			Event:  instance.IDPOIDCConfigAddedEventType,
				// 			Reduce: p.reduceOldOIDCConfigAdded,
				// 		},
				// TODO
				// 		{
				// 			Event:  instance.IDPOIDCConfigChangedEventType,
				// 			Reduce: p.reduceOldOIDCConfigChanged,
				// 		},
				// TODO
				// 		{
				// 			Event:  instance.IDPJWTConfigAddedEventType,
				// 			Reduce: p.reduceOldJWTConfigAdded,
				// 		},
				// TODO
				// 		{
				// 			Event:  instance.IDPJWTConfigChangedEventType,
				// 			Reduce: p.reduceOldJWTConfigChanged,
				// 		},
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
				// 		{
				// 			Event:  instance.GitHubEnterpriseIDPAddedEventType,
				// 			Reduce: p.reduceGitHubEnterpriseIDPAdded,
				// 		},
				// 		{
				// 			Event:  instance.GitHubEnterpriseIDPChangedEventType,
				// 			Reduce: p.reduceGitHubEnterpriseIDPChanged,
				// 		},
				// 		{
				// 			Event:  instance.GitLabIDPAddedEventType,
				// 			Reduce: p.reduceGitLabIDPAdded,
				// 		},
				// 		{
				// 			Event:  instance.GitLabIDPChangedEventType,
				// 			Reduce: p.reduceGitLabIDPChanged,
				// 		},
				// 		{
				// 			Event:  instance.GitLabSelfHostedIDPAddedEventType,
				// 			Reduce: p.reduceGitLabSelfHostedIDPAdded,
				// 		},
				// 		{
				// 			Event:  instance.GitLabSelfHostedIDPChangedEventType,
				// 			Reduce: p.reduceGitLabSelfHostedIDPChanged,
				// 		},
				// 		{
				// 			Event:  instance.GoogleIDPAddedEventType,
				// 			Reduce: p.reduceGoogleIDPAdded,
				// 		},
				// 		{
				// 			Event:  instance.GoogleIDPChangedEventType,
				// 			Reduce: p.reduceGoogleIDPChanged,
				// 		},
				// 		{
				// 			Event:  instance.LDAPIDPAddedEventType,
				// 			Reduce: p.reduceLDAPIDPAdded,
				// 		},
				// 		{
				// 			Event:  instance.LDAPIDPChangedEventType,
				// 			Reduce: p.reduceLDAPIDPChanged,
				// 		},
				// 		{
				// 			Event:  instance.AppleIDPAddedEventType,
				// 			Reduce: p.reduceAppleIDPAdded,
				// 		},
				// 		{
				// 			Event:  instance.AppleIDPChangedEventType,
				// 			Reduce: p.reduceAppleIDPChanged,
				// 		},
				// 		{
				// 			Event:  instance.SAMLIDPAddedEventType,
				// 			Reduce: p.reduceSAMLIDPAdded,
				// 		},
				// 		{
				// 			Event:  instance.SAMLIDPChangedEventType,
				// 			Reduce: p.reduceSAMLIDPChanged,
				// 		},
				// 		{
				// 			Event:  instance.IDPConfigRemovedEventType,
				// 			Reduce: p.reduceIDPConfigRemoved,
				// 		},
				// 		{
				// 			Event:  instance.IDPRemovedEventType,
				// 			Reduce: p.reduceIDPRemoved,
				// 		},
				// 		{
				// 			Event:  instance.InstanceRemovedEventType,
				// 			Reduce: reduceInstanceRemovedHelper(IDPTemplateInstanceIDCol),
				// 		},
				// 	},
				// },
				// {
				// 	Aggregate: org.AggregateType,
				// 	EventReducers: []handler.EventReducer{
				// 		{
				// 			Event:  org.OAuthIDPAddedEventType,
				// 			Reduce: p.reduceOAuthIDPAdded,
				// 		},
				// 		{
				// 			Event:  org.OAuthIDPChangedEventType,
				// 			Reduce: p.reduceOAuthIDPChanged,
				// 		},
				// 		{
				// 			Event:  org.OIDCIDPAddedEventType,
				// 			Reduce: p.reduceOIDCIDPAdded,
				// 		},
				// 		{
				// 			Event:  org.OIDCIDPChangedEventType,
				// 			Reduce: p.reduceOIDCIDPChanged,
				// 		},
				// 		{
				// 			Event:  org.OIDCIDPMigratedAzureADEventType,
				// 			Reduce: p.reduceOIDCIDPMigratedAzureAD,
				// 		},
				// 		{
				// 			Event:  org.OIDCIDPMigratedGoogleEventType,
				// 			Reduce: p.reduceOIDCIDPMigratedGoogle,
				// 		},
				// 		{
				// 			Event:  org.JWTIDPAddedEventType,
				// 			Reduce: p.reduceJWTIDPAdded,
				// 		},
				// 		{
				// 			Event:  org.JWTIDPChangedEventType,
				// 			Reduce: p.reduceJWTIDPChanged,
				// 		},
				// 		{
				// 			Event:  org.IDPConfigAddedEventType,
				// 			Reduce: p.reduceOldConfigAdded,
				// 		},
				// 		{
				// 			Event:  org.IDPConfigChangedEventType,
				// 			Reduce: p.reduceOldConfigChanged,
				// 		},
				// 		{
				// 			Event:  org.IDPOIDCConfigAddedEventType,
				// 			Reduce: p.reduceOldOIDCConfigAdded,
				// 		},
				// 		{
				// 			Event:  org.IDPOIDCConfigChangedEventType,
				// 			Reduce: p.reduceOldOIDCConfigChanged,
				// 		},
				// 		{
				// 			Event:  org.IDPJWTConfigAddedEventType,
				// 			Reduce: p.reduceOldJWTConfigAdded,
				// 		},
				// 		{
				// 			Event:  org.IDPJWTConfigChangedEventType,
				// 			Reduce: p.reduceOldJWTConfigChanged,
				// 		},
				// 		{
				// 			Event:  org.AzureADIDPAddedEventType,
				// 			Reduce: p.reduceAzureADIDPAdded,
				// 		},
				// 		{
				// 			Event:  org.AzureADIDPChangedEventType,
				// 			Reduce: p.reduceAzureADIDPChanged,
				// 		},
				// 		{
				// 			Event:  org.GitHubIDPAddedEventType,
				// 			Reduce: p.reduceGitHubIDPAdded,
				// 		},
				// 		{
				// 			Event:  org.GitHubIDPChangedEventType,
				// 			Reduce: p.reduceGitHubIDPChanged,
				// 		},
				// 		{
				// 			Event:  org.GitHubEnterpriseIDPAddedEventType,
				// 			Reduce: p.reduceGitHubEnterpriseIDPAdded,
				// 		},
				// 		{
				// 			Event:  org.GitHubEnterpriseIDPChangedEventType,
				// 			Reduce: p.reduceGitHubEnterpriseIDPChanged,
				// 		},
				// 		{
				// 			Event:  org.GitLabIDPAddedEventType,
				// 			Reduce: p.reduceGitLabIDPAdded,
				// 		},
				// 		{
				// 			Event:  org.GitLabIDPChangedEventType,
				// 			Reduce: p.reduceGitLabIDPChanged,
				// 		},
				// 		{
				// 			Event:  org.GitLabSelfHostedIDPAddedEventType,
				// 			Reduce: p.reduceGitLabSelfHostedIDPAdded,
				// 		},
				// 		{
				// 			Event:  org.GitLabSelfHostedIDPChangedEventType,
				// 			Reduce: p.reduceGitLabSelfHostedIDPChanged,
				// 		},
				// 		{
				// 			Event:  org.GoogleIDPAddedEventType,
				// 			Reduce: p.reduceGoogleIDPAdded,
				// 		},
				// 		{
				// 			Event:  org.GoogleIDPChangedEventType,
				// 			Reduce: p.reduceGoogleIDPChanged,
				// 		},
				// 		{
				// 			Event:  org.LDAPIDPAddedEventType,
				// 			Reduce: p.reduceLDAPIDPAdded,
				// 		},
				// 		{
				// 			Event:  org.LDAPIDPChangedEventType,
				// 			Reduce: p.reduceLDAPIDPChanged,
				// 		},
				// 		{
				// 			Event:  org.AppleIDPAddedEventType,
				// 			Reduce: p.reduceAppleIDPAdded,
				// 		},
				// 		{
				// 			Event:  org.AppleIDPChangedEventType,
				// 			Reduce: p.reduceAppleIDPChanged,
				// 		},
				// 		{
				// 			Event:  org.SAMLIDPAddedEventType,
				// 			Reduce: p.reduceSAMLIDPAdded,
				// 		},
				// 		{
				// 			Event:  org.SAMLIDPChangedEventType,
				// 			Reduce: p.reduceSAMLIDPChanged,
				// 		},
				// 		{
				// 			Event:  org.IDPConfigRemovedEventType,
				// 			Reduce: p.reduceIDPConfigRemoved,
				// 		},
				// 		{
				// 			Event:  org.IDPRemovedEventType,
				// 			Reduce: p.reduceIDPRemoved,
				// 		},
				// 		{
				// 			Event:  org.OrgRemovedEventType,
				// 			Reduce: p.reduceOwnerRemoved,
				// 		},
			},
		},
	}
}

func (p *idpTemplateRelationalProjection) reduceOAuthIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	// var idpEvent idp.OAuthIDPAddedEvent
	// var idpOwnerType domain.IdentityProviderType
	// switch e := event.(type) {
	// case *org.OAuthIDPAddedEvent:
	// 	idpEvent = e.OAuthIDPAddedEvent
	// 	idpOwnerType = domain.IdentityProviderTypeOrg
	// case *instance.OAuthIDPAddedEvent:
	// 	idpEvent = e.OAuthIDPAddedEvent
	// 	idpOwnerType = domain.IdentityProviderTypeSystem
	// default:
	// }

	e, ok := event.(*instance.OAuthIDPAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ap9ihb", "reduce.wrong.event.type %v", []eventstore.EventType{org.OAuthIDPAddedEventType, instance.OAuthIDPAddedEventType})
	}

	oauth := domain.OAuth{
		ClientID:              e.ClientID,
		ClientSecret:          e.ClientSecret,
		AuthorizationEndpoint: e.AuthorizationEndpoint,
		TokenEndpoint:         e.TokenEndpoint,
		UserEndpoint:          e.UserEndpoint,
		Scopes:                e.Scopes,
		IDAttribute:           e.IDAttribute,
		UsePKCE:               e.UsePKCE,
	}

	payload, err := json.Marshal(oauth)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, e.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
				handler.NewCol(IDPTemplateNameCol, e.Name),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeOAuth.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, e.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, e.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, e.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, e.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, domain.IDPAutoLinkingOption(e.AutoLinkingOption).String()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
				handler.NewCol(CreatedAt, e.CreationDate()),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceOAuthIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	// var idpEvent idp.OAuthIDPChangedEvent
	// switch e := event.(type) {
	// case *org.OAuthIDPChangedEvent:
	// 	idpEvent = e.OAuthIDPChangedEvent
	// case *instance.OAuthIDPChangedEvent:
	// 	idpEvent = e.OAuthIDPChangedEvent
	// default:
	// 	return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OAuthIDPChangedEventType, instance.OAuthIDPChangedEventType})
	// }

	e, ok := event.(*instance.OAuthIDPChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OAuthIDPChangedEventType, instance.OAuthIDPChangedEventType})
	}

	oauth, err := p.idpRepo.GetOAuth(context.Background(), p.idpRepo.IDCondition(e.ID), e.Agg.InstanceID, nil)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(e.Name, e.OptionChanges, &columns)

	payload := &oauth.OAuth
	payloadChanged := reduceOAuthIDPRelationalChangedColumns(payload, &e.OAuthIDPChangedEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, e.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceOIDCIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	// var idpEvent idp.OIDCIDPAddedEvent
	// var idpOwnerType domain.IdentityProviderType
	// switch e := event.(type) {
	// case *org.OIDCIDPAddedEvent:
	// 	idpEvent = e.OIDCIDPAddedEvent
	// 	idpOwnerType = domain.IdentityProviderTypeOrg
	// case *instance.OIDCIDPAddedEvent:
	// 	idpEvent = e.OIDCIDPAddedEvent
	// 	idpOwnerType = domain.IdentityProviderTypeSystem
	// default:
	// 	return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-9s02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPAddedEventType, instance.OIDCIDPAddedEventType})
	// }

	e, ok := event.(*instance.OIDCIDPAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-9s02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPAddedEventType, instance.OIDCIDPAddedEventType})
	}

	payload, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, e.ID),
				handler.NewCol(CreatedAt, e.CreationDate()),
				// handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				// handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				// handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, e.Name),
				// handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeOIDC.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, e.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, e.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, e.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, e.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, domain.IDPAutoLinkingOption(e.AutoLinkingOption).String()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceOIDCIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	// var idpEvent idp.OIDCIDPChangedEvent
	// switch e := event.(type) {
	// case *org.OIDCIDPChangedEvent:
	// 	idpEvent = e.OIDCIDPChangedEvent
	// case *instance.OIDCIDPChangedEvent:
	// 	idpEvent = e.OIDCIDPChangedEvent
	// default:
	// 	return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPChangedEventType, instance.OIDCIDPChangedEventType})
	// }

	e, ok := event.(*instance.OIDCIDPChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPChangedEventType, instance.OIDCIDPChangedEventType})
	}

	oidc, err := p.idpRepo.GetOIDC(context.Background(), p.idpRepo.IDCondition(e.ID), e.Agg.InstanceID, nil)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(e.Name, e.OptionChanges, &columns)

	// ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	// ops = append(ops,
	// 	handler.AddUpdateStatement(
	// 		reduceIDPRelationalChangedTemplateColumns(e.Name, e.OptionChanges),
	// 		[]handler.Condition{
	// 			handler.NewCond(IDPTemplateIDCol, e.ID),
	// 			handler.NewCond(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
	// 		},
	// 	),
	// )
	payload := &oidc.OIDC
	payloadChanged := reduceOIDCIDPRelationalChangedColumns(payload, &e.OIDCIDPChangedEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	// return handler.NewMultiStatement(
	// 	&e,
	// 	ops...,
	// ), nil
	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, e.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceOIDCIDPRelationalMigratedAzureAD(event eventstore.Event) (*handler.Statement, error) {
	// var idpEvent idp.OIDCIDPMigratedAzureADEvent
	// switch e := event.(type) {
	// case *org.OIDCIDPMigratedAzureADEvent:
	// 	idpEvent = e.OIDCIDPMigratedAzureADEvent
	// case *instance.OIDCIDPMigratedAzureADEvent:
	// 	idpEvent = e.OIDCIDPMigratedAzureADEvent
	// default:
	// 	return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPMigratedAzureADEventType, instance.OIDCIDPMigratedAzureADEventType})
	// }

	e, ok := event.(*instance.OIDCIDPMigratedAzureADEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPMigratedAzureADEventType, instance.OIDCIDPMigratedAzureADEventType})
	}

	azure := domain.Azure{
		ClientID:        e.ClientID,
		ClientSecret:    e.ClientSecret,
		Scopes:          e.Scopes,
		Tenant:          e.Tenant,
		IsEmailVerified: e.IsEmailVerified,
	}

	payload, err := json.Marshal(azure)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateNameCol, e.Name),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeAzure.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, e.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, e.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, e.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, e.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, domain.IDPAutoLinkingOption(e.AutoLinkingOption).String()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, e.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceOIDCIDPRelationalMigratedGoogle(event eventstore.Event) (*handler.Statement, error) {
	// var idpEvent idp.OIDCIDPMigratedGoogleEvent
	// switch e := event.(type) {
	// case *org.OIDCIDPMigratedGoogleEvent:
	// 	idpEvent = e.OIDCIDPMigratedGoogleEvent
	// case *instance.OIDCIDPMigratedGoogleEvent:
	// 	idpEvent = e.OIDCIDPMigratedGoogleEvent
	// default:
	// 	return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPMigratedGoogleEventType, instance.OIDCIDPMigratedGoogleEventType})
	// }

	e, ok := event.(*instance.OIDCIDPMigratedGoogleEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPMigratedGoogleEventType, instance.OIDCIDPMigratedGoogleEventType})
	}

	azure := domain.Google{
		ClientID:     e.ClientID,
		ClientSecret: e.ClientSecret,
		Scopes:       e.Scopes,
	}

	payload, err := json.Marshal(azure)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateNameCol, e.Name),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGoogle.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, e.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, e.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, e.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, e.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, domain.IDPAutoLinkingOption(e.AutoLinkingOption).String()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, e.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceJWTIDPReducedAdded(event eventstore.Event) (*handler.Statement, error) {
	// var e idp.JWTIDPAddedEvent
	// var idpOwnerType domain.IdentityProviderType
	// switch e := event.(type) {
	// case *org.JWTIDPAddedEvent:
	// 	idpEvent = e.JWTIDPAddedEvent
	// 	idpOwnerType = domain.IdentityProviderTypeOrg
	// case *instance.JWTIDPAddedEvent:
	// 	idpEvent = e.JWTIDPAddedEvent
	// 	idpOwnerType = domain.IdentityProviderTypeSystem
	// default:
	// 	return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-xopi2s", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPAddedEventType, instance.JWTIDPAddedEventType})
	// }

	e, ok := event.(*instance.JWTIDPAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-xopi2s", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPAddedEventType, instance.JWTIDPAddedEventType})
	}

	jwt := domain.JWT{
		JWTEndpoint:  e.JWTEndpoint,
		Issuer:       e.Issuer,
		KeysEndpoint: e.KeysEndpoint,
		HeaderName:   e.HeaderName,
	}

	payload, err := json.Marshal(jwt)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, e.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
				handler.NewCol(IDPTemplateNameCol, e.Name),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeJWT.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, e.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, e.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, e.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, e.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, domain.IDPAutoLinkingOption(e.AutoLinkingOption).String()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
				handler.NewCol(CreatedAt, e.CreationDate()),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceJWTIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	// var idpEvent idp.JWTIDPChangedEvent
	// switch e := event.(type) {
	// case *org.JWTIDPChangedEvent:
	// 	idpEvent = e.JWTIDPChangedEvent
	// case *instance.JWTIDPChangedEvent:
	// 	idpEvent = e.JWTIDPChangedEvent
	// default:
	// 	return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPChangedEventType, instance.JWTIDPChangedEventType})
	// }

	e, ok := event.(*instance.JWTIDPChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPChangedEventType, instance.JWTIDPChangedEventType})
	}

	jwt, err := p.idpRepo.GetJWT(context.Background(), p.idpRepo.IDCondition(e.ID), e.Agg.InstanceID, nil)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(e.Name, e.OptionChanges, &columns)

	payload := &jwt.JWT
	payloadChanged := reduceJWTIDPRelationalChangedColumns(payload, &e.JWTIDPChangedEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, e.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
	), nil
}

// func (p *idpTemplateRelationalProjection) reduceOldConfigAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.IDPConfigAddedEvent
// 	var idpOwnerType domain.IdentityProviderType
// 	switch e := event.(type) {
// 	case *org.IDPConfigAddedEvent:
// 		idpEvent = e.IDPConfigAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeOrg
// 	case *instance.IDPConfigAddedEvent:
// 		idpEvent = e.IDPConfigAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeSystem
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ADfeg", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigAddedEventType, instance.IDPConfigAddedEventType})
// 	}

// 	return handler.NewCreateStatement(
// 		event,
// 		[]handler.Column{
// 			handler.NewCol(IDPTemplateIDCol, idpEvent.ConfigID),
// 			handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
// 			handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
// 			handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
// 			handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
// 			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
// 			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
// 			handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
// 			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeUnspecified),
// 			handler.NewCol(IDPTemplateIsCreationAllowedCol, true),
// 			handler.NewCol(IDPTemplateIsLinkingAllowedCol, true),
// 			handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.AutoRegister),
// 			handler.NewCol(IDPTemplateIsAutoUpdateCol, false),
// 			handler.NewCol(IDPTemplateAutoLinkingCol, domain.AutoLinkingOptionUnspecified),
// 		},
// 	), nil
// }

// func (p *idpTemplateProjection) reduceOldConfigChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.IDPConfigChangedEvent
// 	switch e := event.(type) {
// 	case *org.IDPConfigChangedEvent:
// 		idpEvent = e.IDPConfigChangedEvent
// 	case *instance.IDPConfigChangedEvent:
// 		idpEvent = e.IDPConfigChangedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SAfg2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigChangedEventType, instance.IDPConfigChangedEventType})
// 	}

// 	cols := make([]handler.Column, 0, 4)
// 	if idpEvent.Name != nil {
// 		cols = append(cols, handler.NewCol(IDPTemplateNameCol, *idpEvent.Name))
// 	}
// 	if idpEvent.AutoRegister != nil {
// 		cols = append(cols, handler.NewCol(IDPTemplateIsAutoCreationCol, *idpEvent.AutoRegister))
// 	}
// 	cols = append(cols,
// 		handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
// 		handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
// 	)

// 	return handler.NewUpdateStatement(
// 		event,
// 		cols,
// 		[]handler.Condition{
// 			handler.NewCond(IDPTemplateIDCol, idpEvent.ConfigID),
// 			handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 		},
// 	), nil
// }

// func (p *idpTemplateProjection) reduceOldOIDCConfigAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.OIDCConfigAddedEvent
// 	switch e := event.(type) {
// 	case *org.IDPOIDCConfigAddedEvent:
// 		idpEvent = e.OIDCConfigAddedEvent
// 	case *instance.IDPOIDCConfigAddedEvent:
// 		idpEvent = e.OIDCConfigAddedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ASFdq2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigAddedEventType, instance.IDPOIDCConfigAddedEventType})
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		handler.AddUpdateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
// 				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeOIDC),
// 			},
// 			[]handler.Condition{
// 				handler.NewCond(IDPTemplateIDCol, idpEvent.IDPConfigID),
// 				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(OIDCIDCol, idpEvent.IDPConfigID),
// 				handler.NewCol(OIDCInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(OIDCIssuerCol, idpEvent.Issuer),
// 				handler.NewCol(OIDCClientIDCol, idpEvent.ClientID),
// 				handler.NewCol(OIDCClientSecretCol, idpEvent.ClientSecret),
// 				handler.NewCol(OIDCScopesCol, database.TextArray[string](idpEvent.Scopes)),
// 				handler.NewCol(OIDCIDTokenMappingCol, true),
// 				handler.NewCol(OIDCUsePKCECol, false),
// 			},
// 			handler.WithTableSuffix(IDPTemplateOIDCSuffix),
// 		),
// 	), nil
// }

// func (p *idpTemplateProjection) reduceOldOIDCConfigChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.OIDCConfigChangedEvent
// 	switch e := event.(type) {
// 	case *org.IDPOIDCConfigChangedEvent:
// 		idpEvent = e.OIDCConfigChangedEvent
// 	case *instance.IDPOIDCConfigChangedEvent:
// 		idpEvent = e.OIDCConfigChangedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPChangedEventType, instance.OIDCIDPChangedEventType})
// 	}

// 	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
// 	ops = append(ops,
// 		handler.AddUpdateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
// 			},
// 			[]handler.Condition{
// 				handler.NewCond(IDPTemplateIDCol, idpEvent.IDPConfigID),
// 				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),
// 	)
// 	oidcCols := make([]handler.Column, 0, 4)
// 	if idpEvent.ClientID != nil {
// 		oidcCols = append(oidcCols, handler.NewCol(OIDCClientIDCol, *idpEvent.ClientID))
// 	}
// 	if idpEvent.ClientSecret != nil {
// 		oidcCols = append(oidcCols, handler.NewCol(OIDCClientSecretCol, *idpEvent.ClientSecret))
// 	}
// 	if idpEvent.Issuer != nil {
// 		oidcCols = append(oidcCols, handler.NewCol(OIDCIssuerCol, *idpEvent.Issuer))
// 	}
// 	if idpEvent.Scopes != nil {
// 		oidcCols = append(oidcCols, handler.NewCol(OIDCScopesCol, database.TextArray[string](idpEvent.Scopes)))
// 	}
// 	if len(oidcCols) > 0 {
// 		ops = append(ops,
// 			handler.AddUpdateStatement(
// 				oidcCols,
// 				[]handler.Condition{
// 					handler.NewCond(OIDCIDCol, idpEvent.IDPConfigID),
// 					handler.NewCond(OIDCInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				},
// 				handler.WithTableSuffix(IDPTemplateOIDCSuffix),
// 			),
// 		)
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		ops...,
// 	), nil
// }

// func (p *idpTemplateProjection) reduceOldJWTConfigAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.JWTConfigAddedEvent
// 	switch e := event.(type) {
// 	case *org.IDPJWTConfigAddedEvent:
// 		idpEvent = e.JWTConfigAddedEvent
// 	case *instance.IDPJWTConfigAddedEvent:
// 		idpEvent = e.JWTConfigAddedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ASFdq2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigAddedEventType, instance.IDPJWTConfigAddedEventType})
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		handler.AddUpdateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
// 				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeJWT),
// 			},
// 			[]handler.Condition{
// 				handler.NewCond(IDPTemplateIDCol, idpEvent.IDPConfigID),
// 				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(JWTIDCol, idpEvent.IDPConfigID),
// 				handler.NewCol(JWTInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(JWTIssuerCol, idpEvent.Issuer),
// 				handler.NewCol(JWTEndpointCol, idpEvent.JWTEndpoint),
// 				handler.NewCol(JWTKeysEndpointCol, idpEvent.KeysEndpoint),
// 				handler.NewCol(JWTHeaderNameCol, idpEvent.HeaderName),
// 			},
// 			handler.WithTableSuffix(IDPTemplateJWTSuffix),
// 		),
// 	), nil
// }

// func (p *idpTemplateProjection) reduceOldJWTConfigChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.JWTConfigChangedEvent
// 	switch e := event.(type) {
// 	case *org.IDPJWTConfigChangedEvent:
// 		idpEvent = e.JWTConfigChangedEvent
// 	case *instance.IDPJWTConfigChangedEvent:
// 		idpEvent = e.JWTConfigChangedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPChangedEventType, instance.JWTIDPChangedEventType})
// 	}

// 	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
// 	ops = append(ops,
// 		handler.AddUpdateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
// 			},
// 			[]handler.Condition{
// 				handler.NewCond(IDPTemplateIDCol, idpEvent.IDPConfigID),
// 				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),
// 	)
// 	jwtCols := make([]handler.Column, 0, 4)
// 	if idpEvent.JWTEndpoint != nil {
// 		jwtCols = append(jwtCols, handler.NewCol(JWTEndpointCol, *idpEvent.JWTEndpoint))
// 	}
// 	if idpEvent.KeysEndpoint != nil {
// 		jwtCols = append(jwtCols, handler.NewCol(JWTKeysEndpointCol, *idpEvent.KeysEndpoint))
// 	}
// 	if idpEvent.HeaderName != nil {
// 		jwtCols = append(jwtCols, handler.NewCol(JWTHeaderNameCol, *idpEvent.HeaderName))
// 	}
// 	if idpEvent.Issuer != nil {
// 		jwtCols = append(jwtCols, handler.NewCol(JWTIssuerCol, *idpEvent.Issuer))
// 	}
// 	if len(jwtCols) > 0 {
// 		ops = append(ops,
// 			handler.AddUpdateStatement(
// 				jwtCols,
// 				[]handler.Condition{
// 					handler.NewCond(JWTIDCol, idpEvent.IDPConfigID),
// 					handler.NewCond(JWTInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				},
// 				handler.WithTableSuffix(IDPTemplateJWTSuffix),
// 			),
// 		)
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		ops...,
// 	), nil
// }

func (p *idpTemplateRelationalProjection) reduceAzureADIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	// var idpEvent idp.AzureADIDPAddedEvent
	// var idpOwnerType domain.IdentityProviderType
	// switch e := event.(type) {
	// case *org.AzureADIDPAddedEvent:
	// 	idpEvent = e.AzureADIDPAddedEvent
	// 	idpOwnerType = domain.IdentityProviderTypeOrg
	// case *instance.AzureADIDPAddedEvent:
	// 	idpEvent = e.AzureADIDPAddedEvent
	// 	idpOwnerType = domain.IdentityProviderTypeSystem
	// default:
	// 	return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.AzureADIDPAddedEventType, instance.AzureADIDPAddedEventType})
	// }

	e, ok := event.(*instance.AzureADIDPAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.AzureADIDPAddedEventType, instance.AzureADIDPAddedEventType})
	}

	azure := domain.Azure{
		ClientID:        e.ClientID,
		ClientSecret:    e.ClientSecret,
		Scopes:          e.Scopes,
		Tenant:          e.Tenant,
		IsEmailVerified: e.IsEmailVerified,
	}

	payload, err := json.Marshal(azure)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, e.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateNameCol, e.Name),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeAzure.String()),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, e.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, e.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, e.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, e.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, domain.IDPAutoLinkingOption(e.AutoLinkingOption).String()),
				handler.NewCol(CreatedAt, e.CreationDate()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceAzureADIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	// var idpEvent idp.AzureADIDPChangedEvent
	// switch e := event.(type) {
	// case *org.AzureADIDPChangedEvent:
	// 	idpEvent = e.AzureADIDPChangedEvent
	// case *instance.AzureADIDPChangedEvent:
	// 	idpEvent = e.AzureADIDPChangedEvent
	// default:
	// 	return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.AzureADIDPChangedEventType, instance.AzureADIDPChangedEventType})
	// }

	e, ok := event.(*instance.AzureADIDPChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.AzureADIDPChangedEventType, instance.AzureADIDPChangedEventType})
	}

	oauth, err := p.idpRepo.GetOAzureAD(context.Background(), p.idpRepo.IDCondition(e.ID), e.Agg.InstanceID, nil)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(e.Name, e.OptionChanges, &columns)

	payload := &oauth.Azure
	payloadChanged := reduceAzureADIDPRelationalChangedColumns(payload, &e.AzureADIDPChangedEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, e.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *idpTemplateRelationalProjection) reduceGitHubIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	// var idpEvent idp.GitHubIDPAddedEvent
	// var idpOwnerType domain.IdentityProviderType
	// switch e := event.(type) {
	// case *org.GitHubIDPAddedEvent:
	// 	idpEvent = e.GitHubIDPAddedEvent
	// 	idpOwnerType = domain.IdentityProviderTypeOrg
	// case *instance.GitHubIDPAddedEvent:
	// 	idpEvent = e.GitHubIDPAddedEvent
	// 	idpOwnerType = domain.IdentityProviderTypeSystem
	// default:
	// 	return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPAddedEventType, instance.GitHubIDPAddedEventType})
	// }

	e, ok := event.(*instance.GitHubIDPAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPAddedEventType, instance.GitHubIDPAddedEventType})
	}

	github := domain.Github{
		ClientID:     e.ClientID,
		ClientSecret: e.ClientSecret,
		Scopes:       e.Scopes,
	}

	payload, err := json.Marshal(github)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, e.ID),
				handler.NewCol(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateNameCol, e.Name),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGithub.String()),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
				handler.NewCol(IDPRelationalAllowCreationCol, e.IsCreationAllowed),
				handler.NewCol(IDPRelationalAllowLinkingCol, e.IsLinkingAllowed),
				handler.NewCol(IDPRelationalAllowAutoCreationCol, e.IsAutoCreation),
				handler.NewCol(IDPRelationalAllowAutoUpdateCol, e.IsAutoUpdate),
				handler.NewCol(IDPRelationalAllowAutoLinkingCol, domain.IDPAutoLinkingOption(e.AutoLinkingOption).String()),
				handler.NewCol(CreatedAt, e.CreationDate()),
				handler.NewCol(IDPRelationalPayloadCol, payload),
			},
		),
	), nil

	// return handler.NewMultiStatement(
	// 	&idpEvent,
	// 	handler.AddCreateStatement(
	// 		[]handler.Column{
	// 			handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
	// 			handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
	// 			handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
	// 			handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
	// 			handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
	// 			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
	// 			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
	// 			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
	// 			handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
	// 			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGithub),
	// 			handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
	// 			handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
	// 			handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
	// 			handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
	// 			handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
	// 		},
	// 	),
	// 	handler.AddCreateStatement(
	// 		[]handler.Column{
	// 			handler.NewCol(GitHubIDCol, idpEvent.ID),
	// 			handler.NewCol(GitHubInstanceIDCol, idpEvent.Aggregate().InstanceID),
	// 			handler.NewCol(GitHubClientIDCol, idpEvent.ClientID),
	// 			handler.NewCol(GitHubClientSecretCol, idpEvent.ClientSecret),
	// 			handler.NewCol(GitHubScopesCol, database.TextArray[string](idpEvent.Scopes)),
	// 		},
	// 		handler.WithTableSuffix(IDPTemplateGitHubSuffix),
	// 	),
	// ), nil
}

// func (p *idpTemplateProjection) reduceGitHubEnterpriseIDPAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.GitHubEnterpriseIDPAddedEvent
// 	var idpOwnerType domain.IdentityProviderType
// 	switch e := event.(type) {
// 	case *org.GitHubEnterpriseIDPAddedEvent:
// 		idpEvent = e.GitHubEnterpriseIDPAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeOrg
// 	case *instance.GitHubEnterpriseIDPAddedEvent:
// 		idpEvent = e.GitHubEnterpriseIDPAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeSystem
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Sf3g2a", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubEnterpriseIDPAddedEventType, instance.GitHubEnterpriseIDPAddedEventType})
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
// 				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
// 				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
// 				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
// 				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
// 				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
// 				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGitHubEnterprise),
// 				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
// 				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
// 				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
// 				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
// 				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
// 			},
// 		),
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(GitHubEnterpriseIDCol, idpEvent.ID),
// 				handler.NewCol(GitHubEnterpriseInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(GitHubEnterpriseClientIDCol, idpEvent.ClientID),
// 				handler.NewCol(GitHubEnterpriseClientSecretCol, idpEvent.ClientSecret),
// 				handler.NewCol(GitHubEnterpriseAuthorizationEndpointCol, idpEvent.AuthorizationEndpoint),
// 				handler.NewCol(GitHubEnterpriseTokenEndpointCol, idpEvent.TokenEndpoint),
// 				handler.NewCol(GitHubEnterpriseUserEndpointCol, idpEvent.UserEndpoint),
// 				handler.NewCol(GitHubEnterpriseScopesCol, database.TextArray[string](idpEvent.Scopes)),
// 			},
// 			handler.WithTableSuffix(IDPTemplateGitHubEnterpriseSuffix),
// 		),
// 	), nil
// }

func (p *idpTemplateRelationalProjection) reduceGitHubIDPRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	// var idpEvent idp.GitHubIDPChangedEvent
	// switch e := event.(type) {
	// case *org.GitHubIDPChangedEvent:
	// 	idpEvent = e.GitHubIDPChangedEvent
	// case *instance.GitHubIDPChangedEvent:
	// 	idpEvent = e.GitHubIDPChangedEvent
	// default:
	// 	return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPChangedEventType, instance.GitHubIDPChangedEventType})
	// }

	e, ok := event.(*instance.GitHubIDPChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPChangedEventType, instance.GitHubIDPChangedEventType})
	}

	github, err := p.idpRepo.GetGithub(context.Background(), p.idpRepo.IDCondition(e.ID), e.Agg.InstanceID, nil)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 7)
	reduceIDPRelationalChangedTemplateColumns(e.Name, e.OptionChanges, &columns)

	payload := &github.Github
	payloadChanged := reduceGitHubIDPRelationalChangedColumns(payload, &e.GitHubIDPChangedEvent)
	if payloadChanged {
		payload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		columns = append(columns, handler.NewCol(IDPRelationalPayloadCol, payload))
	}

	fmt.Printf("@@ >>>>>>>>>>>>>>>>>>>>>>>>>>>> *e.Name = %+v\n", *e.Name)
	fmt.Println("@@ >>>>>>>>>>>>>>>>>>>>>>>>>>>> UPDATE GITHUB")

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, e.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
	), nil

	// ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	// ops = append(ops,
	// 	handler.AddUpdateStatement(
	// 		reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
	// 		[]handler.Condition{
	// 			handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
	// 			handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
	// 		},
	// 	),
	// )
	// githubCols := reduceGitHubIDPChangedColumns(idpEvent)
	// if len(githubCols) > 0 {
	// 	ops = append(ops,
	// 		handler.AddUpdateStatement(
	// 			githubCols,
	// 			[]handler.Condition{
	// 				handler.NewCond(GitHubIDCol, idpEvent.ID),
	// 				handler.NewCond(GitHubInstanceIDCol, idpEvent.Aggregate().InstanceID),
	// 			},
	// 			handler.WithTableSuffix(IDPTemplateGitHubSuffix),
	// 		),
	// 	)
	// }

	// return handler.NewMultiStatement(
	// 	&idpEvent,
	// 	ops...,
	// ), nil
}

// func (p *idpTemplateProjection) reduceGitHubEnterpriseIDPChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.GitHubEnterpriseIDPChangedEvent
// 	switch e := event.(type) {
// 	case *org.GitHubEnterpriseIDPChangedEvent:
// 		idpEvent = e.GitHubEnterpriseIDPChangedEvent
// 	case *instance.GitHubEnterpriseIDPChangedEvent:
// 		idpEvent = e.GitHubEnterpriseIDPChangedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SDg3g", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubEnterpriseIDPChangedEventType, instance.GitHubEnterpriseIDPChangedEventType})
// 	}

// 	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
// 	ops = append(ops,
// 		handler.AddUpdateStatement(
// 			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
// 			[]handler.Condition{
// 				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
// 				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),
// 	)
// 	githubCols := reduceGitHubEnterpriseIDPChangedColumns(idpEvent)
// 	if len(githubCols) > 0 {
// 		ops = append(ops,
// 			handler.AddUpdateStatement(
// 				githubCols,
// 				[]handler.Condition{
// 					handler.NewCond(GitHubEnterpriseIDCol, idpEvent.ID),
// 					handler.NewCond(GitHubEnterpriseInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				},
// 				handler.WithTableSuffix(IDPTemplateGitHubEnterpriseSuffix),
// 			),
// 		)
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		ops...,
// 	), nil
// }

// func (p *idpTemplateProjection) reduceGitLabIDPAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.GitLabIDPAddedEvent
// 	var idpOwnerType domain.IdentityProviderType
// 	switch e := event.(type) {
// 	case *org.GitLabIDPAddedEvent:
// 		idpEvent = e.GitLabIDPAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeOrg
// 	case *instance.GitLabIDPAddedEvent:
// 		idpEvent = e.GitLabIDPAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeSystem
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabIDPAddedEventType, instance.GitLabIDPAddedEventType})
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
// 				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
// 				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
// 				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
// 				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
// 				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
// 				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGitLab),
// 				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
// 				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
// 				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
// 				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
// 				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
// 			},
// 		),
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(GitLabIDCol, idpEvent.ID),
// 				handler.NewCol(GitLabInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(GitLabClientIDCol, idpEvent.ClientID),
// 				handler.NewCol(GitLabClientSecretCol, idpEvent.ClientSecret),
// 				handler.NewCol(GitLabScopesCol, database.TextArray[string](idpEvent.Scopes)),
// 			},
// 			handler.WithTableSuffix(IDPTemplateGitLabSuffix),
// 		),
// 	), nil
// }

// func (p *idpTemplateProjection) reduceGitLabIDPChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.GitLabIDPChangedEvent
// 	switch e := event.(type) {
// 	case *org.GitLabIDPChangedEvent:
// 		idpEvent = e.GitLabIDPChangedEvent
// 	case *instance.GitLabIDPChangedEvent:
// 		idpEvent = e.GitLabIDPChangedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabIDPChangedEventType, instance.GitLabIDPChangedEventType})
// 	}

// 	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
// 	ops = append(ops,
// 		handler.AddUpdateStatement(
// 			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
// 			[]handler.Condition{
// 				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
// 				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),
// 	)
// 	gitlabCols := reduceGitLabIDPChangedColumns(idpEvent)
// 	if len(gitlabCols) > 0 {
// 		ops = append(ops,
// 			handler.AddUpdateStatement(
// 				gitlabCols,
// 				[]handler.Condition{
// 					handler.NewCond(GitLabIDCol, idpEvent.ID),
// 					handler.NewCond(GitLabInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				},
// 				handler.WithTableSuffix(IDPTemplateGitLabSuffix),
// 			),
// 		)
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		ops...,
// 	), nil
// }

// func (p *idpTemplateProjection) reduceGitLabSelfHostedIDPAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.GitLabSelfHostedIDPAddedEvent
// 	var idpOwnerType domain.IdentityProviderType
// 	switch e := event.(type) {
// 	case *org.GitLabSelfHostedIDPAddedEvent:
// 		idpEvent = e.GitLabSelfHostedIDPAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeOrg
// 	case *instance.GitLabSelfHostedIDPAddedEvent:
// 		idpEvent = e.GitLabSelfHostedIDPAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeSystem
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SAF3gw", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabSelfHostedIDPAddedEventType, instance.GitLabSelfHostedIDPAddedEventType})
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
// 				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
// 				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
// 				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
// 				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
// 				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
// 				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGitLabSelfHosted),
// 				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
// 				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
// 				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
// 				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
// 				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
// 			},
// 		),
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(GitLabSelfHostedIDCol, idpEvent.ID),
// 				handler.NewCol(GitLabSelfHostedInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(GitLabSelfHostedIssuerCol, idpEvent.Issuer),
// 				handler.NewCol(GitLabSelfHostedClientIDCol, idpEvent.ClientID),
// 				handler.NewCol(GitLabSelfHostedClientSecretCol, idpEvent.ClientSecret),
// 				handler.NewCol(GitLabSelfHostedScopesCol, database.TextArray[string](idpEvent.Scopes)),
// 			},
// 			handler.WithTableSuffix(IDPTemplateGitLabSelfHostedSuffix),
// 		),
// 	), nil
// }

// func (p *idpTemplateProjection) reduceGitLabSelfHostedIDPChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.GitLabSelfHostedIDPChangedEvent
// 	switch e := event.(type) {
// 	case *org.GitLabSelfHostedIDPChangedEvent:
// 		idpEvent = e.GitLabSelfHostedIDPChangedEvent
// 	case *instance.GitLabSelfHostedIDPChangedEvent:
// 		idpEvent = e.GitLabSelfHostedIDPChangedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SAf3g2", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabSelfHostedIDPChangedEventType, instance.GitLabSelfHostedIDPChangedEventType})
// 	}

// 	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
// 	ops = append(ops,
// 		handler.AddUpdateStatement(
// 			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
// 			[]handler.Condition{
// 				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
// 				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),
// 	)
// 	gitlabCols := reduceGitLabSelfHostedIDPChangedColumns(idpEvent)
// 	if len(gitlabCols) > 0 {
// 		ops = append(ops,
// 			handler.AddUpdateStatement(
// 				gitlabCols,
// 				[]handler.Condition{
// 					handler.NewCond(GitLabSelfHostedIDCol, idpEvent.ID),
// 					handler.NewCond(GitLabSelfHostedInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				},
// 				handler.WithTableSuffix(IDPTemplateGitLabSelfHostedSuffix),
// 			),
// 		)
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		ops...,
// 	), nil
// }

// func (p *idpTemplateProjection) reduceGoogleIDPAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.GoogleIDPAddedEvent
// 	var idpOwnerType domain.IdentityProviderType
// 	switch e := event.(type) {
// 	case *org.GoogleIDPAddedEvent:
// 		idpEvent = e.GoogleIDPAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeOrg
// 	case *instance.GoogleIDPAddedEvent:
// 		idpEvent = e.GoogleIDPAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeSystem
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ap9ihb", "reduce.wrong.event.type %v", []eventstore.EventType{org.GoogleIDPAddedEventType, instance.GoogleIDPAddedEventType})
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
// 				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
// 				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
// 				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
// 				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
// 				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
// 				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGoogle),
// 				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
// 				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
// 				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
// 				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
// 				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
// 			},
// 		),
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(GoogleIDCol, idpEvent.ID),
// 				handler.NewCol(GoogleInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(GoogleClientIDCol, idpEvent.ClientID),
// 				handler.NewCol(GoogleClientSecretCol, idpEvent.ClientSecret),
// 				handler.NewCol(GoogleScopesCol, database.TextArray[string](idpEvent.Scopes)),
// 			},
// 			handler.WithTableSuffix(IDPTemplateGoogleSuffix),
// 		),
// 	), nil
// }

// func (p *idpTemplateProjection) reduceGoogleIDPChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.GoogleIDPChangedEvent
// 	switch e := event.(type) {
// 	case *org.GoogleIDPChangedEvent:
// 		idpEvent = e.GoogleIDPChangedEvent
// 	case *instance.GoogleIDPChangedEvent:
// 		idpEvent = e.GoogleIDPChangedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GoogleIDPChangedEventType, instance.GoogleIDPChangedEventType})
// 	}

// 	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
// 	ops = append(ops,
// 		handler.AddUpdateStatement(
// 			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
// 			[]handler.Condition{
// 				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
// 				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),
// 	)
// 	googleCols := reduceGoogleIDPChangedColumns(idpEvent)
// 	if len(googleCols) > 0 {
// 		ops = append(ops,
// 			handler.AddUpdateStatement(
// 				googleCols,
// 				[]handler.Condition{
// 					handler.NewCond(GoogleIDCol, idpEvent.ID),
// 					handler.NewCond(GoogleInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				},
// 				handler.WithTableSuffix(IDPTemplateGoogleSuffix),
// 			),
// 		)
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		ops...,
// 	), nil
// }

// func (p *idpTemplateProjection) reduceLDAPIDPAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.LDAPIDPAddedEvent
// 	var idpOwnerType domain.IdentityProviderType
// 	switch e := event.(type) {
// 	case *org.LDAPIDPAddedEvent:
// 		idpEvent = e.LDAPIDPAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeOrg
// 	case *instance.LDAPIDPAddedEvent:
// 		idpEvent = e.LDAPIDPAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeSystem
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-9s02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPAddedEventType, instance.LDAPIDPAddedEventType})
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
// 				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
// 				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
// 				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
// 				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
// 				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
// 				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeLDAP),
// 				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
// 				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
// 				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
// 				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
// 				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
// 			},
// 		),
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(LDAPIDCol, idpEvent.ID),
// 				handler.NewCol(LDAPInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(LDAPServersCol, database.TextArray[string](idpEvent.Servers)),
// 				handler.NewCol(LDAPStartTLSCol, idpEvent.StartTLS),
// 				handler.NewCol(LDAPBaseDNCol, idpEvent.BaseDN),
// 				handler.NewCol(LDAPBindDNCol, idpEvent.BindDN),
// 				handler.NewCol(LDAPBindPasswordCol, idpEvent.BindPassword),
// 				handler.NewCol(LDAPUserBaseCol, idpEvent.UserBase),
// 				handler.NewCol(LDAPUserObjectClassesCol, database.TextArray[string](idpEvent.UserObjectClasses)),
// 				handler.NewCol(LDAPUserFiltersCol, database.TextArray[string](idpEvent.UserFilters)),
// 				handler.NewCol(LDAPTimeoutCol, idpEvent.Timeout),
// 				handler.NewCol(LDAPRootCACol, idpEvent.RootCA),
// 				handler.NewCol(LDAPIDAttributeCol, idpEvent.IDAttribute),
// 				handler.NewCol(LDAPFirstNameAttributeCol, idpEvent.FirstNameAttribute),
// 				handler.NewCol(LDAPLastNameAttributeCol, idpEvent.LastNameAttribute),
// 				handler.NewCol(LDAPDisplayNameAttributeCol, idpEvent.DisplayNameAttribute),
// 				handler.NewCol(LDAPNickNameAttributeCol, idpEvent.NickNameAttribute),
// 				handler.NewCol(LDAPPreferredUsernameAttributeCol, idpEvent.PreferredUsernameAttribute),
// 				handler.NewCol(LDAPEmailAttributeCol, idpEvent.EmailAttribute),
// 				handler.NewCol(LDAPEmailVerifiedAttributeCol, idpEvent.EmailVerifiedAttribute),
// 				handler.NewCol(LDAPPhoneAttributeCol, idpEvent.PhoneAttribute),
// 				handler.NewCol(LDAPPhoneVerifiedAttributeCol, idpEvent.PhoneVerifiedAttribute),
// 				handler.NewCol(LDAPPreferredLanguageAttributeCol, idpEvent.PreferredLanguageAttribute),
// 				handler.NewCol(LDAPAvatarURLAttributeCol, idpEvent.AvatarURLAttribute),
// 				handler.NewCol(LDAPProfileAttributeCol, idpEvent.ProfileAttribute),
// 			},
// 			handler.WithTableSuffix(IDPTemplateLDAPSuffix),
// 		),
// 	), nil
// }

// func (p *idpTemplateProjection) reduceLDAPIDPChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.LDAPIDPChangedEvent
// 	switch e := event.(type) {
// 	case *org.LDAPIDPChangedEvent:
// 		idpEvent = e.LDAPIDPChangedEvent
// 	case *instance.LDAPIDPChangedEvent:
// 		idpEvent = e.LDAPIDPChangedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPChangedEventType, instance.LDAPIDPChangedEventType})
// 	}

// 	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
// 	ops = append(ops,
// 		handler.AddUpdateStatement(
// 			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
// 			[]handler.Condition{
// 				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
// 				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),
// 	)

// 	ldapCols := reduceLDAPIDPChangedColumns(idpEvent)
// 	if len(ldapCols) > 0 {
// 		ops = append(ops,
// 			handler.AddUpdateStatement(
// 				ldapCols,
// 				[]handler.Condition{
// 					handler.NewCond(LDAPIDCol, idpEvent.ID),
// 					handler.NewCond(LDAPInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				},
// 				handler.WithTableSuffix(IDPTemplateLDAPSuffix),
// 			),
// 		)
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		ops...,
// 	), nil
// }

// func (p *idpTemplateProjection) reduceSAMLIDPAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.SAMLIDPAddedEvent
// 	var idpOwnerType domain.IdentityProviderType
// 	switch e := event.(type) {
// 	case *org.SAMLIDPAddedEvent:
// 		idpEvent = e.SAMLIDPAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeOrg
// 	case *instance.SAMLIDPAddedEvent:
// 		idpEvent = e.SAMLIDPAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeSystem
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-9s02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.SAMLIDPAddedEventType, instance.SAMLIDPAddedEventType})
// 	}

// 	columns := []handler.Column{
// 		handler.NewCol(SAMLIDCol, idpEvent.ID),
// 		handler.NewCol(SAMLInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 		handler.NewCol(SAMLMetadataCol, idpEvent.Metadata),
// 		handler.NewCol(SAMLKeyCol, idpEvent.Key),
// 		handler.NewCol(SAMLCertificateCol, idpEvent.Certificate),
// 		handler.NewCol(SAMLBindingCol, idpEvent.Binding),
// 		handler.NewCol(SAMLWithSignedRequestCol, idpEvent.WithSignedRequest),
// 		handler.NewCol(SAMLTransientMappingAttributeName, idpEvent.TransientMappingAttributeName),
// 		handler.NewCol(SAMLFederatedLogoutEnabled, idpEvent.FederatedLogoutEnabled),
// 	}
// 	if idpEvent.NameIDFormat != nil {
// 		columns = append(columns, handler.NewCol(SAMLNameIDFormatCol, *idpEvent.NameIDFormat))
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
// 				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
// 				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
// 				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
// 				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
// 				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
// 				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeSAML),
// 				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
// 				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
// 				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
// 				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
// 				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
// 			},
// 		),
// 		handler.AddCreateStatement(
// 			columns,
// 			handler.WithTableSuffix(IDPTemplateSAMLSuffix),
// 		),
// 	), nil
// }

// func (p *idpTemplateProjection) reduceSAMLIDPChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.SAMLIDPChangedEvent
// 	switch e := event.(type) {
// 	case *org.SAMLIDPChangedEvent:
// 		idpEvent = e.SAMLIDPChangedEvent
// 	case *instance.SAMLIDPChangedEvent:
// 		idpEvent = e.SAMLIDPChangedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-o7c0fii4ad", "reduce.wrong.event.type %v", []eventstore.EventType{org.SAMLIDPChangedEventType, instance.SAMLIDPChangedEventType})
// 	}

// 	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
// 	ops = append(ops,
// 		handler.AddUpdateStatement(
// 			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
// 			[]handler.Condition{
// 				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
// 				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),
// 	)

// 	SAMLCols := reduceSAMLIDPChangedColumns(idpEvent)
// 	if len(SAMLCols) > 0 {
// 		ops = append(ops,
// 			handler.AddUpdateStatement(
// 				SAMLCols,
// 				[]handler.Condition{
// 					handler.NewCond(SAMLIDCol, idpEvent.ID),
// 					handler.NewCond(SAMLInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				},
// 				handler.WithTableSuffix(IDPTemplateSAMLSuffix),
// 			),
// 		)
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		ops...,
// 	), nil
// }

// func (p *idpTemplateProjection) reduceAppleIDPAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.AppleIDPAddedEvent
// 	var idpOwnerType domain.IdentityProviderType
// 	switch e := event.(type) {
// 	case *org.AppleIDPAddedEvent:
// 		idpEvent = e.AppleIDPAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeOrg
// 	case *instance.AppleIDPAddedEvent:
// 		idpEvent = e.AppleIDPAddedEvent
// 		idpOwnerType = domain.IdentityProviderTypeSystem
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SFvg3", "reduce.wrong.event.type %v", []eventstore.EventType{org.AppleIDPAddedEventType /*, instance.AppleIDPAddedEventType*/})
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
// 				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
// 				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
// 				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
// 				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
// 				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
// 				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
// 				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeApple),
// 				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
// 				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
// 				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
// 				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
// 				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
// 			},
// 		),
// 		handler.AddCreateStatement(
// 			[]handler.Column{
// 				handler.NewCol(AppleIDCol, idpEvent.ID),
// 				handler.NewCol(AppleInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				handler.NewCol(AppleClientIDCol, idpEvent.ClientID),
// 				handler.NewCol(AppleTeamIDCol, idpEvent.TeamID),
// 				handler.NewCol(AppleKeyIDCol, idpEvent.KeyID),
// 				handler.NewCol(ApplePrivateKeyCol, idpEvent.PrivateKey),
// 				handler.NewCol(AppleScopesCol, database.TextArray[string](idpEvent.Scopes)),
// 			},
// 			handler.WithTableSuffix(IDPTemplateAppleSuffix),
// 		),
// 	), nil
// }

// func (p *idpTemplateProjection) reduceAppleIDPChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.AppleIDPChangedEvent
// 	switch e := event.(type) {
// 	case *org.AppleIDPChangedEvent:
// 		idpEvent = e.AppleIDPChangedEvent
// 	case *instance.AppleIDPChangedEvent:
// 		idpEvent = e.AppleIDPChangedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-GBez3", "reduce.wrong.event.type %v", []eventstore.EventType{org.AppleIDPChangedEventType /*, instance.AppleIDPChangedEventType*/})
// 	}

// 	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
// 	ops = append(ops,
// 		handler.AddUpdateStatement(
// 			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
// 			[]handler.Condition{
// 				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
// 				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			},
// 		),
// 	)
// 	appleCols := reduceAppleIDPChangedColumns(idpEvent)
// 	if len(appleCols) > 0 {
// 		ops = append(ops,
// 			handler.AddUpdateStatement(
// 				appleCols,
// 				[]handler.Condition{
// 					handler.NewCond(AppleIDCol, idpEvent.ID),
// 					handler.NewCond(AppleInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 				},
// 				handler.WithTableSuffix(IDPTemplateAppleSuffix),
// 			),
// 		)
// 	}

// 	return handler.NewMultiStatement(
// 		&idpEvent,
// 		ops...,
// 	), nil
// }

// func (p *idpTemplateProjection) reduceIDPConfigRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idpconfig.IDPConfigRemovedEvent
// 	switch e := event.(type) {
// 	case *org.IDPConfigRemovedEvent:
// 		idpEvent = e.IDPConfigRemovedEvent
// 	case *instance.IDPConfigRemovedEvent:
// 		idpEvent = e.IDPConfigRemovedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SAFet", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigRemovedEventType, instance.IDPConfigRemovedEventType})
// 	}

// 	return handler.NewDeleteStatement(
// 		&idpEvent,
// 		[]handler.Condition{
// 			handler.NewCond(IDPTemplateIDCol, idpEvent.ConfigID),
// 			handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 		},
// 	), nil
// }

// func (p *idpTemplateProjection) reduceIDPRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	var idpEvent idp.RemovedEvent
// 	switch e := event.(type) {
// 	case *org.IDPRemovedEvent:
// 		idpEvent = e.RemovedEvent
// 	case *instance.IDPRemovedEvent:
// 		idpEvent = e.RemovedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-xbcvwin2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPRemovedEventType, instance.IDPRemovedEventType})
// 	}

// 	return handler.NewDeleteStatement(
// 		&idpEvent,
// 		[]handler.Condition{
// 			handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
// 			handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 		},
// 	), nil
// }

// func (p *idpTemplateProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	e, ok := event.(*org.OrgRemovedEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Jp0D2K", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
// 	}

// 	return handler.NewDeleteStatement(
// 		e,
// 		[]handler.Condition{
// 			handler.NewCond(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
// 			handler.NewCond(IDPTemplateResourceOwnerCol, e.Aggregate().ID),
// 		},
// 	), nil
// }

// func reduceIDPChangedTemplateColumns(name *string, creationDate time.Time, sequence uint64, optionChanges idp.OptionChanges) []handler.Column {
// 	cols := make([]handler.Column, 0, 7)
// 	if name != nil {
// 		cols = append(cols, handler.NewCol(IDPTemplateNameCol, *name))
// 	}
// 	if optionChanges.IsCreationAllowed != nil {
// 		cols = append(cols, handler.NewCol(IDPTemplateIsCreationAllowedCol, *optionChanges.IsCreationAllowed))
// 	}
// 	if optionChanges.IsLinkingAllowed != nil {
// 		cols = append(cols, handler.NewCol(IDPTemplateIsLinkingAllowedCol, *optionChanges.IsLinkingAllowed))
// 	}
// 	if optionChanges.IsAutoCreation != nil {
// 		cols = append(cols, handler.NewCol(IDPTemplateIsAutoCreationCol, *optionChanges.IsAutoCreation))
// 	}
// 	if optionChanges.IsAutoUpdate != nil {
// 		cols = append(cols, handler.NewCol(IDPTemplateIsAutoUpdateCol, *optionChanges.IsAutoUpdate))
// 	}
// 	if optionChanges.AutoLinkingOption != nil {
// 		cols = append(cols, handler.NewCol(IDPTemplateAutoLinkingCol, *optionChanges.AutoLinkingOption))
// 	}
// 	return append(cols,
// 		handler.NewCol(IDPTemplateChangeDateCol, creationDate),
// 		handler.NewCol(IDPTemplateSequenceCol, sequence),
// 	)
// }

// func reduceOAuthIDPChangedColumns(idpEvent idp.OAuthIDPChangedEvent) []handler.Column {
// 	oauthCols := make([]handler.Column, 0, 7)
// 	if idpEvent.ClientID != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(OAuthClientIDCol, *idpEvent.ClientID))
// 	}
// 	if idpEvent.ClientSecret != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(OAuthClientSecretCol, *idpEvent.ClientSecret))
// 	}
// 	if idpEvent.AuthorizationEndpoint != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(OAuthAuthorizationEndpointCol, *idpEvent.AuthorizationEndpoint))
// 	}
// 	if idpEvent.TokenEndpoint != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(OAuthTokenEndpointCol, *idpEvent.TokenEndpoint))
// 	}
// 	if idpEvent.UserEndpoint != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(OAuthUserEndpointCol, *idpEvent.UserEndpoint))
// 	}
// 	if idpEvent.Scopes != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(OAuthScopesCol, database.TextArray[string](idpEvent.Scopes)))
// 	}
// 	if idpEvent.IDAttribute != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(OAuthIDAttributeCol, *idpEvent.IDAttribute))
// 	}
// 	if idpEvent.UsePKCE != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(OAuthUsePKCECol, *idpEvent.UsePKCE))
// 	}
// 	return oauthCols
// }

// func reduceOIDCIDPChangedColumns(idpEvent idp.OIDCIDPChangedEvent) []handler.Column {
// 	oidcCols := make([]handler.Column, 0, 5)
// 	if idpEvent.ClientID != nil {
// 		oidcCols = append(oidcCols, handler.NewCol(OIDCClientIDCol, *idpEvent.ClientID))
// 	}
// 	if idpEvent.ClientSecret != nil {
// 		oidcCols = append(oidcCols, handler.NewCol(OIDCClientSecretCol, *idpEvent.ClientSecret))
// 	}
// 	if idpEvent.Issuer != nil {
// 		oidcCols = append(oidcCols, handler.NewCol(OIDCIssuerCol, *idpEvent.Issuer))
// 	}
// 	if idpEvent.Scopes != nil {
// 		oidcCols = append(oidcCols, handler.NewCol(OIDCScopesCol, database.TextArray[string](idpEvent.Scopes)))
// 	}
// 	if idpEvent.IsIDTokenMapping != nil {
// 		oidcCols = append(oidcCols, handler.NewCol(OIDCIDTokenMappingCol, *idpEvent.IsIDTokenMapping))
// 	}
// 	if idpEvent.UsePKCE != nil {
// 		oidcCols = append(oidcCols, handler.NewCol(OIDCUsePKCECol, *idpEvent.UsePKCE))
// 	}
// 	return oidcCols
// }

// func reduceJWTIDPChangedColumns(idpEvent idp.JWTIDPChangedEvent) []handler.Column {
// 	jwtCols := make([]handler.Column, 0, 4)
// 	if idpEvent.JWTEndpoint != nil {
// 		jwtCols = append(jwtCols, handler.NewCol(JWTEndpointCol, *idpEvent.JWTEndpoint))
// 	}
// 	if idpEvent.KeysEndpoint != nil {
// 		jwtCols = append(jwtCols, handler.NewCol(JWTKeysEndpointCol, *idpEvent.KeysEndpoint))
// 	}
// 	if idpEvent.HeaderName != nil {
// 		jwtCols = append(jwtCols, handler.NewCol(JWTHeaderNameCol, *idpEvent.HeaderName))
// 	}
// 	if idpEvent.Issuer != nil {
// 		jwtCols = append(jwtCols, handler.NewCol(JWTIssuerCol, *idpEvent.Issuer))
// 	}
// 	return jwtCols
// }

// func reduceAzureADIDPChangedColumns(idpEvent idp.AzureADIDPChangedEvent) []handler.Column {
// 	azureADCols := make([]handler.Column, 0, 5)
// 	if idpEvent.ClientID != nil {
// 		azureADCols = append(azureADCols, handler.NewCol(AzureADClientIDCol, *idpEvent.ClientID))
// 	}
// 	if idpEvent.ClientSecret != nil {
// 		azureADCols = append(azureADCols, handler.NewCol(AzureADClientSecretCol, *idpEvent.ClientSecret))
// 	}
// 	if idpEvent.Scopes != nil {
// 		azureADCols = append(azureADCols, handler.NewCol(AzureADScopesCol, database.TextArray[string](idpEvent.Scopes)))
// 	}
// 	if idpEvent.Tenant != nil {
// 		azureADCols = append(azureADCols, handler.NewCol(AzureADTenantCol, *idpEvent.Tenant))
// 	}
// 	if idpEvent.IsEmailVerified != nil {
// 		azureADCols = append(azureADCols, handler.NewCol(AzureADIsEmailVerified, *idpEvent.IsEmailVerified))
// 	}
// 	return azureADCols
// }

// func reduceGitHubIDPChangedColumns(idpEvent idp.GitHubIDPChangedEvent) []handler.Column {
// 	oauthCols := make([]handler.Column, 0, 3)
// 	if idpEvent.ClientID != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(GitHubClientIDCol, *idpEvent.ClientID))
// 	}
// 	if idpEvent.ClientSecret != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(GitHubClientSecretCol, *idpEvent.ClientSecret))
// 	}
// 	if idpEvent.Scopes != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(GitHubScopesCol, database.TextArray[string](idpEvent.Scopes)))
// 	}
// 	return oauthCols
// }

// func reduceGitHubEnterpriseIDPChangedColumns(idpEvent idp.GitHubEnterpriseIDPChangedEvent) []handler.Column {
// 	oauthCols := make([]handler.Column, 0, 6)
// 	if idpEvent.ClientID != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseClientIDCol, *idpEvent.ClientID))
// 	}
// 	if idpEvent.ClientSecret != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseClientSecretCol, *idpEvent.ClientSecret))
// 	}
// 	if idpEvent.AuthorizationEndpoint != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseAuthorizationEndpointCol, *idpEvent.AuthorizationEndpoint))
// 	}
// 	if idpEvent.TokenEndpoint != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseTokenEndpointCol, *idpEvent.TokenEndpoint))
// 	}
// 	if idpEvent.UserEndpoint != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseUserEndpointCol, *idpEvent.UserEndpoint))
// 	}
// 	if idpEvent.Scopes != nil {
// 		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseScopesCol, database.TextArray[string](idpEvent.Scopes)))
// 	}
// 	return oauthCols
// }

// func reduceGitLabIDPChangedColumns(idpEvent idp.GitLabIDPChangedEvent) []handler.Column {
// 	gitlabCols := make([]handler.Column, 0, 3)
// 	if idpEvent.ClientID != nil {
// 		gitlabCols = append(gitlabCols, handler.NewCol(GitLabClientIDCol, *idpEvent.ClientID))
// 	}
// 	if idpEvent.ClientSecret != nil {
// 		gitlabCols = append(gitlabCols, handler.NewCol(GitLabClientSecretCol, *idpEvent.ClientSecret))
// 	}
// 	if idpEvent.Scopes != nil {
// 		gitlabCols = append(gitlabCols, handler.NewCol(GitLabScopesCol, database.TextArray[string](idpEvent.Scopes)))
// 	}
// 	return gitlabCols
// }

// func reduceGitLabSelfHostedIDPChangedColumns(idpEvent idp.GitLabSelfHostedIDPChangedEvent) []handler.Column {
// 	gitlabCols := make([]handler.Column, 0, 4)
// 	if idpEvent.Issuer != nil {
// 		gitlabCols = append(gitlabCols, handler.NewCol(GitLabSelfHostedIssuerCol, *idpEvent.Issuer))
// 	}
// 	if idpEvent.ClientID != nil {
// 		gitlabCols = append(gitlabCols, handler.NewCol(GitLabSelfHostedClientIDCol, *idpEvent.ClientID))
// 	}
// 	if idpEvent.ClientSecret != nil {
// 		gitlabCols = append(gitlabCols, handler.NewCol(GitLabSelfHostedClientSecretCol, *idpEvent.ClientSecret))
// 	}
// 	if idpEvent.Scopes != nil {
// 		gitlabCols = append(gitlabCols, handler.NewCol(GitLabSelfHostedScopesCol, database.TextArray[string](idpEvent.Scopes)))
// 	}
// 	return gitlabCols
// }

// func reduceGoogleIDPChangedColumns(idpEvent idp.GoogleIDPChangedEvent) []handler.Column {
// 	googleCols := make([]handler.Column, 0, 3)
// 	if idpEvent.ClientID != nil {
// 		googleCols = append(googleCols, handler.NewCol(GoogleClientIDCol, *idpEvent.ClientID))
// 	}
// 	if idpEvent.ClientSecret != nil {
// 		googleCols = append(googleCols, handler.NewCol(GoogleClientSecretCol, *idpEvent.ClientSecret))
// 	}
// 	if idpEvent.Scopes != nil {
// 		googleCols = append(googleCols, handler.NewCol(GoogleScopesCol, database.TextArray[string](idpEvent.Scopes)))
// 	}
// 	return googleCols
// }

// func reduceLDAPIDPChangedColumns(idpEvent idp.LDAPIDPChangedEvent) []handler.Column {
// 	ldapCols := make([]handler.Column, 0, 22)
// 	if idpEvent.Servers != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPServersCol, database.TextArray[string](idpEvent.Servers)))
// 	}
// 	if idpEvent.StartTLS != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPStartTLSCol, *idpEvent.StartTLS))
// 	}
// 	if idpEvent.BaseDN != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPBaseDNCol, *idpEvent.BaseDN))
// 	}
// 	if idpEvent.BindDN != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPBindDNCol, *idpEvent.BindDN))
// 	}
// 	if idpEvent.BindPassword != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPBindPasswordCol, idpEvent.BindPassword))
// 	}
// 	if idpEvent.UserBase != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPUserBaseCol, *idpEvent.UserBase))
// 	}
// 	if idpEvent.UserObjectClasses != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPUserObjectClassesCol, database.TextArray[string](idpEvent.UserObjectClasses)))
// 	}
// 	if idpEvent.UserFilters != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPUserFiltersCol, database.TextArray[string](idpEvent.UserFilters)))
// 	}
// 	if idpEvent.Timeout != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPTimeoutCol, *idpEvent.Timeout))
// 	}
// 	if idpEvent.RootCA != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPRootCACol, idpEvent.RootCA))
// 	}
// 	if idpEvent.IDAttribute != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPIDAttributeCol, *idpEvent.IDAttribute))
// 	}
// 	if idpEvent.FirstNameAttribute != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPFirstNameAttributeCol, *idpEvent.FirstNameAttribute))
// 	}
// 	if idpEvent.LastNameAttribute != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPLastNameAttributeCol, *idpEvent.LastNameAttribute))
// 	}
// 	if idpEvent.DisplayNameAttribute != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPDisplayNameAttributeCol, *idpEvent.DisplayNameAttribute))
// 	}
// 	if idpEvent.NickNameAttribute != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPNickNameAttributeCol, *idpEvent.NickNameAttribute))
// 	}
// 	if idpEvent.PreferredUsernameAttribute != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPPreferredUsernameAttributeCol, *idpEvent.PreferredUsernameAttribute))
// 	}
// 	if idpEvent.EmailAttribute != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPEmailAttributeCol, *idpEvent.EmailAttribute))
// 	}
// 	if idpEvent.EmailVerifiedAttribute != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPEmailVerifiedAttributeCol, *idpEvent.EmailVerifiedAttribute))
// 	}
// 	if idpEvent.PhoneAttribute != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPPhoneAttributeCol, *idpEvent.PhoneAttribute))
// 	}
// 	if idpEvent.PhoneVerifiedAttribute != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPPhoneVerifiedAttributeCol, *idpEvent.PhoneVerifiedAttribute))
// 	}
// 	if idpEvent.PreferredLanguageAttribute != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPPreferredLanguageAttributeCol, *idpEvent.PreferredLanguageAttribute))
// 	}
// 	if idpEvent.AvatarURLAttribute != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPAvatarURLAttributeCol, *idpEvent.AvatarURLAttribute))
// 	}
// 	if idpEvent.ProfileAttribute != nil {
// 		ldapCols = append(ldapCols, handler.NewCol(LDAPProfileAttributeCol, *idpEvent.ProfileAttribute))
// 	}
// 	return ldapCols
// }

// func reduceAppleIDPChangedColumns(idpEvent idp.AppleIDPChangedEvent) []handler.Column {
// 	appleCols := make([]handler.Column, 0, 5)
// 	if idpEvent.ClientID != nil {
// 		appleCols = append(appleCols, handler.NewCol(AppleClientIDCol, *idpEvent.ClientID))
// 	}
// 	if idpEvent.TeamID != nil {
// 		appleCols = append(appleCols, handler.NewCol(AppleTeamIDCol, *idpEvent.TeamID))
// 	}
// 	if idpEvent.KeyID != nil {
// 		appleCols = append(appleCols, handler.NewCol(AppleKeyIDCol, *idpEvent.KeyID))
// 	}
// 	if idpEvent.PrivateKey != nil {
// 		appleCols = append(appleCols, handler.NewCol(ApplePrivateKeyCol, *idpEvent.PrivateKey))
// 	}
// 	if idpEvent.Scopes != nil {
// 		appleCols = append(appleCols, handler.NewCol(AppleScopesCol, database.TextArray[string](idpEvent.Scopes)))
// 	}
// 	return appleCols
// }

// func reduceSAMLIDPChangedColumns(idpEvent idp.SAMLIDPChangedEvent) []handler.Column {
// 	SAMLCols := make([]handler.Column, 0, 5)
// 	if idpEvent.Metadata != nil {
// 		SAMLCols = append(SAMLCols, handler.NewCol(SAMLMetadataCol, idpEvent.Metadata))
// 	}
// 	if idpEvent.Key != nil {
// 		SAMLCols = append(SAMLCols, handler.NewCol(SAMLKeyCol, idpEvent.Key))
// 	}
// 	if idpEvent.Certificate != nil {
// 		SAMLCols = append(SAMLCols, handler.NewCol(SAMLCertificateCol, idpEvent.Certificate))
// 	}
// 	if idpEvent.Binding != nil {
// 		SAMLCols = append(SAMLCols, handler.NewCol(SAMLBindingCol, *idpEvent.Binding))
// 	}
// 	if idpEvent.WithSignedRequest != nil {
// 		SAMLCols = append(SAMLCols, handler.NewCol(SAMLWithSignedRequestCol, *idpEvent.WithSignedRequest))
// 	}
// 	if idpEvent.NameIDFormat != nil {
// 		SAMLCols = append(SAMLCols, handler.NewCol(SAMLNameIDFormatCol, *idpEvent.NameIDFormat))
// 	}
// 	if idpEvent.TransientMappingAttributeName != nil {
// 		SAMLCols = append(SAMLCols, handler.NewCol(SAMLTransientMappingAttributeName, *idpEvent.TransientMappingAttributeName))
// 	}
// 	if idpEvent.FederatedLogoutEnabled != nil {
// 		SAMLCols = append(SAMLCols, handler.NewCol(SAMLFederatedLogoutEnabled, *idpEvent.FederatedLogoutEnabled))
// 	}
// 	return SAMLCols
// }

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
		*cols = append(*cols, handler.NewCol(IDPRelationalAllowAutoLinkingCol, domain.IDPAutoLinkingOption(*optionChanges.AutoLinkingOption).String()))
	}
}

func reduceOAuthIDPRelationalChangedColumns(payload *domain.OAuth, idpEvent *idp.OAuthIDPChangedEvent) bool {
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

func reduceOIDCIDPRelationalChangedColumns(payload *domain.OIDC, idpEvent *idp.OIDCIDPChangedEvent) bool {
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

func reduceJWTIDPRelationalChangedColumns(payload *domain.JWT, idpEvent *idp.JWTIDPChangedEvent) bool {
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

func reduceAzureADIDPRelationalChangedColumns(payload *domain.Azure, idpEvent *idp.AzureADIDPChangedEvent) bool {
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

func reduceGitHubIDPRelationalChangedColumns(payload *domain.Github, idpEvent *idp.GitHubIDPChangedEvent) bool {
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
