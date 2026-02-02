// NextJS instrumentation entry point
// This file is automatically loaded by NextJS and must export a `register` function
//
// Exports all three OTEL signals:
// - Traces: via OTLP HTTP push (auto-configured from OTEL_EXPORTER_OTLP_ENDPOINT)
// - Metrics: via OTLP HTTP push AND Prometheus scraping endpoint
// - Logs: via OTLP HTTP push (automatic with Winston instrumentation)
//
// Standard OTEL environment variables (automatically handled by SDK):
// - OTEL_SERVICE_NAME: Service name for telemetry
// - OTEL_EXPORTER_OTLP_ENDPOINT: OTLP collector endpoint
// - OTEL_EXPORTER_OTLP_HEADERS: Headers for OTLP requests (key=value,key2=value2)
// - OTEL_EXPORTER_OTLP_PROTOCOL: Protocol (http/protobuf, http/json, grpc)
// - OTEL_RESOURCE_ATTRIBUTES: Additional resource attributes
//
// Custom environment variables:
// - OTEL_LOG_LEVEL: Diagnostic log level (debug, info, warn, error)
// - OTEL_METRICS_PORT: Port for Prometheus metrics scraping (default: 9464)
// - OTEL_METRICS_EXPORT_INTERVAL_MS: Metrics export interval for OTLP (default: 30000)
//
// GCP Cloud Run environment variables (auto-detected for resource attributes):
// - K_SERVICE, K_REVISION, K_CONFIGURATION
// - GOOGLE_CLOUD_PROJECT, GOOGLE_CLOUD_REGION

import type { LoggerProvider } from "@opentelemetry/sdk-logs";
import type { MetricReader } from "@opentelemetry/sdk-metrics";

let _loggerProvider: LoggerProvider | null = null;

/**
 * Detect GCP Cloud Run resource attributes from environment variables
 */
function getGcpResourceAttributes(): Record<string, string> {
  const attributes: Record<string, string> = {};

  const projectId = process.env.GOOGLE_CLOUD_PROJECT || process.env.GCLOUD_PROJECT;
  if (projectId) {
    attributes["cloud.provider"] = "gcp";
    attributes["cloud.account.id"] = projectId;
  }

  const serviceName = process.env.K_SERVICE;
  if (serviceName) {
    attributes["cloud.platform"] = "gcp_cloud_run";
    attributes["faas.name"] = serviceName;
  }

  const revision = process.env.K_REVISION;
  if (revision) {
    attributes["faas.version"] = revision;
    attributes["faas.instance"] = revision;
  }

  const configuration = process.env.K_CONFIGURATION;
  if (configuration) {
    attributes["cloud.resource_id"] = configuration;
  }

  const region = process.env.GOOGLE_CLOUD_REGION || process.env.CLOUD_RUN_REGION;
  if (region) {
    attributes["cloud.region"] = region;
  }

  return attributes;
}

/**
 * Initialize OpenTelemetry with:
 * - Traces: OTLP HTTP push (if OTEL_EXPORTER_OTLP_ENDPOINT is set)
 * - Metrics: OTLP HTTP push + Prometheus scraping
 * - Logs: OTLP HTTP push (via Winston instrumentation)
 *
 * OTLP exporters automatically read configuration from standard environment variables:
 * - OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_EXPORTER_OTLP_HEADERS, etc.
 */
