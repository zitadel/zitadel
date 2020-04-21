package model

import (
	"net"

	"github.com/caos/zitadel/internal/user_agent/model"
)

type BrowserInfo struct {
	UserAgent      string
	AcceptLanguage string
	RemoteIP       net.IP
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
