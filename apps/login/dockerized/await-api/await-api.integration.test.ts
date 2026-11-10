import * as path from "path";
import { fileURLToPath } from "node:url";
import type { StartedNetwork } from "testcontainers";
import { GenericContainer, Network, Wait, type StartedTestContainer } from "testcontainers";
import { afterAll, beforeAll, describe, expect, it } from "vitest";

const TEST_DIR = path.dirname(fileURLToPath(import.meta.url));
const LOGIN_APP_DIR = path.join(TEST_DIR, "../..");
const LOGIN_IMAGE_TAG = `zitadel-login-test:${process.env.GITHUB_SHA || "local"}`;

const DOCKER_TIMEOUT = 180_000;

/**
 * Integration tests for ZITADEL_API_AWAITINITIALCONN support.
 *
 * When this environment variable is set (e.g. "30s", "5m"), the login
 * container polls the ZITADEL API's /debug/ready endpoint before starting
 * the Next.js server. This replaces the wait4x init container that was
 * previously used in the Helm chart.
 *
 * The test suite uses a mock Zitadel server that returns 503 on /debug/ready
 * for a configurable duration before switching to 200. Three scenarios are
 * verified: the API becomes ready within the timeout, the timeout expires
 * before the API is ready, and the variable is not set at all.
 */
describe("ZITADEL_API_AWAITINITIALCONN", () => {
  beforeAll(async () => {
    await GenericContainer.fromDockerfile(LOGIN_APP_DIR).build(LOGIN_IMAGE_TAG);
  }, DOCKER_TIMEOUT);

  describe("when the API becomes ready before the timeout", () => {
    let network: StartedNetwork;
    let mockZitadel: StartedTestContainer;
    let container: StartedTestContainer;
    let appUrl: string;

    beforeAll(async () => {
      network = await new Network().start();

      mockZitadel = await new GenericContainer("node:20-alpine")
        .withNetwork(network)
        .withNetworkAliases("mock-zitadel")
        .withCopyFilesToContainer([
          { source: path.join(TEST_DIR, "mock-zitadel-server.js"), target: "/app/server.js" },
        ])
        .withEnvironment({ PORT: "8080", READY_AFTER_MS: "5000" })
        .withCommand(["node", "/app/server.js"])
        .withExposedPorts(8080)
        .withWaitStrategy(Wait.forLogMessage("Mock Zitadel server listening"))
        .start();

      container = await new GenericContainer(LOGIN_IMAGE_TAG)
        .withNetwork(network)
        .withExposedPorts(3000)
        .withEnvironment({
          ZITADEL_API_URL: "http://mock-zitadel:8080",
          ZITADEL_SERVICE_USER_TOKEN: "test-token",
          ZITADEL_API_AWAITINITIALCONN: "60s",
        })
        .withStartupTimeout(120_000)
        .withWaitStrategy(Wait.forHttp("/ui/v2/login/healthy", 3000))
        .start();

      appUrl = `http://${container.getHost()}:${container.getMappedPort(3000)}`;
    }, DOCKER_TIMEOUT);

    afterAll(async () => {
      await container?.stop().catch(() => {});
      await mockZitadel?.stop().catch(() => {});
      await network?.stop().catch(() => {});
    }, DOCKER_TIMEOUT);

    it("starts successfully after the API becomes ready", async () => {
      const response = await fetch(`${appUrl}/ui/v2/login/healthy`);
      expect(response.status).toBe(200);
    });
  });

  describe("when the API does not become ready before the timeout", () => {
    let network: StartedNetwork;
    let mockZitadel: StartedTestContainer;
    let container: StartedTestContainer;

    beforeAll(async () => {
      network = await new Network().start();

      mockZitadel = await new GenericContainer("node:20-alpine")
        .withNetwork(network)
        .withNetworkAliases("mock-zitadel")
        .withCopyFilesToContainer([
          { source: path.join(TEST_DIR, "mock-zitadel-server.js"), target: "/app/server.js" },
        ])
        .withEnvironment({ PORT: "8080", READY_AFTER_MS: "600000" })
        .withCommand(["node", "/app/server.js"])
        .withExposedPorts(8080)
        .withWaitStrategy(Wait.forLogMessage("Mock Zitadel server listening"))
        .start();

      container = await new GenericContainer(LOGIN_IMAGE_TAG)
        .withNetwork(network)
        .withExposedPorts(3000)
        .withEnvironment({
          ZITADEL_API_URL: "http://mock-zitadel:8080",
          ZITADEL_SERVICE_USER_TOKEN: "test-token",
          ZITADEL_API_AWAITINITIALCONN: "5s",
        })
        .withStartupTimeout(30_000)
        .withWaitStrategy(Wait.forLogMessage(/Timed out after 5s/))
        .start();
    }, DOCKER_TIMEOUT);

    afterAll(async () => {
      await container?.stop().catch(() => {});
      await mockZitadel?.stop().catch(() => {});
      await network?.stop().catch(() => {});
    }, DOCKER_TIMEOUT);

    it("logs the timeout error before exiting", async () => {
      const stream = await container.logs();
      const logs = await new Promise<string>((resolve) => {
        let output = "";
        stream.on("data", (chunk: Buffer) => { output += chunk.toString(); });
        stream.on("end", () => resolve(output));
        setTimeout(() => { stream.destroy(); resolve(output); }, 5_000);
      });
      expect(logs).toContain("Timed out after 5s waiting for API readiness");
    });
  });

  describe("when the variable is not set", () => {
    let network: StartedNetwork;
    let mockZitadel: StartedTestContainer;
    let container: StartedTestContainer;
    let appUrl: string;

    beforeAll(async () => {
      network = await new Network().start();

      mockZitadel = await new GenericContainer("node:20-alpine")
        .withNetwork(network)
        .withNetworkAliases("mock-zitadel")
        .withCopyFilesToContainer([
          { source: path.join(TEST_DIR, "mock-zitadel-server.js"), target: "/app/server.js" },
        ])
        .withEnvironment({ PORT: "8080", READY_AFTER_MS: "0" })
        .withCommand(["node", "/app/server.js"])
        .withExposedPorts(8080)
        .withWaitStrategy(Wait.forLogMessage("Mock Zitadel server listening"))
        .start();

      container = await new GenericContainer(LOGIN_IMAGE_TAG)
        .withNetwork(network)
        .withExposedPorts(3000)
        .withEnvironment({
          ZITADEL_API_URL: "http://mock-zitadel:8080",
          ZITADEL_SERVICE_USER_TOKEN: "test-token",
        })
        .withStartupTimeout(120_000)
        .withWaitStrategy(Wait.forHttp("/ui/v2/login/healthy", 3000))
        .start();

      appUrl = `http://${container.getHost()}:${container.getMappedPort(3000)}`;
    }, DOCKER_TIMEOUT);

    afterAll(async () => {
      await container?.stop().catch(() => {});
      await mockZitadel?.stop().catch(() => {});
      await network?.stop().catch(() => {});
    }, DOCKER_TIMEOUT);

    it("starts without waiting", async () => {
      const response = await fetch(`${appUrl}/ui/v2/login/healthy`);
      expect(response.status).toBe(200);
    });
  });
});
