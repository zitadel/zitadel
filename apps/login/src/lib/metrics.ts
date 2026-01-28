import "server-only";
import { metrics, Counter, Histogram, UpDownCounter } from "@opentelemetry/api";

// Custom metrics interface
interface LoginMetrics {
  requestCounter: Counter | null;
  requestDuration: Histogram | null;
  requestSize: Histogram | null;
  responseSize: Histogram | null;
  responseStatusCounter: Counter | null;
  activeRequests: UpDownCounter | null;
  authAttempts: Counter | null;
  authSuccesses: Counter | null;
  authFailures: Counter | null;
  sessionCreationDuration: Histogram | null;
}

// Custom metrics holder
let metricsInitialized = false;
const loginMetrics: LoginMetrics = {
  requestCounter: null,
  requestDuration: null,
  requestSize: null,
  responseSize: null,
  responseStatusCounter: null,
  activeRequests: null,
  authAttempts: null,
  authSuccesses: null,
  authFailures: null,
  sessionCreationDuration: null,
};

function ensureMetricsInitialized(): void {
  if (metricsInitialized) return;

  const serviceName = process.env.OTEL_SERVICE_NAME || "zitadel-login";
  const meter = metrics.getMeter(serviceName);

  loginMetrics.requestCounter = meter.createCounter(
    "http_server_requests_total",
    {
      description: "Total number of HTTP requests",
    },
  );

  loginMetrics.requestDuration = meter.createHistogram(
    "http_server_request_duration_seconds",
    {
      description: "HTTP request duration in seconds",
      unit: "s",
    },
  );

  loginMetrics.requestSize = meter.createHistogram(
    "http_server_request_size_bytes",
    {
      description: "Size of HTTP request bodies in bytes",
      unit: "By",
    },
  );

  loginMetrics.responseSize = meter.createHistogram(
    "http_server_response_size_bytes",
    {
      description: "Size of HTTP response bodies in bytes",
      unit: "By",
    },
  );

  loginMetrics.responseStatusCounter = meter.createCounter(
    "http_server_response_status_total",
    {
      description: "Count of HTTP responses per status code",
    },
  );

  loginMetrics.activeRequests = meter.createUpDownCounter(
    "http_server_active_requests",
    {
      description: "Number of active HTTP requests",
    },
  );

  loginMetrics.authAttempts = meter.createCounter("login_auth_attempts_total", {
    description: "Total authentication attempts",
  });

  loginMetrics.authSuccesses = meter.createCounter(
    "login_auth_successes_total",
    {
      description: "Total successful authentications",
    },
  );

  loginMetrics.authFailures = meter.createCounter("login_auth_failures_total", {
    description: "Total failed authentications",
  });

  loginMetrics.sessionCreationDuration = meter.createHistogram(
    "login_session_creation_duration_seconds",
    {
      description: "Session creation duration in seconds",
      unit: "s",
    },
  );

  metricsInitialized = true;
}

/**
 * Record an authentication attempt
 * @param method - The authentication method (password, passkey, idp, etc.)
 * @param organization - The organization ID (optional)
 */
export function recordAuthAttempt(
  method: string,
  organization?: string,
): void {
  ensureMetricsInitialized();
  loginMetrics.authAttempts?.add(1, {
    method,
    organization: organization || "unknown",
  });
}

/**
 * Record a successful authentication
 * @param method - The authentication method (password, passkey, idp, etc.)
 * @param organization - The organization ID (optional)
 */
export function recordAuthSuccess(
  method: string,
  organization?: string,
): void {
  ensureMetricsInitialized();
  loginMetrics.authSuccesses?.add(1, {
    method,
    organization: organization || "unknown",
  });
}

/**
 * Record a failed authentication
 * @param method - The authentication method (password, passkey, idp, etc.)
 * @param reason - The failure reason
 * @param organization - The organization ID (optional)
 */
export function recordAuthFailure(
  method: string,
  reason: string,
  organization?: string,
): void {
  ensureMetricsInitialized();
  loginMetrics.authFailures?.add(1, {
    method,
    reason,
    organization: organization || "unknown",
  });
}

/**
 * Record the start of an HTTP request
 */
export function recordRequestStart(): void {
  ensureMetricsInitialized();
  loginMetrics.activeRequests?.add(1);
  loginMetrics.requestCounter?.add(1);
}

/**
 * Record the end of an HTTP request
 * @param startTime - The start time in milliseconds
 * @param statusCode - The HTTP status code
 * @param route - The route path
 * @param requestBytes - Size of request body (optional)
 * @param responseBytes - Size of response body (optional)
 */
export function recordRequestEnd(
  startTime: number,
  statusCode: number,
  route: string,
  requestBytes?: number,
  responseBytes?: number,
): void {
  ensureMetricsInitialized();
  loginMetrics.activeRequests?.add(-1);

  const labels = {
    status_code: statusCode.toString(),
    route,
  };

  const duration = (Date.now() - startTime) / 1000;
  loginMetrics.requestDuration?.record(duration, labels);

  // Record response status
  loginMetrics.responseStatusCounter?.add(1, {
    status_code: statusCode.toString(),
  });

  // Record sizes if provided
  if (requestBytes !== undefined) {
    loginMetrics.requestSize?.record(requestBytes, labels);
  }
  if (responseBytes !== undefined) {
    loginMetrics.responseSize?.record(responseBytes, labels);
  }
}

/**
 * Record session creation duration
 * @param startTime - The start time in milliseconds
 * @param success - Whether the session creation was successful
 */
export function recordSessionCreationDuration(
  startTime: number,
  success: boolean,
): void {
  ensureMetricsInitialized();
  const duration = (Date.now() - startTime) / 1000;
  loginMetrics.sessionCreationDuration?.record(duration, {
    success: success.toString(),
  });
}
