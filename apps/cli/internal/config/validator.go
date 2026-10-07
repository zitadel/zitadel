package config

import (
	"fmt"
	"net/url"
	"strings"
)

// Validate ensures the context has the required fields and well-formed URLs.
func (c *Context) Validate() error {
	if c.Instance == "" {
		return fmt.Errorf("instance URL cannot be empty")
	}

	// Basic URL validation
	rawURL := c.Instance
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid instance URL %q: %w", c.Instance, err)
	}
	if u.Host == "" {
		return fmt.Errorf("invalid instance URL %q: missing host", c.Instance)
	}

	switch c.AuthMethod {
	case "pat":
		if c.PAT == "" {
			return fmt.Errorf("auth method is 'pat' but no PAT is configured")
		}
	case "interactive":
		if c.ClientID == "" {
			return fmt.Errorf("auth method is 'interactive' but no Client ID is configured")
		}
		if c.Token == "" {
			return fmt.Errorf("no access token available, please run 'zitadel-cli login' again")
		}
	default:
		return fmt.Errorf("unknown auth method %q", c.AuthMethod)
	}

	return nil
}