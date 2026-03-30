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
  // Only run OpenTelemetry in the Node.js environment.
  // OTel is opt-in: set OTEL_ENABLED=true to enable it.
  // This prevents unnecessary initialization (and build tracing issues) in
  // environments like Vercel where OTel is not configured.
  if (process.env.NEXT_RUNTIME === "nodejs" && process.env.OTEL_ENABLED === "true") {
    const { registerNode } = await import("./instrumentation.node");
    _loggerProvider = await registerNode();
  }
}

export function getLoggerProvider(): LoggerProvider | null {
  return _loggerProvider;
}
