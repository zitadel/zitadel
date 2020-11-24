package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"github.com/caos/logging"

	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	client_middleware "github.com/caos/zitadel/internal/api/grpc/client/middleware"
	http_util "github.com/caos/zitadel/internal/api/http"
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

const (
	defaultGatewayPort = "8080"
	mimeWildcard       = "*/*"
)

var (
	DefaultJSONMarshaler = &runtime.JSONPb{OrigName: false, EmitDefaults: false}

	DefaultServeMuxOptions = func(customHeaders ...string) []runtime.ServeMuxOption {
		return []runtime.ServeMuxOption{
			runtime.WithMarshalerOption(DefaultJSONMarshaler.ContentType(), DefaultJSONMarshaler),
			runtime.WithMarshalerOption(mimeWildcard, DefaultJSONMarshaler),
			runtime.WithMarshalerOption(runtime.MIMEWildcard, DefaultJSONMarshaler),
			runtime.WithIncomingHeaderMatcher(DefaultHeaderMatcher(customHeaders...)),
			runtime.WithOutgoingHeaderMatcher(runtime.DefaultHeaderMatcher),
		}
	}

	DefaultHeaderMatcher = func(customHeaders ...string) runtime.HeaderMatcherFunc {
		return func(header string) (string, bool) {
			for _, customHeader := range customHeaders {
				if strings.HasPrefix(strings.ToLower(header), customHeader) {
					return header, true
				}
			}
			return runtime.DefaultHeaderMatcher(header)
		}
	}
)

type Gateway interface {
	RegisterGateway() GatewayFunc
	GatewayPathPrefix() string
}

type GatewayFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

//optional extending interfaces of Gateway below

type gatewayCustomServeMuxOptions interface {
	GatewayServeMuxOptions() []runtime.ServeMuxOption
}

type grpcGatewayCustomInterceptor interface {
	GatewayHTTPInterceptor(http.Handler) http.Handler
}

type gatewayCustomCallOptions interface {
	GatewayCallOptions() []grpc.DialOption
}

type GatewayHandler struct {
	mux           *http.ServeMux
	serverPort    string
	gatewayPort   string
	customHeaders []string
}

func CreateGatewayHandler(config grpc_util.Config) *GatewayHandler {
	return &GatewayHandler{
		mux:           http.NewServeMux(),
		serverPort:    config.ServerPort,
		gatewayPort:   config.GatewayPort,
		customHeaders: config.CustomHeaders,
	}
}

//RegisterGateway registers a handler (Gateway interface) on defined port
//Gateway interface may be extended with optional implementation of interfaces (gatewayCustomServeMuxOptions, ...)
func (g *GatewayHandler) RegisterGateway(ctx context.Context, gateway Gateway) {
	handler := createGateway(ctx, gateway, g.serverPort, g.customHeaders...)
	prefix := gateway.GatewayPathPrefix()
	g.RegisterHandler(prefix, handler)
}

func (g *GatewayHandler) RegisterHandler(prefix string, handler http.Handler) {
	http_util.RegisterHandler(g.mux, prefix, handler)
}

func (g *GatewayHandler) Serve(ctx context.Context) {
	http_util.Serve(ctx, g.mux, g.gatewayPort, "api")
}

func createGateway(ctx context.Context, g Gateway, port string, customHeaders ...string) http.Handler {
	mux := createMux(g, customHeaders...)
	opts := createDialOptions(g)
	err := g.RegisterGateway()(ctx, mux, http_util.Endpoint(port), opts)
	logging.Log("SERVE-7B7G0E").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Panic("failed to register grpc gateway")
	return addInterceptors(mux, g)
}

func createMux(g Gateway, customHeaders ...string) *runtime.ServeMux {
	muxOptions := DefaultServeMuxOptions(customHeaders...)
	if customOpts, ok := g.(gatewayCustomServeMuxOptions); ok {
		muxOptions = customOpts.GatewayServeMuxOptions()
	}
	return runtime.NewServeMux(muxOptions...)
}

func createDialOptions(g Gateway) []grpc.DialOption {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(client_middleware.DefaultTracingClient()),
	}

	if customOpts, ok := g.(gatewayCustomCallOptions); ok {
		opts = append(opts, customOpts.GatewayCallOptions()...)
	}
	return opts
}

func addInterceptors(handler http.Handler, g Gateway) http.Handler {
	handler = http_mw.DefaultMetricsyHandler(handler)
	handler = http_mw.DefaultTelemetryHandler(handler)
	handler = http_mw.NoCacheInterceptor(handler)
	if interceptor, ok := g.(grpcGatewayCustomInterceptor); ok {
		handler = interceptor.GatewayHTTPInterceptor(handler)
	}
	return http_mw.CORSInterceptor(handler)
}

func gatewayPort(port string) string {
	if port == "" {
		return defaultGatewayPort
	}
	return port
}
