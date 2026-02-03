import { describe, it, expect, beforeAll, afterAll } from "vitest";
import { execSync } from "child_process";
import * as fs from "fs";
import * as path from "path";

const TEST_DIR = path.dirname(new URL(import.meta.url).pathname);
const OUTPUT_DIR = path.join(TEST_DIR, "output");
const COMPOSE_FILE = path.join(TEST_DIR, "docker-compose.test.yml");

const DOCKER_TIMEOUT = 180000;
const TEST_TIMEOUT = 30000;

// Dynamic port discovery - populated in beforeAll
let APP_URL = "";
let PROMETHEUS_URL = "";
let COLLECTOR_HEALTH_URL = "";
let MOCK_ZITADEL_URL = "";

/**
 * Gets the dynamically assigned host port for a service/container port.
 * Uses `docker compose port` to discover the actual bound port.
 */
function getHostPort(service: string, containerPort: number): number {
  const output = execSync(`docker compose -f ${COMPOSE_FILE} port ${service} ${containerPort}`, {
    cwd: TEST_DIR,
    encoding: "utf-8",
  }).trim();
  // Output format: 0.0.0.0:12345 or [::]:12345
  const match = output.match(/:(\d+)$/);
  if (!match) {
    throw new Error(`Failed to parse port from: ${output}`);
  }
  return parseInt(match[1], 10);
}

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
      if (response.ok) {
        console.log(`[waitForService] ${url} is ready (attempt ${i + 1})`);
        return true;
      }
      console.log(`[waitForService] ${url} returned ${response.status} (attempt ${i + 1})`);
    } catch (error) {
      console.log(`[waitForService] ${url} failed: ${error instanceof Error ? error.message : error} (attempt ${i + 1})`);
      await new Promise((r) => setTimeout(r, delayMs));
    }
  }
  console.log(`[waitForService] ${url} timed out after ${maxAttempts} attempts`);
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
      const stat = fs.statSync(filePath);
      if (stat.size > 0) {
        console.log(`[waitForFile] ${path.basename(filePath)} ready (${stat.size} bytes, attempt ${i + 1})`);
        return true;
      }
      console.log(`[waitForFile] ${path.basename(filePath)} exists but empty (attempt ${i + 1})`);
    } catch {
      if (i === 0 || i % 10 === 0) {
        console.log(`[waitForFile] ${path.basename(filePath)} not found (attempt ${i + 1})`);
      }
    }
    await new Promise((r) => setTimeout(r, delayMs));
  }
  console.log(`[waitForFile] ${path.basename(filePath)} timed out after ${maxAttempts} attempts`);
  return false;
}

/**
 * Reads the traces.json file exported by the OTLP collector.
 * Contains all trace spans in OTLP JSON format.
 */
function readTraces(): string {
  try {
    return fs.readFileSync(path.join(OUTPUT_DIR, "traces.json"), "utf-8");
  } catch {
    return "";
  }
}

/**
 * Reads the metrics.json file exported by the OTLP collector.
 * Contains all metrics in OTLP JSON format.
 */
function readMetrics(): string {
  try {
    return fs.readFileSync(path.join(OUTPUT_DIR, "metrics.json"), "utf-8");
  } catch {
    return "";
  }
}

/**
 * Reads the logs.json file exported by the OTLP collector.
 * Contains all log records in OTLP JSON format.
 */
