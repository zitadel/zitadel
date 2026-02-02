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
const MOCK_ZITADEL_URL = "http://localhost:7432";

const DOCKER_TIMEOUT = 180000;
const TEST_TIMEOUT = 30000;

/**
 * Polls a health endpoint until it returns HTTP 200.
 * Used to wait for Docker services to become ready.
 *
 * @param url - The health check URL to poll
 * @param maxAttempts - Maximum number of polling attempts (default: 30)
 * @param delayMs - Delay between attempts in milliseconds (default: 2000)
 * @returns true if service became healthy, false if max attempts exceeded
 */
async function waitForService(url: string, maxAttempts = 30, delayMs = 2000): Promise<boolean> {
  for (let i = 0; i < maxAttempts; i++) {
    try {
      const response = await fetch(url);
      if (response.ok) return true;
    } catch {
      await new Promise((r) => setTimeout(r, delayMs));
    }
  }
  return false;
}

/**
 * Polls for a file to exist and contain data.
 * Used to wait for OTLP collector to write exported telemetry.
 *
 * @param filePath - Absolute path to the file
 * @param maxAttempts - Maximum number of polling attempts (default: 30)
 * @param delayMs - Delay between attempts in milliseconds (default: 1000)
 * @returns true if file exists with content, false if max attempts exceeded
 */
async function waitForFile(filePath: string, maxAttempts = 30, delayMs = 1000): Promise<boolean> {
  for (let i = 0; i < maxAttempts; i++) {
    try {
      if (fs.statSync(filePath).size > 0) return true;
    } catch {
      await new Promise((r) => setTimeout(r, delayMs));
    }
  }
  return false;
}

/**
 * Reads the traces.json file exported by the OTLP collector.
 * Contains all trace spans in OTLP JSON format.
 */
function readTraces(): string {
  return fs.readFileSync(path.join(OUTPUT_DIR, "traces.json"), "utf-8");
}

/**
 * Reads the metrics.json file exported by the OTLP collector.
 * Contains all metrics in OTLP JSON format.
 */
function readMetrics(): string {
  return fs.readFileSync(path.join(OUTPUT_DIR, "metrics.json"), "utf-8");
}

/**
 * Reads the logs.json file exported by the OTLP collector.
 * Contains all log records in OTLP JSON format.
 */
function readLogs(): string {
  return fs.readFileSync(path.join(OUTPUT_DIR, "logs.json"), "utf-8");
}

/**
 * Fetches metrics from the Prometheus scraping endpoint.
 * Returns metrics in Prometheus text exposition format.
 */
async function fetchPrometheusMetrics(): Promise<string> {
  const response = await fetch(PROMETHEUS_URL);
  return response.text();
}

/**
 * Retrieves HTTP headers captured by the mock Zitadel server.
 * Used to verify trace context propagation to backend gRPC calls.
 */
async function fetchCapturedHeaders(): Promise<Array<{ traceparent: string | null; url: string; method: string }>> {
  const filePath = path.join(OUTPUT_DIR, "captured-headers.json");
  try {
    return JSON.parse(fs.readFileSync(filePath, "utf-8"));
  } catch {
    const response = await fetch(`${MOCK_ZITADEL_URL}/captured-headers`);
    return response.json();
  }
}

