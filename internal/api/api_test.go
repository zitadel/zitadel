package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/ui/login"
)

func TestRegisterDefaultRoutes(t *testing.T) {
	t.Parallel()

	router := mux.NewRouter()
	registerDefaultRoutes(router)

	tests := []struct {
		name         string
		method       string
		path         string
		wantStatus   int
		wantLocation string
	}{
		{
			name:       "favicon returns not found",
			method:     http.MethodGet,
			path:       "/favicon.ico",
			wantStatus: http.StatusNotFound,
		},
		{
			name:         "root redirects to login",
			method:       http.MethodGet,
			path:         "/",
			wantStatus:   http.StatusFound,
			wantLocation: login.HandlerPrefix,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)
			assert.Equal(t, tt.wantStatus, res.Code)
			if tt.wantLocation != "" {
				assert.Equal(t, tt.wantLocation, res.Header().Get("Location"))
			}
		})
	}
}

func TestIsGRPCRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		method      string
		contentType string
		want        bool
	}{
		{
			name:        "grpc content type",
			method:      http.MethodPost,
			contentType: "application/grpc",
			want:        true,
		},
		{
			name:        "grpc proto content type",
			method:      http.MethodPost,
			contentType: "application/grpc+proto",
			want:        true,
		},
		{
			name:        "grpc json content type with charset",
			method:      http.MethodPost,
			contentType: "application/grpc+json; charset=utf-8",
			want:        true,
		},
		{
			name:        "non grpc content type",
			method:      http.MethodPost,
			contentType: "application/json",
			want:        false,
		},
		{
			name:        "non post grpc request",
			method:      http.MethodGet,
			contentType: "application/grpc",
			want:        false,
		},
		{
			name:        "missing content type",
			method:      http.MethodPost,
			contentType: "",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/zitadel.org.v2.OrganizationService/ListOrganizations", nil)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			assert.Equal(t, tt.want, isGRPCRequest(req))
		})
	}
}
