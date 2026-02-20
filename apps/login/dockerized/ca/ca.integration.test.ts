import { describe, it, expect, beforeAll, afterAll } from "vitest";
import {
  GenericContainer,
  type StartedTestContainer,
  Network,
  Wait,
} from "testcontainers";
import type { StartedNetwork } from "testcontainers";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";
import { generateCertificates } from "./utils/tls.ts";

const TEST_DIR = path.dirname(new URL(import.meta.url).pathname);
const LOGIN_APP_DIR = path.join(TEST_DIR, "../..");

const SMOCKER_MOCK_PORT = 443;
const SMOCKER_ADMIN_PORT = 8081;

/**
 * Integration tests for custom CA certificate support in the login application.
 *
 * These tests verify that the Dockerfile correctly configures Node.js to trust
 * custom Certificate Authorities via the SSL_CERT_FILE environment variable and
 * the NODE_OPTIONS=--use-openssl-ca flag.
 *
 * We use Smocker as a TLS server and test CA loading by making curl requests
 * from inside the container - this verifies TLS at the transport layer without
 * needing gRPC/HTTP2 support.
 */
const LOGIN_IMAGE_TAG = `zitadel-login-ca-test:${Date.now()}`;

async function testTlsConnection(
  container: StartedTestContainer,
  url: string,
): Promise<{ exitCode: number; output: string }> {
  return container.exec([
    "node",
    "-e",
    `fetch('${url}').then(() => process.exit(0)).catch(() => process.exit(1))`,
  ]);
}

describe("Custom CA Certificate Integration", () => {
  let network: StartedNetwork;
  let mockServer: StartedTestContainer;
  let certsDir: string;

  beforeAll(async () => {
    certsDir = fs.mkdtempSync(path.join(os.tmpdir(), "ca-test-certs-"));

    const certs = generateCertificates();
    fs.writeFileSync(path.join(certsDir, "ca.crt"), certs.ca.cert);
    fs.writeFileSync(path.join(certsDir, "server.key"), certs.server.key);
    fs.writeFileSync(path.join(certsDir, "server.crt"), certs.server.cert);

    network = await new Network().start();

    mockServer = await new GenericContainer("ghcr.io/smocker-dev/smocker")
      .withNetwork(network)
      .withNetworkAliases("mock-zitadel")
      .withBindMounts([{ source: certsDir, target: "/certs", mode: "ro" }])
      .withEnvironment({
        SMOCKER_MOCK_SERVER_LISTEN_PORT: "443",
        SMOCKER_CONFIG_LISTEN_PORT: "8081",
        SMOCKER_TLS_ENABLE: "true",
        SMOCKER_TLS_CERT_FILE: "/certs/server.crt",
        SMOCKER_TLS_PRIVATE_KEY_FILE: "/certs/server.key",
      })
      .withExposedPorts(SMOCKER_MOCK_PORT, SMOCKER_ADMIN_PORT)
      .withStartupTimeout(120_000)
      .withWaitStrategy(Wait.forLogMessage(/Starting mock server/))
      .start();

    try {
      await GenericContainer.fromDockerfile(LOGIN_APP_DIR).build(LOGIN_IMAGE_TAG);
    } catch (err) {
      throw new Error(`Failed to build login Docker image: ${err instanceof Error ? err.message : err}`);
    }
  }, 180_000);

  afterAll(async () => {
    await mockServer?.stop().catch(() => {});
    await network?.stop().catch(() => {});
  });

  describe("when custom CA certificate is provided", () => {
    let container: StartedTestContainer;

    beforeAll(async () => {
      container = await new GenericContainer(LOGIN_IMAGE_TAG)
        .withNetwork(network)
        .withExposedPorts(3000)
        .withEnvironment({
          SSL_CERT_FILE: "/etc/ssl/certs/custom-ca.crt",
        })
        .withCopyFilesToContainer([
          {
            source: path.join(certsDir, "ca.crt"),
            target: "/etc/ssl/certs/custom-ca.crt",
          },
        ])
        .withStartupTimeout(180_000)
        .withWaitStrategy(Wait.forHttp("/ui/v2/login/healthy", 3000))
        .start();
    }, 180_000);

    afterAll(async () => {
      await container?.stop().catch(() => {});
    });

    it("can establish TLS connection to server with custom CA", async () => {
      const result = await testTlsConnection(container, "https://mock-zitadel/");
      expect(result.exitCode).toBe(0);
    });

    it("has NODE_OPTIONS set with --use-openssl-ca", async () => {
      const result = await container.exec(["printenv", "NODE_OPTIONS"]);
      expect(result.output.trim()).toContain("--use-openssl-ca");
    });

    it("has SSL_CERT_FILE pointing to custom CA", async () => {
      const result = await container.exec(["printenv", "SSL_CERT_FILE"]);
      expect(result.output.trim()).toBe("/etc/ssl/certs/custom-ca.crt");
    });

    it("has CA certificate accessible in container", async () => {
      const result = await container.exec([
        "sh",
        "-c",
        "test -f /etc/ssl/certs/custom-ca.crt && echo exists",
      ]);
      expect(result.output.trim()).toContain("exists");
    });
  });

  describe("when custom CA certificate is not provided", () => {
    let container: StartedTestContainer;

    beforeAll(async () => {
      container = await new GenericContainer(LOGIN_IMAGE_TAG)
        .withNetwork(network)
        .withExposedPorts(3000)
        .withStartupTimeout(180_000)
        .withWaitStrategy(Wait.forHttp("/ui/v2/login/healthy", 3000))
        .start();
    }, 180_000);

    afterAll(async () => {
      await container?.stop().catch(() => {});
    });

    it("cannot establish TLS connection without custom CA", async () => {
      const result = await testTlsConnection(container, "https://mock-zitadel/");
      expect(result.exitCode).not.toBe(0);
    });

    it("uses system CA store by default", async () => {
      const result = await container.exec(["printenv", "SSL_CERT_FILE"]);
      expect(result.output.trim()).toBe("/etc/ssl/certs/ca-certificates.crt");
    });

    it("does not have custom CA certificate mounted", async () => {
      const result = await container.exec([
        "sh",
        "-c",
        "test -f /etc/ssl/certs/custom-ca.crt && echo exists || echo missing",
      ]);
      expect(result.output.trim()).toContain("missing");
    });
  });
});
