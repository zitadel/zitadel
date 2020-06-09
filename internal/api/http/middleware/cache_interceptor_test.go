package middleware

import (
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
		name   string
		fields fields
		want   string
	}{
		{
			"no-store",
			fields{
				NoStore: true,
			},
			"no-store",
		},
		{
			"private and max-age",
			fields{
				Cacheability: CacheabilityPrivate,
				MaxAge:       1 * time.Hour,
				SharedMaxAge: 1 * time.Hour,
			},
			"private, max-age=3600",
		},
		{
			"public, no-cache, proy-revalidate",
			fields{
				Cacheability: CacheabilityPublic,
				NoCache:      true,
				Revalidation: RevalidationProxy,
			},
			"public, no-cache, proxy-revalidate",
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
			got := recorder.Result().Header.Get("cache-control")
			assert.Equal(t, tt.want, got)
		})
	}
}
