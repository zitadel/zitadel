package command

import (
	"reflect"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
)

type OAuthIDPWriteModel struct {
	eventstore.WriteModel

	Name                  string
	ID                    string
	ClientID              string
	ClientSecret          *crypto.CryptoValue
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserEndpoint          string
	Scopes                []string
	idp.Options

	State domain.IDPConfigState //TODO: ?
}

func (wm *OAuthIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.OAuthIDPAddedEvent:
			wm.ID = e.ID
			wm.ClientID = e.ClientID
			wm.ClientSecret = e.ClientSecret
			wm.State = domain.IDPConfigStateActive
		case *idp.OAuthIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceChangedEvent(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OAuthIDPWriteModel) reduceChangedEvent(e *idp.OAuthIDPChangedEvent) {
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.AuthorizationEndpoint != nil {
		wm.AuthorizationEndpoint = *e.AuthorizationEndpoint
	}
	if e.TokenEndpoint != nil {
		wm.TokenEndpoint = *e.TokenEndpoint
	}
	if e.UserEndpoint != nil {
		wm.UserEndpoint = *e.UserEndpoint
	}
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *OAuthIDPWriteModel) NewChanges(
	name,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	authorizationEndpoint,
	tokenEndpoint,
	userEndpoint string,
	scopes []string,
	options idp.Options,
) ([]idp.OAuthIDPChanges, error) {
	changes := make([]idp.OAuthIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeOAuthClientSecret(clientSecret))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeOAuthClientID(clientID))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeOAuthName(name))
	}
	if wm.AuthorizationEndpoint != authorizationEndpoint {
		changes = append(changes, idp.ChangeOAuthAuthorizationEndpoint(authorizationEndpoint))
	}
	if wm.TokenEndpoint != tokenEndpoint {
		changes = append(changes, idp.ChangeOAuthTokenEndpoint(tokenEndpoint))
	}
	if wm.UserEndpoint != userEndpoint {
		changes = append(changes, idp.ChangeOAuthUserEndpoint(userEndpoint))
	}
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeOAuthScopes(scopes))
	}
	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeOAuthOptions(opts))
	}
	return changes, nil
}

type OIDCIDPWriteModel struct {
	eventstore.WriteModel

	Name         string
	Issuer       string
	ID           string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       []string
	idp.Options

	State domain.IDPConfigState //TODO: ?
}

func (wm *OIDCIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.OIDCIDPAddedEvent:
			wm.ID = e.ID
			wm.ClientID = e.ClientID
			wm.ClientSecret = e.ClientSecret
			wm.State = domain.IDPConfigStateActive
		case *idp.OIDCIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceChangedEvent(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OIDCIDPWriteModel) reduceChangedEvent(e *idp.OIDCIDPChangedEvent) {
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.Issuer != nil {
		wm.Issuer = *e.Issuer
	}
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *OIDCIDPWriteModel) NewChanges(
	name,
	issuer,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) ([]idp.OIDCIDPChanges, error) {
	changes := make([]idp.OIDCIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeOIDCClientSecret(clientSecret))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeOIDCClientID(clientID))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeOIDCName(name))
	}
	if wm.Issuer != issuer {
		changes = append(changes, idp.ChangeOIDCIssuer(issuer))
	}
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeOIDCScopes(scopes))
	}
	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeOIDCOptions(opts))
	}
	return changes, nil
}

type GoogleIDPWriteModel struct {
	eventstore.WriteModel

	ID           string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	idp.Options

	State domain.IDPConfigState //TODO: ?
}

func (wm *GoogleIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.GoogleIDPAddedEvent:
			wm.ID = e.ID
			wm.ClientID = e.ClientID
			wm.ClientSecret = e.ClientSecret
			wm.State = domain.IDPConfigStateActive
		case *idp.GoogleIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceChangedEvent(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GoogleIDPWriteModel) reduceChangedEvent(e *idp.GoogleIDPChangedEvent) {
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *GoogleIDPWriteModel) NewChanges(
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	options idp.Options,
) ([]idp.GoogleIDPChanges, error) {
	changes := make([]idp.GoogleIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeGoogleClientSecret(clientSecret))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeGoogleClientID(clientID))
	}

	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeGoogleOptions(opts))
	}
	return changes, nil
}

type GitHubIDPWriteModel struct {
	eventstore.WriteModel

	ID           string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	idp.Options

	State domain.IDPConfigState //TODO: ?
}

func (wm *GitHubIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.GitHubIDPAddedEvent:
			wm.ID = e.ID
			wm.ClientID = e.ClientID
			wm.ClientSecret = e.ClientSecret
			wm.State = domain.IDPConfigStateActive
		case *idp.GitHubIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceChangedEvent(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GitHubIDPWriteModel) reduceChangedEvent(e *idp.GitHubIDPChangedEvent) {
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *GitHubIDPWriteModel) NewChanges(
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	options idp.Options,
) ([]idp.OAuthIDPChanges, error) {
	changes := make([]idp.OAuthIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeOAuthClientSecret(clientSecret))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeOAuthClientID(clientID))
	}

	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeOAuthOptions(opts))
	}
	return changes, nil
}
