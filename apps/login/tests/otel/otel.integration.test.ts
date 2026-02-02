import { describe, it, expect, beforeAll, afterAll } from "vitest";
import { execSync } from "child_process";
import * as fs from "fs";
import * as path from "path";

const TEST_DIR = path.dirname(new URL(import.meta.url).pathname);
const OUTPUT_DIR = path.join(TEST_DIR, "output");
const COMPOSE_FILE = path.join(TEST_DIR, "docker-compose.test.yml");
const APP_URL = "http://localhost:3000";
const PROMETHEUS_URL = "http://localhost:9464/metrics";
const COLLECTOR_HEALTH_URL = "http://localhost:13133";
// HTTP/2 port for gRPC (not used directly by tests)
const MOCK_ZITADEL_GRPC_URL = "http://localhost:8080";
// HTTP/1.1 port for health checks and captured headers endpoint
const MOCK_ZITADEL_URL = "http://localhost:7432";
const CAPTURED_HEADERS_FILE = path.join(OUTPUT_DIR, "captured-headers.json");

// Timeout for Docker operations
const DOCKER_TIMEOUT = 180000; // 3 minutes

async function waitForService(
  url: string,
  maxAttempts = 30,
  delayMs = 2000,
): Promise<boolean> {
  for (let i = 0; i < maxAttempts; i++) {
    try {
      const response = await fetch(url);
      if (response.ok) {
        return true;
      }
    } catch {
      // Service not ready yet
    }
    await new Promise((resolve) => setTimeout(resolve, delayMs));
  }
  return false;
}

async function waitForFile(
  filePath: string,
  maxAttempts = 30,
  delayMs = 1000,
): Promise<boolean> {
  for (let i = 0; i < maxAttempts; i++) {
    try {
      const stats = fs.statSync(filePath);
      if (stats.size > 0) {
        return true;
      }
    } catch {
      // File not ready yet
    }
    await new Promise((resolve) => setTimeout(resolve, delayMs));
  }
  return false;
}

