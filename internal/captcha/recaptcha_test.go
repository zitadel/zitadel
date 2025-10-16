package captcha

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecaptchaV2_Initialize(t *testing.T) {
	tests := []struct {
		name    string
		config  any
		wantErr bool
	}{
		{
			name: "valid config",
			config: RecaptchaV2Config{
				SiteKey:   "test-site-key",
				SecretKey: "test-secret-key",
			},
			wantErr: false,
		},
		{
			name:    "invalid config type",
			config:  "invalid",
			wantErr: true,
		},
		{
			name: "empty keys",
			config: RecaptchaV2Config{
				SiteKey:   "",
				SecretKey: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RecaptchaV2{}
			err := r.Initialize(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecaptchaV2.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRecaptchaV2_GetToken(t *testing.T) {
	tests := []struct {
		name    string
		form    string
		want    string
		wantErr bool
	}{
		{
			name:    "valid token",
			form:    "g-recaptcha-response=valid-token",
			want:    "valid-token",
			wantErr: false,
		},
		{
			name:    "missing token",
			form:    "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RecaptchaV2{}
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.form))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			got, err := r.GetToken(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecaptchaV2.GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecaptchaV2.GetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecaptchaV2_ValidateToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RecaptchaV2{
				config: RecaptchaV2Config{
					SecretKey: "test-secret",
				},
			}
			_, err := r.ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecaptchaV2.ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
