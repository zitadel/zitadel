package oidc

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/go-jose/go-jose/v4"
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
		signer         jose.Signer
		expectedParams map[string]string
	}{
		{
			testName:     "basic with only redirectURI",
			logoutURIStr: "https://example.com/logout",
			redirectURI:  "https://client/cb",
			expectedParams: map[string]string{
				"post_logout_redirect": "https://client/cb",
				"logout_token":         "", // presence checked separately
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
				"logout_token":         "", // presence checked separately
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
				"logout_token":         "", // presence checked separately
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
				"logout_token":         "", // presence checked separately
			},
		},
		{
			testName:     "base with trailing slash",
			logoutURIStr: "https://example.com/logout/",
			redirectURI:  "https://client/cb",
			expectedParams: map[string]string{
				"post_logout_redirect": "https://client/cb",
				"logout_token":         "", // presence checked separately
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
			got, err := buildLoginV2LogoutURL(logoutURI, tc.redirectURI, tc.logoutHint, tc.uiLocales, signer)

			// Then
			require.NoError(t, err)
			gotURL, err := url.Parse(got)
			require.NoError(t, err)
			require.NotContains(t, gotURL.String(), "/logout/")

			q := gotURL.Query()
			// Ensure no unexpected params
			require.Len(t, q, len(tc.expectedParams))

			for k, v := range tc.expectedParams {
				if k == LoginLogoutTokenParam {
					assertLogoutToken(t, q.Get(k), &logoutTokenPayload{
						PostLogoutRedirectURI: tc.redirectURI,
						LogoutHint:            tc.logoutHint,
						UILocales:             tc.uiLocales,
					})
					continue
				}
				assert.Equal(t, v, q.Get(k))
			}
		})
	}
}

func assertLogoutToken(t *testing.T, token string, payload *logoutTokenPayload) {
	signature, err := jose.ParseSigned(token, []jose.SignatureAlgorithm{jose.RS256})
	require.NoError(t, err)
	logoutToken := new(logoutTokenPayload)
	err = json.Unmarshal(signature.UnsafePayloadWithoutVerification(), logoutToken)
	require.NoError(t, err)
	assert.Equal(t, payload, logoutToken)
}

var (
	privKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	signer     = func() jose.Signer {
		signer, _ := jose.NewSigner(
			jose.SigningKey{
				Algorithm: jose.RS256,
				Key:       privKey,
			},
			(&jose.SignerOptions{}).WithType("JWT"),
		)
		return signer
	}()
)
