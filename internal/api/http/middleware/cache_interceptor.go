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
	MaxAge       time.Duration
	SharedMaxAge time.Duration
	//Cacheability Cacheability
	//Revalidation CacheRevalidation
	Public       bool
	Cacheability Cacheability
	NoStore      bool
	NoTransform  bool

	serialized map[string][]string
}
type Cacheability string

const (
	CacheabilityNotSet  Cacheability = ""
	CacheabilityNoCache              = "no-cache"
	CacheabilityNoStore              = "no-store"
)

var (
	NeverCacheOptions = &Cache{
		Cacheability: CacheabilityNoStore,
	}
	AssetOptions = &Cache{
		MaxAge: 365 * 24 * time.Hour,
		Public: true,
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
	control := make([]string, 2, 6)
	pragma := false
	maxAge := c.MaxAge

	control[0] = "private"
	if c.Public {
		control[0] = "public"
	}
	control[1] = fmt.Sprintf("max-age=%v", maxAge.Seconds())
	if maxAge == 0 {
		maxAge = -time.Hour
	}
	expires := time.Now().UTC().Add(maxAge).Format(http.TimeFormat)

	if c.Cacheability != CacheabilityNotSet {
		control = append(control, string(c.Cacheability))
		pragma = true
	}

	w.Header().Set(api.CacheControl, strings.Join(control, ","))
	w.Header().Set(api.Expires, expires)
	if pragma {
		w.Header().Set(api.Pragma, "no-cache")
	}
}
