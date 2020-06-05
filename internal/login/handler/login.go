package handler

import (
	"context"
	"net"
	"net/http"

	"github.com/caos/logging"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/form"

	_ "github.com/caos/zitadel/internal/login/statik"
)

type Login struct {
	endpoint            string
	router              *mux.Router
	renderer            *Renderer
	parser              *form.Parser
	authRepo            *eventsourcing.EsRepository
	zitadelURL          string
	oidcAuthCallbackURL string
}

type Config struct {
	Port                string
	OidcAuthCallbackURL string
	ZitadelURL          string
	LanguageCookieName  string
	DefaultLanguage     language.Tag
}

const (
	login = "LOGIN"
)

func StartLogin(ctx context.Context, config Config, authRepo *eventsourcing.EsRepository) {
	login := &Login{
		endpoint:            config.Port,
		oidcAuthCallbackURL: config.OidcAuthCallbackURL,
		zitadelURL:          config.ZitadelURL,
		authRepo:            authRepo,
	}
	statikFS, err := fs.NewWithNamespace("login")
	logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start listener")

	login.router = CreateRouter(login, statikFS)
	login.renderer = CreateRenderer(statikFS, config.LanguageCookieName, config.DefaultLanguage)
	login.parser = form.NewParser()
	login.Listen(ctx)
}

func (l *Login) Listen(ctx context.Context) {
	if l.endpoint == "" {
		l.endpoint = ":80"
	} else {
		l.endpoint = ":" + l.endpoint
	}

	defer logging.LogWithFields("APP-xUZof", "port", l.endpoint).Info("html is listening")
	httpListener, err := net.Listen("tcp", l.endpoint)
	logging.Log("CONFI-W5q2O").OnError(err).Panic("unable to start listener")

	httpServer := &http.Server{
		Handler: l.router,
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

func setContext(ctx context.Context, resourceOwner string) context.Context {
	data := auth.CtxData{
		UserID: login,
		OrgID:  resourceOwner,
	}
	return auth.SetCtxData(ctx, data)
}