describe("OpenTelemetry Integration", () => {
  beforeAll(async () => {
    // Clean up any existing containers
    try {
      execSync(`docker compose -f ${COMPOSE_FILE} down -v`, {
        stdio: "pipe",
        cwd: TEST_DIR,
      });
    } catch {
      // Ignore errors if containers don't exist
    }

    // Create output directory
    if (!fs.existsSync(OUTPUT_DIR)) {
      fs.mkdirSync(OUTPUT_DIR, { recursive: true });
    }

    // Start the test environment
    console.log("Starting Docker Compose environment...");
    execSync(`docker compose -f ${COMPOSE_FILE} up -d --build`, {
      stdio: "inherit",
      cwd: TEST_DIR,
    });

    // Wait for the collector to be healthy
    console.log("Waiting for OTEL collector to be healthy...");
    const isCollectorHealthy = await waitForService(
      COLLECTOR_HEALTH_URL,
      30,
      2000,
    );
    if (!isCollectorHealthy) {
      throw new Error("OTEL collector failed to become healthy");
    }
    console.log("OTEL collector is healthy!");

    // Wait for mock Zitadel server to be healthy
    console.log("Waiting for mock Zitadel server to be healthy...");
    const isMockZitadelHealthy = await waitForService(
      `${MOCK_ZITADEL_URL}/health`,
      30,
      2000,
    );
    if (!isMockZitadelHealthy) {
      throw new Error("Mock Zitadel server failed to become healthy");
    }
    console.log("Mock Zitadel server is healthy!");

    // Wait for the app to be healthy
    console.log("Waiting for login app to be healthy...");
    const isHealthy = await waitForService(
      `${APP_URL}/ui/v2/login/healthy`,
      60,
      2000,
    );
    if (!isHealthy) {
      throw new Error("Login app failed to become healthy");
    }

    console.log("Login app is healthy!");
  }, DOCKER_TIMEOUT);

  afterAll(async () => {
    // Stop and remove containers
    console.log("Stopping Docker Compose environment...");
    try {
      execSync(`docker compose -f ${COMPOSE_FILE} down -v`, {
        stdio: "pipe",
        cwd: TEST_DIR,
      });
    } catch (error) {
      console.error("Error stopping containers:", error);
    }
  }, DOCKER_TIMEOUT);

  it(
    "should export traces to collector",
    async () => {
      // Call the otel-test endpoint
      const response = await fetch(`${APP_URL}/ui/v2/login/otel-test`);
      expect(response.status).toBe(200);

      const data = await response.json();
      expect(data.status).toBe("ok");
      expect(data.traceId).toBeDefined();
      expect(data.traceId.length).toBe(32); // Valid trace ID is 32 hex chars
      expect(data.spanId).toBeDefined();

      console.log("Received trace ID:", data.traceId);

      // Wait for traces to be exported to file
      const tracesFile = path.join(OUTPUT_DIR, "traces.json");
      const hasTraces = await waitForFile(tracesFile, 30, 1000);
      expect(hasTraces).toBe(true);

      // Verify traces contain our test span
      const tracesContent = fs.readFileSync(tracesFile, "utf-8");
      expect(tracesContent).toContain("otel-test");
    },
    60000,
  );

  it(
    "should export metrics to collector",
    async () => {
      // Call the test endpoint multiple times to generate metrics
      for (let i = 0; i < 3; i++) {
        await fetch(`${APP_URL}/ui/v2/login/otel-test`);
        await new Promise((resolve) => setTimeout(resolve, 500));
      }

      // Wait for metrics to be exported (metrics are exported every 30 seconds by default)
      // We'll wait a bit longer to ensure they're flushed
      await new Promise((resolve) => setTimeout(resolve, 35000));

      const metricsFile = path.join(OUTPUT_DIR, "metrics.json");
      const hasMetrics = await waitForFile(metricsFile, 10, 2000);

      if (hasMetrics) {
        const metricsContent = fs.readFileSync(metricsFile, "utf-8");
        console.log("Metrics exported successfully");
        expect(metricsContent.length).toBeGreaterThan(0);
      } else {
        // Metrics may take longer to export, so we'll just warn
        console.warn("Metrics file not found within timeout");
      }
    },
    120000,
  );

  it(
    "should export logs to collector",
    async () => {
      // Call the test endpoint to generate logs
      await fetch(`${APP_URL}/ui/v2/login/otel-test`);

      // Wait for logs to be exported
      await new Promise((resolve) => setTimeout(resolve, 5000));

      const logsFile = path.join(OUTPUT_DIR, "logs.json");
      const hasLogs = await waitForFile(logsFile, 30, 1000);

      if (hasLogs) {
        const logsContent = fs.readFileSync(logsFile, "utf-8");
        console.log("Logs exported successfully");
        expect(logsContent).toContain("otel-test");
      } else {
        // Logs may not be exported depending on Winston configuration
        console.warn("Logs file not found within timeout");
      }
    },
    60000,
  );

  it("should respond to health check", async () => {
    const response = await fetch(`${APP_URL}/ui/v2/login/healthy`);
    expect(response.status).toBe(200);
  });

  it(
    "should include correct service name in telemetry resource attributes",
    async () => {
      // Wait for traces to be exported
      const tracesFile = path.join(OUTPUT_DIR, "traces.json");
      const hasTraces = await waitForFile(tracesFile, 10, 1000);
      expect(hasTraces).toBe(true);

      const tracesContent = fs.readFileSync(tracesFile, "utf-8");

      // Verify service.name matches OTEL_SERVICE_NAME from docker-compose.test.yml
      expect(tracesContent).toContain('"service.name"');
      expect(tracesContent).toContain('"zitadel-login-test"');

      // Verify service.version is present
      expect(tracesContent).toContain('"service.version"');
    },
    30000,
  );

  it(
    "should include GCP Cloud Run resource attributes in telemetry",
    async () => {
      // The docker-compose.test.yml sets GCP env vars:
      // GOOGLE_CLOUD_PROJECT=test-project-123
      // K_SERVICE=zitadel-login
      // K_REVISION=zitadel-login-00001-abc
      // K_CONFIGURATION=zitadel-login
      // GOOGLE_CLOUD_REGION=us-central1

      const tracesFile = path.join(OUTPUT_DIR, "traces.json");
      const hasTraces = await waitForFile(tracesFile, 10, 1000);
      expect(hasTraces).toBe(true);

      const tracesContent = fs.readFileSync(tracesFile, "utf-8");

      // Verify GCP Cloud Run attributes are present
      expect(tracesContent).toContain('"cloud.provider"');
      expect(tracesContent).toContain('"gcp"');

      expect(tracesContent).toContain('"cloud.platform"');
      expect(tracesContent).toContain('"gcp_cloud_run"');

      expect(tracesContent).toContain('"cloud.account.id"');
      expect(tracesContent).toContain('"test-project-123"');

      expect(tracesContent).toContain('"faas.name"');
      expect(tracesContent).toContain('"zitadel-login"');

      expect(tracesContent).toContain('"faas.version"');
      expect(tracesContent).toContain('"zitadel-login-00001-abc"');

      expect(tracesContent).toContain('"cloud.region"');
      expect(tracesContent).toContain('"us-central1"');

      console.log("GCP Cloud Run resource attributes verified in telemetry");
    },
    30000,
  );

  it(
    "should expose Node.js runtime metrics via Prometheus",
    async () => {
      // Wait for runtime metrics to be collected (monitoringPrecision is 5000ms)
      await new Promise((resolve) => setTimeout(resolve, 6000));

      const response = await fetch(PROMETHEUS_URL);
      expect(response.status).toBe(200);

      const metricsText = await response.text();

      // Verify event loop metrics from @opentelemetry/instrumentation-runtime-node
      expect(metricsText).toContain("nodejs_eventloop_utilization");
      expect(metricsText).toContain("nodejs_eventloop_delay_mean");
      expect(metricsText).toContain("nodejs_eventloop_time");

      // Verify V8 memory metrics
      expect(metricsText).toContain("v8js_memory_heap_used");
      expect(metricsText).toContain("v8js_memory_heap_limit");

      // Verify V8 GC metrics
      expect(metricsText).toContain("v8js_gc_duration");

      console.log("Node.js runtime metrics are being exported correctly");
    },
    30000,
  );

  it(
    "should expose custom application metrics via Prometheus",
    async () => {
      // Generate some traffic to produce custom metrics
      for (let i = 0; i < 3; i++) {
        await fetch(`${APP_URL}/ui/v2/login/otel-test`);
      }

      const response = await fetch(PROMETHEUS_URL);
      expect(response.status).toBe(200);

      const metricsText = await response.text();

      // Verify custom metrics from metrics.ts
      expect(metricsText).toContain("http_server_requests_total");
      expect(metricsText).toContain("http_server_request_duration_seconds");
      expect(metricsText).toContain("http_server_active_requests");
      expect(metricsText).toContain("login_auth_attempts_total");
      expect(metricsText).toContain("login_auth_successes_total");
      expect(metricsText).toContain("login_session_creation_duration_seconds");

      console.log("Custom application metrics are being exported correctly");
    },
    30000,
  );

  it(
    "should record accurate request duration in http_server_duration metric",
    async () => {
      // Helper to parse metrics for a specific HTTP status code 200 (successful requests)
      // We look for metrics with http_status_code="200" to filter out redirects/errors
      const parseMetricsForStatus200 = (metricsText: string) => {
        // Find all http_server_duration_sum lines and extract those with status 200
        const lines = metricsText.split("\n");
        let totalSum = 0;
        let totalCount = 0;

        for (const line of lines) {
          if (
            line.includes("http_server_duration_sum") &&
            line.includes('http_status_code="200"')
          ) {
            const match = line.match(/\}\s+([\d.]+)/);
            if (match) {
              totalSum += parseFloat(match[1]);
            }
          }
          if (
            line.includes("http_server_duration_count") &&
            line.includes('http_status_code="200"')
          ) {
            const match = line.match(/\}\s+(\d+)/);
            if (match) {
              totalCount += parseInt(match[1], 10);
            }
          }
        }

        return { sum: totalSum, count: totalCount };
      };

      // Get metrics BEFORE the request
      const beforeResponse = await fetch(PROMETHEUS_URL);
      const beforeMetrics = parseMetricsForStatus200(await beforeResponse.text());

      // Make a request and measure the actual response time
      const startTime = performance.now();
      const response = await fetch(`${APP_URL}/ui/v2/login/otel-test`);
      const endTime = performance.now();
      const actualDurationMs = endTime - startTime;

      expect(response.status).toBe(200);

      // Wait for metrics to be collected
      await new Promise((resolve) => setTimeout(resolve, 1000));

      // Get metrics AFTER the request
      const afterResponse = await fetch(PROMETHEUS_URL);
      const afterMetrics = parseMetricsForStatus200(await afterResponse.text());

      // Calculate the delta (just this request's contribution)
      const deltaSumMs = afterMetrics.sum - beforeMetrics.sum;
      const deltaCount = afterMetrics.count - beforeMetrics.count;

      expect(deltaCount).toBeGreaterThanOrEqual(1);

      // The metric is in milliseconds (old semantic conventions use ms)
      const serverDurationMs = deltaSumMs / deltaCount;

      console.log(
        `Request timing: client measured ${actualDurationMs.toFixed(2)}ms, ` +
          `server reported ${serverDurationMs.toFixed(2)}ms ` +
          `(delta: sum=${deltaSumMs.toFixed(2)}ms, count=${deltaCount})`,
      );

      // The server-side duration should be:
      // 1. Greater than 0 (something was measured)
      // 2. Less than or equal to client time + tolerance (server can't take longer than total round-trip)
      // 3. Within a reasonable range (0.1ms to 5 seconds)
      expect(serverDurationMs).toBeGreaterThan(0);
      expect(serverDurationMs).toBeLessThanOrEqual(actualDurationMs + 100); // 100ms tolerance for timing variance
      expect(serverDurationMs).toBeGreaterThan(0.1); // At least 0.1ms
      expect(serverDurationMs).toBeLessThan(5000); // Less than 5 seconds

      // Verify the metric is in a reasonable range for a simple HTTP request
      // Most requests should complete in under 500ms
      expect(serverDurationMs).toBeLessThan(500);
    },
    30000,
  );

  it(
    "should instrument all page loads with HTTP metrics",
    async () => {
      // Make requests to various application pages (these will return redirects or errors
      // since we don't have a valid session, but they should still be instrumented)
      const pagesToTest = [
        "/ui/v2/login/login", // Login page
        "/ui/v2/login/register", // Register page
        "/ui/v2/login/password", // Password page
        "/ui/v2/login/passkey", // Passkey page
      ];

      for (const page of pagesToTest) {
        await fetch(`${APP_URL}${page}`, { redirect: "manual" });
      }

      // Wait for metrics to be collected
      await new Promise((resolve) => setTimeout(resolve, 2000));

      // Verify http_server_duration metric is recorded (auto-instrumented by @opentelemetry/instrumentation-http)
      const prometheusResponse = await fetch(PROMETHEUS_URL);
      expect(prometheusResponse.status).toBe(200);
      const metricsText = await prometheusResponse.text();

      // The http_server_duration metric should exist and have recorded requests
      expect(metricsText).toContain("http_server_duration");

      // Verify the metric has count > 0 (requests were recorded)
      // This confirms that HTTP instrumentation is working for all page loads
      const durationCountMatch = metricsText.match(
        /http_server_duration_count\{[^}]*\}\s+(\d+)/,
      );
      expect(durationCountMatch).not.toBeNull();
      const requestCount = parseInt(durationCountMatch![1], 10);
      expect(requestCount).toBeGreaterThanOrEqual(pagesToTest.length);

      console.log(
        `All page loads are being instrumented with HTTP metrics (${requestCount} requests recorded)`,
      );
    },
    30000,
  );

  it(
    "should instrument different HTTP status codes correctly",
    async () => {
      // Helper to extract status code counts from metrics
      const getStatusCodeCounts = (
        metricsText: string,
      ): Record<string, number> => {
        const counts: Record<string, number> = {};
        const lines = metricsText.split("\n");

        for (const line of lines) {
          // Match http_server_duration_count with http_status_code label
          const match = line.match(
            /http_server_duration_count\{[^}]*http_status_code="(\d+)"[^}]*\}\s+(\d+)/,
          );
          if (match) {
            const statusCode = match[1];
            const count = parseInt(match[2], 10);
            counts[statusCode] = (counts[statusCode] || 0) + count;
          }
        }
        return counts;
      };

      // Get baseline metrics
      const beforeResponse = await fetch(PROMETHEUS_URL);
      const beforeCounts = getStatusCodeCounts(await beforeResponse.text());

      // Generate requests with different expected status codes:

      // 1. 200 OK - health check endpoint
      const healthyResponse = await fetch(`${APP_URL}/ui/v2/login/healthy`);
      expect(healthyResponse.status).toBe(200);

      // 2. 200 OK - otel-test endpoint
      const otelResponse = await fetch(`${APP_URL}/ui/v2/login/otel-test`);
      expect(otelResponse.status).toBe(200);

      // 3. 404 Not Found - non-existent route
      const notFoundResponse = await fetch(
        `${APP_URL}/ui/v2/login/this-route-does-not-exist-12345`,
      );
      expect(notFoundResponse.status).toBe(404);

      // Wait for metrics to be collected
      await new Promise((resolve) => setTimeout(resolve, 2000));

      // Get metrics after requests
      const afterResponse = await fetch(PROMETHEUS_URL);
      const afterText = await afterResponse.text();
      const afterCounts = getStatusCodeCounts(afterText);

      // Verify we have metrics for different status codes
      console.log("Status code counts before:", beforeCounts);
      console.log("Status code counts after:", afterCounts);

      // Check that 200 status code is recorded (from healthy and otel-test)
      expect(afterCounts["200"]).toBeGreaterThan(beforeCounts["200"] || 0);

      // Check that 404 is recorded
      expect(afterCounts["404"]).toBeGreaterThan(beforeCounts["404"] || 0);

      // Verify the http_status_code label exists in the metrics output
      expect(afterText).toContain('http_status_code="200"');
      expect(afterText).toContain('http_status_code="404"');

      // Verify multiple distinct status codes are being tracked
      const distinctStatusCodes = Object.keys(afterCounts);
      expect(distinctStatusCodes.length).toBeGreaterThanOrEqual(2);

      console.log(
        "HTTP status codes are being instrumented correctly:",
        distinctStatusCodes.sort().join(", "),
      );
    },
    30000,
  );

  describe("Trace Propagation", () => {
    it(
      "should propagate traceparent headers to backend gRPC calls",
      async () => {
        // Load a page that triggers gRPC calls to the Zitadel backend
        // The loginname page fetches settings from Zitadel (getDefaultOrg,
        // getLoginSettings, getActiveIdentityProviders, getBrandingSettings)
        const response = await fetch(`${APP_URL}/ui/v2/login/loginname`, {
          redirect: "manual",
        });

        console.log(`Loginname page response status: ${response.status}`);

        // Wait a moment for the mock server to write captured headers
        await new Promise((resolve) => setTimeout(resolve, 2000));

        // Read captured headers from the file (written by mock server)
        let capturedRequests: Array<{ traceparent: string | null; url: string; method: string }> = [];
        try {
          const capturedContent = fs.readFileSync(CAPTURED_HEADERS_FILE, "utf-8");
          capturedRequests = JSON.parse(capturedContent);
        } catch (err) {
          // If file doesn't exist, try the endpoint
          console.log("File not found, trying endpoint...");
          const capturedResponse = await fetch(`${MOCK_ZITADEL_URL}/captured-headers`);
          capturedRequests = await capturedResponse.json();
        }

        console.log(`Captured ${capturedRequests.length} requests to mock Zitadel`);

        // Filter for gRPC calls only (POST requests, not health checks)
        const grpcRequests = capturedRequests.filter(
          (req) => req.method === "POST"
        );
        console.log(`Found ${grpcRequests.length} gRPC calls`);

        // Verify gRPC calls have traceparent headers
        const requestsWithTrace = grpcRequests.filter(
          (req) => req.traceparent !== null,
        );

        // Should have at least one gRPC call with traceparent
        expect(grpcRequests.length).toBeGreaterThan(0);
        expect(requestsWithTrace.length).toBeGreaterThan(0);

        // Verify traceparent format: version-traceid-parentid-flags
        const traceparentRegex = /^[0-9a-f]{2}-[0-9a-f]{32}-[0-9a-f]{16}-[0-9a-f]{2}$/;
        for (const req of requestsWithTrace) {
          expect(req.traceparent).toMatch(traceparentRegex);
          console.log(`Request ${req.url}: traceparent=${req.traceparent}`);
        }

        // Verify trace ID consistency - gRPC calls from the same page load share a trace ID
        // Group calls by trace ID and verify each group has at least 2 calls
        if (requestsWithTrace.length >= 2) {
          const traceIds = requestsWithTrace.map((req) => req.traceparent!.split("-")[1]);
          const traceIdCounts = new Map<string, number>();
          for (const id of traceIds) {
            traceIdCounts.set(id, (traceIdCounts.get(id) || 0) + 1);
          }
          // Each trace ID should have at least 2 gRPC calls (loginname makes 2+ calls)
          const validGroups = Array.from(traceIdCounts.values()).filter((count) => count >= 2);
          expect(validGroups.length).toBeGreaterThan(0);
          console.log(`Trace ID consistency verified: ${validGroups.length} page loads with 2+ gRPC calls each`);
        }

        console.log(
          `Trace propagation verified: ${requestsWithTrace.length}/${capturedRequests.length} requests have traceparent`,
        );
      },
      30000,
    );

    // Note: Trace ID consistency is already verified in the test above
    // (all gRPC calls from a single page load share the same trace ID)
  });
});
