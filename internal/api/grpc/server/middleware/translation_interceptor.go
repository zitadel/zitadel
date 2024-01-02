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
			translator, translatorError := getTranslator(ctx)
			if translatorError != nil {
				return resp, err
			}
			translateFields(ctx, loc, translator)
		}
		if err != nil {
			translator, translatorError := getTranslator(ctx)
			if translatorError != nil {
				return resp, err
			}
			err = translateError(ctx, err, translator)
		}
		return resp, err
	}
}

func getTranslator(ctx context.Context) (*i18n.Translator, error) {
	translator, err := i18n.NewZitadelTranslator(authz.GetInstance(ctx).DefaultLanguage())
	if err != nil {
		logging.New().WithError(err).Error("could not load translator")
	}
	return translator, err
}
