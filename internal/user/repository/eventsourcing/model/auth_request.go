package model

import (
	"encoding/json"
	"net"

	"github.com/caos/logging"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type AuthRequest struct {
	ID                  string `json:"id,omitempty"`
	UserAgentID         string `json:"userAgentID,omitempty"`
	SelectedIDPConfigID string `json:"selectedIDPConfigID,omitempty"`
	*BrowserInfo
}

type BrowserInfo struct {
	UserAgent      string `json:"userAgent,omitempty"`
	AcceptLanguage string `json:"acceptLanguage,omitempty"`
	RemoteIP       net.IP `json:"remoteIP,omitempty"`
}

func (a *AuthRequest) SetData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-T5df6").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-yGmhh", "could not unmarshal event")
	}
	return nil
}
