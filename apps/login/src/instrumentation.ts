// OpenTelemetry instrumentation for Next.js
// Configures traces, metrics, and logs export via OTLP and Prometheus

import type { LoggerProvider } from "@opentelemetry/sdk-logs";
import type { MetricReader } from "@opentelemetry/sdk-metrics";

let _loggerProvider: LoggerProvider | null = null;

/** Detect GCP Cloud Run resource attributes from environment variables */
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

export async function register(): Promise<void> {
  if (process.env.NEXT_RUNTIME !== "nodejs") {
    return;
  }

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
  const { containerDetector } = await import("@opentelemetry/resource-detector-container");

  const otlpEnabled = !!process.env.OTEL_EXPORTER_OTLP_ENDPOINT;
  const serviceName = process.env.OTEL_SERVICE_NAME || "zitadel-login";
  const metricsPort = parseInt(process.env.OTEL_METRICS_PORT || "9464", 10);

  if (process.env.OTEL_LOG_LEVEL === "debug") {
    diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.DEBUG);
  }

  const gcpAttributes = getGcpResourceAttributes();
  const gcpDetected = Object.keys(gcpAttributes).length > 0;

  let containerAttributes: Record<string, string> = {};
  try {
    const detectedContainer = containerDetector.detect();
    const attrs = await detectedContainer.attributes;
    if (attrs && typeof attrs === "object") {
      for (const [key, value] of Object.entries(attrs)) {
        if (typeof value === "string") {
          containerAttributes[key] = value;
        }
      }
    }
  } catch {
    // Fails silently outside containerized environments
  }
  const containerDetected = Object.keys(containerAttributes).length > 0;

  const resource = resourceFromAttributes({
    [ATTR_SERVICE_NAME]: serviceName,
    [ATTR_SERVICE_VERSION]: process.env.npm_package_version || "1.0.0",
    "deployment.environment": process.env.NODE_ENV || "development",
    "process.runtime.name": "nodejs",
    "process.runtime.version": process.version,
    ...(process.env.HOSTNAME && { "host.name": process.env.HOSTNAME }),
    ...gcpAttributes,
    ...containerAttributes,
  });

  console.log(
    `OTEL: Initializing (service: ${serviceName}, OTLP: ${otlpEnabled ? "enabled" : "disabled"}, Prometheus: :${metricsPort}/metrics, GCP: ${gcpDetected}, Container: ${containerDetected})`,
  );

  // Traces
  const spanProcessors = otlpEnabled ? [new BatchSpanProcessor(new OTLPTraceExporter())] : [];

  const tracerProvider = new NodeTracerProvider({ resource, spanProcessors });
  tracerProvider.register();

  if (otlpEnabled) {
    console.log("OTEL: Trace OTLP exporter initialized");
  }

  // Propagation
  propagation.setGlobalPropagator(
    new CompositePropagator({
      propagators: [new W3CTraceContextPropagator(), new W3CBaggagePropagator()],
    }),
  );

  // Metrics
  const metricReaders: MetricReader[] = [];
  metricReaders.push(new PrometheusExporter({ port: metricsPort }));
  console.log(`OTEL: Prometheus metrics available on port ${metricsPort} at /metrics`);

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

  // Logs
  const logRecordProcessors = otlpEnabled ? [new BatchLogRecordProcessor(new OTLPLogExporter())] : [];

  const loggerProvider = new LoggerProvider({ resource, processors: logRecordProcessors });
  _loggerProvider = loggerProvider;
  logs.setGlobalLoggerProvider(loggerProvider);

  if (otlpEnabled) {
    console.log("OTEL: Log OTLP exporter initialized");
  }

  // Instrumentations
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

  // Shutdown
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

export function getLoggerProvider(): LoggerProvider | null {
  return _loggerProvider;
}
