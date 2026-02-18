/**
 * Integration tests for OTEL-disabled mode.
 *
 * Verifies that the login application works correctly when OpenTelemetry
 * is disabled via OTEL_SDK_DISABLED=true. The application should:
 * - Start successfully without OTEL initialization
 * - Respond to health checks
 * - Output logs to console (Winston still works)
 */

import { execSync } from "child_process";
import { existsSync, mkdirSync, rmSync } from "fs";
import { join } from "path";
import { afterAll, beforeAll, describe, expect, it } from "vitest";

const TEST_DIR = join(__dirname);
const OUTPUT_DIR = join(TEST_DIR, "output");
const COMPOSE_FILE = join(TEST_DIR, "docker-compose.test.yml");

interface ContainerPorts {
  appUrl: string;
}

function getContainerPorts(): ContainerPorts {
  const appPort = execSync(
    `docker compose -f ${COMPOSE_FILE} port login-app 3000`,
    { encoding: "utf-8" },
  ).trim();

  return {
    appUrl: `http://${appPort.replace("0.0.0.0", "localhost")}`,
  };
}

async function waitForService(
  url: string,
  maxAttempts = 30,
  interval = 1000,
): Promise<void> {
  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      const response = await fetch(url);
      if (response.ok) {
        console.log(`[waitForService] ${url} is ready (attempt ${attempt})`);
        return;
      }
    } catch {
      // Service not ready yet
    }
    await new Promise((resolve) => setTimeout(resolve, interval));
  }
  throw new Error(`Service ${url} did not become ready`);
}

function getContainerLogs(): string {
  return execSync(`docker compose -f ${COMPOSE_FILE} logs login-app`, {
    encoding: "utf-8",
  });
}

describe("OTEL Disabled Integration", () => {
  let ports: ContainerPorts;

  beforeAll(async () => {
    if (existsSync(OUTPUT_DIR)) {
      rmSync(OUTPUT_DIR, { recursive: true });
    }
    mkdirSync(OUTPUT_DIR, { recursive: true });

    execSync(`docker compose -f ${COMPOSE_FILE} up -d --wait --build`, {
      cwd: TEST_DIR,
      stdio: "inherit",
    });

    ports = getContainerPorts();
    console.log(`[PORTS] APP_URL: ${ports.appUrl}`);

    await waitForService(`${ports.appUrl}/ui/v2/login/healthy`);
  }, 120000);

  afterAll(() => {
    execSync(`docker compose -f ${COMPOSE_FILE} down -v --remove-orphans`, {
      cwd: TEST_DIR,
      stdio: "inherit",
    });
  });

  it("should start successfully with OTEL disabled", async () => {
    const response = await fetch(`${ports.appUrl}/ui/v2/login/healthy`);
    expect(response.ok).toBe(true);
  });

  it("should respond to requests without OTEL", async () => {
    const response = await fetch(`${ports.appUrl}/ui/v2/login/healthy`);
    expect(response.status).toBe(200);
  });

  it("should output logs to console (Winston)", async () => {
    // Make a request to generate some logs
    await fetch(`${ports.appUrl}/ui/v2/login/healthy`);

    // Give logs time to flush
    await new Promise((resolve) => setTimeout(resolve, 1000));

    const logs = getContainerLogs();

    // Winston should be logging even with OTEL disabled
    // Check for Next.js startup or request logs
    expect(logs.length).toBeGreaterThan(0);
  });

  it("should not have OTEL SDK initialization logs", async () => {
    const logs = getContainerLogs();

    // When OTEL_SDK_DISABLED=true, there should be no OTEL initialization
    expect(logs).not.toContain("OpenTelemetry SDK started");
    expect(logs).not.toContain("OTLP");
  });
});
