package middleware

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"
)

func newZitadelTranslator(defaultLanguage language.Tag) *i18n.Translator {
	return translatorFromNamespace("zitadel", defaultLanguage)
}

func translatorFromNamespace(namespace string, defaultLanguage language.Tag) *i18n.Translator {
	dir, err := fs.NewWithNamespace(namespace)
	logging.LogWithFields("ERROR-7usEW", "namespace", namespace).OnError(err).Panic("unable to get namespace")

	translator, err := i18n.NewTranslator(dir, i18n.TranslatorConfig{DefaultLanguage: defaultLanguage})
	logging.Log("ERROR-Sk8sf").OnError(err).Panic("unable to get i18n translator")

	return translator
}
