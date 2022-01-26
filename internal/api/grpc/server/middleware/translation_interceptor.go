package middleware

import (
	"context"

	"github.com/caos/zitadel/internal/query"
	"google.golang.org/grpc"

	_ "github.com/caos/zitadel/internal/statik"
)

func TranslationHandler(query *query.Queries) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	translator := newZitadelTranslator(query.GetDefaultLanguage(context.TODO()))

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if loc, ok := resp.(localizers); ok && resp != nil {
			translateFields(ctx, loc, translator)
		}
		if err != nil {
			err = translateError(ctx, err, translator)
		}
		return resp, err
	}
}
