package model

import (
	"encoding/json"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
	"github.com/lib/pq"
	"reflect"
)

type OIDCIDPConfig struct {
	es_models.ObjectRoot
	IDPConfigID  string              `json:"idpConfigId"`
	ClientID     string              `json:"clientId"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Issuer       string              `json:"issuer,omitempty"`
	Scopes       pq.StringArray      `json:"scopes,omitempty"`
}

func (c *OIDCIDPConfig) Changes(changed *OIDCIDPConfig) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	changes["idpConfigId"] = c.IDPConfigID
	if c.ClientID != changed.ClientID {
		changes["clientId"] = changed.ClientID
	}
	if c.ClientSecret != nil {
		changes["clientSecret"] = changed.ClientSecret
	}
	if c.Issuer != changed.Issuer {
		changes["issuer"] = changed.Issuer
	}
	if !reflect.DeepEqual(c.Scopes, changed.Scopes) {
		changes["scopes"] = changed.Scopes
	}
	return changes
}

func OIDCIDPConfigFromModel(config *model.OIDCIDPConfig) *OIDCIDPConfig {
	return &OIDCIDPConfig{
		ObjectRoot:   config.ObjectRoot,
		IDPConfigID:  config.IDPConfigID,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Issuer:       config.Issuer,
		Scopes:       config.Scopes,
	}
}

func OIDCIDPConfigToModel(config *OIDCIDPConfig) *model.OIDCIDPConfig {
	return &model.OIDCIDPConfig{
		ObjectRoot:   config.ObjectRoot,
		IDPConfigID:  config.IDPConfigID,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Issuer:       config.Issuer,
		Scopes:       config.Scopes,
	}
}

func (iam *Iam) appendAddOIDCIDPConfigEvent(event *es_models.Event) error {
	config := new(OIDCIDPConfig)
	err := config.SetData(event)
	if err != nil {
		return err
	}
	config.ObjectRoot.CreationDate = event.CreationDate
	if i, idpConfig := GetIDPConfig(iam.IDPs, config.IDPConfigID); idpConfig != nil {
		iam.IDPs[i].Type = int32(model.IDPConfigTypeOIDC)
		iam.IDPs[i].OIDCIDPConfig = config
	}
	return nil
}

func (iam *Iam) appendChangeOIDCIDPConfigEvent(event *es_models.Event) error {
	config := new(OIDCIDPConfig)
	err := config.SetData(event)
	if err != nil {
		return err
	}

	if i, idpConfig := GetIDPConfig(iam.IDPs, config.IDPConfigID); idpConfig != nil {
		iam.IDPs[i].OIDCIDPConfig.SetData(event)
	}
	return nil
}

func (o *OIDCIDPConfig) SetData(event *es_models.Event) error {
	o.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, o); err != nil {
		logging.Log("EVEN-Msh8s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