function readLogs(): string {
  try {
    return fs.readFileSync(path.join(OUTPUT_DIR, "logs.json"), "utf-8");
  } catch {
    return "";
  }
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

    // Clean output directory to ensure fresh data
    if (fs.existsSync(OUTPUT_DIR)) {
      fs.rmSync(OUTPUT_DIR, { recursive: true });
    }
    fs.mkdirSync(OUTPUT_DIR, { recursive: true });
    execSync(`docker compose -f ${COMPOSE_FILE} up -d --build`, { stdio: "inherit", cwd: TEST_DIR });

    // Discover dynamically assigned ports
    const appPort = getHostPort("login-app", 3000);
    const prometheusPort = getHostPort("login-app", 9464);
    const collectorHealthPort = getHostPort("otel-collector", 13133);
    const mockZitadelHealthPort = getHostPort("mock-zitadel", 7432);

    APP_URL = `http://localhost:${appPort}`;
    PROMETHEUS_URL = `http://localhost:${prometheusPort}/metrics`;
    COLLECTOR_HEALTH_URL = `http://localhost:${collectorHealthPort}`;
    MOCK_ZITADEL_URL = `http://localhost:${mockZitadelHealthPort}`;

    console.log(`[PORTS] APP_URL: ${APP_URL}`);
    console.log(`[PORTS] PROMETHEUS_URL: ${PROMETHEUS_URL}`);
    console.log(`[PORTS] COLLECTOR_HEALTH_URL: ${COLLECTOR_HEALTH_URL}`);
    console.log(`[PORTS] MOCK_ZITADEL_URL: ${MOCK_ZITADEL_URL}`);

    // Debug: show container status and logs
    try {
      const psOutput = execSync(`docker compose -f ${COMPOSE_FILE} ps -a`, { cwd: TEST_DIR, encoding: "utf-8" });
      console.log("[DEBUG] Container status:\n", psOutput);
      const logsOutput = execSync(`docker compose -f ${COMPOSE_FILE} logs otel-collector`, { cwd: TEST_DIR, encoding: "utf-8", stdio: ["pipe", "pipe", "pipe"] });
      console.log("[DEBUG] Collector logs:\n", logsOutput);
    } catch (e) {
      console.log("[DEBUG] Failed to get container info:", e instanceof Error ? e.message : e);
    }

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

    // Wait for batch processors to flush, collector to write files,
    // and Docker volume sync (can be slow on macOS)
    await new Promise((r) => setTimeout(r, 15000));

    // Wait for telemetry files to be written
    const tracesReady = await waitForFile(path.join(OUTPUT_DIR, "traces.json"));
    if (!tracesReady) throw new Error("Traces file not found");
    const logsReady = await waitForFile(path.join(OUTPUT_DIR, "logs.json"));
    if (!logsReady) throw new Error("Logs file not found");
    const metricsReady = await waitForFile(path.join(OUTPUT_DIR, "metrics.json"), 60, 1000);
    if (!metricsReady) throw new Error("Metrics file not found");
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

      it("contains custom test.attribute", () => {
        expect(readTraces()).toContain('"test.attribute"');
      }, TEST_TIMEOUT);

      it("contains custom test.timestamp", () => {
        expect(readTraces()).toContain('"test.timestamp"');
      }, TEST_TIMEOUT);
    });

    describe("Span Status", () => {
      it("contains status code OK", () => {
        expect(readTraces()).toContain('"code":1');
      }, TEST_TIMEOUT);

      it("contains status code ERROR", () => {
        expect(readTraces()).toContain('"code":2');
      }, TEST_TIMEOUT);

      it("contains error message on failed spans", () => {
        expect(readTraces()).toContain('"message"');
      }, TEST_TIMEOUT);
    });

    describe("Span Hierarchy", () => {
      it("contains parent span IDs", () => {
        expect(readTraces()).toContain('"parentSpanId"');
      }, TEST_TIMEOUT);

      it("has child spans linked to parent", () => {
        const traces = readTraces();
        const parentMatches = traces.match(/"parentSpanId":"([a-f0-9]{16})"/g) || [];
        expect(parentMatches.length).toBeGreaterThan(0);
      }, TEST_TIMEOUT);
    });

    describe("Span Filtering", () => {
      it("HTTP instrumentation excludes /healthy", () => {
        const traces = readTraces();
        // Check that no trace contains /healthy path (with dynamic port)
        expect(traces).not.toContain("/ui/v2/login/healthy");
      }, TEST_TIMEOUT);

      it("includes /otel-test in traces", () => {
        const traces = readTraces();
        expect(traces).toContain("otel-test");
      }, TEST_TIMEOUT);

      it("includes /loginname in traces", () => {
        const traces = readTraces();
        expect(traces).toContain("loginname");
      }, TEST_TIMEOUT);
    });

    describe("Span Events", () => {
      it("contains exception events", () => {
        expect(readTraces()).toContain('"name":"exception"');
      }, TEST_TIMEOUT);

      it("contains exception.type attribute", () => {
        expect(readTraces()).toContain('"exception.type"');
      }, TEST_TIMEOUT);

      it("contains exception.message attribute", () => {
        expect(readTraces()).toContain('"exception.message"');
      }, TEST_TIMEOUT);

      it("contains exception.stacktrace attribute", () => {
        expect(readTraces()).toContain('"exception.stacktrace"');
      }, TEST_TIMEOUT);
    });

    describe("gRPC Spans", () => {
      it("contains rpc.system attribute", () => {
        expect(readTraces()).toContain('"rpc.system"');
      }, TEST_TIMEOUT);

      it("contains rpc.service attribute", () => {
        expect(readTraces()).toContain('"rpc.service"');
      }, TEST_TIMEOUT);

      it("contains rpc.method attribute", () => {
        expect(readTraces()).toContain('"rpc.method"');
      }, TEST_TIMEOUT);

      it("contains server.address attribute", () => {
        expect(readTraces()).toContain('"server.address"');
      }, TEST_TIMEOUT);

      it("contains rpc.grpc.status_code attribute", () => {
        expect(readTraces()).toContain('"rpc.grpc.status_code"');
      }, TEST_TIMEOUT);

      it("identifies grpc as rpc system", () => {
        expect(readTraces()).toContain('"grpc"');
      }, TEST_TIMEOUT);
    });

    describe("HTTP Status Codes", () => {
      it("records HTTP 200 status", () => {
        expect(readTraces()).toContain('"intValue":"200"');
      }, TEST_TIMEOUT);

      it("records HTTP 404 status", () => {
        expect(readTraces()).toContain('"intValue":"404"');
      }, TEST_TIMEOUT);

      it("records OK status text", () => {
        expect(readTraces()).toContain('"OK"');
      }, TEST_TIMEOUT);

      it("records NOT FOUND status text", () => {
        expect(readTraces()).toContain('"NOT FOUND"');
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

      it("contains deployment.environment.name", () => {
        expect(readTraces()).toContain('"deployment.environment.name"');
      }, TEST_TIMEOUT);

      it("contains process.runtime.name", () => {
        expect(readTraces()).toContain('"nodejs"');
      }, TEST_TIMEOUT);

      it("contains process.runtime.version", () => {
        expect(readTraces()).toContain('"process.runtime.version"');
      }, TEST_TIMEOUT);

      it("contains container.id when running in container", () => {
        expect(readTraces()).toContain('"container.id"');
      }, TEST_TIMEOUT);

      it("contains host.name", () => {
        expect(readTraces()).toContain('"host.name"');
      }, TEST_TIMEOUT);

      it("contains host.arch", () => {
        expect(readTraces()).toContain('"host.arch"');
      }, TEST_TIMEOUT);

      it("contains process.pid", () => {
        expect(readTraces()).toContain('"process.pid"');
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

      it("exposes login_auth_successes_total", async () => {
        expect(await fetchPrometheusMetrics()).toContain("login_auth_successes_total");
      }, TEST_TIMEOUT);

      it("exposes login_session_creation_duration", async () => {
        expect(await fetchPrometheusMetrics()).toContain("login_session_creation_duration");
      }, TEST_TIMEOUT);

      it("exposes http_server_response_status_total", async () => {
        expect(await fetchPrometheusMetrics()).toContain("http_server_response_status_total");
      }, TEST_TIMEOUT);

      it("exposes http_server_request_duration_seconds", async () => {
        expect(await fetchPrometheusMetrics()).toContain("http_server_request_duration_seconds");
      }, TEST_TIMEOUT);
    });

    describe("Custom Test Metrics", () => {
      it("exposes otel_test_requests_total", async () => {
        expect(await fetchPrometheusMetrics()).toContain("otel_test_requests_total");
      }, TEST_TIMEOUT);

      it("exposes otel_test_duration_seconds", async () => {
        expect(await fetchPrometheusMetrics()).toContain("otel_test_duration_seconds");
      }, TEST_TIMEOUT);
    });

    describe("Metric Labels", () => {
      it("includes method label on auth metrics", async () => {
        expect(await fetchPrometheusMetrics()).toContain('method="test"');
      }, TEST_TIMEOUT);

      it("includes organization label on auth metrics", async () => {
        expect(await fetchPrometheusMetrics()).toContain('organization="test-org"');
      }, TEST_TIMEOUT);

      it("includes route label on request metrics", async () => {
        expect(await fetchPrometheusMetrics()).toContain('route="/otel-test"');
      }, TEST_TIMEOUT);

      it("includes success label on session metrics", async () => {
        expect(await fetchPrometheusMetrics()).toContain('success="true"');
      }, TEST_TIMEOUT);
    });

    describe("Resource Attributes", () => {
      it("contains service.name", () => {
        expect(readMetrics()).toContain('"zitadel-login-test"');
      }, TEST_TIMEOUT);

      it("contains deployment.environment.name", () => {
        expect(readMetrics()).toContain('"deployment.environment.name"');
      }, TEST_TIMEOUT);

      it("contains process.runtime.name", () => {
        expect(readMetrics()).toContain('"nodejs"');
      }, TEST_TIMEOUT);

      it("contains process.runtime.version", () => {
        expect(readMetrics()).toContain('"process.runtime.version"');
      }, TEST_TIMEOUT);

      it("contains container.id when running in container", () => {
        expect(readMetrics()).toContain('"container.id"');
      }, TEST_TIMEOUT);

      it("contains host.name", () => {
        expect(readMetrics()).toContain('"host.name"');
      }, TEST_TIMEOUT);

      it("contains host.arch", () => {
        expect(readMetrics()).toContain('"host.arch"');
      }, TEST_TIMEOUT);

      it("contains process.pid", () => {
        expect(readMetrics()).toContain('"process.pid"');
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

    describe("Log Content", () => {
      it("contains OTEL test log message", () => {
        expect(readLogs()).toContain("OTEL test endpoint called");
      }, TEST_TIMEOUT);

      it("contains test completed log message", () => {
        expect(readLogs()).toContain("OTEL test completed successfully");
      }, TEST_TIMEOUT);

      it("contains INFO severity level", () => {
        const logs = readLogs();
        expect(logs).toContain('"severityText":"info"');
      }, TEST_TIMEOUT);

      it("contains severity number 9 for INFO", () => {
        expect(readLogs()).toContain('"severityNumber":9');
      }, TEST_TIMEOUT);
    });

    describe("Log-Trace Correlation", () => {
      it("includes traceId in log records", () => {
        expect(readLogs()).toContain('"traceId"');
      }, TEST_TIMEOUT);

      it("includes spanId in log records", () => {
        expect(readLogs()).toContain('"spanId"');
      }, TEST_TIMEOUT);

      it("has valid trace ID format", () => {
        expect(readLogs()).toMatch(/"traceId":"[a-f0-9]{32}"/);
      }, TEST_TIMEOUT);

      it("has valid span ID format", () => {
        expect(readLogs()).toMatch(/"spanId":"[a-f0-9]{16}"/);
      }, TEST_TIMEOUT);

      it("includes trace_id attribute", () => {
        expect(readLogs()).toContain('"trace_id"');
      }, TEST_TIMEOUT);

      it("includes span_id attribute", () => {
        expect(readLogs()).toContain('"span_id"');
      }, TEST_TIMEOUT);

      it("includes trace_flags attribute", () => {
        expect(readLogs()).toContain('"trace_flags"');
      }, TEST_TIMEOUT);
    });

    describe("Log Attributes", () => {
      it("includes context attribute", () => {
        expect(readLogs()).toContain('"context"');
      }, TEST_TIMEOUT);

      it("includes service attribute", () => {
        expect(readLogs()).toContain('"service"');
      }, TEST_TIMEOUT);

      it("includes timestamp attribute", () => {
        expect(readLogs()).toContain('"timestamp"');
      }, TEST_TIMEOUT);

      it("includes custom test attribute", () => {
        expect(readLogs()).toContain('"test"');
      }, TEST_TIMEOUT);

      it("includes duration attribute", () => {
        expect(readLogs()).toContain('"duration"');
      }, TEST_TIMEOUT);
    });

    describe("Resource Attributes", () => {
      it("contains service.name", () => {
        expect(readLogs()).toContain('"zitadel-login-test"');
      }, TEST_TIMEOUT);

      it("contains deployment.environment.name", () => {
        expect(readLogs()).toContain('"deployment.environment.name"');
      }, TEST_TIMEOUT);

      it("contains process.runtime.name", () => {
        expect(readLogs()).toContain('"nodejs"');
      }, TEST_TIMEOUT);

      it("contains process.runtime.version", () => {
        expect(readLogs()).toContain('"process.runtime.version"');
      }, TEST_TIMEOUT);

      it("contains container.id when running in container", () => {
        expect(readLogs()).toContain('"container.id"');
      }, TEST_TIMEOUT);

      it("contains host.name", () => {
        expect(readLogs()).toContain('"host.name"');
      }, TEST_TIMEOUT);

      it("contains host.arch", () => {
        expect(readLogs()).toContain('"host.arch"');
      }, TEST_TIMEOUT);

      it("contains process.pid", () => {
        expect(readLogs()).toContain('"process.pid"');
      }, TEST_TIMEOUT);
    });
  });
});
