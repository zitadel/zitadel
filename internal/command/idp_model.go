package command

import (
	"reflect"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
)

type LDAPIDPWriteModel struct {
	eventstore.WriteModel

	ID                  string
	Name                string
	Host                string
	Port                string
	TLS                 bool
	BaseDN              string
	UserObjectClass     string
	UserUniqueAttribute string
	Admin               string
	Password            *crypto.CryptoValue
	idp.LDAPAttributes
	idp.Options

	State domain.IDPState
}

func (wm *LDAPIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.LDAPIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceAddeddEvent(e)
		case *idp.LDAPIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceChangedEvent(e)
		case *idp.RemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.State = domain.IDPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *LDAPIDPWriteModel) reduceAddeddEvent(e *idp.LDAPIDPAddedEvent) {
	wm.Name = e.Name
	wm.Host = e.Host
	wm.Port = e.Port
	wm.TLS = e.TLS
	wm.BaseDN = e.BaseDN
	wm.UserObjectClass = e.UserObjectClass
	wm.UserUniqueAttribute = e.UserUniqueAttribute
	wm.Admin = e.Admin
	wm.Password = e.Password
	wm.LDAPAttributes = e.LDAPAttributes
	wm.Options = e.Options
	wm.State = domain.IDPStateActive
}

func (wm *LDAPIDPWriteModel) reduceChangedEvent(e *idp.LDAPIDPChangedEvent) {
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.Host != nil {
		wm.Host = *e.Host
	}
	if e.Port != nil {
		wm.Port = *e.Port
	}
	if e.TLS != nil {
		wm.TLS = *e.TLS
	}
	if e.BaseDN != nil {
		wm.BaseDN = *e.BaseDN
	}
	if e.UserObjectClass != nil {
		wm.UserObjectClass = *e.UserObjectClass
	}
	if e.UserUniqueAttribute != nil {
		wm.UserUniqueAttribute = *e.UserUniqueAttribute
	}
	if e.Admin != nil {
		wm.Admin = *e.Admin
	}
	if e.Password != nil {
		wm.Password = e.Password
	}
	wm.LDAPAttributes.ReduceChanges(e.LDAPAttributeChanges)
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *LDAPIDPWriteModel) NewChanges(
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
) ([]idp.LDAPIDPChanges, error) {
	changes := make([]idp.LDAPIDPChanges, 0)
	var cryptedPassword *crypto.CryptoValue
	var err error
	if password != "" {
		cryptedPassword, err = crypto.Crypt([]byte(password), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeLDAPPassword(cryptedPassword))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeLDAPName(name))
	}
	if wm.Host != host {
		changes = append(changes, idp.ChangeLDAPHost(host))
	}
	if wm.Port != port {
		changes = append(changes, idp.ChangeLDAPPort(port))
	}
	if wm.TLS != tls {
		changes = append(changes, idp.ChangeLDAPTLS(tls))
	}
	if wm.BaseDN != baseDN {
		changes = append(changes, idp.ChangeLDAPBaseDN(baseDN))
	}
	if wm.UserObjectClass != userObjectClass {
		changes = append(changes, idp.ChangeLDAPUserObjectClass(userObjectClass))
	}
	if wm.UserUniqueAttribute != userUniqueAttribute {
		changes = append(changes, idp.ChangeLDAPUserUniqueAttribute(userUniqueAttribute))
	}
	if wm.Admin != admin {
		changes = append(changes, idp.ChangeLDAPAdmin(admin))
	}
	attrs := wm.LDAPAttributes.Changes(attributes)
	if !attrs.IsZero() {
		changes = append(changes, idp.ChangeLDAPAttributes(attrs))
	}
	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeLDAPOptions(opts))
	}
	return changes, nil
}

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

	State domain.IDPState
}

