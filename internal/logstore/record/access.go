package record

import (
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc/codes"

	zitadel_http "github.com/zitadel/zitadel/internal/api/http"
)

type AccessLog struct {
	LogDate        time.Time      `json:"logDate"`
	Protocol       AccessProtocol `json:"protocol"`
	RequestURL     string         `json:"requestUrl"`
	ResponseStatus uint32         `json:"responseStatus"`
	// RequestHeaders and ResponseHeaders are plain maps so varying implementations
	// between HTTP and gRPC don't interfere with each other
	RequestHeaders  map[string][]string `json:"requestHeaders"`
	ResponseHeaders map[string][]string `json:"responseHeaders"`
	InstanceID      string              `json:"instanceId"`
	ProjectID       string              `json:"projectId"`
	RequestedDomain string              `json:"requestedDomain"`
	RequestedHost   string              `json:"requestedHost"`
	normalized      bool                `json:"-"`
}

type AccessProtocol uint8

const (
	GRPC AccessProtocol = iota
	HTTP

	redacted = "[REDACTED]"
)

func (a AccessLog) IsAuthenticated() bool {
	if !a.normalized {
		panic("access log not normalized, Normalize() must be called before IsAuthenticated()")
	}
	// TODO: Is it possible to maliciously produce usage on public endpoints like this, just by adding an auth header?
	_, hasHTTPAuthHeader := a.RequestHeaders[strings.ToLower(zitadel_http.Authorization)]
	return hasHTTPAuthHeader &&
		!strings.HasPrefix(a.RequestURL, "/zitadel.system.v1.SystemService/") &&
		!strings.HasPrefix(a.RequestURL, "/system/v1/") &&
		(a.Protocol == HTTP &&
			a.ResponseStatus != http.StatusForbidden &&
			a.ResponseStatus != http.StatusInternalServerError &&
			a.ResponseStatus != http.StatusTooManyRequests) ||
		(a.Protocol == GRPC &&
			a.ResponseStatus != uint32(codes.PermissionDenied) &&
			a.ResponseStatus != uint32(codes.Internal) &&
			a.ResponseStatus != uint32(codes.ResourceExhausted))
}

func (a AccessLog) Normalize() *AccessLog {
	a.RequestedDomain = cutString(a.RequestedDomain, 200)
	a.RequestURL = cutString(a.RequestURL, 200)
	a.RequestHeaders = normalizeHeaders(a.RequestHeaders, strings.ToLower(zitadel_http.Authorization), "grpcgateway-authorization", "cookie", "grpcgateway-cookie")
	a.ResponseHeaders = normalizeHeaders(a.ResponseHeaders, "set-cookie")
	a.normalized = true
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
