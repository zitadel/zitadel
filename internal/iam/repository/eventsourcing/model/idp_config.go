package model

import (
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

func GetIDPConfif(idps []*IDPConfig, id string) (int, *IDPConfig) {
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
