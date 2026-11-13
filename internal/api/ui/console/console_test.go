package console

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/record"
)

func TestStart_environmentJSON(t *testing.T) {
	tests := []struct {
		name           string
		externalSecure bool
		requestHost    string
		headers        map[string]string
		wantAPI        string
	}{
		{
			name:           "request host",
			externalSecure: true,
			requestHost:    "zitadel.example.com",
			wantAPI:        "https://zitadel.example.com",
		},
		{
			name:        "x-forwarded-host differs from request host",
			requestHost: "zitadel.internal:8080",
			headers: map[string]string{
				"X-Forwarded-Host": "custom.example.com",
			},
			wantAPI: "http://custom.example.com",
		},
		{
			name:        "forwarded header differs from request host",
			requestHost: "zitadel.internal:8080",
			headers: map[string]string{
				"Forwarded": "host=custom.example.com;proto=https",
			},
			wantAPI: "https://custom.example.com",
		},
		{
			name:           "public host differs from instance host",
			externalSecure: true,
			requestHost:    "zitadel.internal:8080",
			headers: map[string]string{
				"x-zitadel-instance-host": "zitadel.example.com",
				"x-zitadel-public-host":   "custom.example.com",
			},
			wantAPI: "https://custom.example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passthrough := func(next http.Handler) http.Handler { return next }
			issuer := func(r *http.Request) string { return "" }
			accessInterceptor := middleware.NewAccessInterceptor(&logstore.Service[*record.AccessLog]{}, nil, &middleware.AccessConfig{})
			handler, err := Start(Config{}, issuer, passthrough, passthrough, accessInterceptor, "")
			require.NoError(t, err)
			handler = middleware.WithOrigin(tt.externalSecure, "host", ":authority", []string{"x-zitadel-instance-host"}, []string{"x-zitadel-public-host"})(handler)

			req := httptest.NewRequest(http.MethodGet, envRequestPath, nil)
			req.Host = tt.requestHost
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

			var env struct {
				API string `json:"api"`
			}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &env))
			assert.Equal(t, tt.wantAPI, env.API)
		})
	}
}
