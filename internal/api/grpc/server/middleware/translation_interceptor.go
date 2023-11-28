package middleware

import (
	"context"

	"github.com/zitadel/logging"
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/i18n"
	_ "github.com/zitadel/zitadel/internal/statik"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func TranslationHandler() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		ctx, span := tracing.NewSpan(ctx)
		defer func() { span.EndWithError(err) }()

		if loc, ok := resp.(localizers); ok && resp != nil {
			translator, translatorError := newZitadelTranslator(authz.GetInstance(ctx).DefaultLanguage(), i18n.SupportedLanguages())
			if translatorError != nil {
				logging.New().WithError(translatorError).Error("could not load translator")
				return resp, err
			}
			translateFields(ctx, loc, translator)
		}
		if err != nil {
			translator, translatorError := newZitadelTranslator(authz.GetInstance(ctx).DefaultLanguage(), i18n.SupportedLanguages())
			if translatorError != nil {
				logging.New().WithError(translatorError).Error("could not load translator")
				return resp, err
			}
			err = translateError(ctx, err, translator)
		}
		return resp, err
	}
}
