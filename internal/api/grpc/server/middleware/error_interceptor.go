package middleware

import (
	"context"

	"golang.org/x/text/language"

	"google.golang.org/grpc"

	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	_ "github.com/caos/zitadel/internal/statik"
)

func ErrorHandler(defaultLanguage language.Tag) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	translator := newZitadelTranslator(defaultLanguage)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		return resp, grpc_util.CaosToGRPCError(err, ctx, translator)
	}
}
