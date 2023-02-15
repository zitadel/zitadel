package access

import (
	"strings"
	"time"

	zitadel_http "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/logstore"
)

var _ logstore.LogRecord = (*Record)(nil)

type Record struct {
	LogDate        time.Time `json:"logDate"`
	Protocol       Protocol  `json:"protocol"`
	RequestURL     string    `json:"requestUrl"`
	ResponseStatus uint32    `json:"responseStatus"`
	// RequestHeaders are plain maps so varying implementation
	// between HTTP and gRPC don't interfere with each other
	RequestHeaders map[string][]string `json:"requestHeaders"`
	// ResponseHeaders are plain maps so varying implementation
	// between HTTP and gRPC don't interfere with each other
	ResponseHeaders map[string][]string `json:"responseHeaders"`
	InstanceID      string              `json:"instanceId"`
	ProjectID       string              `json:"projectId"`
	RequestedDomain string              `json:"requestedDomain"`
	RequestedHost   string              `json:"requestedHost"`
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
func normalizeHeaders(header map[string][]string, redactKeysLower ...string) {
	lowerKeys(header)
	redactKeys(header, redactKeysLower...)
	pruneKeys(header)
}

func lowerKeys(header map[string][]string) {
	for k, v := range header {
		delete(header, k)
		header[strings.ToLower(k)] = v
	}
}

func redactKeys(header map[string][]string, redactKeysLower ...string) {
	for _, redactKey := range redactKeysLower {
		if _, ok := header[redactKey]; ok {
			header[redactKey] = []string{redacted}
		}
	}
}

func pruneKeys(header map[string][]string) {
	for key, value := range header {
		valueItems := make([]string, 0, maxValuesPerKey)
		for i, valueItem := range value {
			// Max 10 header values per key
			if i > maxValuesPerKey {
				break
			}
			// Max 200 value length
			valueItems[i] = cutString(valueItem, 200)
		}
		header[key] = valueItems
	}
}

func cutString(str string, pos int) string {
	if len(str) <= pos {
		return str
	}
	return str[:pos]
}
