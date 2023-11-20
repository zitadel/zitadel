package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/i18n"
	"testing"
)

func TestMain(m *testing.M) {
	i18n.LoadMockedSupportedLanguages(domain.StringsToLanguages([]string{"en", "de", "es", "fr", "it", "ru"}))
	m.Run()
}
