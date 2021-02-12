package model

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/project/model"
)

type OIDCConfig struct {
	es_models.ObjectRoot
	Version                  int32               `json:"oidcVersion,omitempty"`
	AppID                    string              `json:"appId"`
	ClientID                 string              `json:"clientId,omitempty"`
	ClientSecret             *crypto.CryptoValue `json:"clientSecret,omitempty"`
	RedirectUris             []string            `json:"redirectUris,omitempty"`
	ResponseTypes            []int32             `json:"responseTypes,omitempty"`
	GrantTypes               []int32             `json:"grantTypes,omitempty"`
	ApplicationType          int32               `json:"applicationType,omitempty"`
	AuthMethodType           int32               `json:"authMethodType,omitempty"`
	PostLogoutRedirectUris   []string            `json:"postLogoutRedirectUris,omitempty"`
	DevMode                  bool                `json:"devMode,omitempty"`
	AccessTokenType          int32               `json:"accessTokenType,omitempty"`
	AccessTokenRoleAssertion bool                `json:"accessTokenRoleAssertion,omitempty"`
	IDTokenRoleAssertion     bool                `json:"idTokenRoleAssertion,omitempty"`
	IDTokenUserinfoAssertion bool                `json:"idTokenUserinfoAssertion,omitempty"`
	ClockSkew                time.Duration       `json:"clockSkew,omitempty"`
	ClientKeys               []*ClientKey        `json:"-"`
}

func (c *OIDCConfig) Changes(changed *OIDCConfig) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	changes["appId"] = c.AppID
	if !reflect.DeepEqual(c.RedirectUris, changed.RedirectUris) {
		changes["redirectUris"] = changed.RedirectUris
	}
	if !reflect.DeepEqual(c.ResponseTypes, changed.ResponseTypes) {
		changes["responseTypes"] = changed.ResponseTypes
	}
	if !reflect.DeepEqual(c.GrantTypes, changed.GrantTypes) {
		changes["grantTypes"] = changed.GrantTypes
	}
	if c.ApplicationType != changed.ApplicationType {
		changes["applicationType"] = changed.ApplicationType
	}
	if c.AuthMethodType != changed.AuthMethodType {
		changes["authMethodType"] = changed.AuthMethodType
	}
	if c.Version != changed.Version {
		changes["oidcVersion"] = changed.Version
	}
	if !reflect.DeepEqual(c.PostLogoutRedirectUris, changed.PostLogoutRedirectUris) {
		changes["postLogoutRedirectUris"] = changed.PostLogoutRedirectUris
	}
	if c.DevMode != changed.DevMode {
		changes["devMode"] = changed.DevMode
	}
	if c.AccessTokenType != changed.AccessTokenType {
		changes["accessTokenType"] = changed.AccessTokenType
	}
	if c.AccessTokenRoleAssertion != changed.AccessTokenRoleAssertion {
		changes["accessTokenRoleAssertion"] = changed.AccessTokenRoleAssertion
	}
	if c.IDTokenRoleAssertion != changed.IDTokenRoleAssertion {
		changes["idTokenRoleAssertion"] = changed.IDTokenRoleAssertion
	}
	if c.IDTokenUserinfoAssertion != changed.IDTokenUserinfoAssertion {
		changes["idTokenUserinfoAssertion"] = changed.IDTokenUserinfoAssertion
	}
	if c.ClockSkew != changed.ClockSkew {
		changes["clockSkew"] = changed.ClockSkew
	}
	return changes
}

func OIDCConfigFromModel(config *model.OIDCConfig) *OIDCConfig {
	responseTypes := make([]int32, len(config.ResponseTypes))
	for i, rt := range config.ResponseTypes {
		responseTypes[i] = int32(rt)
	}
	grantTypes := make([]int32, len(config.GrantTypes))
	for i, rt := range config.GrantTypes {
		grantTypes[i] = int32(rt)
	}
	return &OIDCConfig{
		ObjectRoot:               config.ObjectRoot,
		AppID:                    config.AppID,
		Version:                  int32(config.OIDCVersion),
		ClientID:                 config.ClientID,
		ClientSecret:             config.ClientSecret,
		RedirectUris:             config.RedirectUris,
		ResponseTypes:            responseTypes,
		GrantTypes:               grantTypes,
		ApplicationType:          int32(config.ApplicationType),
		AuthMethodType:           int32(config.AuthMethodType),
		PostLogoutRedirectUris:   config.PostLogoutRedirectUris,
		DevMode:                  config.DevMode,
		AccessTokenType:          int32(config.AccessTokenType),
		AccessTokenRoleAssertion: config.AccessTokenRoleAssertion,
		IDTokenRoleAssertion:     config.IDTokenRoleAssertion,
		IDTokenUserinfoAssertion: config.IDTokenUserinfoAssertion,
		ClockSkew:                config.ClockSkew,
	}
}

