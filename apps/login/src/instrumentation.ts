// Next.js server instrumentation: runs once on the server at startup.
// OpenTelemetry tracing is initialized here for the Node.js runtime only.

export async function register() {
  // Only run in the Node.js runtime; never in the Edge runtime
  if (process.env.NEXT_RUNTIME === 'nodejs') {
    // Initialize OpenTelemetry tracing
    await import('./otel-tracing');
    // Initialize console patching (if enabled)
    await import('./logger');
  }
}
