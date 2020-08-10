package model

import (
	"encoding/json"
	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
)

type IdpConfig struct {
	es_models.ObjectRoot
	IDPConfigID   string         `json:"idpConfigId"`
	State         int32          `json:"-"`
	Name          string         `json:"name,omitempty"`
	Type          int32          `json:"idpType,omitempty"`
	LogoSrc       string         `json:"logoSrc,omitempty"`
	OIDCIDPConfig *OidcIdpConfig `json:"-"`
}

type IdpConfigID struct {
	es_models.ObjectRoot
	IdpConfigID string `json:"idpConfigId"`
}

func GetIdpConfig(idps []*IdpConfig, id string) (int, *IdpConfig) {
	for i, idp := range idps {
		if idp.IDPConfigID == id {
			return i, idp
		}
	}
	return -1, nil
}

func (c *IdpConfig) Changes(changed *IdpConfig) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	changes["idpConfigId"] = c.IDPConfigID
	if changed.Name != "" && c.Name != changed.Name {
		changes["name"] = changed.Name
	}
	if changed.LogoSrc != "" && c.LogoSrc != changed.LogoSrc {
		changes["logoSrc"] = changed.LogoSrc
	}
	return changes
}

func IdpConfigsToModel(idps []*IdpConfig) []*model.IdpConfig {
	convertedIDPConfigs := make([]*model.IdpConfig, len(idps))
	for i, idp := range idps {
		convertedIDPConfigs[i] = IdpConfigToModel(idp)
	}
	return convertedIDPConfigs
}

func IdpConfigsFromModel(idps []*model.IdpConfig) []*IdpConfig {
	convertedIDPConfigs := make([]*IdpConfig, len(idps))
	for i, idp := range idps {
		convertedIDPConfigs[i] = IdpConfigFromModel(idp)
	}
	return convertedIDPConfigs
}

func IdpConfigFromModel(idp *model.IdpConfig) *IdpConfig {
	converted := &IdpConfig{
		ObjectRoot:  idp.ObjectRoot,
		IDPConfigID: idp.IDPConfigID,
		Name:        idp.Name,
		State:       int32(idp.State),
		Type:        int32(idp.Type),
		LogoSrc:     idp.LogoSrc,
	}
	if idp.OIDCConfig != nil {
		converted.OIDCIDPConfig = OidcIdpConfigFromModel(idp.OIDCConfig)
	}
	return converted
}

func IdpConfigToModel(idp *IdpConfig) *model.IdpConfig {
	converted := &model.IdpConfig{
		ObjectRoot:  idp.ObjectRoot,
		IDPConfigID: idp.IDPConfigID,
		Name:        idp.Name,
		LogoSrc:     idp.LogoSrc,
		State:       model.IdpConfigState(idp.State),
		Type:        model.IdpConfigType(idp.Type),
	}
	if idp.OIDCIDPConfig != nil {
		converted.OIDCConfig = OidcIdpConfigToModel(idp.OIDCIDPConfig)
	}
	return converted
}

func (iam *Iam) appendAddIdpConfigEvent(event *es_models.Event) error {
	idp := new(IdpConfig)
	err := idp.SetData(event)
	if err != nil {
		return err
	}
	idp.ObjectRoot.CreationDate = event.CreationDate
	iam.IDPs = append(iam.IDPs, idp)
	return nil
}

func (iam *Iam) appendChangeIdpConfigEvent(event *es_models.Event) error {
	idp := new(IdpConfig)
	err := idp.SetData(event)
	if err != nil {
		return err
	}
	if i, a := GetIdpConfig(iam.IDPs, idp.IDPConfigID); a != nil {
		iam.IDPs[i].SetData(event)
	}
	return nil
}

func (iam *Iam) appendRemoveIdpConfigEvent(event *es_models.Event) error {
	idp := new(IdpConfig)
	err := idp.SetData(event)
	if err != nil {
		return err
	}
	if i, a := GetIdpConfig(iam.IDPs, idp.IDPConfigID); a != nil {
		iam.IDPs[i] = iam.IDPs[len(iam.IDPs)-1]
		iam.IDPs[len(iam.IDPs)-1] = nil
		iam.IDPs = iam.IDPs[:len(iam.IDPs)-1]
	}
	return nil
}

func (iam *Iam) appendIdpConfigStateEvent(event *es_models.Event, state model.IdpConfigState) error {
	idp := new(IdpConfig)
	err := idp.SetData(event)
	if err != nil {
		return err
	}

	if i, a := GetIdpConfig(iam.IDPs, idp.IDPConfigID); a != nil {
		a.State = int32(state)
		iam.IDPs[i] = a
	}
	return nil
}

func (c *IdpConfig) SetData(event *es_models.Event) error {
	c.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, c); err != nil {
		logging.Log("EVEN-Msj9w").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
