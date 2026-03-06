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
  // Only run OpenTelemetry in the Node.js environment
  if (process.env.NEXT_RUNTIME === "nodejs") {
    // Disable by default in local development to avoid unnecessary overhead
    if (process.env.NODE_ENV === "development" && process.env.OTEL_SDK_DISABLED !== "false") {
      return;
    }

    // Explicit check for disabled env variable
    if (process.env.OTEL_SDK_DISABLED === "true") return;

    const { registerNode } = await import("./instrumentation.node");
    _loggerProvider = await registerNode();
  }
}

export function getLoggerProvider(): LoggerProvider | null {
  return _loggerProvider;
}
