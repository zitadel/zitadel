package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_serializeHeaders(t *testing.T) {
	type fields struct {
		Cacheability Cacheability
		NoCache      bool
		NoStore      bool
		MaxAge       time.Duration
		SharedMaxAge time.Duration
		NoTransform  bool
		Revalidation Revalidation
	}
	tests := []struct {
		name        string
		fields      fields
		wantControl string
		wantExpires string
		wantPragma  string
	}{
		{
			"no-store",
			fields{
				NoStore: true,
			},
			"no-store",
			time.Now().UTC().Add(-1 * time.Hour).Format(http.TimeFormat),
			"no-cache",
		},
		{
			"private and max-age",
			fields{
				Cacheability: CacheabilityPrivate,
				MaxAge:       1 * time.Hour,
				SharedMaxAge: 1 * time.Hour,
			},
			"private, max-age=3600",
			time.Now().UTC().Add(1 * time.Hour).Format(http.TimeFormat),
			"",
		},
		{
			"public, no-cache, proxy-revalidate",
			fields{
				Cacheability: CacheabilityPublic,
				NoCache:      true,
				Revalidation: RevalidationProxy,
			},
			"public, max-age=0, no-cache, proxy-revalidate",
			time.Now().UTC().Add(-1 * time.Hour).Format(http.TimeFormat),
			"no-cache",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c := &Cache{
				Cacheability: tt.fields.Cacheability,
				NoCache:      tt.fields.NoCache,
				NoStore:      tt.fields.NoStore,
				MaxAge:       tt.fields.MaxAge,
				SharedMaxAge: tt.fields.SharedMaxAge,
				NoTransform:  tt.fields.NoTransform,
				Revalidation: tt.fields.Revalidation,
			}
			c.serializeHeaders(recorder)
			cc := recorder.Result().Header.Get("cache-control")
			assert.Equal(t, tt.wantControl, cc)
			exp := recorder.Result().Header.Get("expires")
			assert.Equal(t, tt.wantExpires, exp)
			pragma := recorder.Result().Header.Get("pragma")
			assert.Equal(t, tt.wantPragma, pragma)
		})
	}
}
