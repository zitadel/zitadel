package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	zhttp "github.com/zitadel/zitadel/internal/api/http"
)

func TestContentTypeMiddleware(t *testing.T) {
	tests := []struct {
		name              string
		contentTypeHeader string
		acceptHeader      string
		wantErr           bool
	}{
		{
			name:              "valid",
			contentTypeHeader: "application/scim+json",
			acceptHeader:      "application/scim+json",
			wantErr:           false,
		},
		{
			name:              "invalid content type",
			contentTypeHeader: "application/octet-stream",
			acceptHeader:      "application/json",
			wantErr:           true,
		},
		{
			name:              "invalid accept",
			contentTypeHeader: "application/json",
			acceptHeader:      "application/octet-stream",
			wantErr:           true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.acceptHeader != "" {
				req.Header.Set(zhttp.Accept, tt.acceptHeader)
			}

			if tt.contentTypeHeader != "" {
				req.Header.Set(zhttp.ContentType, tt.contentTypeHeader)
			}

			err := ContentTypeMiddleware(func(w http.ResponseWriter, r *http.Request) error {
				return nil
			})(httptest.NewRecorder(), req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentTypeMiddleware() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateContentType(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		want        bool
	}{
		{
			name:        "empty",
			contentType: "",
			want:        true,
		},
		{
			name:        "json",
			contentType: "application/json",
			want:        true,
		},
		{
			name:        "scim",
			contentType: "application/scim+json",
			want:        true,
		},
		{
			name:        "json utf-8",
			contentType: "application/json; charset=utf-8",
			want:        true,
		},
		{
			name:        "scim utf-8",
			contentType: "application/scim+json; charset=utf-8",
			want:        true,
		},
		{
			name:        "unknown content type",
			contentType: "application/octet-stream",
			want:        false,
		},
		{
			name:        "unknown charset",
			contentType: "application/scim+json; charset=utf-16",
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateContentType(tt.contentType); got != tt.want {
				t.Errorf("validateContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
