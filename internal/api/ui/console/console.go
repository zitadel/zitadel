package console

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/op"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
)

type Config struct {
	ShortCache middleware.CacheConfig
	LongCache  middleware.CacheConfig
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

func Start(config Config, externalSecure bool, issuer op.IssuerFromRequest, instanceHandler func(http.Handler) http.Handler) (http.Handler, error) {
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
	security := middleware.SecurityHeaders(csp(), nil)

	handler := mux.NewRouter()
	handler.Use(cache, security)
	handler.Handle(envRequestPath, middleware.TelemetryHandler()(instanceHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		instance := authz.GetInstance(r.Context())
		if instance.InstanceID() == "" {
			http.Error(w, "empty instanceID", http.StatusInternalServerError)
			return
		}
		url := http_util.BuildOrigin(r.Host, externalSecure)
		environmentJSON, err := createEnvironmentJSON(url, issuer(r), instance.ConsoleClientID())
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to marshal env for console: %v", err), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(environmentJSON)
		logging.OnError(err).Error("error serving environment.json")
	}))))
	handler.SkipClean(true).PathPrefix("").Handler(http.FileServer(&spaHandler{http.FS(fSys)}))
	return handler, nil
}

func csp() *middleware.CSP {
	csp := middleware.DefaultSCP
	csp.StyleSrc = csp.StyleSrc.AddInline()
	csp.ScriptSrc = csp.ScriptSrc.AddEval()
	csp.ConnectSrc = csp.ConnectSrc.AddOwnHost()
	csp.ImgSrc = csp.ImgSrc.AddOwnHost().AddScheme("blob")
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
