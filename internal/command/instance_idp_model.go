package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceOAuthIDPWriteModel struct {
	OAuthIDPWriteModel
}

func NewOAuthInstanceIDPWriteModel(instanceID, id string) *InstanceOAuthIDPWriteModel {
	return &InstanceOAuthIDPWriteModel{
		OAuthIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceOAuthIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.OAuthIDPAddedEvent:
			wm.OAuthIDPWriteModel.AppendEvents(&e.OAuthIDPAddedEvent)
		case *instance.OAuthIDPChangedEvent:
			wm.OAuthIDPWriteModel.AppendEvents(&e.OAuthIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.OAuthIDPWriteModel.AppendEvents(&e.RemovedEvent)
		}
	}
}

func (wm *InstanceOAuthIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.OAuthIDPAddedEventType,
			instance.OAuthIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Builder()
}

func (wm *InstanceOAuthIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	authorizationEndpoint,
	tokenEndpoint,
	userEndpoint,
	idAttribute string,
	scopes []string,
	options idp.Options,
) (*instance.OAuthIDPChangedEvent, error) {

	changes, err := wm.OAuthIDPWriteModel.NewChanges(
		name,
		clientID,
		clientSecretString,
		secretCrypto,
		authorizationEndpoint,
		tokenEndpoint,
		userEndpoint,
		idAttribute,
		scopes,
		options,
	)
	if err != nil || len(changes) == 0 {
		return nil, err
	}
	return instance.NewOAuthIDPChangedEvent(ctx, aggregate, id, changes)
}

type InstanceOIDCIDPWriteModel struct {
	OIDCIDPWriteModel
}

func NewOIDCInstanceIDPWriteModel(instanceID, id string) *InstanceOIDCIDPWriteModel {
	return &InstanceOIDCIDPWriteModel{
		OIDCIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceOIDCIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.OIDCIDPAddedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCIDPAddedEvent)
		case *instance.OIDCIDPChangedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.RemovedEvent)

			// old events
		case *instance.IDPConfigAddedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *instance.IDPConfigChangedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.IDPConfigChangedEvent)
		case *instance.IDPOIDCConfigAddedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCConfigAddedEvent)
		case *instance.IDPOIDCConfigChangedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCConfigChangedEvent)
		case *instance.IDPConfigRemovedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.OIDCIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceOIDCIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.OIDCIDPAddedEventType,
			instance.OIDCIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Or(). // old events
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.IDPConfigAddedEventType,
			instance.IDPConfigChangedEventType,
			instance.IDPOIDCConfigAddedEventType,
			instance.IDPOIDCConfigChangedEventType,
			instance.IDPConfigRemovedEventType,
		).
		EventData(map[string]interface{}{"idpConfigId": wm.ID}).
		Builder()
}

func (wm *InstanceOIDCIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	issuer,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	idTokenMapping bool,
	options idp.Options,
) (*instance.OIDCIDPChangedEvent, error) {

	changes, err := wm.OIDCIDPWriteModel.NewChanges(
		name,
		issuer,
		clientID,
		clientSecretString,
		secretCrypto,
		scopes,
		idTokenMapping,
		options,
	)
	if err != nil || len(changes) == 0 {
		return nil, err
	}
	return instance.NewOIDCIDPChangedEvent(ctx, aggregate, id, changes)
}

type InstanceJWTIDPWriteModel struct {
	JWTIDPWriteModel
}

func NewJWTInstanceIDPWriteModel(instanceID, id string) *InstanceJWTIDPWriteModel {
	return &InstanceJWTIDPWriteModel{
		JWTIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceJWTIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.JWTIDPAddedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTIDPAddedEvent)
		case *instance.JWTIDPChangedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.RemovedEvent)

			// old events
		case *instance.IDPConfigAddedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *instance.IDPConfigChangedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.IDPConfigChangedEvent)
		case *instance.IDPJWTConfigAddedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTConfigAddedEvent)
		case *instance.IDPJWTConfigChangedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTConfigChangedEvent)
		case *instance.IDPConfigRemovedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.JWTIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceJWTIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.JWTIDPAddedEventType,
			instance.JWTIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Or(). // old events
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.IDPConfigAddedEventType,
			instance.IDPConfigChangedEventType,
			instance.IDPJWTConfigAddedEventType,
			instance.IDPJWTConfigChangedEventType,
			instance.IDPConfigRemovedEventType,
		).
		EventData(map[string]interface{}{"idpConfigId": wm.ID}).
		Builder()
}

