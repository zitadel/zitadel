package model

import (
	"encoding/json"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"net"

	"github.com/caos/zitadel/internal/auth_request/model"
)

type AuthRequest struct {
	ID                  string `json:"id,omitempty"`
	UserAgentID         string `json:"userAgentID,omitempty"`
	SelectedIDPConfigID string `json:"selectedIDPConfigID,omitempty"`
	*BrowserInfo
}

func AuthRequestFromModel(request *model.AuthRequest) *AuthRequest {
	req := &AuthRequest{
		ID:                  request.ID,
		UserAgentID:         request.AgentID,
		SelectedIDPConfigID: request.SelectedIDPConfigID,
	}
	if request.BrowserInfo != nil {
		req.BrowserInfo = BrowserInfoFromModel(request.BrowserInfo)
	}
	return req
}

func AuthRequestToModel(request *AuthRequest) *model.AuthRequest {
	req := &model.AuthRequest{
		ID:                  request.ID,
		AgentID:             request.UserAgentID,
		SelectedIDPConfigID: request.SelectedIDPConfigID,
	}
	if request.BrowserInfo != nil {
		req.BrowserInfo = BrowserInfoToModel(request.BrowserInfo)
	}
	return req
}

type BrowserInfo struct {
	UserAgent      string `json:"userAgent,omitempty"`
	AcceptLanguage string `json:"acceptLanguage,omitempty"`
	RemoteIP       net.IP `json:"remoteIP,omitempty"`
}

func BrowserInfoFromModel(info *model.BrowserInfo) *BrowserInfo {
	return &BrowserInfo{
		UserAgent:      info.UserAgent,
		AcceptLanguage: info.AcceptLanguage,
		RemoteIP:       info.RemoteIP,
	}
}

func BrowserInfoToModel(info *BrowserInfo) *model.BrowserInfo {
	return &model.BrowserInfo{
		UserAgent:      info.UserAgent,
		AcceptLanguage: info.AcceptLanguage,
		RemoteIP:       info.RemoteIP,
	}
}
func (a *AuthRequest) SetData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-T5df6").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-yGmhh", "could not unmarshal event")
	}
	return nil
}
