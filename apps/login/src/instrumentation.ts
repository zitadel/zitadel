// NextJS instrumentation entry point
// This file is automatically loaded by NextJS and must export a `register` function
//
// Exports all three OTEL signals:
// - Traces: via OTLP HTTP push
// - Metrics: via OTLP HTTP push AND Prometheus scraping endpoint
// - Logs: via OTLP HTTP push (automatic with Winston instrumentation)
//
// Supported environment variables:
// - OTEL_SERVICE_NAME: Service name (default: "zitadel-login")
// - OTEL_EXPORTER_OTLP_ENDPOINT: OTLP collector endpoint (required for OTLP export)
// - OTEL_EXPORTER_OTLP_HEADERS: Headers for OTLP requests (key=value,key2=value2)
// - OTEL_LOG_LEVEL: Log level (debug, info, warn, error)
// - OTEL_METRICS_EXPORT_INTERVAL_MS: Metrics export interval for OTLP (default: 30000)
// - OTEL_METRICS_PORT: Port for Prometheus metrics scraping (default: 9464)
// - OTEL_SEMCONV_STABILITY_OPT_IN: Semantic convention stability opt-in (for future use)
//
// GCP Cloud Run environment variables (auto-detected):
// - K_SERVICE, K_REVISION, K_CONFIGURATION
// - GOOGLE_CLOUD_PROJECT
// - GOOGLE_CLOUD_REGION
//
// Note: Current instrumentation uses old semantic conventions (pre-v1.23):
// - http.server.duration: milliseconds (ms)
// When @opentelemetry/instrumentation-http is upgraded to 0.54.0+, set OTEL_SEMCONV_STABILITY_OPT_IN=http
// to use stable conventions with http.server.request.duration in seconds.

// Note: All OTEL imports are done dynamically in register() to avoid
// webpack bundling issues with Node.js-only modules like 'http'

import type { LoggerProvider } from "@opentelemetry/sdk-logs";
import type { MetricReader } from "@opentelemetry/sdk-metrics";

// Store provider for external access
let _loggerProvider: LoggerProvider | null = null;

/**
 * Parse OTLP headers from environment variable
 */
function parseOtlpHeaders(): Record<string, string> {
  const headersEnv = process.env.OTEL_EXPORTER_OTLP_HEADERS;
  if (!headersEnv) return {};

  const headers: Record<string, string> = {};
  headersEnv.split(",").forEach((header) => {
    const [key, ...valueParts] = header.split("=");
    const value = valueParts.join("=");
    if (key && value) {
      headers[key.trim()] = value.trim();
    }
  });
  return headers;
}

/**
 * Detect GCP Cloud Run resource attributes from environment variables
 */
function getGcpResourceAttributes(): Record<string, string> {
  const attributes: Record<string, string> = {};

  const projectId =
    process.env.GOOGLE_CLOUD_PROJECT || process.env.GCLOUD_PROJECT;
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

  const region =
    process.env.GOOGLE_CLOUD_REGION || process.env.CLOUD_RUN_REGION;
  if (region) {
    attributes["cloud.region"] = region;
  }

  return attributes;
}

/**
 * Get host and process attributes
 */
function getHostAttributes(): Record<string, string> {
  const attributes: Record<string, string> = {};

  const hostname = process.env.HOSTNAME;
  if (hostname) {
    attributes["host.name"] = hostname;
  }

  attributes["process.runtime.name"] = "nodejs";
  attributes["process.runtime.version"] = process.version;

  return attributes;
}

/**
 * Initialize OpenTelemetry with:
 * - Traces: OTLP HTTP push
 * - Metrics: OTLP HTTP push + Prometheus scraping
 * - Logs: OTLP HTTP push (via Winston instrumentation)
 */
