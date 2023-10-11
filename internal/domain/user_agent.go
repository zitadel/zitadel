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
