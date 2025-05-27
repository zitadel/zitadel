package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zitadel/logging"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	client_middleware "github.com/zitadel/zitadel/internal/api/grpc/client/middleware"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

const (
	mimeWildcard = "*/*"
	UnknownPath  = "UNKNOWN_PATH"
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

	httpErrorHandler = runtime.RoutingErrorHandlerFunc(
		func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, httpStatus int) {
			if httpStatus != http.StatusMethodNotAllowed {
				runtime.DefaultRoutingErrorHandler(ctx, mux, marshaler, w, r, httpStatus)
				return
			}

			// Use HTTPStatusError to customize the DefaultHTTPErrorHandler status code
			err := &runtime.HTTPStatusError{
				HTTPStatus: httpStatus,
				Err:        status.Error(codes.Unimplemented, http.StatusText(httpStatus)),
			}

			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
		},
	)

	// we need the errorHandler to set the request URI pattern in case of an error
	errorHandler = runtime.ErrorHandlerFunc(
		func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
			setRequestURIPattern(ctx)
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
		})

	serveMuxOptions = func(hostHeaders []string) []runtime.ServeMuxOption {
		return []runtime.ServeMuxOption{
			runtime.WithMarshalerOption(jsonMarshaler.ContentType(nil), jsonMarshaler),
			runtime.WithMarshalerOption(mimeWildcard, jsonMarshaler),
			runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonMarshaler),
			runtime.WithIncomingHeaderMatcher(headerMatcher(hostHeaders)),
			runtime.WithOutgoingHeaderMatcher(runtime.DefaultHeaderMatcher),
			runtime.WithForwardResponseOption(responseForwarder),
			runtime.WithRoutingErrorHandler(httpErrorHandler),
			runtime.WithErrorHandler(errorHandler),
		}
	}

	headerMatcher = func(hostHeaders []string) runtime.HeaderMatcherFunc {
		customHeaders = slices.Compact(append(customHeaders, hostHeaders...))
		return func(header string) (string, bool) {
			for _, customHeader := range customHeaders {
				if strings.HasPrefix(strings.ToLower(header), customHeader) {
					return header, true
				}
			}
			return runtime.DefaultHeaderMatcher(header)
		}
	}

	responseForwarder = func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
		setRequestURIPattern(ctx)
		t, ok := resp.(CustomHTTPResponse)
		if ok {
			// TODO: find a way to return a location header if needed w.Header().Set("location", t.Location())
			w.WriteHeader(t.CustomHTTPCode())
		}
		return nil
	}
)

type Gateway struct {
	mux               *runtime.ServeMux
	connection        *grpc.ClientConn
	accessInterceptor *http_mw.AccessInterceptor
}

func (g *Gateway) Handler() http.Handler {
	return addInterceptors(g.mux, g.accessInterceptor)
}

type CustomHTTPResponse interface {
	CustomHTTPCode() int
}

type RegisterGatewayFunc func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

func CreateGatewayWithPrefix(
	ctx context.Context,
	g WithGatewayPrefix,
	port uint16,
	hostHeaders []string,
	accessInterceptor *http_mw.AccessInterceptor,
	tlsConfig *tls.Config,
) (http.Handler, string, error) {
	runtimeMux := runtime.NewServeMux(serveMuxOptions(hostHeaders)...)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(grpcCredentials(tlsConfig)),
		grpc.WithChainUnaryInterceptor(
			client_middleware.UnaryActivityClientInterceptor(),
		),
		grpc.WithStatsHandler(client_middleware.DefaultTracingClient()),
	}
	connection, err := dial(ctx, port, opts)
	if err != nil {
		return nil, "", err
	}
	err = g.RegisterGateway()(ctx, runtimeMux, connection)
	if err != nil {
		return nil, "", fmt.Errorf("failed to register grpc gateway: %w", err)
	}
	return addInterceptors(runtimeMux, accessInterceptor), g.GatewayPathPrefix(), nil
}

