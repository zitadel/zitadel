package schemas

import (
	"encoding/json"
	"net/url"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type HttpURL url.URL

func ParseHTTPURL(rawURL string) (*HttpURL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "SCIM-htturl1", "HTTP URL expected, got %v", parsedURL.Scheme)
	}

	return (*HttpURL)(parsedURL), nil
}

func (u *HttpURL) UnmarshalJSON(data []byte) error {
	var urlStr string
	if err := json.Unmarshal(data, &urlStr); err != nil {
		return err
	}

	parsedURL, err := ParseHTTPURL(urlStr)
	if err != nil {
		return err
	}

	*u = *parsedURL
	return nil
}

func (u *HttpURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *HttpURL) String() string {
	if u == nil {
		return ""
	}

	return (*url.URL)(u).String()
}
