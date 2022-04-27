package model

import (
	"net"
	"net/http"

	http_util "github.com/zitadel/zitadel/internal/api/http"
)

type BrowserInfo struct {
	UserAgent      string
	AcceptLanguage string
	RemoteIP       net.IP
}

func BrowserInfoFromRequest(r *http.Request) *BrowserInfo {
	return &BrowserInfo{
		UserAgent:      r.Header.Get(http_util.UserAgentHeader),
		AcceptLanguage: r.Header.Get(http_util.AcceptLanguage),
		RemoteIP:       http_util.RemoteIPFromRequest(r),
	}
}

func (i *BrowserInfo) IsValid() bool {
	return i.UserAgent != "" &&
		i.AcceptLanguage != "" &&
		i.RemoteIP != nil && !i.RemoteIP.IsUnspecified()
}
