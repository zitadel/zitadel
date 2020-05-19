package login

import (
	"context"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/gorilla/mux"
	"golang.org/x/text/language"
)

type Config struct {
	Port                string
	StaticDir           string
	OidcAuthCallbackURL string
	CitadelURL          string
	LanguageCookieName  string
	DefaultLanguage     language.Tag
}

type Login struct {
	endpoint string
	router   *mux.Router
	//renderer            *Renderer
	//parser              *form.Parser
	//service             *service.ExternalService
	//citadelURL          string
	//oidcAuthCallbackURL string
	//userAgentHandler    *auth.UserAgentHandler
}

func Start(ctx context.Context, config Config, systemDefaults sd.SystemDefaults) {
}
