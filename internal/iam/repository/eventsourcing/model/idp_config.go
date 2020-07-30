package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
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
