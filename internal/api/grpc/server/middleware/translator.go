package middleware

import (
	"context"
	"errors"

	"github.com/caos/logging"
	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/i18n"
)

type localizers interface {
	Localizers() []Localizer
}
type Localizer interface {
	LocalizationKey() string
	SetLocalizedMessage(string)
}

func translateFields(ctx context.Context, object localizers, translator *i18n.Translator) {
	if translator == nil || object == nil {
		return
	}
	for _, field := range object.Localizers() {
		field.SetLocalizedMessage(translator.LocalizeFromCtx(ctx, field.LocalizationKey(), nil))
	}
}

func translateError(ctx context.Context, err error, translator *i18n.Translator) error {
	if translator == nil || err == nil {
		return err
	}
	caosErr := new(caos_errs.CaosError)
	if errors.As(err, &caosErr) {
		caosErr.SetMessage(translator.LocalizeFromCtx(ctx, caosErr.GetMessage(), nil))
	}
	return err
}

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
