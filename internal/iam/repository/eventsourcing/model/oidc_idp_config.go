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

type OidcIdpConfig struct {
	es_models.ObjectRoot
	IdpConfigID  string              `json:"idpConfigId"`
	ClientID     string              `json:"clientId"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Issuer       string              `json:"issuer,omitempty"`
	Scopes       pq.StringArray      `json:"scopes,omitempty"`
}

func (c *OidcIdpConfig) Changes(changed *OidcIdpConfig) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	changes["idpConfigId"] = c.IdpConfigID
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

func OidcIdpConfigFromModel(config *model.OidcIdpConfig) *OidcIdpConfig {
	return &OidcIdpConfig{
		ObjectRoot:   config.ObjectRoot,
		IdpConfigID:  config.IDPConfigID,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Issuer:       config.Issuer,
		Scopes:       config.Scopes,
	}
}

func OidcIdpConfigToModel(config *OidcIdpConfig) *model.OidcIdpConfig {
	return &model.OidcIdpConfig{
		ObjectRoot:   config.ObjectRoot,
		IDPConfigID:  config.IdpConfigID,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Issuer:       config.Issuer,
		Scopes:       config.Scopes,
	}
}

func (iam *Iam) appendAddOidcIdpConfigEvent(event *es_models.Event) error {
	config := new(OidcIdpConfig)
	err := config.SetData(event)
	if err != nil {
		return err
	}
	config.ObjectRoot.CreationDate = event.CreationDate
	if i, a := GetIdpConfig(iam.IDPs, config.IdpConfigID); a != nil {
		iam.IDPs[i].Type = int32(model.IDPConfigTypeOIDC)
		iam.IDPs[i].OIDCIDPConfig = config
	}
	return nil
}

func (iam *Iam) appendChangeOidcIdpConfigEvent(event *es_models.Event) error {
	config := new(OidcIdpConfig)
	err := config.SetData(event)
	if err != nil {
		return err
	}

	if i, a := GetIdpConfig(iam.IDPs, config.IdpConfigID); a != nil {
		iam.IDPs[i].OIDCIDPConfig.SetData(event)
	}
	return nil
}

func (o *OidcIdpConfig) SetData(event *es_models.Event) error {
	o.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, o); err != nil {
		logging.Log("EVEN-Msh8s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
