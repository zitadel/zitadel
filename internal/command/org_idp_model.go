package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type OrgLDAPIDPWriteModel struct {
	LDAPIDPWriteModel
}

func NewLDAPOrgIDPWriteModel(orgID, id string) *OrgLDAPIDPWriteModel {
	return &OrgLDAPIDPWriteModel{
		LDAPIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgLDAPIDPWriteModel) Reduce() error {
	return wm.LDAPIDPWriteModel.Reduce()
}

func (wm *OrgLDAPIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.LDAPIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.LDAPIDPWriteModel.AppendEvents(&e.LDAPIDPAddedEvent)
		case *org.LDAPIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.LDAPIDPWriteModel.AppendEvents(&e.LDAPIDPChangedEvent)
		case *org.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.LDAPIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.LDAPIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgLDAPIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.LDAPIDPAddedEventType,
			org.LDAPIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		Builder()
}

func (wm *OrgLDAPIDPWriteModel) NewChangedEvent(
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
) (*org.LDAPIDPChangedEvent, error) {

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
	changeEvent, err := org.NewLDAPIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgOAuthIDPWriteModel struct {
	OAuthIDPWriteModel
}

func NewOAuthOrgIDPWriteModel(orgID, id string) *OrgOAuthIDPWriteModel {
	return &OrgOAuthIDPWriteModel{
		OAuthIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgOAuthIDPWriteModel) Reduce() error {
	return wm.OAuthIDPWriteModel.Reduce()
}

func (wm *OrgOAuthIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.OAuthIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OAuthIDPWriteModel.AppendEvents(&e.OAuthIDPAddedEvent)
		case *org.OAuthIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OAuthIDPWriteModel.AppendEvents(&e.OAuthIDPChangedEvent)
		case *org.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OAuthIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.OAuthIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgOAuthIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.OAuthIDPAddedEventType,
			org.OAuthIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		Builder()
}

func (wm *OrgOAuthIDPWriteModel) NewChangedEvent(
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
) (*org.OAuthIDPChangedEvent, error) {

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
	changeEvent, err := org.NewOAuthIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgOIDCIDPWriteModel struct {
	OIDCIDPWriteModel
}

func NewOIDCOrgIDPWriteModel(orgID, id string) *OrgOIDCIDPWriteModel {
	return &OrgOIDCIDPWriteModel{
		OIDCIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgOIDCIDPWriteModel) Reduce() error {
	return wm.OIDCIDPWriteModel.Reduce()
}

func (wm *OrgOIDCIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.OIDCIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCIDPAddedEvent)
		case *org.OIDCIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCIDPChangedEvent)
		case *org.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.RemovedEvent)
		case *org.IDPConfigAddedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *org.IDPOIDCConfigAddedEvent:
			if wm.ID != e.IDPConfigID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCConfigAddedEvent)
		case *org.IDPOIDCConfigChangedEvent:
			if wm.ID != e.IDPConfigID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCConfigChangedEvent)
		case *org.IDPConfigRemovedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.OIDCIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgOIDCIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.OIDCIDPAddedEventType,
			org.OIDCIDPChangedEventType,
			org.IDPRemovedEventType,
			org.IDPConfigAddedEventType,
			org.IDPOIDCConfigAddedEventType,
			org.IDPOIDCConfigChangedEventType,
			org.IDPConfigRemovedEventType,
		).
		Builder()
}

func (wm *OrgOIDCIDPWriteModel) NewChangedEvent(
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
) (*org.OIDCIDPChangedEvent, error) {

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
	changeEvent, err := org.NewOIDCIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgJWTIDPWriteModel struct {
	JWTIDPWriteModel
}

func NewJWTOrgIDPWriteModel(orgID, id string) *OrgJWTIDPWriteModel {
	return &OrgJWTIDPWriteModel{
		JWTIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgJWTIDPWriteModel) Reduce() error {
	return wm.JWTIDPWriteModel.Reduce()
}

func (wm *OrgJWTIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.JWTIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTIDPAddedEvent)
		case *org.JWTIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTIDPChangedEvent)
		case *org.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.RemovedEvent)
		case *org.IDPConfigAddedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *org.IDPJWTConfigAddedEvent:
			if wm.ID != e.IDPConfigID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTConfigAddedEvent)
		case *org.IDPJWTConfigChangedEvent:
			if wm.ID != e.IDPConfigID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTConfigChangedEvent)
		case *org.IDPConfigRemovedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.JWTIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgJWTIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.JWTIDPAddedEventType,
			org.JWTIDPChangedEventType,
			org.IDPRemovedEventType,
			org.IDPConfigAddedEventType,
			org.IDPJWTConfigAddedEventType,
			org.IDPJWTConfigChangedEventType,
			org.IDPConfigRemovedEventType,
		).
		Builder()
}

