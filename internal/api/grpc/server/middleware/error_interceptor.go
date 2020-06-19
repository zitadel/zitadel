package middleware

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"

	"google.golang.org/grpc"

	grpc_util "github.com/caos/zitadel/internal/api/grpc"
)

func ErrorHandler(defaultLanguage language.Tag) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//TODO: Add correct namespace
	dir, err := fs.NewWithNamespace("login")
	logging.Log("ERROR-7usEW").OnError(err).Panic("unable to get zitadel namespace")

	i18n, err := i18n.NewTranslator(dir, i18n.TranslatorConfig{DefaultLanguage: defaultLanguage})
	if err != nil {
		logging.Log("ERROR-Sk8sf").OnError(err).Panic("unable to get i18n translator")
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		return resp, grpc_util.CaosToGRPCError(err, ctx, i18n)
	}
}