export async function register(): Promise<void> {
  // Only run on server side (Node.js runtime)
  if (process.env.NEXT_RUNTIME !== "nodejs") {
    return;
  }

  // Note: Stable HTTP semantic conventions (OTEL_SEMCONV_STABILITY_OPT_IN=http) require
  // @opentelemetry/instrumentation-http v0.54.0+. Current version uses old conventions.
  // See: https://opentelemetry.io/docs/specs/semconv/http/http-metrics/

  // Dynamic imports to avoid webpack bundling issues
  const { NodeTracerProvider } = await import("@opentelemetry/sdk-trace-node");
  const { BatchSpanProcessor } = await import("@opentelemetry/sdk-trace-base");
  const { OTLPTraceExporter } = await import("@opentelemetry/exporter-trace-otlp-http");
  const { MeterProvider, PeriodicExportingMetricReader } = await import("@opentelemetry/sdk-metrics");
  const { OTLPMetricExporter } = await import("@opentelemetry/exporter-metrics-otlp-http");
  // Use webpackIgnore to prevent webpack from trying to bundle this Node.js-only module
  const { PrometheusExporter } = await import(/* webpackIgnore: true */ "@opentelemetry/exporter-prometheus");
  const { LoggerProvider, BatchLogRecordProcessor } = await import("@opentelemetry/sdk-logs");
  const { OTLPLogExporter } = await import("@opentelemetry/exporter-logs-otlp-http");
  const { resourceFromAttributes } = await import("@opentelemetry/resources");
  const { ATTR_SERVICE_NAME, ATTR_SERVICE_VERSION } = await import("@opentelemetry/semantic-conventions");
  const { metrics, diag, DiagConsoleLogger, DiagLogLevel, propagation } = await import("@opentelemetry/api");
  const { logs } = await import("@opentelemetry/api-logs");
  const { registerInstrumentations } = await import("@opentelemetry/instrumentation");
  const { CompositePropagator, W3CTraceContextPropagator, W3CBaggagePropagator } = await import(
    "@opentelemetry/core"
  );
  const { HttpInstrumentation } = await import("@opentelemetry/instrumentation-http");
  const { WinstonInstrumentation } = await import("@opentelemetry/instrumentation-winston");
  const { RuntimeNodeInstrumentation } = await import("@opentelemetry/instrumentation-runtime-node");

  const endpoint = process.env.OTEL_EXPORTER_OTLP_ENDPOINT;
  const serviceName = process.env.OTEL_SERVICE_NAME || "zitadel-login";
  const metricsPort = parseInt(process.env.OTEL_METRICS_PORT || "9464", 10);

  // Enable OTEL diagnostics in debug mode
  const logLevel = process.env.OTEL_LOG_LEVEL || "info";
  if (logLevel === "debug") {
    diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.DEBUG);
  }

  const headers = parseOtlpHeaders();
  const gcpAttributes = getGcpResourceAttributes();
  const hostAttributes = getHostAttributes();
  const gcpDetected = Object.keys(gcpAttributes).length > 0;

  const resource = resourceFromAttributes({
    [ATTR_SERVICE_NAME]: serviceName,
    [ATTR_SERVICE_VERSION]: process.env.npm_package_version || "1.0.0",
    "deployment.environment": process.env.NODE_ENV || "development",
    ...gcpAttributes,
    ...hostAttributes,
  });

  console.log(
    `OTEL: Initializing (service: ${serviceName}, OTLP: ${endpoint || "disabled"}, Prometheus: :${metricsPort}/metrics, GCP: ${gcpDetected})`,
  );

  // ============ TRACES ============
  const spanProcessors = endpoint
    ? [
        new BatchSpanProcessor(
          new OTLPTraceExporter({
            url: `${endpoint}/v1/traces`,
            headers,
          }),
        ),
      ]
    : [];

  if (endpoint) {
    console.log("OTEL: Trace OTLP exporter initialized");
  }

  const tracerProvider = new NodeTracerProvider({ resource, spanProcessors });
  tracerProvider.register();
  console.log("OTEL: Trace provider registered");

  // ============ PROPAGATION ============
  // Set up W3C trace context and baggage propagation for cross-service tracing
  propagation.setGlobalPropagator(
    new CompositePropagator({
      propagators: [new W3CTraceContextPropagator(), new W3CBaggagePropagator()],
    }),
  );
  console.log("OTEL: W3C trace context and baggage propagation enabled");

  // ============ METRICS ============
  const metricReaders: MetricReader[] = [];

  // Prometheus exporter for scraping (always enabled)
  const prometheusExporter = new PrometheusExporter({ port: metricsPort });
  metricReaders.push(prometheusExporter);
  console.log(`OTEL: Prometheus metrics available on port ${metricsPort} at /metrics`);

  // OTLP exporter for push (if endpoint configured)
  if (endpoint) {
    const metricExporter = new OTLPMetricExporter({
      url: `${endpoint}/v1/metrics`,
      headers,
    });
    const metricsInterval = parseInt(
      process.env.OTEL_METRICS_EXPORT_INTERVAL_MS || "30000",
      10,
    );
    metricReaders.push(
      new PeriodicExportingMetricReader({
        exporter: metricExporter,
        exportIntervalMillis: metricsInterval,
      }),
    );
    console.log(`OTEL: Metrics OTLP exporter initialized (interval: ${metricsInterval}ms)`);
  }

  const meterProvider = new MeterProvider({ resource, readers: metricReaders });
  metrics.setGlobalMeterProvider(meterProvider);
  console.log("OTEL: Meter provider registered");

  // ============ LOGS ============
  const logRecordProcessors = endpoint
    ? [
        new BatchLogRecordProcessor(
          new OTLPLogExporter({
            url: `${endpoint}/v1/logs`,
            headers,
          }),
        ),
      ]
    : [];

  if (endpoint) {
    console.log("OTEL: Log OTLP exporter initialized");
  }

  const loggerProvider = new LoggerProvider({ resource, processors: logRecordProcessors });
  _loggerProvider = loggerProvider;
  logs.setGlobalLoggerProvider(loggerProvider);
  console.log("OTEL: Log provider registered");

  // ============ INSTRUMENTATIONS ============
  registerInstrumentations({
    tracerProvider,
    meterProvider,
    instrumentations: [
      new HttpInstrumentation({
        ignoreIncomingRequestHook: (request) => {
          const ignoredPaths = [
            "/healthy",
            "/_next/",
            "/favicon.ico",
            "/__nextjs",
            "/metrics",
          ];
          return ignoredPaths.some((path) => request.url?.includes(path));
        },
      }),
      // Winston instrumentation: auto-injects trace context and sends logs to OTEL
      new WinstonInstrumentation({
        disableLogSending: false,
        disableLogCorrelation: false,
      }),
      // Runtime instrumentation: Node.js runtime metrics (GC, event loop, memory, etc.)
      new RuntimeNodeInstrumentation({
        monitoringPrecision: 5000, // Collect metrics every 5 seconds
      }),
    ],
  });

  console.log("OTEL: HTTP, Winston, and Runtime instrumentation registered");

  // ============ GRACEFUL SHUTDOWN ============
  const shutdown = async () => {
    console.log("OTEL: Shutting down...");
    try {
      await Promise.all([
        tracerProvider?.shutdown(),
        meterProvider?.shutdown(),
        loggerProvider?.shutdown(),
      ]);
      console.log("OTEL: Shutdown complete");
    } catch (error) {
      console.error("OTEL: Error during shutdown", error);
    }
  };

  process.on("SIGTERM", shutdown);
  process.on("SIGINT", shutdown);

  console.log("OTEL: Instrumentation complete");
}

/**
 * Get the OTEL LoggerProvider for direct log emission if needed
 */
export function getLoggerProvider(): LoggerProvider | null {
  return _loggerProvider;
}
