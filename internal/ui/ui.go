package ui

import (
	"context"
	"net/http"

	sentryhttp "github.com/getsentry/sentry-go/http"

	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/ui/console"
	"github.com/caos/zitadel/internal/ui/login"
)

const (
	uiname = "ui"
)

type Config struct {
	Port    string
	Login   login.Config
	Console console.Config
}

type UI struct {
	port string
	mux  *http.ServeMux
}

func Create(config Config) *UI {
	return &UI{
		port: config.Port,
		mux:  http.NewServeMux(),
	}
}

func (u *UI) RegisterHandler(prefix string, handler http.Handler) {
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	http_util.RegisterHandler(u.mux, prefix, sentryHandler.Handle(handler))
}

func (u *UI) Start(ctx context.Context) {
	http_util.Serve(ctx, u.mux, u.port, uiname)
}
