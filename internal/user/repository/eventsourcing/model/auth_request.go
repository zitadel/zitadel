package model

import (
	"net"

	"github.com/caos/zitadel/internal/auth_request/model"
)

type AuthRequest struct {
	UserAgentID string `json:"userAgentID,omitempty"`
	*BrowserInfo
}

func AuthRequestFromModel(request *model.AuthRequest) *AuthRequest {
	return &AuthRequest{
		UserAgentID: request.ObjectRoot.AggregateID,
		BrowserInfo: BrowserInfoFromModel(request.BrowserInfo),
	}
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

//
//func AuthRequestToModel(request *AuthRequest) *model.AuthRequest {
//	return &model.AuthRequest{
//		UserAgentID: request.UserAgentID,
//		BrowserInfo: BrowserInfoToModel(request.BrowserInfo),
//	}
//}
