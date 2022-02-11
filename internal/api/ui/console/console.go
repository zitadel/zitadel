package console

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/caos/logging"

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

const (
	envRequestPath    = "/assets/environment.json"
	consoleDefaultDir = "./console/"
	HandlerPrefix     = "/ui/console"
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

func Start(config Config, domain, url, issuer, clientID string) (http.Handler, error) {
	environmentJSON, err := createEnvironmentJSON(url, issuer, clientID)
	logging.Log("CONSO-tMAsY").OnError(err).Fatal("unable to marshal env")

	consoleDir := consoleDefaultDir
	if config.ConsoleOverwriteDir != "" {
		consoleDir = config.ConsoleOverwriteDir
	}
	consoleHTTPDir := http.Dir(consoleDir)

	cache := assetsCacheInterceptorIgnoreManifest(
		config.ShortCache.MaxAge,
		config.ShortCache.SharedMaxAge,
		config.LongCache.MaxAge,
		config.LongCache.SharedMaxAge,
	)
	security := middleware.SecurityHeaders(csp(domain), nil)

	handler := &http.ServeMux{}
	handler.Handle("/", cache(security(http.FileServer(&spaHandler{consoleHTTPDir}))))
	handler.Handle(envRequestPath, cache(security(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(environmentJSON)
		logging.Log("CONSOLE-sdet2").OnError(err).Error("error serving environment.json")
	}))))
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

func createEnvironmentJSON(url, issuer, clientID string) ([]byte, error) {
	environment := struct {
		AuthServiceUrl         string `json:"authServiceUrl,omitempty"`
		MgmtServiceUrl         string `json:"mgmtServiceUrl,omitempty"`
		AdminServiceUrl        string `json:"adminServiceUrl,omitempty"`
		SubscriptionServiceUrl string `json:"subscriptionServiceUrl,omitempty"`
		AssetServiceUrl        string `json:"assetServiceUrl,omitempty"`
		Issuer                 string `json:"issuer,omitempty"`
		ClientID               string `json:"clientid,omitempty"`
	}{
		AuthServiceUrl:         url,
		MgmtServiceUrl:         url,
		AdminServiceUrl:        url,
		SubscriptionServiceUrl: url,
		AssetServiceUrl:        url,
		Issuer:                 issuer,
		ClientID:               clientID,
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
