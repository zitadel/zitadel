import { describe, it, expect, beforeAll, afterAll } from "vitest";
import { GenericContainer, type StartedTestContainer, Network, Wait } from "testcontainers";
import type { StartedNetwork } from "testcontainers";
import * as fs from "fs";
import * as path from "path";

const TEST_DIR = path.dirname(new URL(import.meta.url).pathname);
const OUTPUT_DIR = path.join(TEST_DIR, "output");
const LOGIN_APP_DIR = path.join(TEST_DIR, "../..");

const DOCKER_TIMEOUT = 180000;
const TEST_TIMEOUT = 30000;

const LOGIN_IMAGE_TAG = "zitadel-login-otel-test:latest";

function readTraces(): string {
  try {
    return fs.readFileSync(path.join(OUTPUT_DIR, "traces.json"), "utf-8");
  } catch {
    return "";
  }
}

function readMetrics(): string {
  try {
    return fs.readFileSync(path.join(OUTPUT_DIR, "metrics.json"), "utf-8");
  } catch {
    return "";
  }
}

function readLogs(): string {
  try {
    return fs.readFileSync(path.join(OUTPUT_DIR, "logs.json"), "utf-8");
  } catch {
    return "";
  }
}

function readCapturedHeaders(): Array<{ traceparent: string | null; url: string; method: string }> {
  try {
    return JSON.parse(fs.readFileSync(path.join(OUTPUT_DIR, "captured-headers.json"), "utf-8"));
  } catch {
    return [];
  }
}

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

async function getContainerLogs(container: StartedTestContainer, timeoutMs = 2000): Promise<string> {
  const stream = await container.logs();
  return new Promise((resolve) => {
    let output = "";
    const timeout = setTimeout(() => {
      stream.destroy();
      resolve(output);
    }, timeoutMs);
    stream.on("data", (chunk: Buffer) => {
      output += chunk.toString();
    });
    stream.on("end", () => {
      clearTimeout(timeout);
      resolve(output);
    });
    stream.on("error", () => {
      clearTimeout(timeout);
      resolve(output);
    });
  });
}

