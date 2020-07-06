package middleware

import (
	"context"

	"github.com/caos/zitadel/internal/i18n"
	"golang.org/x/text/language"

	"google.golang.org/grpc"

	_ "github.com/caos/zitadel/internal/statik"
)

func TranslationHandler(defaultLanguage language.Tag) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	translator := newZitadelTranslator(defaultLanguage)

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
