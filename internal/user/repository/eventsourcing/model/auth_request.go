package model

import (
	"net"

	"github.com/zitadel/logging"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
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

func (a *AuthRequest) SetData(event eventstore.Event) error {
	if err := event.Unmarshal(a); err != nil {
		logging.Log("EVEN-T5df6").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-yGmhh", "could not unmarshal event")
	}
	return nil
}
