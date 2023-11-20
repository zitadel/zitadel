package command

import (
	"testing"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/i18n"
)

func TestMain(m *testing.M) {
	i18n.LoadMockedSupportedLanguages(language.English, language.German)
	m.Run()
}