func OIDCConfigToModel(config *OIDCConfig) *model.OIDCConfig {
	responseTypes := make([]model.OIDCResponseType, len(config.ResponseTypes))
	for i, rt := range config.ResponseTypes {
		responseTypes[i] = model.OIDCResponseType(rt)
	}
	grantTypes := make([]model.OIDCGrantType, len(config.GrantTypes))
	for i, rt := range config.GrantTypes {
		grantTypes[i] = model.OIDCGrantType(rt)
	}
	oidcConfig := &model.OIDCConfig{
		ObjectRoot:               config.ObjectRoot,
		AppID:                    config.AppID,
		OIDCVersion:              model.OIDCVersion(config.Version),
		ClientID:                 config.ClientID,
		ClientSecret:             config.ClientSecret,
		RedirectUris:             config.RedirectUris,
		ResponseTypes:            responseTypes,
		GrantTypes:               grantTypes,
		ApplicationType:          model.OIDCApplicationType(config.ApplicationType),
		AuthMethodType:           model.OIDCAuthMethodType(config.AuthMethodType),
		PostLogoutRedirectUris:   config.PostLogoutRedirectUris,
		DevMode:                  config.DevMode,
		AccessTokenType:          model.OIDCTokenType(config.AccessTokenType),
		AccessTokenRoleAssertion: config.AccessTokenRoleAssertion,
		IDTokenRoleAssertion:     config.IDTokenRoleAssertion,
		IDTokenUserinfoAssertion: config.IDTokenUserinfoAssertion,
		ClockSkew:                config.ClockSkew,
		ClientKeys:               ClientKeysToModel(config.ClientKeys),
	}
	oidcConfig.FillCompliance()
	return oidcConfig
}

func (p *Project) appendAddOIDCConfigEvent(event *es_models.Event) error {
	config := new(OIDCConfig)
	err := config.setData(event)
	if err != nil {
		return err
	}
	config.ObjectRoot.CreationDate = event.CreationDate
	if i, a := GetApplication(p.Applications, config.AppID); a != nil {
		p.Applications[i].Type = int32(model.AppTypeOIDC)
		p.Applications[i].OIDCConfig = config
	}
	return nil
}

func (p *Project) appendChangeOIDCConfigEvent(event *es_models.Event) error {
	config := new(OIDCConfig)
	err := config.setData(event)
	if err != nil {
		return err
	}

	if i, a := GetApplication(p.Applications, config.AppID); a != nil {
		return p.Applications[i].OIDCConfig.setData(event)
	}
	return nil
}

func (p *Project) appendAddClientKeyEvent(event *es_models.Event) error {
	key := new(ClientKey)
	err := key.SetData(event)
	if err != nil {
		return err
	}

	if i, a := GetApplication(p.Applications, key.ApplicationID); a != nil {
		if a.OIDCConfig != nil {
			p.Applications[i].OIDCConfig.ClientKeys = append(p.Applications[i].OIDCConfig.ClientKeys, key)
		}
		if a.APIConfig != nil {
			p.Applications[i].APIConfig.ClientKeys = append(p.Applications[i].APIConfig.ClientKeys, key)
		}
	}
	return nil
}

