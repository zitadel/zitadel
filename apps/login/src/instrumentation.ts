import type { LoggerProvider } from "@opentelemetry/api-logs";

/**
 * @module instrumentation
 *
 * OpenTelemetry instrumentation for Next.js applications using NodeSDK.
 * Set OTEL_SDK_DISABLED=true to disable instrumentation entirely.
 *
 * Configuration via standard OTEL environment variables:
 *
 * ### Exporters (all default to "otlp")
 * - OTEL_TRACES_EXPORTER: trace exporter (otlp, console, none)
 * - OTEL_METRICS_EXPORTER: metric exporter (otlp, prometheus, console, none)
 * - OTEL_LOGS_EXPORTER: log exporter (otlp, console, none)
 * - OTEL_EXPORTER_OTLP_ENDPOINT: base endpoint for OTLP signals
 * - OTEL_EXPORTER_PROMETHEUS_PORT: Prometheus scrape port (default 9464)
 *
 * ### Batch Processors
 * - OTEL_BSP_SCHEDULE_DELAY: span export delay (default 5000ms)
 * - OTEL_METRIC_EXPORT_INTERVAL: metric push interval (default 60000ms)
 *
 * ### Resource
 * - OTEL_SERVICE_NAME: service.name attribute
 * - OTEL_RESOURCE_ATTRIBUTES: additional attributes (key=value,key=value)
 *
 * @see https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/
 */

let _loggerProvider: LoggerProvider | null = null;

export async function register(): Promise<void> {
  if (process.env.NEXT_RUNTIME !== "nodejs") return;

  const [
    { NodeSDK },
    { envDetector, processDetector, hostDetector, resourceFromAttributes },
    { containerDetector },
    { gcpDetector },
    { logs },
    { getNodeAutoInstrumentations },
  ] = await Promise.all([
    import("@opentelemetry/sdk-node"),
    import("@opentelemetry/resources"),
    import("@opentelemetry/resource-detector-container"),
    import("@opentelemetry/resource-detector-gcp"),
    import("@opentelemetry/api-logs"),
    import("@opentelemetry/auto-instrumentations-node"),
  ]);

  const sdk = new NodeSDK({
    resource: resourceFromAttributes({
      "deployment.environment.name": process.env.NODE_ENV || "development",
    }),
    resourceDetectors: [
      envDetector,
      processDetector,
      hostDetector,
      containerDetector,
      gcpDetector,
    ],
    instrumentations: [
      getNodeAutoInstrumentations({
        "@opentelemetry/instrumentation-http": {
          ignoreIncomingRequestHook: (req) =>
            ["/healthy", "/_next/", "/favicon.ico", "/__nextjs", "/metrics"].some((p) =>
              req.url?.includes(p),
            ),
        },
        "@opentelemetry/instrumentation-fs": { enabled: false },
        "@opentelemetry/instrumentation-winston": { disableLogSending: false },
      }),
    ],
  });

  sdk.start();

  _loggerProvider = logs.getLoggerProvider();

  const shutdown = () => sdk.shutdown().catch(console.error);
  process.once("SIGTERM", shutdown);
  process.once("SIGINT", shutdown);
}

export function getLoggerProvider(): LoggerProvider | null {
  return _loggerProvider;
}
