package login

import (
	"context"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
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

func Start(ctx context.Context, config Config, systemDefaults sd.SystemDefaults) {
}
