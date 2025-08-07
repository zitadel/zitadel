package http

import (
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
