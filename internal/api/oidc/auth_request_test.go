package oidc

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestBuildLoginV2LogoutURL(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName       string
		logoutURIStr   string
		redirectURI    string
		logoutHint     string
		uiLocales      []language.Tag
		expectedParams map[string]string
	}{
		{
			testName:     "basic with only redirectURI",
			logoutURIStr: "https://example.com/logout",
			redirectURI:  "https://client/cb",
			expectedParams: map[string]string{
				"post_logout_redirect": "https://client/cb",
			},
		},
		{
			testName:     "with logout hint",
			logoutURIStr: "https://example.com/logout",
			redirectURI:  "https://client/cb",
			logoutHint:   "user@example.com",
			expectedParams: map[string]string{
				"post_logout_redirect": "https://client/cb",
				"logout_hint":          "user@example.com",
			},
		},
		{
			testName:     "with ui_locales",
			logoutURIStr: "https://example.com/logout",
			redirectURI:  "https://client/cb",
			uiLocales:    []language.Tag{language.English, language.Italian},
			expectedParams: map[string]string{
				"post_logout_redirect": "https://client/cb",
				"ui_locales":           "en it",
			},
		},
		{
			testName:     "with all params",
			logoutURIStr: "https://example.com/logout",
			redirectURI:  "https://client/cb",
			logoutHint:   "logoutme",
			uiLocales:    []language.Tag{language.Make("de-CH"), language.Make("fr")},
			expectedParams: map[string]string{
				"post_logout_redirect": "https://client/cb",
				"logout_hint":          "logoutme",
				"ui_locales":           "de-CH fr",
			},
		},
		{
			testName:     "base with trailing slash",
			logoutURIStr: "https://example.com/logout/",
			redirectURI:  "https://client/cb",
			expectedParams: map[string]string{
				"post_logout_redirect": "https://client/cb",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// t.Parallel()

			// Given
			logoutURI, err := url.Parse(tc.logoutURIStr)
			require.NoError(t, err)

			// When
			got := buildLoginV2LogoutURL(logoutURI, tc.redirectURI, tc.logoutHint, tc.uiLocales)

			// Then
			gotURL, err := url.Parse(got)
			require.NoError(t, err)
			require.NotContains(t, gotURL.String(), "/logout/")

			q := gotURL.Query()
			// Ensure no unexpected params
			require.Len(t, q, len(tc.expectedParams))

			for k, v := range tc.expectedParams {
				assert.Equal(t, v, q.Get(k))
			}
		})
	}
}
