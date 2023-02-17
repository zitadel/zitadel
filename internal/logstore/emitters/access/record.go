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
	a.RequestHeaders = normalizeHeaders(a.RequestHeaders, strings.ToLower(zitadel_http.Authorization), "grpcgateway-authorization", "cookie", "grpcgateway-cookie")
	a.ResponseHeaders = normalizeHeaders(a.ResponseHeaders, "set-cookie")
	return &a
}

// normalizeHeaders lowers all header keys and redacts secrets
func normalizeHeaders(header map[string][]string, redactKeysLower ...string) map[string][]string {
	return pruneKeys(redactKeys(lowerKeys(header), redactKeysLower...))
}

func lowerKeys(header map[string][]string) map[string][]string {
	lower := make(map[string][]string, len(header))
	for k, v := range header {
		lower[strings.ToLower(k)] = v
	}
	return lower
}

func redactKeys(header map[string][]string, redactKeysLower ...string) map[string][]string {
	redactedKeys := make(map[string][]string, len(header))
	for k, v := range header {
		redactedKeys[k] = v
	}
	for _, redactKey := range redactKeysLower {
		if _, ok := redactedKeys[redactKey]; ok {
			redactedKeys[redactKey] = []string{redacted}
		}
	}
	return redactedKeys
}

const maxValuesPerKey = 10

func pruneKeys(header map[string][]string) map[string][]string {
	prunedKeys := make(map[string][]string, len(header))
	for key, value := range header {
		valueItems := make([]string, 0, maxValuesPerKey)
		for i, valueItem := range value {
			// Max 10 header values per key
			if i > maxValuesPerKey {
				break
			}
			// Max 200 value length
			valueItems = append(valueItems, cutString(valueItem, 200))
		}
		prunedKeys[key] = valueItems
	}
	return prunedKeys
}

func cutString(str string, pos int) string {
	if len(str) <= pos {
		return str
	}
	return str[:pos-1]
}
