package instrumentation

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/rs/xid"
	slogmulti "github.com/samber/slog-multi"
	slogctx "github.com/veqryn/slog-context"
	slogotel "github.com/veqryn/slog-context/otel"
	"github.com/zitadel/sloggcp"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type LogFormat int

//go:generate enumer -type=LogFormat -trimprefix=LogFormat -text -linecomment -transform=snake
const (
	// Empty line comment sets empty string of unspecified value
	LogFormatUndefined LogFormat = iota //
	LogFormatDisabled
	LogFormatText
	LogFormatJSON
	// JSON formatted logs compatible with Google Cloud Platform
	LogFormatGCP
	LogFormatGCPErrorReporting
)

type LogConfig struct {
	Level     slog.Level
	Streams   []Stream
	Format    LogFormat
	AddSource bool
	Errors    ErrorConfig
	Exporter  ExporterConfig
}

type ErrorConfig struct {
	ReportLocation bool
	StackTrace     bool
}

// setLogger configures the global slog logger.
// Logs are sent to [os.Stderr] and/or the [log.LoggerProvider].
//
// When present in the context, each line emitted contains:
//   - Trace ID
//   - Span ID
//   - Service name
//   - Sanitized request hosts from [http_util.DomainCtx]
//   - Request details, such as method and path.
func setLogger(provider *log.LoggerProvider, cfg LogConfig) {
	options := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     cfg.Level,
	}
	var stdErrHandler slog.Handler
	switch cfg.Format {
	case LogFormatUndefined:
		return // uses default slog logger
	case LogFormatDisabled:
		stdErrHandler = slog.DiscardHandler
	case LogFormatText:
		stdErrHandler = slog.NewTextHandler(os.Stderr, options)
	case LogFormatJSON:
		stdErrHandler = slog.NewJSONHandler(os.Stderr, options)
	case LogFormatGCP:
		options.ReplaceAttr = sloggcp.ReplaceAttr
		stdErrHandler = slog.NewJSONHandler(os.Stderr, options)
	case LogFormatGCPErrorReporting:
		zerrors.GCPErrorReportingEnabled(true)
		options.ReplaceAttr = replaceErrAttr
		stdErrHandler = sloggcp.NewErrorReportingHandler(os.Stderr, options)
	}

	zerrors.EnableReportLocation(cfg.Errors.ReportLocation)
	zerrors.EnableStackTrace(cfg.Errors.StackTrace)

	stdErrHandler = slogctx.NewHandler(
		stdErrHandler,
		&slogctx.HandlerOptions{
			Prependers: []slogctx.AttrExtractor{
				domainExtractor,
				requestExtractor,
				slogotel.ExtractTraceSpanID,
				slogctx.ExtractPrepended,
			},
		},
	)

	otelHandler := otelslog.NewHandler(
		Name,
		otelslog.WithLoggerProvider(provider),
	)

	logger := slog.New(slogmulti.Fanout(stdErrHandler, otelHandler))
	logger.Info("structured logger configured", "config_level", cfg.Level, "format", cfg.Format)
	slog.SetDefault(logger)
}

const (
	ProtocolHttp    = "http"
	ProtocolConnect = "connect"
	ProtocolGrpc    = "grpc"
)

// SetHttpRequestDetails adds static details to each context aware log entry.
func SetHttpRequestDetails(ctx context.Context, service string, request *http.Request) context.Context {
	now := time.Now()
	return context.WithValue(ctx, ctxKey{}, &requestDetails{
		protocol:    ProtocolHttp,
		service:     service,
		http_method: request.Method,
		path:        request.URL.Path,
		requestID:   xid.NewWithTime(now),
		start:       now,
	})
}

// SetConnectRequestDetails adds static details to each context aware log entry.
func SetConnectRequestDetails(ctx context.Context, request connect.AnyRequest) context.Context {
	now := time.Now()
	spec := request.Spec()
	return context.WithValue(ctx, ctxKey{}, &requestDetails{
		protocol:    ProtocolConnect,
		service:     serviceFromRPCMethod(spec.Procedure),
		http_method: request.HTTPMethod(),
		path:        spec.Procedure,
		requestID:   xid.NewWithTime(now),
		start:       now,
	})
}

func SetGrpcRequestDetails(ctx context.Context, info *grpc.UnaryServerInfo) context.Context {
	now := time.Now()
	return context.WithValue(ctx, ctxKey{}, &requestDetails{
		protocol:    ProtocolGrpc,
		service:     serviceFromRPCMethod(info.FullMethod),
		http_method: http.MethodPost, // gRPC always uses POST
		path:        info.FullMethod,
		requestID:   xid.NewWithTime(now),
		start:       now,
	})
}

type ctxKey struct{}

type requestDetails struct {
	protocol    string
	service     string
	http_method string
	path        string
	requestID   xid.ID
	start       time.Time
}

func (r *requestDetails) attrs() []slog.Attr {
	attrs := make([]slog.Attr, 0, 6)
	if r.protocol != "" {
		attrs = append(attrs, slog.String("protocol", r.protocol))
	}
	if r.service != "" {
		attrs = append(attrs, slog.String("service", r.service))
	}
	if r.http_method != "" {
		attrs = append(attrs, slog.String("http_method", r.http_method))
	}
	if r.path != "" {
		attrs = append(attrs, slog.String("path", r.path))
	}
	if !r.requestID.IsZero() {
		attrs = append(attrs, slog.String("request_id", r.requestID.String()))
	}
	if !r.start.IsZero() {
		attrs = append(attrs, slog.Duration("duration", time.Since(r.start)))
	}
	return attrs
}

func serviceFromRPCMethod(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return "unknown"
}

// domainExtractor sets the sanitized request hosts from [http_util.DomainCtx] to a log entry.
func domainExtractor(ctx context.Context, _ time.Time, _ slog.Level, _ string) []slog.Attr {
	return []slog.Attr{
		slog.Any("domain", http_util.DomainContext(ctx)),
	}
}

// requestExtractor sets the request details from [requestDetails] to a log entry.
func requestExtractor(ctx context.Context, _ time.Time, _ slog.Level, _ string) []slog.Attr {
	if r, ok := ctx.Value(ctxKey{}).(*requestDetails); ok {
		return r.attrs()
	}
	return nil
}

// replaceErrAttr renames the "err" attribute to the Google Cloud Platform compatible error key.
func replaceErrAttr(groups []string, a slog.Attr) slog.Attr {
	if len(groups) == 0 && a.Key == "err" {
		a.Key = sloggcp.ErrorKey
	}
	return a
}
