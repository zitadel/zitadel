package model

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
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
		ClockSkew:                config.ClockSkew,
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
		p.Applications[i].OIDCConfig.setData(event)
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
