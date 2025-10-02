package captcha

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
)

type fakeRoundTripper struct {
	respBody string
	status   int
}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: f.status,
		Body:       io.NopCloser(strings.NewReader(f.respBody)),
		Header:     make(http.Header),
	}, nil
}

var testCases = []struct {
	name             string
	captchaType      domain.CaptchaType
	captchaSiteKey   string
	captchaSecretKey string
	serverResponse   string
	serverToken      string
	statusCode       int
	expectedErr      bool
}{
	{
		name:             "valid captcha",
		captchaType:      domain.CaptchaTypeReCaptcha,
		captchaSiteKey:   "site-key",
		captchaSecretKey: "secret-key",
		serverResponse:   `{"success": true, "challenge_ts": "2025-01-01T00:00:00Z", "hostname": "example.com"}`,
		serverToken:      "faketoken",
		statusCode:       200,
		expectedErr:      false,
	},
	{
		name:             "invalid captcha",
		captchaType:      domain.CaptchaTypeReCaptcha,
		captchaSiteKey:   "site-key",
		captchaSecretKey: "secret-key",
		serverResponse:   `{"success": false, "challenge_ts": "2025-01-01T00:00:00Z", "hostname": "example.com"}`,
		serverToken:      "faketoken",
		statusCode:       200,
		expectedErr:      true,
	},
	{
		name:             "backend error",
		captchaType:      domain.CaptchaTypeReCaptcha,
		captchaSiteKey:   "site-key",
		captchaSecretKey: "secret-key",
		serverResponse:   `{}`,
		serverToken:      "faketoken",
		statusCode:       502,
		expectedErr:      true,
	},
	{
		name:             "missing site key",
		captchaType:      domain.CaptchaTypeReCaptcha,
		captchaSiteKey:   "",
		captchaSecretKey: "secret-key",
		serverResponse:   `{"success": true, "challenge_ts": "2025-01-01T00:00:00Z", "hostname": "example.com"}`,
		serverToken:      "faketoken",
		statusCode:       200,
		expectedErr:      true,
	},
	{
		name:             "unspecified captcha type",
		captchaType:      0,
		captchaSiteKey:   "site-key",
		captchaSecretKey: "secret-key",
		serverResponse:   `{"success": true, "challenge_ts": "2025-01-01T00:00:00Z", "hostname": "example.com"}`,
		serverToken:      "faketoken",
		statusCode:       200,
		expectedErr:      true,
	},
	{
		name:             "invalid captcha type",
		captchaType:      99,
		captchaSiteKey:   "site-key",
		captchaSecretKey: "secret-key",
		serverResponse:   `{"success": true, "challenge_ts": "2025-01-01T00:00:00Z", "hostname": "example.com"}`,
		serverToken:      "faketoken",
		statusCode:       200,
		expectedErr:      true,
	},
	{
		name:             "missing token",
		captchaType:      domain.CaptchaTypeReCaptcha,
		captchaSiteKey:   "site-key",
		captchaSecretKey: "secret-key",
		serverResponse:   `{"success": true, "challenge_ts": "2025-01-01T00:00:00Z", "hostname": "example.com"}`,
		serverToken:      "",
		statusCode:       200,
		expectedErr:      true,
	},
}

func TestVerifyCaptchaWithFakeClient(t *testing.T) {
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Setup fake RecaptchaV2 with mock HTTP client
			fakeClient := &RecaptchaV2{
				client: &http.Client{
					Transport: &fakeRoundTripper{
						status:   tt.statusCode,
						respBody: tt.serverResponse,
					},
				},
			}

			// Override factory
			captchaFactory[domain.CaptchaTypeReCaptcha] = func() Captcha {
				return fakeClient
			}
			t.Cleanup(func() {
				// Reset factory after test
				captchaFactory[domain.CaptchaTypeReCaptcha] = func() Captcha {
					return &RecaptchaV2{}
				}
			})

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf("g-recaptcha-response=%s", tt.serverToken)))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.ParseForm()

			authReq := &domain.AuthRequest{
				LoginPolicy: &domain.LoginPolicy{
					CaptchaType:      tt.captchaType,
					CaptchaSiteKey:   tt.captchaSiteKey,
					CaptchaSecretKey: tt.captchaSecretKey,
				},
			}

			err := VerifyCaptcha(req, authReq)
			if (err != nil) != tt.expectedErr {
				t.Errorf("VerifyCaptcha() error = %v, wantErr %v", err, tt.expectedErr)
			}
		})
	}
}
