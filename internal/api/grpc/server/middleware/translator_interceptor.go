package middleware

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"

	"google.golang.org/grpc"

	_ "github.com/caos/zitadel/internal/statik"
)

func TranslationHandler(defaultLanguage language.Tag) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	dir, err := fs.NewWithNamespace("zitadel")
	logging.Log("MIDDL-76QHE").OnError(err).Panic("unable to get zitadel namespace")

	translator, err := i18n.NewTranslator(dir, i18n.TranslatorConfig{DefaultLanguage: defaultLanguage})
	if err != nil {
		logging.Log("MIDDL-ELbF7").OnError(err).Panic("unable to get i18n translator")
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if loc, ok := resp.(localizers); ok {
			TranslateFields(ctx, loc, translator)
		}
		return resp, err
	}
}

type localizers interface {
	Localizers() []Localizer
}
type Localizer interface {
	LocalizationKey() string
	SetLocalizedMessage(string)
}

func TranslateFields(ctx context.Context, object localizers, translator *i18n.Translator) {
	if translator == nil || object == nil {
		return
	}
	for _, field := range object.Localizers() {
		field.SetLocalizedMessage(translator.LocalizeFromCtx(ctx, field.LocalizationKey(), nil))
	}
}
