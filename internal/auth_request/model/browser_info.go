package model

import "net"

type BrowserInfo struct {
	UserAgent      string
	AcceptLanguage string
	RemoteIP       net.IP
}

func (i *BrowserInfo) IsValid() bool {
	return i.UserAgent != "" &&
		i.AcceptLanguage != "" &&
		i.RemoteIP != nil && !i.RemoteIP.IsUnspecified()
}
