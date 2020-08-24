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
