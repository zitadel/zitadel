package model

import (
	"encoding/json"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
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

func (o *APIConfig) setData(event *es_models.Event) error {
	o.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, o); err != nil {
		logging.Log("EVEN-d8e3s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
