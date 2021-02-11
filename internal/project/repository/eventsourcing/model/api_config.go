package model

import (
	"encoding/json"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
)

type APIConfig struct {
	es_models.ObjectRoot
	AppID          string              `json:"appId"`
	ClientID       string              `json:"clientId,omitempty"`
	ClientSecret   *crypto.CryptoValue `json:"clientSecret,omitempty"`
	AuthMethodType int32               `json:"authMethodType,omitempty"`
	ClientKeys     []*ClientKey        `json:"-"`
}

func (c *APIConfig) Changes(changed *APIConfig) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	changes["appId"] = c.AppID
	if c.AuthMethodType != changed.AuthMethodType {
		changes["authMethodType"] = changed.AuthMethodType
	}
	return changes
}

func APIConfigFromModel(config *model.APIConfig) *APIConfig {
	return &APIConfig{
		ObjectRoot:     config.ObjectRoot,
		AppID:          config.AppID,
		ClientID:       config.ClientID,
		ClientSecret:   config.ClientSecret,
		AuthMethodType: int32(config.AuthMethodType),
	}
}

func APIConfigToModel(config *APIConfig) *model.APIConfig {
	oidcConfig := &model.APIConfig{
		ObjectRoot:     config.ObjectRoot,
		AppID:          config.AppID,
		ClientID:       config.ClientID,
		ClientSecret:   config.ClientSecret,
		AuthMethodType: model.APIAuthMethodType(config.AuthMethodType),
		ClientKeys:     ClientKeysToModel(config.ClientKeys),
	}
	return oidcConfig
}

func (p *Project) appendAddAPIConfigEvent(event *es_models.Event) error {
	config := new(APIConfig)
	err := config.setData(event)
	if err != nil {
		return err
	}
	config.ObjectRoot.CreationDate = event.CreationDate
	if i, a := GetApplication(p.Applications, config.AppID); a != nil {
		p.Applications[i].Type = int32(model.AppTypeAPI)
		p.Applications[i].APIConfig = config
	}
	return nil
}

func (p *Project) appendChangeAPIConfigEvent(event *es_models.Event) error {
	config := new(APIConfig)
	err := config.setData(event)
	if err != nil {
		return err
	}

	if i, a := GetApplication(p.Applications, config.AppID); a != nil {
		return p.Applications[i].APIConfig.setData(event)
	}
	return nil
}

func (o *APIConfig) setData(event *es_models.Event) error {
	o.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, o); err != nil {
		logging.Log("EVEN-d8e3s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
