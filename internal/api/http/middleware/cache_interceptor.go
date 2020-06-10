package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/caos/zitadel/internal/api"
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
	CacheabilityPublic               = "public"
	CacheabilityPrivate              = "private"
)

type Revalidation string

const (
	RevalidationNotSet Revalidation = ""
	RevalidationMust                = "must-revalidate"
	RevalidationProxy               = "proxy-revalidate"
)

var (
	NeverCacheOptions = &Cache{
		NoStore: true,
	}
	AssetOptions = &Cache{
		Cacheability: CacheabilityPublic,
		MaxAge:       7 * 24 * time.Hour,
	}
)

func DefaultCacheInterceptor(pattern string) (func(http.Handler) http.Handler, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if regex.MatchString(r.URL.Path) {
				AssetsCacheInterceptor(handler).ServeHTTP(w, r)
				return
			}
			NoCacheInterceptor(handler).ServeHTTP(w, r)
		})
	}, nil
}

func NoCacheInterceptor(h http.Handler) http.Handler {
	return CacheInterceptorOpts(h, NeverCacheOptions)
}

func AssetsCacheInterceptor(h http.Handler) http.Handler {
	return CacheInterceptorOpts(h, AssetOptions)
}

func CacheInterceptorOpts(h http.Handler, cache *Cache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cache.serializeHeaders(w)
		h.ServeHTTP(w, req)
	})
}

func (c *Cache) serializeHeaders(w http.ResponseWriter) {
	control := make([]string, 0, 6)
	pragma := false

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
		control = append(control, fmt.Sprintf("no-cache"))
		pragma = true
	}

	if c.NoStore {
		control = append(control, fmt.Sprintf("no-store"))
	}
	if c.NoTransform {
		control = append(control, fmt.Sprintf("no-transform"))
	}

	if c.Revalidation != RevalidationNotSet {
		control = append(control, string(c.Revalidation))
	}

	w.Header().Set(api.CacheControl, strings.Join(control, ", "))
	w.Header().Set(api.Expires, expires)
	if pragma {
		w.Header().Set(api.Pragma, "no-cache")
	}
}
