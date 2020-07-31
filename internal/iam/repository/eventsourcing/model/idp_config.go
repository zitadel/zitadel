package model

import (
	"encoding/json"
	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
)

type IDPConfig struct {
	es_models.ObjectRoot
	IDPConfigID   string         `json:"idpConfigId"`
	State         int32          `json:"-"`
	Name          string         `json:"name,omitempty"`
	Type          int32          `json:"idpType,omitempty"`
	LogoSrc       string         `json:"logoSrc,omitempty"`
	OIDCIDPConfig *OIDCIDPConfig `json:"-"`
}

type IDPConfigID struct {
	es_models.ObjectRoot
	IDPConfigID string `json:"idpConfigId"`
}

func GetIDPConfig(idps []*IDPConfig, id string) (int, *IDPConfig) {
	for i, idp := range idps {
		if idp.IDPConfigID == id {
			return i, idp
		}
	}
	return -1, nil
}

func (c *IDPConfig) Changes(changed *IDPConfig) map[string]interface{} {
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

func IDPConfigsToModel(idps []*IDPConfig) []*model.IDPConfig {
	convertedIDPConfigs := make([]*model.IDPConfig, len(idps))
	for i, idp := range idps {
		convertedIDPConfigs[i] = IDPConfigToModel(idp)
	}
	return convertedIDPConfigs
}

func IDPConfigsFromModel(idps []*model.IDPConfig) []*IDPConfig {
	convertedIDPConfigs := make([]*IDPConfig, len(idps))
	for i, idp := range idps {
		convertedIDPConfigs[i] = IDPConfigFromModel(idp)
	}
	return convertedIDPConfigs
}

func IDPConfigFromModel(idp *model.IDPConfig) *IDPConfig {
	converted := &IDPConfig{
		ObjectRoot:  idp.ObjectRoot,
		IDPConfigID: idp.IDPConfigID,
		Name:        idp.Name,
		State:       int32(idp.State),
		Type:        int32(idp.Type),
		LogoSrc:     idp.LogoSrc,
	}
	if idp.OIDCConfig != nil {
		converted.OIDCIDPConfig = OIDCIDPConfigFromModel(idp.OIDCConfig)
	}
	return converted
}

func IDPConfigToModel(idp *IDPConfig) *model.IDPConfig {
	converted := &model.IDPConfig{
		ObjectRoot:  idp.ObjectRoot,
		IDPConfigID: idp.IDPConfigID,
		Name:        idp.Name,
		LogoSrc:     idp.LogoSrc,
		State:       model.IDPConfigState(idp.State),
		Type:        model.IDPConfigType(idp.Type),
	}
	if idp.OIDCIDPConfig != nil {
		converted.OIDCConfig = OIDCIDPConfigToModel(idp.OIDCIDPConfig)
	}
	return converted
}

func (iam *Iam) appendAddIdpConfigEvent(event *es_models.Event) error {
	idp := new(IDPConfig)
	err := idp.setData(event)
	if err != nil {
		return err
	}
	idp.ObjectRoot.CreationDate = event.CreationDate
	iam.IDPs = append(iam.IDPs, idp)
	return nil
}

func (iam *Iam) appendChangeIdpConfigEvent(event *es_models.Event) error {
	idp := new(IDPConfig)
	err := idp.setData(event)
	if err != nil {
		return err
	}
	if i, a := GetIDPConfig(iam.IDPs, idp.IDPConfigID); a != nil {
		iam.IDPs[i].setData(event)
	}
	return nil
}

func (iam *Iam) appendRemoveIdpConfigEvent(event *es_models.Event) error {
	idp := new(IDPConfig)
	err := idp.setData(event)
	if err != nil {
		return err
	}
	if i, a := GetIDPConfig(iam.IDPs, idp.IDPConfigID); a != nil {
		iam.IDPs[i] = iam.IDPs[len(iam.IDPs)-1]
		iam.IDPs[len(iam.IDPs)-1] = nil
		iam.IDPs = iam.IDPs[:len(iam.IDPs)-1]
	}
	return nil
}

func (iam *Iam) appendIdpConfigStateEvent(event *es_models.Event, state model.IDPConfigState) error {
	idp := new(IDPConfig)
	err := idp.setData(event)
	if err != nil {
		return err
	}

	if i, a := GetIDPConfig(iam.IDPs, idp.IDPConfigID); a != nil {
		a.State = int32(state)
		iam.IDPs[i] = a
	}
	return nil
}

func (c *IDPConfig) setData(event *es_models.Event) error {
	c.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, c); err != nil {
		logging.Log("EVEN-Msj9w").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
