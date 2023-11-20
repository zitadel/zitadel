package command

import (
	"github.com/zitadel/zitadel/internal/i18n"
	"golang.org/x/text/language"
	"testing"
)

func TestMain(m *testing.M) {
	i18n.LoadMockedSupportedLanguages(language.English, language.German)
	m.Run()
}
