package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

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

func (wm *InstanceLDAPIDPWriteModel) Reduce() error {
	return wm.LDAPIDPWriteModel.Reduce()
}

func (wm *InstanceLDAPIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.LDAPIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.LDAPIDPWriteModel.AppendEvents(&e.LDAPIDPAddedEvent)
		case *instance.LDAPIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.LDAPIDPWriteModel.AppendEvents(&e.LDAPIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
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
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewLDAPIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
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

func (wm *InstanceIDPRemoveWriteModel) Reduce() error {
	return wm.IDPRemoveWriteModel.Reduce()
}

func (wm *InstanceIDPRemoveWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.LDAPIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.LDAPIDPAddedEvent)
		case *instance.LDAPIDPChangedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.LDAPIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.RemovedEvent)
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
			instance.LDAPIDPAddedEventType,
			instance.LDAPIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		Builder()
}

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

func (wm *InstanceOAuthIDPWriteModel) Reduce() error {
	return wm.OAuthIDPWriteModel.Reduce()
}

func (wm *InstanceOAuthIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.OAuthIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OAuthIDPWriteModel.AppendEvents(&e.OAuthIDPAddedEvent)
		case *instance.OAuthIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OAuthIDPWriteModel.AppendEvents(&e.OAuthIDPChangedEvent)
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
		).
		Builder()
}

func (wm *InstanceOAuthIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	oldName,
	name,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	authorizationEndpoint,
	tokenEndpoint,
	userEndpoint string,
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
		scopes,
		options,
	)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewOAuthIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
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

func (wm *InstanceOIDCIDPWriteModel) Reduce() error {
	return wm.OIDCIDPWriteModel.Reduce()
}

func (wm *InstanceOIDCIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.OIDCIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCIDPAddedEvent)
		case *instance.OIDCIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.RemovedEvent)
		case *instance.IDPConfigAddedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *instance.IDPOIDCConfigAddedEvent:
			if wm.ID != e.IDPConfigID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCConfigAddedEvent)
		case *instance.IDPOIDCConfigChangedEvent:
			if wm.ID != e.IDPConfigID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCConfigChangedEvent)
		case *instance.IDPConfigRemovedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
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
			instance.IDPConfigAddedEventType,
			instance.IDPOIDCConfigAddedEventType,
			instance.IDPOIDCConfigChangedEventType,
			instance.IDPConfigRemovedEventType,
		).
		Builder()
}

func (wm *InstanceOIDCIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	oldName,
	name,
	issuer,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*instance.OIDCIDPChangedEvent, error) {

	changes, err := wm.OIDCIDPWriteModel.NewChanges(
		name,
		issuer,
		clientID,
		clientSecretString,
		secretCrypto,
		scopes,
		options,
	)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewOIDCIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
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

func (wm *InstanceJWTIDPWriteModel) Reduce() error {
	return wm.JWTIDPWriteModel.Reduce()
}

func (wm *InstanceJWTIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.JWTIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTIDPAddedEvent)
		case *instance.JWTIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.RemovedEvent)
		case *instance.IDPConfigAddedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *instance.IDPJWTConfigAddedEvent:
			if wm.ID != e.IDPConfigID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTConfigAddedEvent)
		case *instance.IDPJWTConfigChangedEvent:
			if wm.ID != e.IDPConfigID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTConfigChangedEvent)
		case *instance.IDPConfigRemovedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
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
			instance.IDPConfigAddedEventType,
			instance.IDPJWTConfigAddedEventType,
			instance.IDPJWTConfigChangedEventType,
			instance.IDPConfigRemovedEventType,
		).
		Builder()
}

func (wm *InstanceJWTIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	oldName,
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
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewJWTIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
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

func (wm *InstanceAzureADIDPWriteModel) Reduce() error {
	return wm.AzureADIDPWriteModel.Reduce()
}

func (wm *InstanceAzureADIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.AzureADIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.AzureADIDPWriteModel.AppendEvents(&e.AzureADIDPAddedEvent)
		case *instance.AzureADIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.AzureADIDPWriteModel.AppendEvents(&e.AzureADIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
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
		Builder()
}

func (wm *InstanceAzureADIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	oldName,
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
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewAzureADIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
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

func (wm *InstanceGitHubIDPWriteModel) Reduce() error {
	return wm.GitHubIDPWriteModel.Reduce()
}

func (wm *InstanceGitHubIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.GitHubIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitHubIDPWriteModel.AppendEvents(&e.GitHubIDPAddedEvent)
		case *instance.GitHubIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitHubIDPWriteModel.AppendEvents(&e.GitHubIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
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
		Builder()
}

func (wm *InstanceGitHubIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*instance.GitHubIDPChangedEvent, error) {

	changes, err := wm.GitHubIDPWriteModel.NewChanges(clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewGitHubIDPChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
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

func (wm *InstanceGitHubEnterpriseIDPWriteModel) Reduce() error {
	return wm.GitHubEnterpriseIDPWriteModel.Reduce()
}

func (wm *InstanceGitHubEnterpriseIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.GitHubEnterpriseIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitHubEnterpriseIDPWriteModel.AppendEvents(&e.GitHubEnterpriseIDPAddedEvent)
		case *instance.GitHubEnterpriseIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitHubEnterpriseIDPWriteModel.AppendEvents(&e.GitHubEnterpriseIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
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
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewGitHubEnterpriseIDPChangedEvent(ctx, aggregate, id, wm.Name, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
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

func (wm *InstanceGitLabIDPWriteModel) Reduce() error {
	return wm.GitLabIDPWriteModel.Reduce()
}

func (wm *InstanceGitLabIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.GitLabIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitLabIDPWriteModel.AppendEvents(&e.GitLabIDPAddedEvent)
		case *instance.GitLabIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitLabIDPWriteModel.AppendEvents(&e.GitLabIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
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
		Builder()
}

func (wm *InstanceGitLabIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*instance.GitLabIDPChangedEvent, error) {

	changes, err := wm.GitLabIDPWriteModel.NewChanges(clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewGitLabIDPChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
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

func (wm *InstanceGitLabSelfHostedIDPWriteModel) Reduce() error {
	return wm.GitLabSelfHostedIDPWriteModel.Reduce()
}

func (wm *InstanceGitLabSelfHostedIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.GitLabSelfHostedIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitLabSelfHostedIDPWriteModel.AppendEvents(&e.GitLabSelfHostedIDPAddedEvent)
		case *instance.GitLabSelfHostedIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitLabSelfHostedIDPWriteModel.AppendEvents(&e.GitLabSelfHostedIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
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
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewGitLabSelfHostedIDPChangedEvent(ctx, aggregate, id, wm.Name, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
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

func (wm *InstanceGoogleIDPWriteModel) Reduce() error {
	return wm.GoogleIDPWriteModel.Reduce()
}

func (wm *InstanceGoogleIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.GoogleIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPAddedEvent)
		case *instance.GoogleIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPChangedEvent)
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
		).
		Builder()
}

func (wm *InstanceGoogleIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*instance.GoogleIDPChangedEvent, error) {

	changes, err := wm.GoogleIDPWriteModel.NewChanges(clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewGoogleIDPChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}
