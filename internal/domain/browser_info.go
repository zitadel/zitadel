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

func (a *AuthRequest) ToUserAgent() *UserAgent {
	agent := &UserAgent{
		FingerprintID: &a.AgentID,
	}
	if a.BrowserInfo == nil {
		return agent
	}
	agent.IP = a.BrowserInfo.RemoteIP
	agent.Description = &a.BrowserInfo.UserAgent
	agent.Header = a.BrowserInfo.Header
	return agent
}
