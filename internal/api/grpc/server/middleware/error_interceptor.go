package middleware

import (
	"context"
	"github.com/caos/zitadel/internal/api/grpc/errors"

	"golang.org/x/text/language"
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/i18n"
	_ "github.com/caos/zitadel/internal/statik"
)

func ErrorHandler(defaultLanguage language.Tag) grpc.UnaryServerInterceptor {
	translator := newZitadelTranslator(defaultLanguage)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return toGRPCError(ctx, req, handler, translator)
	}
}

func toGRPCError(ctx context.Context, req interface{}, handler grpc.UnaryHandler, translator *i18n.Translator) (interface{}, error) {
	resp, err := handler(ctx, req)
	return resp, errors.CaosToGRPCError(ctx, err, translator)
}
