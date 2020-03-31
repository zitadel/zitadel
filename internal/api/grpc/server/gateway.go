package server

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"github.com/caos/logging"

	client_middleware "github.com/caos/zitadel/internal/api/grpc/client/middleware"
	http_util "github.com/caos/zitadel/internal/api/http"
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
)

const (
	defaultGatewayPort = "8080"
	mimeWildcard       = "*/*"
)

var (
	DefaultJSONMarshaler = &runtime.JSONPb{OrigName: false, EmitDefaults: false}

	DefaultServeMuxOptions = []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(DefaultJSONMarshaler.ContentType(), DefaultJSONMarshaler),
		runtime.WithMarshalerOption(mimeWildcard, DefaultJSONMarshaler),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, DefaultJSONMarshaler),
		runtime.WithIncomingHeaderMatcher(runtime.DefaultHeaderMatcher),
		runtime.WithOutgoingHeaderMatcher(runtime.DefaultHeaderMatcher),
	}
)

type Gateway interface {
	GRPCEndpoint() string
	GatewayPort() string
	Gateway() GatewayFunc
}

type GatewayFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

type gatewayCustomServeMuxOptions interface {
	GatewayServeMuxOptions() []runtime.ServeMuxOption
}
type grpcGatewayCustomInterceptor interface {
	GatewayHTTPInterceptor(http.Handler) http.Handler
}

type gatewayCustomCallOptions interface {
	GatewayCallOptions() []grpc.DialOption
}

func StartGateway(ctx context.Context, g Gateway) {
	mux := createMux(ctx, g)
	serveGateway(ctx, mux, gatewayPort(g.GatewayPort()), g)
}

func createMux(ctx context.Context, g Gateway) *runtime.ServeMux {
	muxOptions := DefaultServeMuxOptions
	if customOpts, ok := g.(gatewayCustomServeMuxOptions); ok {
		muxOptions = customOpts.GatewayServeMuxOptions()
	}
	mux := runtime.NewServeMux(muxOptions...)

	opts := []grpc.DialOption{grpc.WithInsecure()}
	opts = append(opts, client_middleware.DefaultTracingStatsClient())

	if customOpts, ok := g.(gatewayCustomCallOptions); ok {
		opts = append(opts, customOpts.GatewayCallOptions()...)
	}
	err := g.Gateway()(ctx, mux, g.GRPCEndpoint(), opts)
	logging.Log("SERVE-7B7G0E").OnError(err).Panic("failed to create mux for grpc gateway")

	return mux
}

func addInterceptors(handler http.Handler, g Gateway) http.Handler {
	handler = http_mw.DefaultTraceHandler(handler)
	if interceptor, ok := g.(grpcGatewayCustomInterceptor); ok {
		handler = interceptor.GatewayHTTPInterceptor(handler)
	}
	return http_mw.CORSInterceptorOpts(http_mw.DefaultCORSOptions, handler)
}

func serveGateway(ctx context.Context, handler http.Handler, port string, g Gateway) {
	server := &http.Server{
		Handler: addInterceptors(handler, g),
	}

	listener := http_util.CreateListener(port)

	go func() {
		<-ctx.Done()
		err := server.Shutdown(ctx)
		logging.Log("SERVE-m7kBlq").OnError(err).Warn("error during graceful shutdown of grpc gateway")
	}()

	go func() {
		err := server.Serve(listener)
		logging.Log("SERVE-tBHR60").OnError(err).Panic("grpc gateway serve failed")
	}()
	logging.LogWithFields("SERVE-KHh0Cb", "port", port).Info("grpc gateway is listening")
}

func gatewayPort(port string) string {
	if port == "" {
		return defaultGatewayPort
	}
	return port
}
