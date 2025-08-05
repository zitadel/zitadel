package domain

import (
	"net"
	httplib "net/http"
)

type UserAgent struct {
	FingerprintID *string        `json:"fingerprint_id,omitempty"`
	IP            net.IP         `json:"ip,omitempty"`
	Description   *string        `json:"description,omitempty"`
	Header        httplib.Header `json:"header,omitempty"`
}

func (ua UserAgent) IsEmpty() bool {
	return ua.FingerprintID == nil && len(ua.IP) == 0 && ua.Description == nil && ua.Header == nil
}

func (ua *UserAgent) GetFingerprintID() string {
	if ua == nil || ua.FingerprintID == nil {
		return ""
	}
	return *ua.FingerprintID
}