describe("OpenTelemetry Integration", () => {
  beforeAll(async () => {
    try {
      execSync(`docker compose -f ${COMPOSE_FILE} down -v`, { stdio: "pipe", cwd: TEST_DIR });
    } catch {
      // Container cleanup may fail if not running
    }

    fs.mkdirSync(OUTPUT_DIR, { recursive: true });
    execSync(`docker compose -f ${COMPOSE_FILE} up -d --build`, { stdio: "inherit", cwd: TEST_DIR });

    const collectorReady = await waitForService(COLLECTOR_HEALTH_URL);
    if (!collectorReady) throw new Error("OTEL collector failed to start");

    const mockReady = await waitForService(`${MOCK_ZITADEL_URL}/health`);
    if (!mockReady) throw new Error("Mock Zitadel failed to start");

    const appReady = await waitForService(`${APP_URL}/ui/v2/login/healthy`, 60);
    if (!appReady) throw new Error("Login app failed to start");

    await fetch(`${APP_URL}/ui/v2/login/otel-test`);
    await fetch(`${APP_URL}/ui/v2/login/loginname`, { redirect: "manual" });
    await fetch(`${APP_URL}/ui/v2/login/healthy`);
    await fetch(`${APP_URL}/ui/v2/login/this-does-not-exist-404`);

    await waitForFile(path.join(OUTPUT_DIR, "traces.json"));
    await new Promise((r) => setTimeout(r, 35000));
  }, DOCKER_TIMEOUT);

  afterAll(async () => {
    try {
      execSync(`docker compose -f ${COMPOSE_FILE} down -v`, { stdio: "pipe", cwd: TEST_DIR });
    } catch {
      // Container cleanup may fail if already stopped
    }
  }, DOCKER_TIMEOUT);

  describe("Traces", () => {
    describe("OTLP Export", () => {
      it("exports spans to collector", () => {
        expect(readTraces().length).toBeGreaterThan(0);
      }, TEST_TIMEOUT);

      it("contains valid trace IDs", () => {
        expect(readTraces()).toMatch(/"traceId":"[a-f0-9]{32}"/);
      }, TEST_TIMEOUT);

      it("contains valid span IDs", () => {
        expect(readTraces()).toMatch(/"spanId":"[a-f0-9]{16}"/);
      }, TEST_TIMEOUT);

      it("contains start timestamps", () => {
        expect(readTraces()).toContain("startTimeUnixNano");
      }, TEST_TIMEOUT);

      it("contains end timestamps", () => {
        expect(readTraces()).toContain("endTimeUnixNano");
      }, TEST_TIMEOUT);

      it("contains span status", () => {
        expect(readTraces()).toContain('"status"');
      }, TEST_TIMEOUT);

      it("contains span kind", () => {
        expect(readTraces()).toContain('"kind"');
      }, TEST_TIMEOUT);

      it("contains custom test spans", () => {
        expect(readTraces()).toContain("test-span");
      }, TEST_TIMEOUT);
    });

    describe("Span Attributes", () => {
      it("contains http.method", () => {
        expect(readTraces()).toContain('"http.method"');
      }, TEST_TIMEOUT);

      it("contains http.url", () => {
        expect(readTraces()).toContain('"http.url"');
      }, TEST_TIMEOUT);

      it("contains http.status_code", () => {
        expect(readTraces()).toContain('"http.status_code"');
      }, TEST_TIMEOUT);

      it("contains http.target", () => {
        expect(readTraces()).toContain('"http.target"');
      }, TEST_TIMEOUT);

      it("contains net.host.name", () => {
        expect(readTraces()).toContain('"net.host.name"');
      }, TEST_TIMEOUT);

      it("contains net.host.port", () => {
        expect(readTraces()).toContain('"net.host.port"');
      }, TEST_TIMEOUT);
    });

    describe("Trace Propagation", () => {
      it("propagates traceparent to gRPC", async () => {
        const requests = await fetchCapturedHeaders();
        const grpcCalls = requests.filter((r) => r.method === "POST" && r.traceparent !== null);
        expect(grpcCalls.length).toBeGreaterThan(0);
      }, TEST_TIMEOUT);

      it("uses W3C traceparent format", async () => {
        const requests = await fetchCapturedHeaders();
        const grpcCalls = requests.filter((r) => r.traceparent !== null);
        grpcCalls.forEach((r) => expect(r.traceparent).toMatch(/^[0-9a-f]{2}-[0-9a-f]{32}-[0-9a-f]{16}-[0-9a-f]{2}$/));
      }, TEST_TIMEOUT);

      it("uses traceparent version 00", async () => {
        const requests = await fetchCapturedHeaders();
        const grpcCalls = requests.filter((r) => r.traceparent !== null);
        grpcCalls.forEach((r) => expect(r.traceparent).toMatch(/^00-/));
      }, TEST_TIMEOUT);

      it("sets sampled flag to 01", async () => {
        const requests = await fetchCapturedHeaders();
        const grpcCalls = requests.filter((r) => r.traceparent !== null);
        grpcCalls.forEach((r) => expect(r.traceparent).toMatch(/-01$/));
      }, TEST_TIMEOUT);

      it("maintains trace ID consistency", async () => {
        const requests = await fetchCapturedHeaders();
        const grpcCalls = requests.filter((r) => r.method === "POST" && r.traceparent !== null);
        const traceIds = grpcCalls.map((r) => r.traceparent!.split("-")[1]);
        const uniqueIds = new Set(traceIds);
        expect(uniqueIds.size).toBeLessThan(traceIds.length);
      }, TEST_TIMEOUT);
    });

    describe("Resource Attributes", () => {
      it("contains service.name", () => {
        expect(readTraces()).toContain('"zitadel-login-test"');
      }, TEST_TIMEOUT);

      it("contains service.version", () => {
        expect(readTraces()).toContain('"service.version"');
      }, TEST_TIMEOUT);

      it("contains deployment.environment", () => {
        expect(readTraces()).toContain('"deployment.environment"');
      }, TEST_TIMEOUT);

      it("contains process.runtime.name", () => {
        expect(readTraces()).toContain('"nodejs"');
      }, TEST_TIMEOUT);

      it("contains process.runtime.version", () => {
        expect(readTraces()).toContain('"process.runtime.version"');
      }, TEST_TIMEOUT);

      it("contains cloud.provider", () => {
        expect(readTraces()).toContain('"gcp"');
      }, TEST_TIMEOUT);

      it("contains cloud.platform", () => {
        expect(readTraces()).toContain('"gcp_cloud_run"');
      }, TEST_TIMEOUT);

      it("contains cloud.account.id", () => {
        expect(readTraces()).toContain('"test-project-123"');
      }, TEST_TIMEOUT);

      it("contains cloud.region", () => {
        expect(readTraces()).toContain('"us-central1"');
      }, TEST_TIMEOUT);

      it("contains cloud.resource_id", () => {
        expect(readTraces()).toContain('"cloud.resource_id"');
      }, TEST_TIMEOUT);

      it("contains faas.name", () => {
        expect(readTraces()).toContain('"faas.name"');
      }, TEST_TIMEOUT);

      it("contains faas.version", () => {
        expect(readTraces()).toContain('"zitadel-login-00001-abc"');
      }, TEST_TIMEOUT);

      it("contains faas.instance", () => {
        expect(readTraces()).toContain('"faas.instance"');
      }, TEST_TIMEOUT);

      it("contains container.id when available", () => {
        const traces = readTraces();
        if (traces.includes('"container.id"')) {
          expect(traces).toMatch(/"container\.id"[^"]*"[a-f0-9]{64}"/);
        }
      }, TEST_TIMEOUT);
    });
  });

  describe("Metrics", () => {
    describe("OTLP Export", () => {
      it("exports metrics to collector", () => {
        expect(readMetrics().length).toBeGreaterThan(0);
      }, TEST_TIMEOUT);

      it("contains timestamps", () => {
        expect(readMetrics()).toContain("timeUnixNano");
      }, TEST_TIMEOUT);

      it("contains data points", () => {
        expect(readMetrics()).toContain("dataPoints");
      }, TEST_TIMEOUT);
    });

    describe("Prometheus Endpoint", () => {
      it("returns HTTP 200", async () => {
        const response = await fetch(PROMETHEUS_URL);
        expect(response.status).toBe(200);
      }, TEST_TIMEOUT);

      it("returns text/plain content type", async () => {
        const response = await fetch(PROMETHEUS_URL);
        expect(response.headers.get("content-type")).toContain("text/plain");
      }, TEST_TIMEOUT);
    });

    describe("Runtime Metrics", () => {
      it("exposes nodejs_eventloop_utilization", async () => {
        expect(await fetchPrometheusMetrics()).toContain("nodejs_eventloop_utilization");
      }, TEST_TIMEOUT);

      it("exposes nodejs_eventloop_delay", async () => {
        expect(await fetchPrometheusMetrics()).toContain("nodejs_eventloop_delay");
      }, TEST_TIMEOUT);

      it("exposes v8js_memory_heap_used", async () => {
        expect(await fetchPrometheusMetrics()).toContain("v8js_memory_heap_used");
      }, TEST_TIMEOUT);

      it("exposes v8js_memory_heap_limit", async () => {
        expect(await fetchPrometheusMetrics()).toContain("v8js_memory_heap_limit");
      }, TEST_TIMEOUT);

      it("exposes v8js_gc_duration", async () => {
        expect(await fetchPrometheusMetrics()).toContain("v8js_gc_duration");
      }, TEST_TIMEOUT);
    });

    describe("HTTP Metrics", () => {
      it("exposes http_server_duration", async () => {
        expect(await fetchPrometheusMetrics()).toContain("http_server_duration");
      }, TEST_TIMEOUT);

      it("exposes http_server_duration_sum", async () => {
        expect(await fetchPrometheusMetrics()).toContain("http_server_duration_sum");
      }, TEST_TIMEOUT);

      it("exposes http_server_duration_count", async () => {
        expect(await fetchPrometheusMetrics()).toContain("http_server_duration_count");
      }, TEST_TIMEOUT);

      it("exposes http_server_duration_bucket", async () => {
        expect(await fetchPrometheusMetrics()).toContain("http_server_duration_bucket");
      }, TEST_TIMEOUT);

      it("records status code 200", async () => {
        expect(await fetchPrometheusMetrics()).toContain('http_status_code="200"');
      }, TEST_TIMEOUT);

      it("records status code 404", async () => {
        expect(await fetchPrometheusMetrics()).toContain('http_status_code="404"');
      }, TEST_TIMEOUT);

      it("records HTTP GET method", async () => {
        expect(await fetchPrometheusMetrics()).toContain('http_method="GET"');
      }, TEST_TIMEOUT);
    });

    describe("Application Metrics", () => {
      it("exposes http_server_requests_total", async () => {
        expect(await fetchPrometheusMetrics()).toContain("http_server_requests_total");
      }, TEST_TIMEOUT);

      it("exposes http_server_active_requests", async () => {
        expect(await fetchPrometheusMetrics()).toContain("http_server_active_requests");
      }, TEST_TIMEOUT);

      it("exposes login_auth_attempts_total", async () => {
        expect(await fetchPrometheusMetrics()).toContain("login_auth_attempts_total");
      }, TEST_TIMEOUT);

      it("exposes login_session_creation_duration", async () => {
        expect(await fetchPrometheusMetrics()).toContain("login_session_creation_duration");
      }, TEST_TIMEOUT);
    });

    describe("Resource Attributes", () => {
      it("contains service.name", () => {
        expect(readMetrics()).toContain('"zitadel-login-test"');
      }, TEST_TIMEOUT);

      it("contains service.version", () => {
        expect(readMetrics()).toContain('"service.version"');
      }, TEST_TIMEOUT);

      it("contains deployment.environment", () => {
        expect(readMetrics()).toContain('"deployment.environment"');
      }, TEST_TIMEOUT);

      it("contains process.runtime.name", () => {
        expect(readMetrics()).toContain('"nodejs"');
      }, TEST_TIMEOUT);

      it("contains process.runtime.version", () => {
        expect(readMetrics()).toContain('"process.runtime.version"');
      }, TEST_TIMEOUT);

      it("contains cloud.provider", () => {
        expect(readMetrics()).toContain('"gcp"');
      }, TEST_TIMEOUT);

      it("contains cloud.platform", () => {
        expect(readMetrics()).toContain('"gcp_cloud_run"');
      }, TEST_TIMEOUT);

      it("contains cloud.account.id", () => {
        expect(readMetrics()).toContain('"test-project-123"');
      }, TEST_TIMEOUT);

      it("contains cloud.region", () => {
        expect(readMetrics()).toContain('"us-central1"');
      }, TEST_TIMEOUT);

      it("contains cloud.resource_id", () => {
        expect(readMetrics()).toContain('"cloud.resource_id"');
      }, TEST_TIMEOUT);

      it("contains faas.name", () => {
        expect(readMetrics()).toContain('"faas.name"');
      }, TEST_TIMEOUT);

      it("contains faas.version", () => {
        expect(readMetrics()).toContain('"zitadel-login-00001-abc"');
      }, TEST_TIMEOUT);

      it("contains faas.instance", () => {
        expect(readMetrics()).toContain('"faas.instance"');
      }, TEST_TIMEOUT);

      it("contains container.id when available", () => {
        const metrics = readMetrics();
        if (metrics.includes('"container.id"')) {
          expect(metrics).toMatch(/"container\.id"[^"]*"[a-f0-9]{64}"/);
        }
      }, TEST_TIMEOUT);
    });
  });

  describe("Logs", () => {
    describe("OTLP Export", () => {
      it("exports logs to collector", () => {
        expect(readLogs().length).toBeGreaterThan(0);
      }, TEST_TIMEOUT);

      it("contains timestamps", () => {
        expect(readLogs()).toContain("timeUnixNano");
      }, TEST_TIMEOUT);

      it("contains severity", () => {
        expect(readLogs()).toContain("severityNumber");
      }, TEST_TIMEOUT);

      it("contains body", () => {
        expect(readLogs()).toContain("body");
      }, TEST_TIMEOUT);
    });

    describe("Resource Attributes", () => {
      it("contains service.name", () => {
        expect(readLogs()).toContain('"zitadel-login-test"');
      }, TEST_TIMEOUT);

      it("contains service.version", () => {
        expect(readLogs()).toContain('"service.version"');
      }, TEST_TIMEOUT);

      it("contains deployment.environment", () => {
        expect(readLogs()).toContain('"deployment.environment"');
      }, TEST_TIMEOUT);

      it("contains process.runtime.name", () => {
        expect(readLogs()).toContain('"nodejs"');
      }, TEST_TIMEOUT);

      it("contains process.runtime.version", () => {
        expect(readLogs()).toContain('"process.runtime.version"');
      }, TEST_TIMEOUT);

      it("contains cloud.provider", () => {
        expect(readLogs()).toContain('"gcp"');
      }, TEST_TIMEOUT);

      it("contains cloud.platform", () => {
        expect(readLogs()).toContain('"gcp_cloud_run"');
      }, TEST_TIMEOUT);

      it("contains cloud.account.id", () => {
        expect(readLogs()).toContain('"test-project-123"');
      }, TEST_TIMEOUT);

      it("contains cloud.region", () => {
        expect(readLogs()).toContain('"us-central1"');
      }, TEST_TIMEOUT);

      it("contains cloud.resource_id", () => {
        expect(readLogs()).toContain('"cloud.resource_id"');
      }, TEST_TIMEOUT);

      it("contains faas.name", () => {
        expect(readLogs()).toContain('"faas.name"');
      }, TEST_TIMEOUT);

      it("contains faas.version", () => {
        expect(readLogs()).toContain('"zitadel-login-00001-abc"');
      }, TEST_TIMEOUT);

      it("contains faas.instance", () => {
        expect(readLogs()).toContain('"faas.instance"');
      }, TEST_TIMEOUT);

      it("contains container.id when available", () => {
        const logs = readLogs();
        if (logs.includes('"container.id"')) {
          expect(logs).toMatch(/"container\.id"[^"]*"[a-f0-9]{64}"/);
        }
      }, TEST_TIMEOUT);
    });
  });
});
