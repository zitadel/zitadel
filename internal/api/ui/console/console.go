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

	"github.com/gorilla/mux"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/op"

	"github.com/zitadel/zitadel/cmd/build"
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

func LoginHintLink(origin, username string) string {
	return origin + HandlerPrefix + "?login_hint=" + username
}

func (i *spaHandler) Open(name string) (http.File, error) {
	ret, err := i.fileSystem.Open(name)
	if !os.IsNotExist(err) || path.Ext(name) != "" {
		return ret, err
	}

	f, err := i.fileSystem.Open("/index.html")
	if err != nil {
		return nil, err
	}
	return &file{File: f}, nil
}

// file wraps the http.File and fs.FileInfo interfaces
// to return the build.Date() as ModTime() of the file
type file struct {
	http.File
	fs.FileInfo
}

func (f *file) ModTime() time.Time {
	return build.Date()
}

func (f *file) Stat() (_ fs.FileInfo, err error) {
	f.FileInfo, err = f.File.Stat()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func Start(config Config, externalSecure bool, issuer op.IssuerFromRequest, callDurationInterceptor, instanceHandler, accessInterceptor func(http.Handler) http.Handler, customerPortal string) (http.Handler, error) {
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

	handler.Use(callDurationInterceptor, instanceHandler, security, accessInterceptor)
	handler.Handle(envRequestPath, middleware.TelemetryHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := http_util.BuildOrigin(r.Host, externalSecure)
		environmentJSON, err := createEnvironmentJSON(url, issuer(r), authz.GetInstance(r.Context()).ConsoleClientID(), customerPortal)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to marshal env for console: %v", err), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(environmentJSON)
		logging.OnError(err).Error("error serving environment.json")
	})))
	handler.SkipClean(true).PathPrefix("").Handler(cache(http.FileServer(&spaHandler{http.FS(fSys)})))
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

func createEnvironmentJSON(api, issuer, clientID, customerPortal string) ([]byte, error) {
	environment := struct {
		API            string `json:"api,omitempty"`
		Issuer         string `json:"issuer,omitempty"`
		ClientID       string `json:"clientid,omitempty"`
		CustomerPortal string `json:"customer_portal,omitempty"`
	}{
		API:            api,
		Issuer:         issuer,
		ClientID:       clientID,
		CustomerPortal: customerPortal,
	}
	return json.Marshal(environment)
}

func assetsCacheInterceptorIgnoreManifest(shortMaxAge, shortSharedMaxAge, longMaxAge, longSharedMaxAge time.Duration) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, file := range shortCacheFiles {
				if r.URL.Path == file || isIndexOrSubPath(r.URL.Path) {
					middleware.AssetsCacheInterceptor(shortMaxAge, shortSharedMaxAge).Handler(handler).ServeHTTP(w, r)
					return
				}
			}
			middleware.AssetsCacheInterceptor(longMaxAge, longSharedMaxAge).Handler(handler).ServeHTTP(w, r)
		})
	}
}

func isIndexOrSubPath(path string) bool {
	//files will have an extension
	return !strings.Contains(path, ".")
}
