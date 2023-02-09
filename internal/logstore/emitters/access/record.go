package access

import (
	"net/http"
	"strings"
	"time"

	zitadel_http "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/logstore"
)

var _ logstore.LogRecord = (*Record)(nil)

type Record struct {
	LogDate         time.Time
	Protocol        Protocol
	RequestURL      string
	ResponseStatus  uint32
	RequestHeaders  http.Header
	ResponseHeaders http.Header
	InstanceID      string
	ProjectID       string
	RequestedDomain string
	RequestedHost   string
}

type Protocol uint8

const (
	GRPC Protocol = iota
	HTTP

	redacted = "[REDACTED]"
)

func (a Record) Normalize() logstore.LogRecord {
	normalizeHeaders(a.RequestHeaders, strings.ToLower(zitadel_http.Authorization), "grpcgateway-authorization", "cookie", "grpcgateway-cookie")
	normalizeHeaders(a.ResponseHeaders, "set-cookie")
	return &a
}

// normalizeHeaders lowers all header keys and redacts secrets
func normalizeHeaders(header http.Header, redactKeysLower ...string) {
	for k, v := range header {
		lowerKey := strings.ToLower(k)
		delete(header, k)
		header[lowerKey] = v
		for _, r := range redactKeysLower {
			if lowerKey == r {
				header[lowerKey] = []string{redacted}
				break
			}
		}
	}
}