describe("OpenTelemetry Integration", () => {
  let network: StartedNetwork;
  let otelCollector: StartedTestContainer;
  let mockZitadel: StartedTestContainer;
  let loginApp: StartedTestContainer;
  let appUrl: string;
  let prometheusUrl: string;

  beforeAll(async () => {
    if (fs.existsSync(OUTPUT_DIR)) {
      fs.rmSync(OUTPUT_DIR, { recursive: true });
    }
    fs.mkdirSync(OUTPUT_DIR, { recursive: true });

    fs.writeFileSync(path.join(OUTPUT_DIR, "traces.json"), "");
    fs.writeFileSync(path.join(OUTPUT_DIR, "logs.json"), "");
    fs.writeFileSync(path.join(OUTPUT_DIR, "metrics.json"), "");

    network = await new Network().start();

    otelCollector = await new GenericContainer("otel/opentelemetry-collector-contrib:0.145.0")
      .withNetwork(network)
      .withNetworkAliases("otel-collector")
      .withUser("0:0")
      .withCommand(["--config=/etc/otel-collector-config.yaml"])
      .withCopyFilesToContainer([
        { source: path.join(TEST_DIR, "otel-collector-config.yaml"), target: "/etc/otel-collector-config.yaml" },
      ])
      .withBindMounts([{ source: OUTPUT_DIR, target: "/tmp/otel" }])
      .withExposedPorts(4318, 8888, 13133)
      .withWaitStrategy(Wait.forHttp("/", 13133))
      .start();

    mockZitadel = await new GenericContainer("node:20-alpine")
      .withNetwork(network)
      .withNetworkAliases("mock-zitadel")
      .withCopyFilesToContainer([
        { source: path.join(TEST_DIR, "mock-zitadel-server.js"), target: "/app/server.js" },
      ])
      .withBindMounts([{ source: OUTPUT_DIR, target: "/tmp/otel" }])
      .withEnvironment({
        PORT: "8080",
        OUTPUT_DIR: "/tmp/otel",
      })
      .withCommand(["node", "/app/server.js"])
      .withExposedPorts(8080, 7432)
      .withWaitStrategy(Wait.forLogMessage("Mock Zitadel HTTP/2 server listening"))
      .start();

    await GenericContainer.fromDockerfile(LOGIN_APP_DIR).build(LOGIN_IMAGE_TAG);

    loginApp = await new GenericContainer(LOGIN_IMAGE_TAG)
      .withNetwork(network)
      .withExposedPorts(3000, 9464)
      .withEnvironment({
        ZITADEL_API_URL: "http://mock-zitadel:8080",
        ZITADEL_SERVICE_USER_TOKEN: "test-token-for-otel-integration-tests",
        OTEL_SERVICE_NAME: "zitadel-login-test",
        OTEL_EXPORTER_OTLP_ENDPOINT: "http://otel-collector:4318",
        OTEL_EXPORTER_OTLP_PROTOCOL: "http/protobuf",
        OTEL_EXPORTER_OTLP_HEADERS: "x-test-header=test-value-123,x-another-header=another-value",
        OTEL_LOG_LEVEL: "info",
        OTEL_METRICS_EXPORTER: "otlp,prometheus",
        OTEL_LOGS_EXPORTER: "otlp",
        OTEL_EXPORTER_PROMETHEUS_PORT: "9464",
        OTEL_BSP_SCHEDULE_DELAY: "1000",
        OTEL_METRIC_EXPORT_INTERVAL: "5000",
        DEBUG: "true",
      })
      .withStartupTimeout(120000)
      .withWaitStrategy(Wait.forHttp("/ui/v2/login/healthy", 3000))
      .start();

    appUrl = `http://${loginApp.getHost()}:${loginApp.getMappedPort(3000)}`;
    prometheusUrl = `http://${loginApp.getHost()}:${loginApp.getMappedPort(9464)}/metrics`;

    console.log(`[PORTS] APP_URL: ${appUrl}`);
    console.log(`[PORTS] PROMETHEUS_URL: ${prometheusUrl}`);

    await fetch(`${appUrl}/ui/v2/login/otel-test`);
    await fetch(`${appUrl}/ui/v2/login/loginname`, { redirect: "manual" });
    await fetch(`${appUrl}/ui/v2/login/healthy`);
    await fetch(`${appUrl}/ui/v2/login/this-does-not-exist-404`);

    await new Promise((r) => setTimeout(r, 3000));

    const tracesReady = await waitForFile(path.join(OUTPUT_DIR, "traces.json"), 60, 1000);
    if (!tracesReady) throw new Error("Traces file not found");
    const logsReady = await waitForFile(path.join(OUTPUT_DIR, "logs.json"), 60, 1000);
    if (!logsReady) throw new Error("Logs file not found");
    const metricsReady = await waitForFile(path.join(OUTPUT_DIR, "metrics.json"), 60, 1000);
    if (!metricsReady) throw new Error("Metrics file not found");
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
        const tracesContent = readTraces();
        const traceObjects = tracesContent
          .split("\n")
          .filter((line) => line.trim())
          .map((line) => JSON.parse(line));

        const allResourceSpans = traceObjects.flatMap((t) => t.resourceSpans || []);

        const httpSpans = allResourceSpans
          .flatMap((rs: { scopeSpans: Array<{ scope: { name: string }; spans: unknown[] }> }) => rs.scopeSpans)
          .filter((ss: { scope: { name: string } }) => ss.scope.name === "@opentelemetry/instrumentation-http")
          .flatMap((ss: { spans: unknown[] }) => ss.spans);
        const healthySpans = httpSpans.filter(
          (s: { attributes?: Array<{ key: string; value: { stringValue?: string } }> }) =>
            s.attributes?.some(
              (a: { key: string; value: { stringValue?: string } }) =>
                a.key === "http.target" && a.value.stringValue?.includes("/healthy"),
            ),
        );
        expect(healthySpans).toHaveLength(0);
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
      it("propagates traceparent to gRPC", () => {
        const requests = readCapturedHeaders();
        const grpcCalls = requests.filter((r) => r.method === "POST" && r.traceparent !== null);
        expect(grpcCalls.length).toBeGreaterThan(0);
      }, TEST_TIMEOUT);

      it("uses W3C traceparent format", () => {
        const requests = readCapturedHeaders();
        const grpcCalls = requests.filter((r) => r.traceparent !== null);
        grpcCalls.forEach((r) => expect(r.traceparent).toMatch(/^[0-9a-f]{2}-[0-9a-f]{32}-[0-9a-f]{16}-[0-9a-f]{2}$/));
      }, TEST_TIMEOUT);

      it("uses traceparent version 00", () => {
        const requests = readCapturedHeaders();
        const grpcCalls = requests.filter((r) => r.traceparent !== null);
        grpcCalls.forEach((r) => expect(r.traceparent).toMatch(/^00-/));
      }, TEST_TIMEOUT);

      it("sets sampled flag to 01", () => {
        const requests = readCapturedHeaders();
        const grpcCalls = requests.filter((r) => r.traceparent !== null);
        grpcCalls.forEach((r) => expect(r.traceparent).toMatch(/-01$/));
      }, TEST_TIMEOUT);

      it("maintains trace ID consistency", () => {
        const requests = readCapturedHeaders();
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
        const response = await fetch(prometheusUrl);
        expect(response.status).toBe(200);
      }, TEST_TIMEOUT);

      it("returns text/plain content type", async () => {
        const response = await fetch(prometheusUrl);
        expect(response.headers.get("content-type")).toContain("text/plain");
      }, TEST_TIMEOUT);
    });

    describe("Runtime Metrics", () => {
      it("exposes nodejs_eventloop_utilization", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("nodejs_eventloop_utilization");
      }, TEST_TIMEOUT);

      it("exposes nodejs_eventloop_delay", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("nodejs_eventloop_delay");
      }, TEST_TIMEOUT);

      it("exposes v8js_memory_heap_used", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("v8js_memory_heap_used");
      }, TEST_TIMEOUT);

      it("exposes v8js_memory_heap_limit", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("v8js_memory_heap_limit");
      }, TEST_TIMEOUT);

      it("exposes v8js_gc_duration", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("v8js_gc_duration");
      }, TEST_TIMEOUT);
    });

    describe("HTTP Metrics", () => {
      it("exposes http_server_duration", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("http_server_duration");
      }, TEST_TIMEOUT);

      it("exposes http_server_duration_sum", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("http_server_duration_sum");
      }, TEST_TIMEOUT);

      it("exposes http_server_duration_count", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("http_server_duration_count");
      }, TEST_TIMEOUT);

      it("exposes http_server_duration_bucket", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("http_server_duration_bucket");
      }, TEST_TIMEOUT);

      it("records status code 200", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain('http_status_code="200"');
      }, TEST_TIMEOUT);

      it("records status code 404", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain('http_status_code="404"');
      }, TEST_TIMEOUT);

      it("records HTTP GET method", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain('http_method="GET"');
      }, TEST_TIMEOUT);
    });

    describe("Application Metrics", () => {
      it("exposes http_server_requests_total", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("http_server_requests_total");
      }, TEST_TIMEOUT);

      it("exposes http_server_active_requests", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("http_server_active_requests");
      }, TEST_TIMEOUT);

      it("exposes login_auth_attempts_total", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("login_auth_attempts_total");
      }, TEST_TIMEOUT);

      it("exposes login_auth_successes_total", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("login_auth_successes_total");
      }, TEST_TIMEOUT);

      it("exposes login_session_creation_duration", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("login_session_creation_duration");
      }, TEST_TIMEOUT);

      it("exposes http_server_response_status_total", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("http_server_response_status_total");
      }, TEST_TIMEOUT);

      it("exposes http_server_request_duration_seconds", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("http_server_request_duration_seconds");
      }, TEST_TIMEOUT);
    });

    describe("Custom Test Metrics", () => {
      it("exposes otel_test_requests_total", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("otel_test_requests_total");
      }, TEST_TIMEOUT);

      it("exposes otel_test_duration_seconds", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain("otel_test_duration_seconds");
      }, TEST_TIMEOUT);
    });

    describe("Metric Labels", () => {
      it("includes method label on auth metrics", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain('method="test"');
      }, TEST_TIMEOUT);

      it("includes organization label on auth metrics", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain('organization="test-org"');
      }, TEST_TIMEOUT);

      it("includes route label on request metrics", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain('route="/otel-test"');
      }, TEST_TIMEOUT);

      it("includes success label on session metrics", async () => {
        const response = await fetch(prometheusUrl);
        expect(await response.text()).toContain('success="true"');
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

describe("OpenTelemetry Disabled", () => {
  let network: StartedNetwork;
  let mockZitadel: StartedTestContainer;
  let loginApp: StartedTestContainer;
  let appUrl: string;

  beforeAll(async () => {
    network = await new Network().start();

    mockZitadel = await new GenericContainer("node:20-alpine")
      .withNetwork(network)
      .withNetworkAliases("mock-zitadel")
      .withCopyFilesToContainer([
        { source: path.join(TEST_DIR, "mock-zitadel-server.js"), target: "/app/server.js" },
      ])
      .withEnvironment({
        PORT: "8080",
      })
      .withCommand(["node", "/app/server.js"])
      .withExposedPorts(8080, 7432)
      .withWaitStrategy(Wait.forLogMessage("Mock Zitadel HTTP/2 server listening"))
      .start();

    loginApp = await new GenericContainer(LOGIN_IMAGE_TAG)
      .withNetwork(network)
      .withExposedPorts(3000, 9464)
      .withEnvironment({
        ZITADEL_API_URL: "http://mock-zitadel:8080",
        ZITADEL_SERVICE_USER_TOKEN: "test-token",
        OTEL_SDK_DISABLED: "true",
      })
      .withStartupTimeout(120000)
      .withWaitStrategy(Wait.forHttp("/ui/v2/login/healthy", 3000))
      .start();

    appUrl = `http://${loginApp.getHost()}:${loginApp.getMappedPort(3000)}`;
    console.log(`[OTEL DISABLED] APP_URL: ${appUrl}`);
  }, DOCKER_TIMEOUT);


  describe("Application Health", () => {
    it("starts successfully with OTEL_SDK_DISABLED=true", async () => {
      const response = await fetch(`${appUrl}/ui/v2/login/healthy`);
      expect(response.ok).toBe(true);
    }, TEST_TIMEOUT);

    it("responds to health checks", async () => {
      const response = await fetch(`${appUrl}/ui/v2/login/healthy`);
      expect(response.status).toBe(200);
    }, TEST_TIMEOUT);
  });

  describe("Metrics Endpoint", () => {
    it("metrics endpoint is not available when OTEL is disabled", async () => {
      const metricsUrl = `http://${loginApp.getHost()}:${loginApp.getMappedPort(9464)}/metrics`;
      try {
        await fetch(metricsUrl);
        expect.fail("Metrics endpoint should not be available");
      } catch (error) {
        expect(error).toBeDefined();
      }
    }, TEST_TIMEOUT);
  });

  describe("Logging", () => {
    it("outputs logs to console via Winston", async () => {
      await fetch(`${appUrl}/ui/v2/login/healthy`);
      await new Promise((resolve) => setTimeout(resolve, 1000));
      const logOutput = await getContainerLogs(loginApp);
      expect(logOutput.length).toBeGreaterThan(0);
    }, TEST_TIMEOUT);

    it("does not contain OTEL SDK initialization logs", async () => {
      const logOutput = await getContainerLogs(loginApp);
      expect(logOutput).not.toContain("OpenTelemetry SDK started");
      expect(logOutput).not.toContain("Exporter");
    }, TEST_TIMEOUT);
  });
});

describe("Log Level Configuration", () => {
  describe("OTEL Disabled", () => {
    describe("LOG_LEVEL=debug", () => {
      let network: StartedNetwork;
      let mockZitadel: StartedTestContainer;
      let loginApp: StartedTestContainer;
      let appUrl: string;

      beforeAll(async () => {
        network = await new Network().start();

        mockZitadel = await new GenericContainer("node:20-alpine")
          .withNetwork(network)
          .withNetworkAliases("mock-zitadel")
          .withCopyFilesToContainer([
            { source: path.join(TEST_DIR, "mock-zitadel-server.js"), target: "/app/server.js" },
          ])
          .withEnvironment({
            PORT: "8080",
          })
          .withCommand(["node", "/app/server.js"])
          .withExposedPorts(8080, 7432)
          .withWaitStrategy(Wait.forLogMessage("Mock Zitadel HTTP/2 server listening"))
          .start();

        loginApp = await new GenericContainer(LOGIN_IMAGE_TAG)
          .withNetwork(network)
          .withExposedPorts(3000)
          .withEnvironment({
            ZITADEL_API_URL: "http://mock-zitadel:8080",
            ZITADEL_SERVICE_USER_TOKEN: "test-token",
            OTEL_SDK_DISABLED: "true",
            LOG_LEVEL: "debug",
          })
          .withStartupTimeout(120000)
          .withWaitStrategy(Wait.forHttp("/ui/v2/login/healthy", 3000))
          .start();

        appUrl = `http://${loginApp.getHost()}:${loginApp.getMappedPort(3000)}`;
        console.log(`[OTEL_DISABLED + LOG_LEVEL=debug] APP_URL: ${appUrl}`);

        await fetch(`${appUrl}/ui/v2/login/healthy`);
        await fetch(`${appUrl}/ui/v2/login/otel-test`);
        await new Promise((r) => setTimeout(r, 1000));
      }, DOCKER_TIMEOUT);

    
      it("starts successfully", async () => {
        const response = await fetch(`${appUrl}/ui/v2/login/healthy`);
        expect(response.ok).toBe(true);
      }, TEST_TIMEOUT);

      it("includes info level logs in console", async () => {
        const logOutput = await getContainerLogs(loginApp);
        expect(logOutput).toContain('"level":"info"');
      }, TEST_TIMEOUT);
    });

    describe("LOG_LEVEL=warn", () => {
      let network: StartedNetwork;
      let mockZitadel: StartedTestContainer;
      let loginApp: StartedTestContainer;
      let appUrl: string;

      beforeAll(async () => {
        network = await new Network().start();

        mockZitadel = await new GenericContainer("node:20-alpine")
          .withNetwork(network)
          .withNetworkAliases("mock-zitadel")
          .withCopyFilesToContainer([
            { source: path.join(TEST_DIR, "mock-zitadel-server.js"), target: "/app/server.js" },
          ])
          .withEnvironment({
            PORT: "8080",
          })
          .withCommand(["node", "/app/server.js"])
          .withExposedPorts(8080, 7432)
          .withWaitStrategy(Wait.forLogMessage("Mock Zitadel HTTP/2 server listening"))
          .start();

        loginApp = await new GenericContainer(LOGIN_IMAGE_TAG)
          .withNetwork(network)
          .withExposedPorts(3000)
          .withEnvironment({
            ZITADEL_API_URL: "http://mock-zitadel:8080",
            ZITADEL_SERVICE_USER_TOKEN: "test-token",
            OTEL_SDK_DISABLED: "true",
            LOG_LEVEL: "warn",
          })
          .withStartupTimeout(120000)
          .withWaitStrategy(Wait.forHttp("/ui/v2/login/healthy", 3000))
          .start();

        appUrl = `http://${loginApp.getHost()}:${loginApp.getMappedPort(3000)}`;
        console.log(`[OTEL_DISABLED + LOG_LEVEL=warn] APP_URL: ${appUrl}`);

        await fetch(`${appUrl}/ui/v2/login/healthy`);
        await new Promise((r) => setTimeout(r, 1000));
      }, DOCKER_TIMEOUT);

    
      it("starts successfully", async () => {
        const response = await fetch(`${appUrl}/ui/v2/login/healthy`);
        expect(response.ok).toBe(true);
      }, TEST_TIMEOUT);

      it("does not include info level logs in console", async () => {
        const logOutput = await getContainerLogs(loginApp);
        expect(logOutput).not.toContain('"level":"info"');
      }, TEST_TIMEOUT);

      it("does not include debug level logs in console", async () => {
        const logOutput = await getContainerLogs(loginApp);
        expect(logOutput).not.toContain('"level":"debug"');
      }, TEST_TIMEOUT);
    });
  });

  describe("OTEL Enabled", () => {
    const OTEL_LOG_OUTPUT_DIR = path.join(OUTPUT_DIR, "loglevel");

    describe("LOG_LEVEL=debug", () => {
      let network: StartedNetwork;
      let otelCollector: StartedTestContainer;
      let mockZitadel: StartedTestContainer;
      let loginApp: StartedTestContainer;
      let appUrl: string;
      let logOutputDir: string;

      beforeAll(async () => {
        logOutputDir = path.join(OTEL_LOG_OUTPUT_DIR, "debug");
        if (fs.existsSync(logOutputDir)) {
          fs.rmSync(logOutputDir, { recursive: true });
        }
        fs.mkdirSync(logOutputDir, { recursive: true });
        fs.writeFileSync(path.join(logOutputDir, "logs.json"), "");

        network = await new Network().start();

        otelCollector = await new GenericContainer("otel/opentelemetry-collector-contrib:0.145.0")
          .withNetwork(network)
          .withNetworkAliases("otel-collector")
          .withUser("0:0")
          .withCommand(["--config=/etc/otel-collector-config.yaml"])
          .withCopyFilesToContainer([
            { source: path.join(TEST_DIR, "otel-collector-config.yaml"), target: "/etc/otel-collector-config.yaml" },
          ])
          .withBindMounts([{ source: logOutputDir, target: "/tmp/otel" }])
          .withExposedPorts(4318, 13133)
          .withWaitStrategy(Wait.forHttp("/", 13133))
          .start();

        mockZitadel = await new GenericContainer("node:20-alpine")
          .withNetwork(network)
          .withNetworkAliases("mock-zitadel")
          .withCopyFilesToContainer([
            { source: path.join(TEST_DIR, "mock-zitadel-server.js"), target: "/app/server.js" },
          ])
          .withEnvironment({
            PORT: "8080",
          })
          .withCommand(["node", "/app/server.js"])
          .withExposedPorts(8080, 7432)
          .withWaitStrategy(Wait.forLogMessage("Mock Zitadel HTTP/2 server listening"))
          .start();

        loginApp = await new GenericContainer(LOGIN_IMAGE_TAG)
          .withNetwork(network)
          .withExposedPorts(3000)
          .withEnvironment({
            ZITADEL_API_URL: "http://mock-zitadel:8080",
            ZITADEL_SERVICE_USER_TOKEN: "test-token",
            OTEL_SERVICE_NAME: "zitadel-login-loglevel-debug",
            OTEL_EXPORTER_OTLP_ENDPOINT: "http://otel-collector:4318",
            OTEL_EXPORTER_OTLP_PROTOCOL: "http/protobuf",
            OTEL_LOG_LEVEL: "debug",
            OTEL_LOGS_EXPORTER: "otlp",
            OTEL_BSP_SCHEDULE_DELAY: "1000",
          })
          .withStartupTimeout(120000)
          .withWaitStrategy(Wait.forHttp("/ui/v2/login/healthy", 3000))
          .start();

        appUrl = `http://${loginApp.getHost()}:${loginApp.getMappedPort(3000)}`;
        console.log(`[OTEL_ENABLED + LOG_LEVEL=debug] APP_URL: ${appUrl}`);

        await fetch(`${appUrl}/ui/v2/login/healthy`);
        await fetch(`${appUrl}/ui/v2/login/otel-test`);
        await new Promise((r) => setTimeout(r, 3000));

        await waitForFile(path.join(logOutputDir, "logs.json"), 30, 1000);
      }, DOCKER_TIMEOUT);

    
      it("starts successfully", async () => {
        const response = await fetch(`${appUrl}/ui/v2/login/healthy`);
        expect(response.ok).toBe(true);
      }, TEST_TIMEOUT);

      it("exports logs to OTLP collector", () => {
        const logs = fs.readFileSync(path.join(logOutputDir, "logs.json"), "utf-8");
        expect(logs.length).toBeGreaterThan(0);
      }, TEST_TIMEOUT);

      it("includes info level logs in OTLP export", () => {
        const logs = fs.readFileSync(path.join(logOutputDir, "logs.json"), "utf-8");
        expect(logs).toContain('"severityText":"info"');
      }, TEST_TIMEOUT);

      it("includes info level logs in console", async () => {
        const logOutput = await getContainerLogs(loginApp);
        expect(logOutput).toContain('"level":"info"');
      }, TEST_TIMEOUT);
    });

    describe("LOG_LEVEL=warn", () => {
      let network: StartedNetwork;
      let otelCollector: StartedTestContainer;
      let mockZitadel: StartedTestContainer;
      let loginApp: StartedTestContainer;
      let appUrl: string;
      let logOutputDir: string;

      beforeAll(async () => {
        logOutputDir = path.join(OTEL_LOG_OUTPUT_DIR, "warn");
        if (fs.existsSync(logOutputDir)) {
          fs.rmSync(logOutputDir, { recursive: true });
        }
        fs.mkdirSync(logOutputDir, { recursive: true });
        fs.writeFileSync(path.join(logOutputDir, "logs.json"), "");

        network = await new Network().start();

        otelCollector = await new GenericContainer("otel/opentelemetry-collector-contrib:0.145.0")
          .withNetwork(network)
          .withNetworkAliases("otel-collector")
          .withUser("0:0")
          .withCommand(["--config=/etc/otel-collector-config.yaml"])
          .withCopyFilesToContainer([
            { source: path.join(TEST_DIR, "otel-collector-config.yaml"), target: "/etc/otel-collector-config.yaml" },
          ])
          .withBindMounts([{ source: logOutputDir, target: "/tmp/otel" }])
          .withExposedPorts(4318, 13133)
          .withWaitStrategy(Wait.forHttp("/", 13133))
          .start();

        mockZitadel = await new GenericContainer("node:20-alpine")
          .withNetwork(network)
          .withNetworkAliases("mock-zitadel")
          .withCopyFilesToContainer([
            { source: path.join(TEST_DIR, "mock-zitadel-server.js"), target: "/app/server.js" },
          ])
          .withEnvironment({
            PORT: "8080",
          })
          .withCommand(["node", "/app/server.js"])
          .withExposedPorts(8080, 7432)
          .withWaitStrategy(Wait.forLogMessage("Mock Zitadel HTTP/2 server listening"))
          .start();

        loginApp = await new GenericContainer(LOGIN_IMAGE_TAG)
          .withNetwork(network)
          .withExposedPorts(3000)
          .withEnvironment({
            ZITADEL_API_URL: "http://mock-zitadel:8080",
            ZITADEL_SERVICE_USER_TOKEN: "test-token",
            OTEL_SERVICE_NAME: "zitadel-login-loglevel-warn",
            OTEL_EXPORTER_OTLP_ENDPOINT: "http://otel-collector:4318",
            OTEL_EXPORTER_OTLP_PROTOCOL: "http/protobuf",
            OTEL_LOG_LEVEL: "warn",
            OTEL_LOGS_EXPORTER: "otlp",
            OTEL_BSP_SCHEDULE_DELAY: "1000",
          })
          .withStartupTimeout(120000)
          .withWaitStrategy(Wait.forHttp("/ui/v2/login/healthy", 3000))
          .start();

        appUrl = `http://${loginApp.getHost()}:${loginApp.getMappedPort(3000)}`;
        console.log(`[OTEL_ENABLED + LOG_LEVEL=warn] APP_URL: ${appUrl}`);

        await fetch(`${appUrl}/ui/v2/login/healthy`);
        await fetch(`${appUrl}/ui/v2/login/otel-test`);
        await new Promise((r) => setTimeout(r, 3000));
      }, DOCKER_TIMEOUT);

    
      it("starts successfully", async () => {
        const response = await fetch(`${appUrl}/ui/v2/login/healthy`);
        expect(response.ok).toBe(true);
      }, TEST_TIMEOUT);

      it("does not include info level logs in console", async () => {
        const logOutput = await getContainerLogs(loginApp);
        expect(logOutput).not.toContain('"level":"info"');
      }, TEST_TIMEOUT);

      it("does not include debug level logs in console", async () => {
        const logOutput = await getContainerLogs(loginApp);
        expect(logOutput).not.toContain('"level":"debug"');
      }, TEST_TIMEOUT);

      it("does not export info level logs to OTLP", () => {
        const logsPath = path.join(logOutputDir, "logs.json");
        if (fs.existsSync(logsPath)) {
          const logs = fs.readFileSync(logsPath, "utf-8");
          expect(logs).not.toContain('"severityText":"info"');
        }
      }, TEST_TIMEOUT);

      it("does not export debug level logs to OTLP", () => {
        const logsPath = path.join(logOutputDir, "logs.json");
        if (fs.existsSync(logsPath)) {
          const logs = fs.readFileSync(logsPath, "utf-8");
          expect(logs).not.toContain('"severityText":"debug"');
        }
      }, TEST_TIMEOUT);
    });
  });
});
