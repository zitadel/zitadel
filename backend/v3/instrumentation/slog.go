package instrumentation

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"connectrpc.com/connect"
	slogmulti "github.com/samber/slog-multi"
	slogctx "github.com/veqryn/slog-context"
	slogotel "github.com/veqryn/slog-context/otel"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"

	http_util "github.com/zitadel/zitadel/internal/api/http"
)

type LogConfig struct {
	Level      slog.Level
	StdErr     bool
	JSONFormat bool
	AddSource  bool
	Exporter   ExporterConfig
}

// setLogger configures the global slog logger.
// Logs are sent to StdErr and/or the [log.LoggerProvider].
//
// When present in the context, each line emitted contains:
//   - Trace ID
//   - Span ID
//   - Service name
//   - Sanitized request hosts from [http_util.DomainCtx]
//   - Request details, such as method and path.
func setLogger(provider *log.LoggerProvider, cfg LogConfig) {
	handlers := []slog.Handler{
		otelslog.NewHandler("zitadel", otelslog.WithLoggerProvider(provider)),
	}

	if cfg.StdErr {
		options := &slog.HandlerOptions{
			AddSource: cfg.AddSource,
			Level:     cfg.Level,
		}
		if cfg.JSONFormat {
			handlers = append(handlers, slog.NewJSONHandler(os.Stderr,
				options,
			))
		} else {
			handlers = append(handlers, slog.NewTextHandler(os.Stderr,
				options,
			))
		}
	}
	handler := slogctx.NewHandler(
		slogmulti.Fanout(handlers...),
		&slogctx.HandlerOptions{
			Prependers: []slogctx.AttrExtractor{
				domainExtractor,
				requestExtractor,
				slogotel.ExtractTraceSpanID,
				slogctx.ExtractPrepended,
			},
		},
	)
	slog.SetDefault(slog.New(handler))
}

// SetHttpRequestDetails adds static details to each context aware log entry.
func SetHttpRequestDetails(ctx context.Context, service string, request *http.Request) context.Context {
	return context.WithValue(ctx, ctxKey{}, httpRequest{
		service: service,
		request: request,
	})
}

// SetConnectRequestDetails adds static details to each context aware log entry.
func SetConnectRequestDetails(ctx context.Context, request connect.AnyRequest) context.Context {
	return context.WithValue(ctx, ctxKey{}, connectRequest{request})
}

func SetGrpcRequestDetails(ctx context.Context, info *grpc.UnaryServerInfo) context.Context {
	return context.WithValue(ctx, ctxKey{}, grpcRequest{info})
}

type ctxKey struct{}

type requestDetails interface {
	attrs() []slog.Attr
}

type httpRequest struct {
	service string
	request *http.Request
}

var _ requestDetails = httpRequest{}

func (r httpRequest) attrs() []slog.Attr {
	return []slog.Attr{
		slog.String("protocol", "http"),
		slog.String("service", r.service),
		slog.String("method", r.request.Method),
		slog.String("path", r.request.URL.Path),
	}
}

type connectRequest struct {
	connect.AnyRequest
}

var _ requestDetails = connectRequest{}

func (r connectRequest) attrs() []slog.Attr {
	spec := r.Spec()
	proc := strings.Split(spec.Procedure, "/")
	var service string
	if len(proc) > 0 {
		service = proc[0]
	}

	return []slog.Attr{
		slog.String("protocol", "connect"),
		slog.String("service", service),
		slog.String("method", r.HTTPMethod()),
		slog.String("procedure", spec.Procedure),
	}
}

type grpcRequest struct {
	info *grpc.UnaryServerInfo
}

var _ requestDetails = grpcRequest{}

func (r grpcRequest) attrs() []slog.Attr {
	proc := strings.Split(r.info.FullMethod, "/")
	var service string
	if len(proc) > 0 {
		service = proc[0]
	}

	return []slog.Attr{
		slog.String("protocol", "grpc"),
		slog.String("service", service),
		slog.String("procedure", r.info.FullMethod),
	}
}

// domainExtractor sets the sanitized request hosts from [http_util.DomainCtx] to a log entry.
func domainExtractor(ctx context.Context, _ time.Time, _ slog.Level, _ string) []slog.Attr {
	return []slog.Attr{
		slog.Any("domain", http_util.DomainContext(ctx)),
	}
}

func requestExtractor(ctx context.Context, _ time.Time, _ slog.Level, _ string) []slog.Attr {
	if r, ok := ctx.Value(ctxKey{}).(requestDetails); ok {
		return r.attrs()
	}
	return nil
}
