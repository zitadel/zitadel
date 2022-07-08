package middleware

import (
	"context"

	"github.com/zitadel/logging"
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"

	_ "github.com/zitadel/zitadel/internal/statik"
)

func TranslationHandler() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		translator, translatorError := newZitadelTranslator(authz.GetInstance(ctx).DefaultLanguage())
		if translatorError != nil {
			logging.New().WithError(translatorError).Error("could not load translator")
			return resp, err
		}
		if loc, ok := resp.(localizers); ok && resp != nil {
			translateFields(ctx, loc, translator)
		}
		if err != nil {
			err = translateError(ctx, err, translator)
		}
		return resp, err
	}
}
