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
	// NotCountable can be used by the logging service to explicitly stating,
	// that the request must not increase the amount of countable (authenticated) requests
	NotCountable bool `json:"-"`
	normalized   bool `json:"-"`
}

type AccessProtocol uint8

const (
	GRPC AccessProtocol = iota
	HTTP

	redacted = "[REDACTED]"
)

var (
	unaccountableEndpoints = []string{
		"/zitadel.system.v1.SystemService/",
		"/zitadel.admin.v1.AdminService/Healthz",
		"/zitadel.management.v1.ManagementService/Healthz",
		"/zitadel.management.v1.ManagementService/GetOIDCInformation",
		"/zitadel.auth.v1.AuthService/Healthz",
	}
)

func (a AccessLog) IsAuthenticated() bool {
	if a.NotCountable {
		return false
	}
	if !a.normalized {
		panic("access log not normalized, Normalize() must be called before IsAuthenticated()")
	}
	_, hasHTTPAuthHeader := a.RequestHeaders[strings.ToLower(zitadel_http.Authorization)]
	// ignore requests, which were unauthorized or do not require an authorization (even if one was sent)
	// also ignore if the limit was already reached or if the server returned an internal error
	// not that endpoints paths are only checked with the gRPC representation as HTTP (gateway) will not log them
	return hasHTTPAuthHeader &&
		(a.Protocol == HTTP &&
			a.ResponseStatus != http.StatusInternalServerError &&
			a.ResponseStatus != http.StatusTooManyRequests &&
			a.ResponseStatus != http.StatusUnauthorized) ||
		(a.Protocol == GRPC &&
			a.ResponseStatus != uint32(codes.Internal) &&
			a.ResponseStatus != uint32(codes.ResourceExhausted) &&
			a.ResponseStatus != uint32(codes.Unauthenticated) &&
			!a.isUnaccountableEndpoint())
}

func (a AccessLog) isUnaccountableEndpoint() bool {
	for _, endpoint := range unaccountableEndpoints {
		if strings.HasPrefix(a.RequestURL, endpoint) {
			return true
		}
	}
	return false
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
