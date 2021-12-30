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
	"github.com/gorilla/mux"
)

type Config struct {
	ConsoleOverwriteDir string
	ShortCache          middleware.CacheConfig
	LongCache           middleware.CacheConfig
	CSPDomain           string
	Environment         Environment
}

type Environment struct {
	AuthServiceUrl         string `json:"authServiceUrl,omitempty"`
	MgmtServiceUrl         string `json:"mgmtServiceUrl,omitempty"`
	AdminServiceUrl        string `json:"adminServiceUrl,omitempty"`
	SubscriptionServiceUrl string `json:"subscriptionServiceUrl,omitempty"`
	AssetServiceUrl        string `json:"assetServiceUrl,omitempty"`
	Issuer                 string `json:"issuer,omitempty"`
	Clientid               string `json:"clientid,omitempty"`
}

type spaHandler struct {
	fileSystem http.FileSystem
}

const (
	envRequestPath    = "/assets/environment.json"
	consoleDefaultDir = "./console/"
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

func New(uiRouter *mux.Router, config Config) {
	consoleDir := consoleDefaultDir
	if config.ConsoleOverwriteDir != "" {
		consoleDir = config.ConsoleOverwriteDir
	}
	consoleHTTPDir := http.Dir(consoleDir)
	cache := AssetsCacheInterceptorIgnoreManifest(
		config.ShortCache.MaxAge.Duration,
		config.ShortCache.SharedMaxAge.Duration,
		config.LongCache.MaxAge.Duration,
		config.LongCache.SharedMaxAge.Duration,
	)
	security := middleware.SecurityHeaders(csp(config.CSPDomain), nil)
	consoleRouter := uiRouter.PathPrefix("/console").Subrouter()
	consoleRouter.PathPrefix(envRequestPath).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		environmentJSON, err := json.Marshal(config.Environment)
		logging.Log("CONSO-tMAsY").OnError(err).Error("unable to marshal env")
		w.Write(environmentJSON)
	}))
	consoleRouter.NewRoute().Handler(http.StripPrefix("/ui/console", cache(security(http.FileServer(&spaHandler{consoleHTTPDir})))))
}

func csp(zitadelDomain string) *middleware.CSP {
	if !strings.HasPrefix(zitadelDomain, "*.") {
		zitadelDomain = "*." + zitadelDomain
	}
	csp := middleware.DefaultSCP
	csp.StyleSrc = csp.StyleSrc.AddInline().AddHost("fonts.googleapis.com").AddHost("maxst.icons8.com") //TODO: host it
	csp.FontSrc = csp.FontSrc.AddHost("fonts.gstatic.com").AddHost("maxst.icons8.com")                  //TODO: host it
	csp.ScriptSrc = csp.ScriptSrc.AddEval()
	csp.ConnectSrc = csp.ConnectSrc.AddHost(zitadelDomain).
		AddHost("fonts.googleapis.com").
		AddHost("fonts.gstatic.com").
		AddHost("maxst.icons8.com") //TODO: host it
	csp.ImgSrc = csp.ImgSrc.AddHost(zitadelDomain).AddScheme("blob")
	return &csp
}

func AssetsCacheInterceptorIgnoreManifest(shortMaxAge, shortSharedMaxAge, longMaxAge, longSharedMaxAge time.Duration) func(http.Handler) http.Handler {
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
