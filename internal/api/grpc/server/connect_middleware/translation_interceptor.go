package connect_middleware

import (
	"context"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/i18n"
	_ "github.com/zitadel/zitadel/internal/statik"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func TranslationHandler() connect.UnaryInterceptorFunc {

	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := handler(ctx, req)
			ctx, span := tracing.NewSpan(ctx)
			defer func() { span.EndWithError(err) }()

			if err != nil {
				translator := getTranslator(ctx)
				return resp, translateError(ctx, err, translator)
			}
			if loc, ok := resp.Any().(localizers); ok {
				translator := getTranslator(ctx)
				translateFields(ctx, loc, translator)
			}
			return resp, nil
		}
	}
}

func getTranslator(ctx context.Context) *i18n.Translator {
	return i18n.NewZitadelTranslator(authz.GetInstance(ctx).DefaultLanguage())
}
