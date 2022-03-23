package model

import (
	"encoding/json"
	"reflect"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
	"github.com/lib/pq"
)

type OIDCIDPConfig struct {
	es_models.ObjectRoot
	IDPConfigID           string              `json:"idpConfigId"`
	ClientID              string              `json:"clientId"`
	ClientSecret          *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Issuer                string              `json:"issuer,omitempty"`
	Scopes                pq.StringArray      `json:"scopes,omitempty"`
	IDPDisplayNameMapping int32               `json:"idpDisplayNameMapping,omitempty"`
	UsernameMapping       int32               `json:"usernameMapping,omitempty"`
}

func (c *OIDCIDPConfig) Changes(changed *OIDCIDPConfig) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	changes["idpConfigId"] = c.IDPConfigID
	if c.ClientID != changed.ClientID {
		changes["clientId"] = changed.ClientID
	}
	if changed.ClientSecret != nil && c.ClientSecret != changed.ClientSecret {
		changes["clientSecret"] = changed.ClientSecret
	}
	if c.Issuer != changed.Issuer {
		changes["issuer"] = changed.Issuer
	}
	if !reflect.DeepEqual(c.Scopes, changed.Scopes) {
		changes["scopes"] = changed.Scopes
	}
	if c.IDPDisplayNameMapping != changed.IDPDisplayNameMapping {
		changes["idpDisplayNameMapping"] = changed.IDPDisplayNameMapping
	}
	if c.UsernameMapping != changed.UsernameMapping {
		changes["usernameMapping"] = changed.UsernameMapping
	}
	return changes
}

func OIDCIDPConfigToModel(config *OIDCIDPConfig) *model.OIDCIDPConfig {
	return &model.OIDCIDPConfig{
		ObjectRoot:            config.ObjectRoot,
		IDPConfigID:           config.IDPConfigID,
		ClientID:              config.ClientID,
		ClientSecret:          config.ClientSecret,
		Issuer:                config.Issuer,
		Scopes:                config.Scopes,
		IDPDisplayNameMapping: model.OIDCMappingField(config.IDPDisplayNameMapping),
		UsernameMapping:       model.OIDCMappingField(config.UsernameMapping),
	}
}

func (o *OIDCIDPConfig) SetData(event *es_models.Event) error {
	o.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, o); err != nil {
		logging.Log("EVEN-Msh8s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