func (wm *OrgJWTIDPWriteModel) NewChangedEvent(
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
) (*org.JWTIDPChangedEvent, error) {

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
	changeEvent, err := org.NewJWTIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgAzureADIDPWriteModel struct {
	AzureADIDPWriteModel
}

func NewAzureADOrgIDPWriteModel(orgID, id string) *OrgAzureADIDPWriteModel {
	return &OrgAzureADIDPWriteModel{
		AzureADIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgAzureADIDPWriteModel) Reduce() error {
	return wm.AzureADIDPWriteModel.Reduce()
}

func (wm *OrgAzureADIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.AzureADIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.AzureADIDPWriteModel.AppendEvents(&e.AzureADIDPAddedEvent)
		case *org.AzureADIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.AzureADIDPWriteModel.AppendEvents(&e.AzureADIDPChangedEvent)
		case *org.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.AzureADIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.AzureADIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgAzureADIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.AzureADIDPAddedEventType,
			org.AzureADIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		Builder()
}

func (wm *OrgAzureADIDPWriteModel) NewChangedEvent(
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
) (*org.AzureADIDPChangedEvent, error) {

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
	changeEvent, err := org.NewAzureADIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgGitHubIDPWriteModel struct {
	GitHubIDPWriteModel
}

func NewGitHubOrgIDPWriteModel(orgID, id string) *OrgGitHubIDPWriteModel {
	return &OrgGitHubIDPWriteModel{
		GitHubIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgGitHubIDPWriteModel) Reduce() error {
	return wm.GitHubIDPWriteModel.Reduce()
}

func (wm *OrgGitHubIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.GitHubIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitHubIDPWriteModel.AppendEvents(&e.GitHubIDPAddedEvent)
		case *org.GitHubIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitHubIDPWriteModel.AppendEvents(&e.GitHubIDPChangedEvent)
		case *org.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitHubIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.GitHubIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgGitHubIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.GitHubIDPAddedEventType,
			org.GitHubIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		Builder()
}

func (wm *OrgGitHubIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*org.GitHubIDPChangedEvent, error) {

	changes, err := wm.GitHubIDPWriteModel.NewChanges(clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := org.NewGitHubIDPChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgGitHubEnterpriseIDPWriteModel struct {
	GitHubEnterpriseIDPWriteModel
}

func NewGitHubEnterpriseOrgIDPWriteModel(orgID, id string) *OrgGitHubEnterpriseIDPWriteModel {
	return &OrgGitHubEnterpriseIDPWriteModel{
		GitHubEnterpriseIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgGitHubEnterpriseIDPWriteModel) Reduce() error {
	return wm.GitHubEnterpriseIDPWriteModel.Reduce()
}

func (wm *OrgGitHubEnterpriseIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.GitHubEnterpriseIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitHubEnterpriseIDPWriteModel.AppendEvents(&e.GitHubEnterpriseIDPAddedEvent)
		case *org.GitHubEnterpriseIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitHubEnterpriseIDPWriteModel.AppendEvents(&e.GitHubEnterpriseIDPChangedEvent)
		case *org.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitHubEnterpriseIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.GitHubEnterpriseIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgGitHubEnterpriseIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.GitHubEnterpriseIDPAddedEventType,
			org.GitHubEnterpriseIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		Builder()
}

func (wm *OrgGitHubEnterpriseIDPWriteModel) NewChangedEvent(
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
) (*org.GitHubEnterpriseIDPChangedEvent, error) {

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
	changeEvent, err := org.NewGitHubEnterpriseIDPChangedEvent(ctx, aggregate, id, wm.Name, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgGitLabIDPWriteModel struct {
	GitLabIDPWriteModel
}

func NewGitLabOrgIDPWriteModel(orgID, id string) *OrgGitLabIDPWriteModel {
	return &OrgGitLabIDPWriteModel{
		GitLabIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgGitLabIDPWriteModel) Reduce() error {
	return wm.GitLabIDPWriteModel.Reduce()
}

func (wm *OrgGitLabIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.GitLabIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitLabIDPWriteModel.AppendEvents(&e.GitLabIDPAddedEvent)
		case *org.GitLabIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitLabIDPWriteModel.AppendEvents(&e.GitLabIDPChangedEvent)
		case *org.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitLabIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.GitLabIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgGitLabIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.GitLabIDPAddedEventType,
			org.GitLabIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		Builder()
}

func (wm *OrgGitLabIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*org.GitLabIDPChangedEvent, error) {

	changes, err := wm.GitLabIDPWriteModel.NewChanges(clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := org.NewGitLabIDPChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgGitLabSelfHostedIDPWriteModel struct {
	GitLabSelfHostedIDPWriteModel
}

func NewGitLabSelfHostedOrgIDPWriteModel(orgID, id string) *OrgGitLabSelfHostedIDPWriteModel {
	return &OrgGitLabSelfHostedIDPWriteModel{
		GitLabSelfHostedIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgGitLabSelfHostedIDPWriteModel) Reduce() error {
	return wm.GitLabSelfHostedIDPWriteModel.Reduce()
}

func (wm *OrgGitLabSelfHostedIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.GitLabSelfHostedIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitLabSelfHostedIDPWriteModel.AppendEvents(&e.GitLabSelfHostedIDPAddedEvent)
		case *org.GitLabSelfHostedIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitLabSelfHostedIDPWriteModel.AppendEvents(&e.GitLabSelfHostedIDPChangedEvent)
		case *org.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GitLabSelfHostedIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.GitLabSelfHostedIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgGitLabSelfHostedIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.GitLabSelfHostedIDPAddedEventType,
			org.GitLabSelfHostedIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		Builder()
}

func (wm *OrgGitLabSelfHostedIDPWriteModel) NewChangedEvent(
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
) (*org.GitLabSelfHostedIDPChangedEvent, error) {

	changes, err := wm.GitLabSelfHostedIDPWriteModel.NewChanges(name, issuer, clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := org.NewGitLabSelfHostedIDPChangedEvent(ctx, aggregate, id, wm.Name, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgGoogleIDPWriteModel struct {
	GoogleIDPWriteModel
}

func NewGoogleOrgIDPWriteModel(orgID, id string) *OrgGoogleIDPWriteModel {
	return &OrgGoogleIDPWriteModel{
		GoogleIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgGoogleIDPWriteModel) Reduce() error {
	return wm.GoogleIDPWriteModel.Reduce()
}

func (wm *OrgGoogleIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.GoogleIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPAddedEvent)
		case *org.GoogleIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPChangedEvent)
		case *org.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GoogleIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.GoogleIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgGoogleIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.GoogleIDPAddedEventType,
			org.GoogleIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		Builder()
}

func (wm *OrgGoogleIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*org.GoogleIDPChangedEvent, error) {

	changes, err := wm.GoogleIDPWriteModel.NewChanges(clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := org.NewGoogleIDPChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgIDPRemoveWriteModel struct {
	IDPRemoveWriteModel
}

func NewOrgIDPRemoveWriteModel(orgID, id string) *OrgIDPRemoveWriteModel {
	return &OrgIDPRemoveWriteModel{
		IDPRemoveWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgIDPRemoveWriteModel) Reduce() error {
	return wm.IDPRemoveWriteModel.Reduce()
}

func (wm *OrgIDPRemoveWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.OAuthIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.OAuthIDPAddedEvent)
		case *org.OIDCIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.OIDCIDPAddedEvent)
		case *org.JWTIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.JWTIDPAddedEvent)
		case *org.AzureADIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.AzureADIDPAddedEvent)
		case *org.GitHubIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GitHubIDPAddedEvent)
		case *org.GitHubEnterpriseIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GitHubEnterpriseIDPAddedEvent)
		case *org.GitLabIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GitLabIDPAddedEvent)
		case *org.GitLabSelfHostedIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GitLabSelfHostedIDPAddedEvent)
		case *org.GoogleIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GoogleIDPAddedEvent)
		case *org.IDPConfigAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *org.IDPConfigRemovedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		case *org.LDAPIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.LDAPIDPAddedEvent)
		case *org.LDAPIDPChangedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.LDAPIDPChangedEvent)
		case *org.IDPRemovedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.IDPRemoveWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgIDPRemoveWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.OAuthIDPAddedEventType,
			org.OIDCIDPAddedEventType,
			org.JWTIDPAddedEventType,
			org.AzureADIDPAddedEventType,
			org.GitHubIDPAddedEventType,
			org.GitHubEnterpriseIDPAddedEventType,
			org.GitLabIDPAddedEventType,
			org.GitLabSelfHostedIDPAddedEventType,
			org.GoogleIDPAddedEventType,
			org.IDPConfigAddedEventType,
			org.LDAPIDPAddedEventType,
			org.LDAPIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		Builder()
}
