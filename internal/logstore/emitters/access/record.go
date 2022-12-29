package access

import (
	"net/http"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/logstore"

	zitadel_http "github.com/zitadel/zitadel/internal/api/http"
)

var _ logstore.LogRecord = (*Record)(nil)

type Record struct {
	Timestamp       time.Time   `json:"ts"`
	Protocol        Protocol    `json:"protocol"`
	RequestURL      string      `json:"requestURL"`
	ResponseStatus  uint32      `json:"responseStatus"`
	RequestHeaders  http.Header `json:"requestHeaders"`
	ResponseHeaders http.Header `json:"responseHeaders"`
	InstanceID      string      `json:"instanceID"`
	ProjectID       string      `json:"projectID"`
	RequestedDomain string      `json:"requestedDomain"`
	RequestedHost   string      `json:"requestedHost"`
}

type Protocol uint8

const (
	GRPC Protocol = iota
	HTTP
	// TODO: GRPC-Web?
	// TODO: HTTPS?

	redacted = "[REDACTED]"
)

func (a *Record) Redact() logstore.LogRecord {
	clone := &(*a)
	redactHeaders(clone.RequestHeaders, strings.ToLower(zitadel_http.Authorization), "cookie")
	redactHeaders(clone.ResponseHeaders, "set-cookie")
	return clone
}

func redactHeaders(header http.Header, redactKeysLower ...string) {
	for k := range header {
		lowerKey := strings.ToLower(k)
		for _, r := range redactKeysLower {
			if lowerKey == r {
				header[k] = []string{redacted}
				break
			}
		}
	}
}