func (wm *OAuthIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.OAuthIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceAddeddEvent(e)
		case *idp.OAuthIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceChangedEvent(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OAuthIDPWriteModel) reduceAddeddEvent(e *idp.OAuthIDPAddedEvent) {
	wm.Name = e.Name
	wm.ClientID = e.ClientID
	wm.ClientSecret = e.ClientSecret
	wm.AuthorizationEndpoint = e.AuthorizationEndpoint
	wm.TokenEndpoint = e.TokenEndpoint
	wm.UserEndpoint = e.UserEndpoint
	wm.Scopes = e.Scopes
	wm.State = domain.IDPStateActive
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
	ID           string
	Issuer       string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       []string
	idp.Options

	State domain.IDPState
}

func (wm *OIDCIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {

		case *idp.OIDCIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Name = e.Name
			wm.Issuer = e.Issuer
			wm.ClientID = e.ClientID
			wm.ClientSecret = e.ClientSecret
			wm.Scopes = e.Scopes
			wm.State = domain.IDPStateActive
		case *idp.OIDCIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceChangedEvent(e)
		case *idpconfig.IDPConfigAddedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
			wm.reduceConfigAddedEvent(e)
		case *idpconfig.IDPConfigChangedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
			wm.reduceConfigChangedEvent(e)
		case *idpconfig.IDPConfigRemovedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
			wm.State = domain.IDPStateRemoved
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

func (wm *OIDCIDPWriteModel) reduceConfigAddedEvent(e *idpconfig.IDPConfigAddedEvent) {
	wm.Name = e.Name
	//rm.StylingType = e.StylingType //TODO: drop?
	wm.Options.IsAutoCreation = e.AutoRegister
	wm.State = domain.IDPStateActive
}

func (wm *OIDCIDPWriteModel) reduceConfigChangedEvent(e *idpconfig.IDPConfigChangedEvent) {
	if e.Name != nil {
		wm.Name = *e.Name
	}
	//if e.StylingType != nil && e.StylingType.Valid() { //TODO: drop?
	//	rm.StylingType = *e.StylingType
	//}
	if e.AutoRegister != nil {
		wm.Options.IsAutoCreation = *e.AutoRegister
	}
}

type JWTIDPWriteModel struct {
	eventstore.WriteModel

	ID           string
	Name         string
	Issuer       string
	JWTEndpoint  string
	KeysEndpoint string
	HeaderName   string
	idp.Options

	State domain.IDPState //TODO: ?
}

func (wm *JWTIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.JWTIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Name = e.Name
			wm.Issuer = e.Issuer
			wm.JWTEndpoint = e.JWTEndpoint
			wm.KeysEndpoint = e.KeysEndpoint
			wm.HeaderName = e.HeaderName
			wm.State = domain.IDPStateActive
		case *idp.JWTIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceChangedEvent(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *JWTIDPWriteModel) reduceChangedEvent(e *idp.JWTIDPChangedEvent) {
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.Issuer != nil {
		wm.Issuer = *e.Issuer
	}
	if e.JWTEndpoint != nil {
		wm.JWTEndpoint = *e.JWTEndpoint
	}
	if e.KeysEndpoint != nil {
		wm.KeysEndpoint = *e.KeysEndpoint
	}
	if e.HeaderName != nil {
		wm.HeaderName = *e.HeaderName
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *JWTIDPWriteModel) NewChanges(
	name,
	issuer,
	jwtEndpoint,
	keysEndpoint,
	headerName string,
	options idp.Options,
) ([]idp.JWTIDPChanges, error) {
	changes := make([]idp.JWTIDPChanges, 0)
	if wm.Name != name {
		changes = append(changes, idp.ChangeJWTName(name))
	}
	if wm.Issuer != issuer {
		changes = append(changes, idp.ChangeJWTIssuer(issuer))
	}
	if wm.JWTEndpoint != jwtEndpoint {
		changes = append(changes, idp.ChangeJWTEndpoint(jwtEndpoint))
	}
	if wm.KeysEndpoint != keysEndpoint {
		changes = append(changes, idp.ChangeJWTKeysEndpoint(keysEndpoint))
	}
	if wm.HeaderName != headerName {
		changes = append(changes, idp.ChangeJWTHeaderName(headerName))
	}
	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeJWTOptions(opts))
	}
	return changes, nil
}

type AzureADIDPWriteModel struct {
	eventstore.WriteModel

	ID              string
	Name            string
	ClientID        string
	ClientSecret    *crypto.CryptoValue
	Scopes          []string
	Tenant          string
	IsEmailVerified bool
	idp.Options

	State domain.IDPState //TODO: ?
}

func (wm *AzureADIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.AzureADIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Name = e.Name
			wm.ClientID = e.ClientID
			wm.ClientSecret = e.ClientSecret
			wm.Scopes = e.Scopes
			wm.Tenant = e.Tenant
			wm.IsEmailVerified = e.IsEmailVerified
			wm.State = domain.IDPStateActive
		case *idp.AzureADIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceChangedEvent(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *AzureADIDPWriteModel) reduceChangedEvent(e *idp.AzureADIDPChangedEvent) {
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
	if e.Tenant != nil {
		wm.Tenant = *e.Tenant
	}
	if e.IsEmailVerified != nil {
		wm.IsEmailVerified = *e.IsEmailVerified
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *AzureADIDPWriteModel) NewChanges(
	name string,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	tenant string,
	isEmailVerified bool,
	options idp.Options,
) ([]idp.AzureADIDPChanges, error) {
	changes := make([]idp.AzureADIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeAzureADClientSecret(clientSecret))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeAzureADName(name))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeAzureADClientID(clientID))
	}
	if wm.Tenant != tenant {
		changes = append(changes, idp.ChangeAzureADTenant(tenant))
	}
	if wm.IsEmailVerified != isEmailVerified {
		changes = append(changes, idp.ChangeAzureADIsEmailVerified(isEmailVerified))
	}
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeAzureADScopes(scopes))
	}

	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeAzureADOptions(opts))
	}
	return changes, nil
}

