package command

import (
	"testing"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/i18n"
)

var (
	SupportedLanguages   = []language.Tag{language.English, language.German}
	OnlyAllowedLanguages = []language.Tag{language.English}
	AllowedLanguage      = language.English
	DisallowedLanguage   = language.German
	UnsupportedLanguage  = language.Spanish
)

func TestMain(m *testing.M) {
	i18n.SupportLanguages(SupportedLanguages...)
	m.Run()
}
