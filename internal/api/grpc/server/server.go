package server

import (
	"context"

	"github.com/caos/logging"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"golang.org/x/text/language"
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
	"github.com/caos/zitadel/internal/api/http"
)

const (
	defaultGrpcPort = "80"
)

type Server interface {
	Gateway
	RegisterServer(*grpc.Server)
	AppName() string
	MethodPrefix() string
	AuthMethods() authz.MethodMapping
}

func CreateServer(verifier *authz.TokenVerifier, authConfig authz.Config, lang language.Tag) *grpc.Server {
	return grpc.NewServer(
		middleware.DefaultTracingStatsServer(),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middleware.ErrorHandler(),
				middleware.AuthorizationInterceptor(verifier, authConfig),
				middleware.TranslationHandler(lang),
				middleware.ValidationHandler(),
				middleware.ServiceHandler(),
			),
		),
	)
}

func Serve(ctx context.Context, server *grpc.Server, port string) {
	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	go func() {
		listener := http.CreateListener(port)
		err := server.Serve(listener)
		logging.Log("SERVE-Ga3e94").OnError(err).Panic("grpc server serve failed")
	}()
	logging.LogWithFields("SERVE-bZ44QM", "port", port).Info("grpc server is listening")
}

func grpcPort(port string) string {
	if port == "" {
		return defaultGrpcPort
	}
	return port
}