export async function register(): Promise<void> {
  if (process.env.NEXT_RUNTIME !== "nodejs") {
    return;
  }

  // Dynamic imports to avoid webpack bundling issues with Node.js-only modules
  const { NodeTracerProvider } = await import("@opentelemetry/sdk-trace-node");
  const { BatchSpanProcessor } = await import("@opentelemetry/sdk-trace-base");
  const { OTLPTraceExporter } = await import("@opentelemetry/exporter-trace-otlp-http");
  const { MeterProvider, PeriodicExportingMetricReader } = await import("@opentelemetry/sdk-metrics");
  const { OTLPMetricExporter } = await import("@opentelemetry/exporter-metrics-otlp-http");
  const { PrometheusExporter } = await import(/* webpackIgnore: true */ "@opentelemetry/exporter-prometheus");
  const { LoggerProvider, BatchLogRecordProcessor } = await import("@opentelemetry/sdk-logs");
  const { OTLPLogExporter } = await import("@opentelemetry/exporter-logs-otlp-http");
  const { resourceFromAttributes } = await import("@opentelemetry/resources");
  const { ATTR_SERVICE_NAME, ATTR_SERVICE_VERSION } = await import("@opentelemetry/semantic-conventions");
  const { metrics, diag, DiagConsoleLogger, DiagLogLevel, propagation } = await import("@opentelemetry/api");
  const { logs } = await import("@opentelemetry/api-logs");
  const { registerInstrumentations } = await import("@opentelemetry/instrumentation");
  const { CompositePropagator, W3CTraceContextPropagator, W3CBaggagePropagator } = await import("@opentelemetry/core");
  const { HttpInstrumentation } = await import("@opentelemetry/instrumentation-http");
  const { WinstonInstrumentation } = await import("@opentelemetry/instrumentation-winston");
  const { RuntimeNodeInstrumentation } = await import("@opentelemetry/instrumentation-runtime-node");

  // Check if OTLP export is enabled (exporters auto-read endpoint from env)
  const otlpEnabled = !!process.env.OTEL_EXPORTER_OTLP_ENDPOINT;
  const serviceName = process.env.OTEL_SERVICE_NAME || "zitadel-login";
  const metricsPort = parseInt(process.env.OTEL_METRICS_PORT || "9464", 10);

  // Enable OTEL diagnostics in debug mode
  if (process.env.OTEL_LOG_LEVEL === "debug") {
    diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.DEBUG);
  }

  const gcpAttributes = getGcpResourceAttributes();
  const gcpDetected = Object.keys(gcpAttributes).length > 0;

  // Build resource with service info and detected attributes
  const resource = resourceFromAttributes({
    [ATTR_SERVICE_NAME]: serviceName,
    [ATTR_SERVICE_VERSION]: process.env.npm_package_version || "1.0.0",
    "deployment.environment": process.env.NODE_ENV || "development",
    "process.runtime.name": "nodejs",
    "process.runtime.version": process.version,
    ...(process.env.HOSTNAME && { "host.name": process.env.HOSTNAME }),
    ...gcpAttributes,
  });

  console.log(
    `OTEL: Initializing (service: ${serviceName}, OTLP: ${otlpEnabled ? "enabled" : "disabled"}, Prometheus: :${metricsPort}/metrics, GCP: ${gcpDetected})`,
  );

  // ============ TRACES ============
  // OTLPTraceExporter automatically reads from OTEL_EXPORTER_OTLP_* env vars
  const spanProcessors = otlpEnabled ? [new BatchSpanProcessor(new OTLPTraceExporter())] : [];

  const tracerProvider = new NodeTracerProvider({ resource, spanProcessors });
  tracerProvider.register();

  if (otlpEnabled) {
    console.log("OTEL: Trace OTLP exporter initialized");
  }

  // ============ PROPAGATION ============
  propagation.setGlobalPropagator(
    new CompositePropagator({
      propagators: [new W3CTraceContextPropagator(), new W3CBaggagePropagator()],
    }),
  );

  // ============ METRICS ============
  const metricReaders: MetricReader[] = [];

  // Prometheus exporter for scraping (always enabled)
  metricReaders.push(new PrometheusExporter({ port: metricsPort }));
  console.log(`OTEL: Prometheus metrics available on port ${metricsPort} at /metrics`);

  // OTLPMetricExporter automatically reads from OTEL_EXPORTER_OTLP_* env vars
  if (otlpEnabled) {
    const metricsInterval = parseInt(process.env.OTEL_METRICS_EXPORT_INTERVAL_MS || "30000", 10);
    metricReaders.push(
      new PeriodicExportingMetricReader({
        exporter: new OTLPMetricExporter(),
        exportIntervalMillis: metricsInterval,
      }),
    );
    console.log(`OTEL: Metrics OTLP exporter initialized (interval: ${metricsInterval}ms)`);
  }

  const meterProvider = new MeterProvider({ resource, readers: metricReaders });
  metrics.setGlobalMeterProvider(meterProvider);

  // ============ LOGS ============
  // OTLPLogExporter automatically reads from OTEL_EXPORTER_OTLP_* env vars
  const logRecordProcessors = otlpEnabled ? [new BatchLogRecordProcessor(new OTLPLogExporter())] : [];

  const loggerProvider = new LoggerProvider({ resource, processors: logRecordProcessors });
  _loggerProvider = loggerProvider;
  logs.setGlobalLoggerProvider(loggerProvider);

  if (otlpEnabled) {
    console.log("OTEL: Log OTLP exporter initialized");
  }

  // ============ INSTRUMENTATIONS ============
  registerInstrumentations({
    tracerProvider,
    meterProvider,
    instrumentations: [
      new HttpInstrumentation({
        ignoreIncomingRequestHook: (request) => {
          const ignoredPaths = ["/healthy", "/_next/", "/favicon.ico", "/__nextjs", "/metrics"];
          return ignoredPaths.some((path) => request.url?.includes(path));
        },
      }),
      new WinstonInstrumentation({
        disableLogSending: false,
        disableLogCorrelation: false,
      }),
      new RuntimeNodeInstrumentation({
        monitoringPrecision: 5000,
      }),
    ],
  });

  console.log("OTEL: Instrumentation complete");

  // ============ GRACEFUL SHUTDOWN ============
  const shutdown = async () => {
    console.log("OTEL: Shutting down...");
    try {
      await Promise.all([tracerProvider.shutdown(), meterProvider.shutdown(), loggerProvider.shutdown()]);
      console.log("OTEL: Shutdown complete");
    } catch (error) {
      console.error("OTEL: Error during shutdown", error);
    }
  };

  process.on("SIGTERM", shutdown);
  process.on("SIGINT", shutdown);
}

/**
 * Get the OTEL LoggerProvider for direct log emission if needed
 */
export function getLoggerProvider(): LoggerProvider | null {
  return _loggerProvider;
}
