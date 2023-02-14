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

// normalizeHeaders lowers all header keys and redacts secrets
func normalizeHeaders(header http.Header, redactKeysLower ...string) {
	for k, v := range header {
		lowerKey := strings.ToLower(k)
		delete(header, k)
		for _, r := range redactKeysLower {
			if lowerKey == r {
				header[lowerKey] = []string{redacted}
				break
			}
		}
		vItems := make([]string, 0)
		for i, vitem := range v {
			// Max 10 header values per key
			if i > 10 {
				break
			}
			// Max 200 value length
			vItems = append(vItems, cutString(vitem, 200))
		}
		header[lowerKey] = vItems
	}
}

func cutString(str string, pos int) string {
	if len(str) <= pos {
		return str
	}
	return str[:pos]
}
