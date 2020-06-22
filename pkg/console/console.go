package console

import (
	"context"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rakyll/statik/fs"

	"github.com/caos/zitadel/internal/api/http/middleware"
	_ "github.com/caos/zitadel/pkg/console/statik"
)

type Config struct {
	Port            string
	EnvOverwriteDir string
	Cache           middleware.CacheConfig
	CSPDomain       string
}

type spaHandler struct {
	fileSystem http.FileSystem
}

const (
	envRequestPath = "/assets/environment.json"
	envDefaultDir  = "/console/"

	manifestFile = "/manifest.webmanifest"
)

func (i *spaHandler) Open(name string) (http.File, error) {
	ret, err := i.fileSystem.Open(name)
	if !os.IsNotExist(err) || path.Ext(name) != "" {
		return ret, err
	}

	return i.fileSystem.Open("/index.html")
}

func Start(ctx context.Context, config Config) error {
	statikFS, err := fs.NewWithNamespace("console")
	if err != nil {
		return err
	}
	envDir := envDefaultDir
	if config.EnvOverwriteDir != "" {
		envDir = config.EnvOverwriteDir
	}
	cache := AssetsCacheInterceptorIgnoreManifest(config.Cache.MaxAge.Duration, config.Cache.SharedMaxAge.Duration)
	security := middleware.SecurityHeaders(csp(config.CSPDomain), nil)
	http.Handle("/", cache(security(http.FileServer(&spaHandler{statikFS}))))
	http.Handle(envRequestPath, http.StripPrefix("/assets", http.FileServer(http.Dir(envDir))))
	return http.ListenAndServe(":"+config.Port, nil)
}

func csp(zitadelDomain string) *middleware.CSP {
	if !strings.HasPrefix(zitadelDomain, "*.") {
		zitadelDomain = "*." + zitadelDomain
	}
	csp := middleware.DefaultSCP
	csp.StyleSrc = csp.StyleSrc.AddInline().AddHost("fonts.googleapis.com").AddHost("maxst.icons8.com") //TODO: host it
	csp.FontSrc = csp.FontSrc.AddHost("fonts.gstatic.com").AddHost("maxst.icons8.com")                  //TODO: host it
	csp.ScriptSrc = csp.ScriptSrc.AddEval()
	csp.ConnectSrc = csp.ConnectSrc.AddHost(zitadelDomain)
	return &csp
}

func AssetsCacheInterceptorIgnoreManifest(maxAge, sharedMaxAge time.Duration) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == manifestFile {
				middleware.NoCacheInterceptor(handler).ServeHTTP(w, r)
				return
			}
			middleware.AssetsCacheInterceptor(maxAge, sharedMaxAge, handler).ServeHTTP(w, r)
		})
	}
}