type GitHubIDPWriteModel struct {
	eventstore.WriteModel

	ID           string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       []string
	idp.Options

	State domain.IDPState //TODO: ?
}

func (wm *GitHubIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.GitHubIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.ClientID = e.ClientID
			wm.ClientSecret = e.ClientSecret
			wm.Scopes = e.Scopes
			wm.State = domain.IDPStateActive
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
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *GitHubIDPWriteModel) NewChanges(
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
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
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeOAuthScopes(scopes))
	}

	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeOAuthOptions(opts))
	}
	return changes, nil
}

type GitHubEnterpriseIDPWriteModel struct {
	eventstore.WriteModel

	ID                    string
	Name                  string
	ClientID              string
	ClientSecret          *crypto.CryptoValue
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserEndpoint          string
	Scopes                []string
	idp.Options

	State domain.IDPState //TODO: ?
}

func (wm *GitHubEnterpriseIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.GitHubEnterpriseIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Name = e.Name
			wm.ClientID = e.ClientID
			wm.ClientSecret = e.ClientSecret
			wm.AuthorizationEndpoint = e.AuthorizationEndpoint
			wm.TokenEndpoint = e.TokenEndpoint
			wm.UserEndpoint = e.UserEndpoint
			wm.Scopes = e.Scopes
			wm.State = domain.IDPStateActive
		case *idp.GitHubEnterpriseIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceChangedEvent(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GitHubEnterpriseIDPWriteModel) reduceChangedEvent(e *idp.GitHubEnterpriseIDPChangedEvent) {
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

func (wm *GitHubEnterpriseIDPWriteModel) NewChanges(
	name,
	clientID string,
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

type GitLabIDPWriteModel struct {
	eventstore.WriteModel

	ID           string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       []string
	idp.Options

	State domain.IDPState //TODO: ?
}

func (wm *GitLabIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.GitLabIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.ClientID = e.ClientID
			wm.ClientSecret = e.ClientSecret
			wm.Scopes = e.Scopes
			wm.State = domain.IDPStateActive
		case *idp.GitLabIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceChangedEvent(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GitLabIDPWriteModel) reduceChangedEvent(e *idp.GitLabIDPChangedEvent) {
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *GitLabIDPWriteModel) NewChanges(
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) ([]idp.GitLabIDPChanges, error) {
	changes := make([]idp.GitLabIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeGitLabClientSecret(clientSecret))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeGitLabClientID(clientID))
	}
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeGitLabScopes(scopes))
	}

	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeGitLabOptions(opts))
	}
	return changes, nil
}

