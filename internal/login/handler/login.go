package handler

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/form"
	"github.com/gorilla/mux"
	"golang.org/x/text/language"
	"net"
	"net/http"
)

type Login struct {
	endpoint string
	router   *mux.Router
	renderer            *Renderer
	parser              *form.Parser
	authRepo 			*eventsourcing.EsRepository
	zitadelURL          string
	oidcAuthCallbackURL string
	//userAgentHandler    *auth.UserAgentHandler
}

type Config struct {
	Port                  string
	StaticDir             string
	OidcAuthCallbackURL   string
	ZitadelURL            string
	LanguageCookieName    string
	DefaultLanguage       language.Tag
	//UserAgentCookieConfig *auth.UserAgentCookieConfig

}

func StartLogin(ctx context.Context, config Config, authRepo *eventsourcing.EsRepository) (err error) {
	login := &Login{
		endpoint: config.Port,
		oidcAuthCallbackURL: config.OidcAuthCallbackURL,
		zitadelURL: config.ZitadelURL,
		authRepo: authRepo,
	}
	login.router = CreateRouter(login, config.StaticDir)
	login.Listen(ctx)
	return err
}

func (login *Login) Listen(ctx context.Context) {
	if login.endpoint == "" {
		login.endpoint = ":80"
	} else {
		login.endpoint = ":" + login.endpoint
	}

	defer logging.LogWithFields("APP-xUZof", "port", login.endpoint).Info("html is listening")
	httpListener, err := net.Listen("tcp", login.endpoint)
	logging.Log("CONFI-W5q2O").OnError(err).Panic("unable to start listener")

	httpServer := &http.Server{
		Handler: login.router,
	}

	go func() {
		<-ctx.Done()
		if err = httpServer.Shutdown(ctx); err != nil {
			logging.Log("APP-mJKTv").WithError(err)
		}
	}()

	go func() {
		err := httpServer.Serve(httpListener)
		logging.Log("APP-oSklt").OnError(err).Panic("unable to start listener")
	}()
}
