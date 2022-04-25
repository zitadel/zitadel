package middleware

import (
	"context"

	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/api/authz"

	_ "github.com/caos/zitadel/internal/statik"
)

func TranslationHandler() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		translator := newZitadelTranslator(authz.GetInstance(ctx).DefaultLanguage())
		if loc, ok := resp.(localizers); ok && resp != nil {
			translateFields(ctx, loc, translator)
		}
		if err != nil {
			err = translateError(ctx, err, translator)
		}
		return resp, err
	}
}