func (wm *InstanceJWTIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	issuer,
	jwtEndpoint,
	keysEndpoint,
	headerName string,
	options idp.Options,
) (*instance.JWTIDPChangedEvent, error) {

	changes, err := wm.JWTIDPWriteModel.NewChanges(
		name,
		issuer,
		jwtEndpoint,
		keysEndpoint,
		headerName,
		options,
	)
	if err != nil || len(changes) == 0 {
		return nil, err
	}
	return instance.NewJWTIDPChangedEvent(ctx, aggregate, id, changes)
}

type InstanceAzureADIDPWriteModel struct {
	AzureADIDPWriteModel
}

func NewAzureADInstanceIDPWriteModel(instanceID, id string) *InstanceAzureADIDPWriteModel {
	return &InstanceAzureADIDPWriteModel{
		AzureADIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceAzureADIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.AzureADIDPAddedEvent:
			wm.AzureADIDPWriteModel.AppendEvents(&e.AzureADIDPAddedEvent)
		case *instance.AzureADIDPChangedEvent:
			wm.AzureADIDPWriteModel.AppendEvents(&e.AzureADIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.AzureADIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.AzureADIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceAzureADIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.AzureADIDPAddedEventType,
			instance.AzureADIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Builder()
}

func (wm *InstanceAzureADIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	tenant string,
	isEmailVerified bool,
	options idp.Options,
) (*instance.AzureADIDPChangedEvent, error) {

	changes, err := wm.AzureADIDPWriteModel.NewChanges(
		name,
		clientID,
		clientSecretString,
		secretCrypto,
		scopes,
		tenant,
		isEmailVerified,
		options,
	)
	if err != nil || len(changes) == 0 {
		return nil, err
	}
	return instance.NewAzureADIDPChangedEvent(ctx, aggregate, id, changes)
}

type InstanceGitHubIDPWriteModel struct {
	GitHubIDPWriteModel
}

func NewGitHubInstanceIDPWriteModel(instanceID, id string) *InstanceGitHubIDPWriteModel {
	return &InstanceGitHubIDPWriteModel{
		GitHubIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceGitHubIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.GitHubIDPAddedEvent:
			wm.GitHubIDPWriteModel.AppendEvents(&e.GitHubIDPAddedEvent)
		case *instance.GitHubIDPChangedEvent:
			wm.GitHubIDPWriteModel.AppendEvents(&e.GitHubIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.GitHubIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.GitHubIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceGitHubIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.GitHubIDPAddedEventType,
			instance.GitHubIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Builder()
}

func (wm *InstanceGitHubIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*instance.GitHubIDPChangedEvent, error) {

	changes, err := wm.GitHubIDPWriteModel.NewChanges(name, clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil || len(changes) == 0 {
		return nil, err
	}
	return instance.NewGitHubIDPChangedEvent(ctx, aggregate, id, changes)
}

type InstanceGitHubEnterpriseIDPWriteModel struct {
	GitHubEnterpriseIDPWriteModel
}

func NewGitHubEnterpriseInstanceIDPWriteModel(instanceID, id string) *InstanceGitHubEnterpriseIDPWriteModel {
	return &InstanceGitHubEnterpriseIDPWriteModel{
		GitHubEnterpriseIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceGitHubEnterpriseIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.GitHubEnterpriseIDPAddedEvent:
			wm.GitHubEnterpriseIDPWriteModel.AppendEvents(&e.GitHubEnterpriseIDPAddedEvent)
		case *instance.GitHubEnterpriseIDPChangedEvent:
			wm.GitHubEnterpriseIDPWriteModel.AppendEvents(&e.GitHubEnterpriseIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.GitHubEnterpriseIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.GitHubEnterpriseIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceGitHubEnterpriseIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.GitHubEnterpriseIDPAddedEventType,
			instance.GitHubEnterpriseIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Builder()
}

func (wm *InstanceGitHubEnterpriseIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	authorizationEndpoint,
	tokenEndpoint,
	userEndpoint string,
	scopes []string,
	options idp.Options,
) (*instance.GitHubEnterpriseIDPChangedEvent, error) {

	changes, err := wm.GitHubEnterpriseIDPWriteModel.NewChanges(
		name,
		clientID,
		clientSecretString,
		secretCrypto,
		authorizationEndpoint,
		tokenEndpoint,
		userEndpoint,
		scopes,
		options,
	)
	if err != nil || len(changes) == 0 {
		return nil, err
	}
	return instance.NewGitHubEnterpriseIDPChangedEvent(ctx, aggregate, id, changes)
}

type InstanceGitLabIDPWriteModel struct {
	GitLabIDPWriteModel
}

func NewGitLabInstanceIDPWriteModel(instanceID, id string) *InstanceGitLabIDPWriteModel {
	return &InstanceGitLabIDPWriteModel{
		GitLabIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceGitLabIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.GitLabIDPAddedEvent:
			wm.GitLabIDPWriteModel.AppendEvents(&e.GitLabIDPAddedEvent)
		case *instance.GitLabIDPChangedEvent:
			wm.GitLabIDPWriteModel.AppendEvents(&e.GitLabIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.GitLabIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.GitLabIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceGitLabIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.GitLabIDPAddedEventType,
			instance.GitLabIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Builder()
}

func (wm *InstanceGitLabIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*instance.GitLabIDPChangedEvent, error) {

	changes, err := wm.GitLabIDPWriteModel.NewChanges(name, clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil || len(changes) == 0 {
		return nil, err
	}
	return instance.NewGitLabIDPChangedEvent(ctx, aggregate, id, changes)
}

type InstanceGitLabSelfHostedIDPWriteModel struct {
	GitLabSelfHostedIDPWriteModel
}

func NewGitLabSelfHostedInstanceIDPWriteModel(instanceID, id string) *InstanceGitLabSelfHostedIDPWriteModel {
	return &InstanceGitLabSelfHostedIDPWriteModel{
		GitLabSelfHostedIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceGitLabSelfHostedIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.GitLabSelfHostedIDPAddedEvent:
			wm.GitLabSelfHostedIDPWriteModel.AppendEvents(&e.GitLabSelfHostedIDPAddedEvent)
		case *instance.GitLabSelfHostedIDPChangedEvent:
			wm.GitLabSelfHostedIDPWriteModel.AppendEvents(&e.GitLabSelfHostedIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.GitLabSelfHostedIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.GitLabSelfHostedIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceGitLabSelfHostedIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.GitLabSelfHostedIDPAddedEventType,
			instance.GitLabSelfHostedIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Builder()
}

func (wm *InstanceGitLabSelfHostedIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	issuer,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*instance.GitLabSelfHostedIDPChangedEvent, error) {

	changes, err := wm.GitLabSelfHostedIDPWriteModel.NewChanges(name, issuer, clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil || len(changes) == 0 {
		return nil, err
	}
	return instance.NewGitLabSelfHostedIDPChangedEvent(ctx, aggregate, id, changes)
}

type InstanceGoogleIDPWriteModel struct {
	GoogleIDPWriteModel
}

func NewGoogleInstanceIDPWriteModel(instanceID, id string) *InstanceGoogleIDPWriteModel {
	return &InstanceGoogleIDPWriteModel{
		GoogleIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceGoogleIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.GoogleIDPAddedEvent:
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPAddedEvent)
		case *instance.GoogleIDPChangedEvent:
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.GoogleIDPWriteModel.AppendEvents(&e.RemovedEvent)
		}
	}
}

func (wm *InstanceGoogleIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.GoogleIDPAddedEventType,
			instance.GoogleIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Builder()
}

func (wm *InstanceGoogleIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*instance.GoogleIDPChangedEvent, error) {

	changes, err := wm.GoogleIDPWriteModel.NewChanges(name, clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil || len(changes) == 0 {
		return nil, err
	}
	return instance.NewGoogleIDPChangedEvent(ctx, aggregate, id, changes)
}

type InstanceLDAPIDPWriteModel struct {
	LDAPIDPWriteModel
}

func NewLDAPInstanceIDPWriteModel(instanceID, id string) *InstanceLDAPIDPWriteModel {
	return &InstanceLDAPIDPWriteModel{
		LDAPIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceLDAPIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.LDAPIDPAddedEvent:
			wm.LDAPIDPWriteModel.AppendEvents(&e.LDAPIDPAddedEvent)
		case *instance.LDAPIDPChangedEvent:
			wm.LDAPIDPWriteModel.AppendEvents(&e.LDAPIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.LDAPIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.LDAPIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceLDAPIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.LDAPIDPAddedEventType,
			instance.LDAPIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Builder()
}

func (wm *InstanceLDAPIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	oldName,
	name,
	host,
	port string,
	tls bool,
	baseDN,
	userObjectClass,
	userUniqueAttribute,
	admin string,
	password string,
	secretCrypto crypto.Crypto,
	attributes idp.LDAPAttributes,
	options idp.Options,
) (*instance.LDAPIDPChangedEvent, error) {

	changes, err := wm.LDAPIDPWriteModel.NewChanges(
		name,
		host,
		port,
		tls,
		baseDN,
		userObjectClass,
		userUniqueAttribute,
		admin,
		password,
		secretCrypto,
		attributes,
		options,
	)
	if err != nil || len(changes) == 0 {
		return nil, err
	}
	return instance.NewLDAPIDPChangedEvent(ctx, aggregate, id, oldName, changes)
}

type InstanceIDPRemoveWriteModel struct {
	IDPRemoveWriteModel
}

func NewInstanceIDPRemoveWriteModel(instanceID, id string) *InstanceIDPRemoveWriteModel {
	return &InstanceIDPRemoveWriteModel{
		IDPRemoveWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceIDPRemoveWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.OAuthIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.OAuthIDPAddedEvent)
		case *instance.OIDCIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.OIDCIDPAddedEvent)
		case *instance.JWTIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.JWTIDPAddedEvent)
		case *instance.AzureADIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.AzureADIDPAddedEvent)
		case *instance.GitHubIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GitHubIDPAddedEvent)
		case *instance.GitHubEnterpriseIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GitHubEnterpriseIDPAddedEvent)
		case *instance.GitLabIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GitLabIDPAddedEvent)
		case *instance.GitLabSelfHostedIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GitLabSelfHostedIDPAddedEvent)
		case *instance.GoogleIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GoogleIDPAddedEvent)
		case *instance.LDAPIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.LDAPIDPAddedEvent)
		case *instance.IDPRemovedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.RemovedEvent)
		case *instance.IDPConfigAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *instance.IDPConfigRemovedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.IDPRemoveWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceIDPRemoveWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.OAuthIDPAddedEventType,
			instance.OIDCIDPAddedEventType,
			instance.JWTIDPAddedEventType,
			instance.AzureADIDPAddedEventType,
			instance.GitHubIDPAddedEventType,
			instance.GitHubEnterpriseIDPAddedEventType,
			instance.GitLabIDPAddedEventType,
			instance.GitLabSelfHostedIDPAddedEventType,
			instance.GoogleIDPAddedEventType,
			instance.LDAPIDPAddedEventType,
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Or(). // old events
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.IDPConfigAddedEventType,
			instance.IDPConfigRemovedEventType,
		).
		EventData(map[string]interface{}{"idpConfigId": wm.ID}).
		Builder()
}
