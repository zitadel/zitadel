package http

import (
	"context"
	"fmt"
	"net/url"
)

func GetOriginFromURLString(s string) (string, error) {
	parsed, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host), nil
}

func IsOriginAllowed(allowList []string, origin string) bool {
	for _, allowed := range allowList {
		if allowed == origin {
			return true
		}
	}
	return false
}

// IsOrigin checks if provided string is an origin (scheme://hostname[:port]) without path, query or fragment
func IsOrigin(rawOrigin string) bool {
	parsedUrl, err := url.Parse(rawOrigin)
	if err != nil {
		return false
	}
	return parsedUrl.Scheme != "" && parsedUrl.Host != "" && parsedUrl.Path == "" && len(parsedUrl.Query()) == 0 && parsedUrl.Fragment == ""
}

func BuildHTTP(hostname string, externalPort uint16, secure bool) string {
	if externalPort == 0 || (externalPort == 443 && secure) || (externalPort == 80 && !secure) {
		return BuildOrigin(hostname, secure)
	}
	return BuildOrigin(fmt.Sprintf("%s:%d", hostname, externalPort), secure)
}

func BuildOrigin(host string, secure bool) string {
	schema := "https"
	if !secure {
		schema = "http"
	}
	return fmt.Sprintf("%s://%s", schema, host)
}

type RequestOrigin struct {
	Full      string // Full is the full origin including scheme and host
	Host      string // Host includes the port if not standard
	Domain    string // Domain is the host without the port
	Scheme    string // Scheme is http or https
	IsDefault bool   // IsDefault is true if the origin is not read from a request but from the runtime config
}

func RequestOriginFromCtx(ctx context.Context) RequestOrigin {
	o, _ := ctx.Value(origin).(RequestOrigin)
	return o
}

func WithRequestOrigin(ctx context.Context, value RequestOrigin) context.Context {
	return context.WithValue(ctx, origin, value)
}

func WithDefaultOrigin(ctx context.Context, externalSecure bool, externalDomain string, externalPort uint16) context.Context {
	defaultOrigin := RequestOrigin{
		Full:      BuildHTTP(externalDomain, externalPort, externalSecure),
		Host:      externalDomain,
		Domain:    externalDomain,
		Scheme:    "http",
		IsDefault: true,
	}
	if externalSecure {
		defaultOrigin.Scheme = "https"
	}
	if externalPort != 80 && externalPort != 443 {
		defaultOrigin.Host = fmt.Sprintf("%s:%d", defaultOrigin.Host, externalPort)
	}
	return WithRequestOrigin(ctx, defaultOrigin)
}
