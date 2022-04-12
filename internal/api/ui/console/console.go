package console

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"

	"github.com/caos/zitadel/internal/api/http/middleware"
)

type Config struct {
	ConsoleOverwriteDir string
	ShortCache          middleware.CacheConfig
	LongCache           middleware.CacheConfig
}

type spaHandler struct {
	fileSystem http.FileSystem
}

var (
	//go:embed static/*
	static embed.FS
)

const (
	envRequestPath = "/assets/environment.json"
	HandlerPrefix  = "/ui/console"
)

var (
	shortCacheFiles = []string{
		"/",
		"/index.html",
		"/manifest.webmanifest",
		"/ngsw.json",
		"/ngsw-worker.js",
		"/safety-worker.js",
		"/worker-basic.min.js",
	}
)

func (i *spaHandler) Open(name string) (http.File, error) {
	ret, err := i.fileSystem.Open(name)
	if !os.IsNotExist(err) || path.Ext(name) != "" {
		return ret, err
	}

	return i.fileSystem.Open("/index.html")
}

func Start(config Config, domain, url, issuer string, instanceHandler func(http.Handler) http.Handler) (http.Handler, error) {
	fSys, err := fs.Sub(static, "static")
	if err != nil {
		return nil, err
	}
	cache := assetsCacheInterceptorIgnoreManifest(
		config.ShortCache.MaxAge,
		config.ShortCache.SharedMaxAge,
		config.LongCache.MaxAge,
		config.LongCache.SharedMaxAge,
	)
	security := middleware.SecurityHeaders(csp(domain), nil)

	handler := &http.ServeMux{}
	handler.Handle("/", cache(security(http.FileServer(&spaHandler{http.FS(fSys)}))))
	handler.Handle(envRequestPath, instanceHandler(cache(security(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		instance := authz.GetInstance(r.Context())
		if instance.InstanceID() == "" {
			http.Error(w, "empty instanceID", http.StatusInternalServerError)
			return
		}
		environmentJSON, err := createEnvironmentJSON(url, issuer, instance.ConsoleClientID())
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to marshal env for console: %v", err), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(environmentJSON)
		logging.OnError(err).Error("error serving environment.json")
	})))))
	return handler, nil
}

func csp(zitadelDomain string) *middleware.CSP {
	if !strings.HasPrefix(zitadelDomain, "*.") {
		zitadelDomain = "*." + zitadelDomain
	}
	csp := middleware.DefaultSCP
	csp.StyleSrc = csp.StyleSrc.AddInline()
	csp.ScriptSrc = csp.ScriptSrc.AddEval()
	csp.ConnectSrc = csp.ConnectSrc.AddHost(zitadelDomain)
	csp.ImgSrc = csp.ImgSrc.AddHost(zitadelDomain).AddScheme("blob")
	return &csp
}

func createEnvironmentJSON(api, issuer, clientID string) ([]byte, error) {
	environment := struct {
		API      string `json:"api,omitempty"`
		Issuer   string `json:"issuer,omitempty"`
		ClientID string `json:"clientid,omitempty"`
	}{
		API:      api,
		Issuer:   issuer,
		ClientID: clientID,
	}
	return json.Marshal(environment)
}

func assetsCacheInterceptorIgnoreManifest(shortMaxAge, shortSharedMaxAge, longMaxAge, longSharedMaxAge time.Duration) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, file := range shortCacheFiles {
				if r.URL.Path == file {
					middleware.AssetsCacheInterceptor(shortMaxAge, shortSharedMaxAge, handler).ServeHTTP(w, r)
					return
				}
			}
			middleware.AssetsCacheInterceptor(longMaxAge, longSharedMaxAge, handler).ServeHTTP(w, r)
			return
		})
	}
}
