package middleware

import (
	"context"
	"errors"

	"github.com/rakyll/statik/fs"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	caos_errs "github.com/zitadel/zitadel/v2/internal/errors"
	"github.com/zitadel/zitadel/v2/internal/i18n"
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

func newZitadelTranslator(defaultLanguage language.Tag) (*i18n.Translator, error) {
	return translatorFromNamespace("zitadel", defaultLanguage)
}

func translatorFromNamespace(namespace string, defaultLanguage language.Tag) (*i18n.Translator, error) {
	dir, err := fs.NewWithNamespace(namespace)
	logging.WithFields("namespace", namespace).OnError(err).Panic("unable to get namespace")

	return i18n.NewTranslator(dir, defaultLanguage, "")
}