type GitLabSelfHostedIDPWriteModel struct {
	eventstore.WriteModel

	ID           string
	Name         string
	Issuer       string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       []string
	idp.Options

	State domain.IDPState //TODO: ?
}

func (wm *GitLabSelfHostedIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.GitLabSelfHostedIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Name = e.Name
			wm.Issuer = e.Issuer
			wm.ClientID = e.ClientID
			wm.ClientSecret = e.ClientSecret
			wm.Scopes = e.Scopes
			wm.State = domain.IDPStateActive
		case *idp.GitLabSelfHostedIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceChangedEvent(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GitLabSelfHostedIDPWriteModel) reduceChangedEvent(e *idp.GitLabSelfHostedIDPChangedEvent) {
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

func (wm *GitLabSelfHostedIDPWriteModel) NewChanges(
	name string,
	issuer string,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) ([]idp.GitLabSelfHostedIDPChanges, error) {
	changes := make([]idp.GitLabSelfHostedIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeGitLabSelfHostedClientSecret(clientSecret))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeGitLabSelfHostedClientID(clientID))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeGitLabSelfHostedName(name))
	}
	if wm.Issuer != issuer {
		changes = append(changes, idp.ChangeGitLabSelfHostedIssuer(issuer))
	}
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeGitLabSelfHostedScopes(scopes))
	}
	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeGitLabSelfHostedOptions(opts))
	}
	return changes, nil
}

type GoogleIDPWriteModel struct {
	eventstore.WriteModel

	ID           string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       []string
	idp.Options

	State domain.IDPState
}

func (wm *GoogleIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.GoogleIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.ClientID = e.ClientID
			wm.ClientSecret = e.ClientSecret
			wm.Scopes = e.Scopes
			wm.State = domain.IDPStateActive
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
	scopes []string,
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
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeGoogleScopes(scopes))
	}

	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeGoogleOptions(opts))
	}
	return changes, nil
}

type IDPRemoveWriteModel struct {
	eventstore.WriteModel

	ID    string
	State domain.IDPState
	name  string
}

func (wm *IDPRemoveWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.OAuthIDPAddedEvent:
			wm.reduceAdded(e.ID, e.Name)
		case *idp.OIDCIDPAddedEvent:
			wm.reduceAdded(e.ID, e.Name)
		case *idp.JWTIDPAddedEvent:
			wm.reduceAdded(e.ID, e.Name)
		case *idp.AzureADIDPAddedEvent:
			wm.reduceAdded(e.ID, e.Name)
		case *idp.GitHubIDPAddedEvent:
			wm.reduceAdded(e.ID, e.Name)
		case *idp.GitHubEnterpriseIDPAddedEvent:
			wm.reduceAdded(e.ID, e.Name)
		case *idp.GitLabIDPAddedEvent:
			wm.reduceAdded(e.ID, "")
		case *idp.GitLabSelfHostedIDPAddedEvent:
			wm.reduceAdded(e.ID, e.Name)
		case *idp.GoogleIDPAddedEvent:
			wm.reduceAdded(e.ID, "")
		case *idp.LDAPIDPAddedEvent:
			wm.reduceAdded(e.ID, e.Name)
		case *idp.LDAPIDPChangedEvent:
			wm.reduceChanged(e.ID, e.Name)
		case *idp.RemovedEvent:
			wm.reduceRemoved(e.ID)
		case *idpconfig.IDPConfigAddedEvent:
			wm.reduceAdded(e.ConfigID, "")
		case *idpconfig.IDPConfigRemovedEvent:
			wm.reduceRemoved(e.ConfigID)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *IDPRemoveWriteModel) reduceAdded(id string, name string) {
	if wm.ID != id {
		return
	}
	wm.State = domain.IDPStateActive
	wm.name = name
}

func (wm *IDPRemoveWriteModel) reduceChanged(id string, name *string) {
	if wm.ID != id || name == nil {
		return
	}
	wm.name = *name
}

func (wm *IDPRemoveWriteModel) reduceRemoved(id string) {
	if wm.ID != id {
		return
	}
	wm.State = domain.IDPStateRemoved
}
