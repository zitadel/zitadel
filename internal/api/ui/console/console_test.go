package console

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/record"
)

func TestStart_environmentJSON(t *testing.T) {
	tests := []struct {
		name         string
		requestHost  string
		instanceHost string
		publicHost   string
		protocol     string
		wantAPI      string
	}{
		{
			name:         "instance host",
			requestHost:  "zitadel.example.com",
			instanceHost: "zitadel.example.com",
			protocol:     "https",
			wantAPI:      "https://zitadel.example.com",
		},
		{
			name:         "public host differs from request host",
			requestHost:  "zitadel.internal:8080",
			instanceHost: "zitadel.example.com",
			publicHost:   "custom.example.com",
			protocol:     "https",
			wantAPI:      "https://custom.example.com",
		},
		{
			name:         "forwarded instance host differs from request host",
			requestHost:  "zitadel.internal:8080",
			instanceHost: "custom.example.com",
			protocol:     "http",
			wantAPI:      "http://custom.example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passthrough := func(next http.Handler) http.Handler { return next }
			issuer := func(r *http.Request) string { return http_util.DomainContext(r.Context()).Origin() }
			accessInterceptor := middleware.NewAccessInterceptor(&logstore.Service[*record.AccessLog]{}, nil, &middleware.AccessConfig{})
			handler, err := Start(Config{}, issuer, passthrough, passthrough, accessInterceptor, "")
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, envRequestPath, nil)
			req.Host = tt.requestHost
			req = req.WithContext(http_util.WithDomainContext(req.Context(), http_util.NewDomainCtx(tt.instanceHost, tt.publicHost, tt.protocol)))
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

			var env struct {
				API    string `json:"api"`
				Issuer string `json:"issuer"`
			}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &env))
			assert.Equal(t, tt.wantAPI, env.API)
			assert.Equal(t, env.Issuer, env.API)
		})
	}
}
