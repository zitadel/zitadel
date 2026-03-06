import type { LoggerProvider } from "@opentelemetry/api-logs";

export async function registerNode(): Promise<LoggerProvider | null> {
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

  // Use console during shutdown as the logger may be unavailable
  const shutdown = () => sdk.shutdown().catch(console.error);
  process.once("SIGTERM", shutdown);
  process.once("SIGINT", shutdown);

  return logs.getLoggerProvider();
}
