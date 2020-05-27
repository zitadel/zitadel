package model

import (
	"net"
	"net/http"

	"github.com/caos/zitadel/internal/api"
	http_util "github.com/caos/zitadel/internal/api/http"
)

type BrowserInfo struct {
	UserAgent      string
	AcceptLanguage string
	RemoteIP       net.IP
}

func BrowserInfoFromRequest(r *http.Request) *BrowserInfo {
	return &BrowserInfo{
		UserAgent:      r.Header.Get(api.UserAgent),
		AcceptLanguage: r.Header.Get(api.AcceptLanguage),
		RemoteIP:       http_util.RemoteIPFromRequest(r),
	}
}

func (i *BrowserInfo) IsValid() bool {
	return i.UserAgent != "" &&
		i.AcceptLanguage != "" &&
		i.RemoteIP != nil && !i.RemoteIP.IsUnspecified()
}
