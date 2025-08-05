package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
)

type Cache struct {
	Cacheability Cacheability
	NoCache      bool
	NoStore      bool
	MaxAge       time.Duration
	SharedMaxAge time.Duration
	NoTransform  bool
	Revalidation Revalidation
}

type Cacheability string

const (
	CacheabilityNotSet  Cacheability = ""
	CacheabilityPublic  Cacheability = "public"
	CacheabilityPrivate Cacheability = "private"
)

type Revalidation string

const (
	RevalidationNotSet Revalidation = ""
	RevalidationMust   Revalidation = "must-revalidate"
	RevalidationProxy  Revalidation = "proxy-revalidate"
)

type CacheConfig struct {
	MaxAge       time.Duration
	SharedMaxAge time.Duration
}

var (
	NeverCacheOptions = &Cache{
		NoStore: true,
	}
	AssetOptions = func(maxAge, SharedMaxAge time.Duration) *Cache {
		return &Cache{
			Cacheability: CacheabilityPublic,
			MaxAge:       maxAge,
			SharedMaxAge: SharedMaxAge,
		}
	}
)

func NoCacheInterceptor() *cacheInterceptor {
	return CacheInterceptorOpts(NeverCacheOptions)
}

func AssetsCacheInterceptor(maxAge, sharedMaxAge time.Duration) *cacheInterceptor {
	return CacheInterceptorOpts(AssetOptions(maxAge, sharedMaxAge))
}

func CacheInterceptorOpts(cache *Cache) *cacheInterceptor {
	return &cacheInterceptor{
		cache: cache,
	}
}

type cacheInterceptor struct {
	cache *Cache
}

func (c *cacheInterceptor) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(&cachingResponseWriter{
			ResponseWriter: w,
			Cache:          c.cache,
		}, r)
	})
}

func (c *cacheInterceptor) HandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(&cachingResponseWriter{
			ResponseWriter: w,
			Cache:          c.cache,
		}, r)
	}
}

type cachingResponseWriter struct {
	http.ResponseWriter
	*Cache
}

func (w *cachingResponseWriter) WriteHeader(code int) {
	if code >= 400 {
		NeverCacheOptions.serializeHeaders(w.ResponseWriter)
		w.ResponseWriter.WriteHeader(code)
		return
	}
	w.Cache.serializeHeaders(w.ResponseWriter)
	w.ResponseWriter.WriteHeader(code)
}

func (c *Cache) serializeHeaders(w http.ResponseWriter) {
	control := make([]string, 0, 6)
	pragma := false

	// Do not overwrite cache-control header if set by business logic.
	if w.Header().Get(http_utils.CacheControl) != "" {
		return
	}

	if c.Cacheability != CacheabilityNotSet {
		control = append(control, string(c.Cacheability))
		control = append(control, fmt.Sprintf("max-age=%v", c.MaxAge.Seconds()))
		if c.SharedMaxAge != c.MaxAge {
			control = append(control, fmt.Sprintf("s-maxage=%v", c.SharedMaxAge.Seconds()))
		}
	}
	maxAge := c.MaxAge
	if maxAge == 0 {
		maxAge = -time.Hour
	}
	expires := time.Now().UTC().Add(maxAge).Format(http.TimeFormat)

	if c.NoCache {
		control = append(control, "no-cache")
		pragma = true
	}

	if c.NoStore {
		control = append(control, "no-store")
		pragma = true
	}
	if c.NoTransform {
		control = append(control, "no-transform")
	}

	if c.Revalidation != RevalidationNotSet {
		control = append(control, string(c.Revalidation))
	}

	w.Header().Set(http_utils.CacheControl, strings.Join(control, ", "))
	w.Header().Set(http_utils.Expires, expires)
	if pragma {
		w.Header().Set(http_utils.Pragma, "no-cache")
	}
}
