package domain

import (
	"net"
	net_http "net/http"

	http_util "github.com/zitadel/zitadel/internal/api/http"
)

type BrowserInfo struct {
	UserAgent      string
	AcceptLanguage string
	RemoteIP       net.IP
	Header         net_http.Header
}

func BrowserInfoFromRequest(r *net_http.Request) *BrowserInfo {
	return &BrowserInfo{
		UserAgent:      r.Header.Get(http_util.UserAgentHeader),
		AcceptLanguage: r.Header.Get(http_util.AcceptLanguage),
		RemoteIP:       http_util.RemoteIPFromRequest(r),
		Header:         r.Header,
	}
}

func (b *BrowserInfo) ToUserAgent() *UserAgent {
	if b == nil {
		return nil
	}
	return &UserAgent{
		FingerprintID: &b.UserAgent,
		IP:            b.RemoteIP,
		Description:   &b.UserAgent,
		Header:        b.Header,
	}
}
