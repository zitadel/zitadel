package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/org"
)

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
			wm.OAuthIDPWriteModel.AppendEvents(&e.OAuthIDPAddedEvent)
		case *org.OAuthIDPChangedEvent:
			wm.OAuthIDPWriteModel.AppendEvents(&e.OAuthIDPChangedEvent)
		case *org.IDPRemovedEvent:
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
		EventData(map[string]interface{}{"id": wm.ID}).
		Builder()
}

func (wm *OrgOAuthIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
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
	return org.NewOAuthIDPChangedEvent(ctx, aggregate, id, changes)
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
			wm.GitHubIDPWriteModel.AppendEvents(&e.GitHubIDPAddedEvent)
		case *org.GitHubIDPChangedEvent:
			wm.GitHubIDPWriteModel.AppendEvents(&e.GitHubIDPChangedEvent)
		case *org.IDPRemovedEvent:
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
		EventData(map[string]interface{}{"id": wm.ID}).
		Builder()
}

func (wm *OrgGitHubIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*org.GitHubIDPChangedEvent, error) {

	changes, err := wm.GitHubIDPWriteModel.NewChanges(name, clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	return org.NewGitHubIDPChangedEvent(ctx, aggregate, id, changes)
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
			wm.GitHubEnterpriseIDPWriteModel.AppendEvents(&e.GitHubEnterpriseIDPAddedEvent)
		case *org.GitHubEnterpriseIDPChangedEvent:
			wm.GitHubEnterpriseIDPWriteModel.AppendEvents(&e.GitHubEnterpriseIDPChangedEvent)
		case *org.IDPRemovedEvent:
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
		EventData(map[string]interface{}{"id": wm.ID}).
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
	return org.NewGitHubEnterpriseIDPChangedEvent(ctx, aggregate, id, changes)
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
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPAddedEvent)
		case *org.GoogleIDPChangedEvent:
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPChangedEvent)
		case *org.IDPRemovedEvent:
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
		EventData(map[string]interface{}{"id": wm.ID}).
		Builder()
}

func (wm *OrgGoogleIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*org.GoogleIDPChangedEvent, error) {

	changes, err := wm.GoogleIDPWriteModel.NewChanges(name, clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	return org.NewGoogleIDPChangedEvent(ctx, aggregate, id, changes)
}

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
			wm.LDAPIDPWriteModel.AppendEvents(&e.LDAPIDPAddedEvent)
		case *org.LDAPIDPChangedEvent:
			wm.LDAPIDPWriteModel.AppendEvents(&e.LDAPIDPChangedEvent)
		case *org.IDPRemovedEvent:
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
		EventData(map[string]interface{}{"id": wm.ID}).
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
	return org.NewLDAPIDPChangedEvent(ctx, aggregate, id, oldName, changes)
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
		case *org.GitHubIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GitHubIDPAddedEvent)
		case *org.GitHubIDPChangedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GitHubIDPChangedEvent)
		case *org.GitHubEnterpriseIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GitHubEnterpriseIDPAddedEvent)
		case *org.GitHubEnterpriseIDPChangedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GitHubEnterpriseIDPChangedEvent)
		case *org.OAuthIDPChangedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.OAuthIDPChangedEvent)
		case *org.GoogleIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GoogleIDPAddedEvent)
		case *org.GoogleIDPChangedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GoogleIDPChangedEvent)
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
			org.OAuthIDPChangedEventType,
			org.GitHubIDPAddedEventType,
			org.GitHubIDPChangedEventType,
			org.GitHubEnterpriseIDPAddedEventType,
			org.GitHubEnterpriseIDPChangedEventType,
			org.GoogleIDPAddedEventType,
			org.GoogleIDPChangedEventType,
			org.LDAPIDPAddedEventType,
			org.LDAPIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Builder()
}
