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
	LogDate         time.Time   `json:"logDate"`
	Protocol        Protocol    `json:"protocol"`
	RequestURL      string      `json:"requestUrl"`
	ResponseStatus  uint32      `json:"responseStatus"`
	RequestHeaders  http.Header `json:"requestHeaders"`
	ResponseHeaders http.Header `json:"responseHeaders"`
	InstanceID      string      `json:"instanceId"`
	ProjectID       string      `json:"projectId"`
	RequestedDomain string      `json:"requestedDomain"`
	RequestedHost   string      `json:"requestedHost"`
}

type Protocol uint8

const (
	GRPC Protocol = iota
	HTTP

	redacted = "[REDACTED]"
)

func (a Record) Normalize() logstore.LogRecord {
	a.RequestedDomain = cutString(a.RequestedDomain, 200)
	a.RequestURL = cutString(a.RequestURL, 200)
	normalizeHeaders(a.RequestHeaders, strings.ToLower(zitadel_http.Authorization), "grpcgateway-authorization", "cookie", "grpcgateway-cookie")
	normalizeHeaders(a.ResponseHeaders, "set-cookie")
	return &a
}

const maxValuesPerKey = 10

// normalizeHeaders lowers all header keys and redacts secrets
func normalizeHeaders(header http.Header, redactKeysLower ...string) {
	// set readacted values where needed
	for _, key := range redactKeysLower {
		if len(header.Values(key)) > 0 {
			header.Set(key, redacted)
		}
	}

	// normalize keys to lowercase and ensure limit
	for key, values := range header {
		lowerKey := strings.ToLower(key)
		header.Del(key)

		for i, value := range values {
			// Max 10 header values per key
			if i > maxValuesPerKey {
				break
			}
			// Max 200 value length
			header.Add(lowerKey, cutString(value, 200))
		}
	}
}

func cutString(str string, pos int) string {
	if len(str) <= pos {
		return str
	}
	return str[:pos]
}
