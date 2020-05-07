package model

import (
	"net"

	"github.com/caos/zitadel/internal/user/model"
)

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
