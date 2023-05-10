package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zitadel/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/encoding/protojson"

	client_middleware "github.com/zitadel/zitadel/internal/api/grpc/client/middleware"
	"github.com/zitadel/zitadel/internal/api/grpc/server/middleware"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/query"
)

const (
	mimeWildcard = "*/*"
)

var (
	customHeaders = []string{
		"x-zitadel-",
	}
	jsonMarshaler = &runtime.JSONPb{
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}

	serveMuxOptions = []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(jsonMarshaler.ContentType(nil), jsonMarshaler),
		runtime.WithMarshalerOption(mimeWildcard, jsonMarshaler),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonMarshaler),
		runtime.WithIncomingHeaderMatcher(headerMatcher),
		runtime.WithOutgoingHeaderMatcher(runtime.DefaultHeaderMatcher),
	}

	headerMatcher = runtime.HeaderMatcherFunc(
		func(header string) (string, bool) {
			for _, customHeader := range customHeaders {
				if strings.HasPrefix(strings.ToLower(header), customHeader) {
					return header, true
				}
			}
			return runtime.DefaultHeaderMatcher(header)
		},
	)
)

type Gateway struct {
	mux           *runtime.ServeMux
	http1HostName string
	connection    *grpc.ClientConn
	cookieHandler *http_utils.CookieHandler
	cookieConfig  *http_mw.AccessConfig
	queries       *query.Queries
}

func (g *Gateway) Handler() http.Handler {
	return addInterceptors(g.mux, g.http1HostName, g.cookieHandler, g.cookieConfig, g.queries)
}

type RegisterGatewayFunc func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

func CreateGatewayWithPrefix(
	ctx context.Context,
	g WithGatewayPrefix,
	port uint16,
	http1HostName string,
	cookieHandler *http_utils.CookieHandler,
	cookieConfig *http_mw.AccessConfig,
	queries *query.Queries,
) (http.Handler, string, error) {
	runtimeMux := runtime.NewServeMux(serveMuxOptions...)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(client_middleware.DefaultTracingClient()),
	}
	connection, err := dial(ctx, port, opts)
	if err != nil {
		return nil, "", err
	}
	err = g.RegisterGateway()(ctx, runtimeMux, connection)
	if err != nil {
		return nil, "", fmt.Errorf("failed to register grpc gateway: %w", err)
	}
	return addInterceptors(runtimeMux, http1HostName, cookieHandler, cookieConfig, queries), g.GatewayPathPrefix(), nil
}

func CreateGateway(ctx context.Context, port uint16, http1HostName string, cookieHandler *http_utils.CookieHandler, cookieConfig *http_mw.AccessConfig) (*Gateway, error) {
	connection, err := dial(ctx,
		port,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(client_middleware.DefaultTracingClient()),
		})
	if err != nil {
		return nil, err
	}
	runtimeMux := runtime.NewServeMux(append(serveMuxOptions, runtime.WithHealthzEndpoint(healthpb.NewHealthClient(connection)))...)
	return &Gateway{
		mux:           runtimeMux,
		http1HostName: http1HostName,
		connection:    connection,
		cookieHandler: cookieHandler,
		cookieConfig:  cookieConfig,
	}, nil
}

func RegisterGateway(ctx context.Context, gateway *Gateway, server Server) error {
	err := server.RegisterGateway()(ctx, gateway.mux, gateway.connection)
	if err != nil {
		return fmt.Errorf("failed to register grpc gateway: %w", err)
	}
	return nil
}

func dial(ctx context.Context, port uint16, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	endpoint := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				logging.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				logging.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()
	return conn, nil
}

func addInterceptors(
	handler http.Handler,
	http1HostName string,
	cookieHandler *http_utils.CookieHandler,
	cookieConfig *http_mw.AccessConfig,
	queries *query.Queries,
) http.Handler {
	handler = http_mw.CallDurationHandler(handler)
	handler = http1Host(handler, http1HostName)
	handler = http_mw.CORSInterceptor(handler)
	handler = http_mw.DefaultTelemetryHandler(handler)
	// For some non-obvious reason, the exhaustedCookieInterceptor sends the SetCookie header
	// only if it follows the http_mw.DefaultTelemetryHandler
	handler = exhaustedCookieInterceptor(handler, cookieHandler, cookieConfig, queries)
	handler = http_mw.DefaultMetricsHandler(handler)
	return handler
}

func http1Host(next http.Handler, http1HostName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, err := http_mw.HostFromRequest(r, http1HostName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		r.Header.Set(middleware.HTTP1Host, host)
		next.ServeHTTP(w, r)
	})
}

func exhaustedCookieInterceptor(
	next http.Handler,
	cookieHandler *http_utils.CookieHandler,
	cookieConfig *http_mw.AccessConfig,
	queries *query.Queries,
) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		next.ServeHTTP(&cookieResponseWriter{
			ResponseWriter: writer,
			cookieHandler:  cookieHandler,
			cookieConfig:   cookieConfig,
			request:        request,
			queries:        queries,
		}, request)
	})
}

type cookieResponseWriter struct {
	http.ResponseWriter
	cookieHandler *http_utils.CookieHandler
	cookieConfig  *http_mw.AccessConfig
	request       *http.Request
	queries       *query.Queries
}

func (r *cookieResponseWriter) WriteHeader(status int) {
	if status >= 200 && status < 300 {
		http_mw.DeleteExhaustedCookie(r.cookieHandler, r.ResponseWriter, r.request, r.cookieConfig)
	}
	if status == http.StatusTooManyRequests {
		http_mw.SetExhaustedCookie(r.cookieHandler, r.ResponseWriter, r.cookieConfig, r.request)
	}
	r.ResponseWriter.WriteHeader(status)
}
