import "server-only";
import { metrics } from "@opentelemetry/api";

const serviceName = process.env.OTEL_SERVICE_NAME || "zitadel-login";
const meter = metrics.getMeter(serviceName);

const requestCounter = meter.createCounter("http_server_requests_total", {
  description: "Total number of HTTP requests",
});

const requestDuration = meter.createHistogram(
  "http_server_request_duration_seconds",
  {
    description: "HTTP request duration in seconds",
    unit: "s",
  },
);

const requestSize = meter.createHistogram("http_server_request_size_bytes", {
  description: "Size of HTTP request bodies in bytes",
  unit: "By",
});

const responseSize = meter.createHistogram("http_server_response_size_bytes", {
  description: "Size of HTTP response bodies in bytes",
  unit: "By",
});

const responseStatusCounter = meter.createCounter(
  "http_server_response_status_total",
  {
    description: "Count of HTTP responses per status code",
  },
);

const activeRequests = meter.createUpDownCounter("http_server_active_requests", {
  description: "Number of active HTTP requests",
});

const authAttempts = meter.createCounter("login_auth_attempts_total", {
  description: "Total authentication attempts",
});

const authSuccesses = meter.createCounter("login_auth_successes_total", {
  description: "Total successful authentications",
});

const authFailures = meter.createCounter("login_auth_failures_total", {
  description: "Total failed authentications",
});

const sessionCreationDuration = meter.createHistogram(
  "login_session_creation_duration_seconds",
  {
    description: "Session creation duration in seconds",
    unit: "s",
  },
);

/**
 * Record an authentication attempt.
 */
export function recordAuthAttempt(method: string, organization?: string): void {
  authAttempts.add(1, {
    method,
    organization: organization || "unknown",
  });
}

/**
 * Record a successful authentication.
 */
export function recordAuthSuccess(method: string, organization?: string): void {
  authSuccesses.add(1, {
    method,
    organization: organization || "unknown",
  });
}

/**
 * Record a failed authentication.
 */
export function recordAuthFailure(
  method: string,
  reason: string,
  organization?: string,
): void {
  authFailures.add(1, {
    method,
    reason,
    organization: organization || "unknown",
  });
}

/**
 * Record the start of an HTTP request.
 */
export function recordRequestStart(): void {
  activeRequests.add(1);
  requestCounter.add(1);
}

/**
 * Record the end of an HTTP request.
 */
export function recordRequestEnd(
  startTime: number,
  statusCode: number,
  route: string,
  requestBytes?: number,
  responseBytes?: number,
): void {
  activeRequests.add(-1);

  const labels = {
    status_code: statusCode.toString(),
    route,
  };

  const duration = (Date.now() - startTime) / 1000;
  requestDuration.record(duration, labels);

  responseStatusCounter.add(1, {
    status_code: statusCode.toString(),
  });

  if (requestBytes !== undefined) {
    requestSize.record(requestBytes, labels);
  }
  if (responseBytes !== undefined) {
    responseSize.record(responseBytes, labels);
  }
}

/**
 * Record session creation duration.
 */
export function recordSessionCreationDuration(
  startTime: number,
  success: boolean,
): void {
  const duration = (Date.now() - startTime) / 1000;
  sessionCreationDuration.record(duration, {
    success: success.toString(),
  });
}
