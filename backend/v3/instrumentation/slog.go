package instrumentation

import (
	"log/slog"
	"os"

	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/sdk/log"
)

type LogConfig struct {
	Level      slog.Level
	StdErr     bool
	JSONFormat bool
	AddSource  bool
	Exporter   ExporterConfig
}

func setLogger(provider *log.LoggerProvider, cfg LogConfig) {
	options := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     cfg.Level,
	}
	handlers := []slog.Handler{
		otelslog.NewHandler("zitadel", otelslog.WithLoggerProvider(provider)),
	}

	if cfg.StdErr {
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

	logger := slog.New(slogmulti.Fanout(handlers...))
	slog.SetDefault(logger)
}