func CreateGateway(
	ctx context.Context,
	port uint16,
	hostHeaders []string,
	accessInterceptor *http_mw.AccessInterceptor,
	tlsConfig *tls.Config,
) (*Gateway, error) {
	connection, err := dial(ctx,
		port,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(grpcCredentials(tlsConfig)),
			grpc.WithChainUnaryInterceptor(
				client_middleware.UnaryActivityClientInterceptor(),
			),
			grpc.WithStatsHandler(client_middleware.DefaultTracingClient()),
		})
	if err != nil {
		return nil, err
	}
	runtimeMux := runtime.NewServeMux(append(serveMuxOptions(hostHeaders), runtime.WithHealthzEndpoint(healthpb.NewHealthClient(connection)))...)
	return &Gateway{
		mux:               runtimeMux,
		connection:        connection,
		accessInterceptor: accessInterceptor,
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
	accessInterceptor *http_mw.AccessInterceptor,
) http.Handler {
	handler = http_mw.CallDurationHandler(handler)
	handler = http_mw.CORSInterceptor(handler)
	handler = http_mw.RobotsTagHandler(handler)
	handler = http_mw.DefaultTelemetryHandler(handler)
	handler = http_mw.ActivityHandler(handler)
	// For some non-obvious reason, the exhaustedCookieInterceptor sends the SetCookie header
	// only if it follows the http_mw.DefaultTelemetryHandler
	handler = exhaustedCookieInterceptor(handler, accessInterceptor)
	handler = http_mw.MetricsHandler([]metrics.MetricType{
		metrics.MetricTypeTotalCount,
		metrics.MetricTypeStatusCode,
	}, http_utils.Probes...)(handler)
	return handler
}

func exhaustedCookieInterceptor(
	next http.Handler,
	accessInterceptor *http_mw.AccessInterceptor,
) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		next.ServeHTTP(&cookieResponseWriter{
			ResponseWriter:    writer,
			accessInterceptor: accessInterceptor,
			request:           request,
		}, request)
	})
}

type cookieResponseWriter struct {
	http.ResponseWriter
	accessInterceptor *http_mw.AccessInterceptor
	request           *http.Request
	headerWritten     bool
}

func (r *cookieResponseWriter) WriteHeader(status int) {
	if status >= 200 && status < 300 {
		r.accessInterceptor.DeleteExhaustedCookie(r.ResponseWriter)
	}
	if status == http.StatusTooManyRequests {
		r.accessInterceptor.SetExhaustedCookie(r.ResponseWriter, r.request)
	}
	r.headerWritten = true
	r.ResponseWriter.WriteHeader(status)
}

func (r *cookieResponseWriter) Write(bytes []byte) (int, error) {
	if !r.headerWritten {
		// If no header was written before the data, the status code is 200 and we can delete the cookie
		r.accessInterceptor.DeleteExhaustedCookie(r.ResponseWriter)
	}
	return r.ResponseWriter.Write(bytes)
}

func grpcCredentials(tlsConfig *tls.Config) credentials.TransportCredentials {
	creds := insecure.NewCredentials()
	if tlsConfig != nil {
		tlsConfigClone := tlsConfig.Clone()
		// We don't want to verify the certificate of the internal grpc server
		// That's up to the client who called the gRPC gateway
		tlsConfigClone.InsecureSkipVerify = true
		creds = credentials.NewTLS(tlsConfigClone)
	}
	return creds
}

func setRequestURIPattern(ctx context.Context) {
	pattern, ok := runtime.HTTPPathPattern(ctx)
	if !ok {
		// As all unmatched paths will be handled by the gateway, any request not matching a pattern,
		// means there's no route to the path.
		// To prevent high cardinality on metrics and tracing, we want to make sure we don't record
		// the actual path as name (it will still be recorded explicitly in the span http info).
		pattern = UnknownPath
	}
	span := trace.SpanFromContext(ctx)
	span.SetName(pattern)
	metrics.SetRequestURIPattern(ctx, pattern)
}