func (p *Project) appendRemoveClientKeyEvent(event *es_models.Event) error {
	key := new(ClientKey)
	err := key.SetData(event)
	if err != nil {
		return err
	}
	if i, a := GetApplication(p.Applications, key.ApplicationID); a != nil {
		if a.OIDCConfig != nil {
			if j, k := GetClientKey(p.Applications[i].OIDCConfig.ClientKeys, key.KeyID); k != nil {
				p.Applications[i].OIDCConfig.ClientKeys[j] = p.Applications[i].OIDCConfig.ClientKeys[len(p.Applications[i].OIDCConfig.ClientKeys)-1]
				p.Applications[i].OIDCConfig.ClientKeys[len(p.Applications[i].OIDCConfig.ClientKeys)-1] = nil
				p.Applications[i].OIDCConfig.ClientKeys = p.Applications[i].OIDCConfig.ClientKeys[:len(p.Applications[i].OIDCConfig.ClientKeys)-1]
			}
		}
		if a.APIConfig != nil {
			if j, k := GetClientKey(p.Applications[i].APIConfig.ClientKeys, key.KeyID); k != nil {
				p.Applications[i].APIConfig.ClientKeys[j] = p.Applications[i].APIConfig.ClientKeys[len(p.Applications[i].APIConfig.ClientKeys)-1]
				p.Applications[i].APIConfig.ClientKeys[len(p.Applications[i].APIConfig.ClientKeys)-1] = nil
				p.Applications[i].APIConfig.ClientKeys = p.Applications[i].APIConfig.ClientKeys[:len(p.Applications[i].APIConfig.ClientKeys)-1]
			}
		}
	}
	return nil
}

func (o *OIDCConfig) setData(event *es_models.Event) error {
	o.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, o); err != nil {
		logging.Log("EVEN-d8e3s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func GetClientKey(keys []*ClientKey, id string) (int, *ClientKey) {
	for i, k := range keys {
		if k.KeyID == id {
			return i, k
		}
	}
	return -1, nil
}

type ClientKey struct {
	es_models.ObjectRoot `json:"-"`
	ApplicationID        string    `json:"applicationID,omitempty"`
	ClientID             string    `json:"clientId,omitempty"`
	KeyID                string    `json:"keyId,omitempty"`
	Type                 int32     `json:"type,omitempty"`
	ExpirationDate       time.Time `json:"expirationDate,omitempty"`
	PublicKey            []byte    `json:"publicKey,omitempty"`
	privateKey           []byte
}

func (key *ClientKey) SetData(event *es_models.Event) error {
	key.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, key); err != nil {
		logging.Log("EVEN-SADdg").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (key *ClientKey) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		err := key.AppendEvent(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (key *ClientKey) AppendEvent(event *es_models.Event) (err error) {
	key.ObjectRoot.AppendEvent(event)
	switch event.Type {
	case ClientKeyAdded:
		err = json.Unmarshal(event.Data, key)
		if err != nil {
			return errors.ThrowInternal(err, "MODEL-Fetg3", "Errors.Internal")
		}
	case ClientKeyRemoved:
		key.ExpirationDate = event.CreationDate
	}
	return err
}

func ClientKeyFromModel(key *model.ClientKey) *ClientKey {
	return &ClientKey{
		ObjectRoot:     key.ObjectRoot,
		ExpirationDate: key.ExpirationDate,
		ApplicationID:  key.ApplicationID,
		ClientID:       key.ClientID,
		KeyID:          key.KeyID,
		Type:           int32(key.Type),
	}
}

func ClientKeysToModel(keys []*ClientKey) []*model.ClientKey {
	clientKeys := make([]*model.ClientKey, len(keys))
	for i, key := range keys {
		clientKeys[i] = ClientKeyToModel(key)
	}
	return clientKeys
}

func ClientKeyToModel(key *ClientKey) *model.ClientKey {
	return &model.ClientKey{
		ObjectRoot:     key.ObjectRoot,
		ExpirationDate: key.ExpirationDate,
		ApplicationID:  key.ApplicationID,
		ClientID:       key.ClientID,
		KeyID:          key.KeyID,
		PrivateKey:     key.privateKey,
		Type:           key_model.AuthNKeyType(key.Type),
	}
}

func (key *ClientKey) GenerateClientKeyPair(keySize int) error {
	privateKey, publicKey, err := crypto.GenerateKeyPair(keySize)
	if err != nil {
		return err
	}
	key.PublicKey, err = crypto.PublicKeyToBytes(publicKey)
	if err != nil {
		return err
	}
	key.privateKey = crypto.PrivateKeyToBytes(privateKey)
	return nil
}
