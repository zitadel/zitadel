package instrumentation

import (
	"context"
	"log/slog"
	"os"
	"slices"
	"sync/atomic"
	"time"

	"github.com/rs/xid"
	slogmulti "github.com/samber/slog-multi"
	slogctx "github.com/veqryn/slog-context"
	slogotel "github.com/veqryn/slog-context/otel"
	old_logging "github.com/zitadel/logging" //nolint:staticcheck
	"github.com/zitadel/sloggcp"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/sdk/log"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type LogFormat int

//go:generate enumer -type=LogFormat -trimprefix=LogFormat -text -linecomment -transform=snake
const (
	// Empty line comment sets empty string of unspecified value
	LogFormatUnspecified LogFormat = iota //
	LogFormatDisabled
	LogFormatText
	LogFormatJSON
	// JSON formatted logs compatible with Google Cloud Platform
	LogFormatGCP
	LogFormatGCPErrorReporting
)

func (f LogFormat) isDisabled() bool {
	return f == LogFormatUnspecified || f == LogFormatDisabled
}

type LogConfig struct {
	Level     slog.Level
	Streams   []Stream
	Mask      MaskConfig
	Format    LogFormat
	AddSource bool
	Errors    ErrorConfig
	Exporter  ExporterConfig
}

func (c *LogConfig) SetLegacyConfig(lc *old_logging.Config) {
	if lc == nil || !c.Format.isDisabled() {
		return
	}
	err := c.Level.UnmarshalText([]byte(lc.Level))
	if err != nil {
		c.Level = slog.LevelInfo
	}
	c.Format, err = LogFormatString(lc.Formatter.Format)
	if err != nil {
		c.Format = LogFormatText
	}
	c.AddSource = lc.AddSource
}

func (c LogConfig) replacer() replacer {
	var replacers []replacer
	if len(c.Mask.Keys) > 0 {
		replacers = append(replacers, c.Mask.replacer())
	}
	if c.Format == LogFormatGCP {
		replacers = append(replacers, sloggcp.ReplaceAttr)
	}
	if c.Format == LogFormatGCPErrorReporting {
		replacers = append(replacers, errReplacer())
	}
	return chainReplacers(replacers...)
}

type MaskConfig struct {
	Keys  []string
	Value string
}

// replacer returns a function that replaces the value of any attribute whose key matches
// one of the configured keys with the configured value.
// If a group key matches, all attribute values in that group are replaced.
// Attribute structure is preserved, e.g. "password" is replaced with "masked" but still appears as "password" in the logs.
func (c MaskConfig) replacer() replacer {
	keys := slices.Clone(c.Keys)
	slices.Sort(keys)
	return func(groups []string, a slog.Attr) slog.Attr {
		// mask all entries in a matching group, e.g. "data.*"
		for _, g := range groups {
			if _, found := slices.BinarySearch(keys, g); found {
				a.Value = slog.StringValue(c.Value)
				return a
			}
		}
		if _, found := slices.BinarySearch(keys, a.Key); found {
			a.Value = slog.StringValue(c.Value)
		}
		return a
	}
}

type ErrorConfig struct {
	ReportLocation bool
	StackTrace     bool
}

var addSourceEnabled atomic.Bool

func IsAddSourceEnabled() bool {
	return addSourceEnabled.Load()
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
		AddSource:   cfg.AddSource,
		Level:       cfg.Level,
		ReplaceAttr: cfg.replacer(),
	}
	addSourceEnabled.Store(cfg.AddSource)

	var stdErrHandler slog.Handler
	switch cfg.Format {
	case LogFormatUnspecified:
		return // uses default slog logger
	case LogFormatDisabled:
		stdErrHandler = slog.DiscardHandler
	case LogFormatText:
		stdErrHandler = slog.NewTextHandler(os.Stderr, options)
	case LogFormatJSON:
		stdErrHandler = slog.NewJSONHandler(os.Stderr, options)
	case LogFormatGCP:
		stdErrHandler = slog.NewJSONHandler(os.Stderr, options)
	case LogFormatGCPErrorReporting:
		zerrors.GCPErrorReportingEnabled(true)
		stdErrHandler = sloggcp.NewErrorReportingHandler(os.Stderr, options)
	}

	EnableStreams(cfg.Streams...)
	zerrors.EnableReportLocation(cfg.Errors.ReportLocation)
	zerrors.EnableStackTrace(cfg.Errors.StackTrace)

	stdErrHandler = slogctx.NewHandler(
		stdErrHandler,
		&slogctx.HandlerOptions{
			Prependers: []slogctx.AttrExtractor{
				instanceExtractor,
				requestIDExtractor,
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

// Instance is a minimal interface for logging the instance ID.
type Instance interface {
	InstanceID() string
}

// SetInstance adds the instance to the context for logging.
func SetInstance(ctx context.Context, instance Instance) context.Context {
	return context.WithValue(ctx, ctxKeyInstance, instance)
}

// SetRequest generates a new [xid.ID] based on the passed request timestamp
// and adds it to the context.
func SetRequestID(ctx context.Context, ts time.Time) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, xid.NewWithTime(ts))
}

type ctxKeyType int

const (
	ctxKeyRequestID ctxKeyType = iota
	ctxKeyInstance
)

// instanceExtractor sets the instance ID from [Instance] to a log entry.
func instanceExtractor(ctx context.Context, _ time.Time, _ slog.Level, _ string) []slog.Attr {
	if instance, ok := ctx.Value(ctxKeyInstance).(Instance); ok {
		return []slog.Attr{
			slog.String("instance", instance.InstanceID()),
		}
	}
	return nil
}

// requestIDExtractor sets the request XID to a log entry.
func requestIDExtractor(ctx context.Context, _ time.Time, _ slog.Level, _ string) []slog.Attr {
	if r, ok := ctx.Value(ctxKeyRequestID).(xid.ID); ok {
		return []slog.Attr{
			slog.String("request_id", r.String()),
		}
	}
	return nil
}

type replacer func(groups []string, a slog.Attr) slog.Attr

func chainReplacers(replacers ...replacer) replacer {
	return func(groups []string, a slog.Attr) slog.Attr {
		for _, r := range replacers {
			a = r(groups, a)
		}
		return a
	}
}

func errReplacer() replacer {
	return func(groups []string, a slog.Attr) slog.Attr {
		if len(groups) == 0 && a.Key == "err" {
			a.Key = sloggcp.ErrorKey
		}
		return a
	}
}
