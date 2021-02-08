package domain

import (
	"net"
	net_http "net/http"

	http_util "github.com/caos/zitadel/internal/api/http"
)

type BrowserInfo struct {
	UserAgent      string
	AcceptLanguage string
	RemoteIP       net.IP
}

func BrowserInfoFromRequest(r *net_http.Request) *BrowserInfo {
	return &BrowserInfo{
		UserAgent:      r.Header.Get(http_util.UserAgentHeader),
		AcceptLanguage: r.Header.Get(http_util.AcceptLanguage),
		RemoteIP:       http_util.RemoteIPFromRequest(r),
	}
}
