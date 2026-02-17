package instrumentation

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	sdk_log "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

func newLoggerProvider(ctx context.Context, cfg ExporterConfig, resource *resource.Resource) (_ *sdk_log.LoggerProvider, err error) {
	var exporter sdk_log.Exporter
	switch cfg.Type {
	case ExporterTypeUnspecified, ExporterTypeNone:
		// no exporter
	case ExporterTypeStdOut, ExporterTypeStdErr:
		exporter, err = logStdOutExporter(cfg)
	case ExporterTypeGRPC:
		exporter, err = logGrpcExporter(ctx, cfg)
	case ExporterTypeHTTP:
		exporter, err = logHttpExporter(ctx, cfg)
	case ExporterTypeGoogle, ExporterTypePrometheus:
		fallthrough // google, prometheus are not supported for logs
	default:
		err = errExporterType(cfg.Type, "logger")
	}
	if err != nil {
		return nil, fmt.Errorf("log exporter: %w", err)
	}
	opts := []sdk_log.LoggerProviderOption{
		sdk_log.WithResource(resource),
	}
	if exporter != nil {
		opts = append(opts, sdk_log.WithProcessor(
			sdk_log.NewBatchProcessor(
				exporter,
				sdk_log.WithExportInterval(cfg.BatchDuration),
			),
		))
	}
	loggerProvider := sdk_log.NewLoggerProvider(opts...)
	return loggerProvider, nil
}

func logStdOutExporter(cfg ExporterConfig) (sdk_log.Exporter, error) {
	options := []stdoutlog.Option{
		stdoutlog.WithPrettyPrint(),
	}
	if cfg.Type == ExporterTypeStdErr {
		options = append(options, stdoutlog.WithWriter(os.Stderr))
	}
	exporter, err := stdoutlog.New(options...)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}

func logGrpcExporter(ctx context.Context, cfg ExporterConfig) (sdk_log.Exporter, error) {
	var grpcOpts []otlploggrpc.Option
	if cfg.Endpoint != "" {
		grpcOpts = append(grpcOpts, otlploggrpc.WithEndpoint(cfg.Endpoint))
	}
	if cfg.Insecure {
		grpcOpts = append(grpcOpts, otlploggrpc.WithInsecure())
	}

	exporter, err := otlploggrpc.New(ctx, grpcOpts...)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}

func logHttpExporter(ctx context.Context, cfg ExporterConfig) (sdk_log.Exporter, error) {
	var httpOpts []otlploghttp.Option
	if cfg.Endpoint != "" {
		httpOpts = append(httpOpts, otlploghttp.WithEndpoint(cfg.Endpoint))
	}
	if cfg.Insecure {
		httpOpts = append(httpOpts, otlploghttp.WithInsecure())
	}

	exporter, err := otlploghttp.New(ctx, httpOpts...)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}
