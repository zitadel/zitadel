package instrumentation

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	slogmulti "github.com/samber/slog-multi"
	slogctx "github.com/veqryn/slog-context"
	slogotel "github.com/veqryn/slog-context/otel"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/sdk/log"

	http_util "github.com/zitadel/zitadel/internal/api/http"
)

type LogConfig struct {
	Level      slog.Level
	StdErr     bool
	JSONFormat bool
	AddSource  bool
	Exporter   ExporterConfig
}

type RequestDetails interface {
	*http.Request
}

// SetRequestDetails adds static details to each context aware log entry.
func SetRequestDetails[R RequestDetails](ctx context.Context, service string, details R) context.Context {
	return context.WithValue(ctx, ctxKey{}, requestDetails[R]{
		service: service,
		details: details,
	})
}

type ctxKey struct{}

type requestDetails[R RequestDetails] struct {
	service string
	details R
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

// domainExtractor sets the sanitized request hosts from [http_util.DomainCtx] to a log entry.
func domainExtractor(ctx context.Context, _ time.Time, _ slog.Level, _ string) []slog.Attr {
	return []slog.Attr{
		slog.Any("domain", http_util.DomainContext(ctx)),
	}
}

func requestExtractor(ctx context.Context, _ time.Time, _ slog.Level, _ string) []slog.Attr {
	switch details := ctx.Value(ctxKey{}).(type) {
	case *http.Request:
		return httpRequestToAttr(details)
	default:
		return nil
	}
}

func httpRequestToAttr(r *http.Request) []slog.Attr {
	return []slog.Attr{
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	}
}
