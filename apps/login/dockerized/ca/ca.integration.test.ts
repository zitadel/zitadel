import { describe, it, expect, beforeAll, afterAll } from "vitest";
import { GenericContainer, type StartedTestContainer, Network, Wait } from "testcontainers";
import type { StartedNetwork } from "testcontainers";
import * as fs from "fs";
import * as path from "path";
import * as esbuild from "esbuild";
import { generateCertificates } from "./utils/tls.ts";

const TEST_DIR = path.dirname(new URL(import.meta.url).pathname);
const OUTPUT_DIR = path.join(TEST_DIR, "output");
const CERTS_DIR = path.join(OUTPUT_DIR, "certs");
const LOGIN_APP_DIR = path.join(TEST_DIR, "../..");

async function bundleMockServer(): Promise<string> {
  const bundleDir = path.join(OUTPUT_DIR, "bundle");
  fs.mkdirSync(bundleDir, { recursive: true });

  const outfile = path.join(bundleDir, "mock-server.js");
  await esbuild.build({
    entryPoints: [path.join(TEST_DIR, "mock-server.ts")],
    bundle: true,
    platform: "node",
    target: "node22",
    outfile,
    format: "cjs",
  });

  return outfile;
}

function readCapturedRequests(): Array<{ method: string; url: string; tlsConnected: boolean }> {
  const filePath = path.join(OUTPUT_DIR, "requests.json");
  if (!fs.existsSync(filePath)) {
    return [];
  }
  return JSON.parse(fs.readFileSync(filePath, "utf-8"));
}

/**
 * Integration tests for custom CA certificate support in the login application.
 *
 * These tests verify that the Dockerfile correctly configures Node.js to trust custom
 * Certificate Authorities via the SSL_CERT_FILE environment variable and the
 * NODE_OPTIONS=--use-openssl-ca flag. This is essential for deployments where the
 * ZITADEL API is served behind a TLS-terminating proxy using internal or self-signed
 * certificates.
 *
 * The test suite spins up a mock TLS server with a self-signed certificate inside a
 * Docker network. The login application container connects to this mock server, and
 * the tests verify that TLS connections succeed when the custom CA is mounted and
 * fail when it is not. The mock server records all incoming requests to a shared
 * volume, allowing the tests to verify that TLS handshakes completed successfully.
 *
 * The Docker network ensures container-to-container communication uses internal DNS
 * resolution, avoiding port conflicts on the host machine and ensuring reliable
 * execution in CI environments.
 */
const LOGIN_IMAGE_TAG = "zitadel-login-otel-test:latest";

describe("Custom CA Certificate Integration", () => {
  let certs: ReturnType<typeof generateCertificates>;
  let network: StartedNetwork;
  let mockServer: StartedTestContainer;
  let mockServerBundle: string;

  beforeAll(async () => {
    if (fs.existsSync(OUTPUT_DIR)) {
      fs.rmSync(OUTPUT_DIR, { recursive: true });
    }
    fs.mkdirSync(CERTS_DIR, { recursive: true });

    certs = generateCertificates();

    fs.writeFileSync(path.join(CERTS_DIR, "ca.crt"), certs.ca.cert);
    fs.writeFileSync(path.join(CERTS_DIR, "server.key"), certs.server.key);
    fs.writeFileSync(path.join(CERTS_DIR, "server.crt"), certs.server.cert);

    mockServerBundle = await bundleMockServer();

    network = await new Network().start();

    mockServer = await new GenericContainer("node:22-alpine")
      .withNetwork(network)
      .withNetworkAliases("mock-zitadel")
      .withCopyFilesToContainer([
        { source: mockServerBundle, target: "/app/server.js" },
        { source: path.join(CERTS_DIR, "server.key"), target: "/certs/server.key" },
        { source: path.join(CERTS_DIR, "server.crt"), target: "/certs/server.crt" },
      ])
      .withBindMounts([{ source: OUTPUT_DIR, target: "/output" }])
      .withCommand(["node", "/app/server.js"])
      .withExposedPorts(443)
      .withWaitStrategy(Wait.forLogMessage("Mock TLS server listening"))
      .start();

    await GenericContainer.fromDockerfile(LOGIN_APP_DIR).build(LOGIN_IMAGE_TAG);
  }, 180000);


  describe("when custom CA certificate is provided", () => {
    let container: StartedTestContainer;
    let appUrl: string;

    beforeAll(async () => {
      if (fs.existsSync(path.join(OUTPUT_DIR, "requests.json"))) {
        fs.unlinkSync(path.join(OUTPUT_DIR, "requests.json"));
      }

      container = await new GenericContainer(LOGIN_IMAGE_TAG)
        .withNetwork(network)
        .withExposedPorts(3000)
        .withEnvironment({
          ZITADEL_API_URL: "https://mock-zitadel",
          ZITADEL_SERVICE_USER_TOKEN: "test-token",
          SSL_CERT_FILE: "/etc/ssl/certs/custom-ca.crt",
        })
        .withCopyFilesToContainer([
          {
            source: path.join(CERTS_DIR, "ca.crt"),
            target: "/etc/ssl/certs/custom-ca.crt",
          },
        ])
        .withStartupTimeout(180000)
        .withWaitStrategy(Wait.forHttp("/ui/v2/login/healthy", 3000))
        .start();

      appUrl = `http://${container.getHost()}:${container.getMappedPort(3000)}`;

      await fetch(`${appUrl}/ui/v2/login/loginname`, { redirect: "manual" });

      for (let i = 0; i < 30; i++) {
        const requests = readCapturedRequests();
        if (requests.length > 0) break;
        await new Promise((r) => setTimeout(r, 1000));
      }
    });

  
    it("connects to mock server over TLS", () => {
      const requests = readCapturedRequests();
      expect(requests.length).toBeGreaterThan(0);
    });

    it("marks requests as TLS connected", () => {
      const requests = readCapturedRequests();
      const tlsRequests = requests.filter((r) => r.tlsConnected === true);
      expect(tlsRequests.length).toBeGreaterThan(0);
    });

    it("returns 200 from healthy endpoint", async () => {
      const response = await fetch(`${appUrl}/ui/v2/login/healthy`);
      expect(response.status).toBe(200);
    });

    it("serves the login page", async () => {
      const response = await fetch(`${appUrl}/ui/v2/login/loginname`, { redirect: "manual" });
      expect([200, 302, 303, 307, 308]).toContain(response.status);
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
    let appUrl: string;
    let requestsBefore: number;

    beforeAll(async () => {
      requestsBefore = readCapturedRequests().length;

      container = await new GenericContainer(LOGIN_IMAGE_TAG)
        .withNetwork(network)
        .withExposedPorts(3000)
        .withEnvironment({
          ZITADEL_API_URL: "https://mock-zitadel",
          ZITADEL_SERVICE_USER_TOKEN: "test-token",
        })
        .withStartupTimeout(180000)
        .withWaitStrategy(Wait.forHttp("/ui/v2/login/healthy", 3000))
        .start();

      appUrl = `http://${container.getHost()}:${container.getMappedPort(3000)}`;

      await fetch(`${appUrl}/ui/v2/login/loginname`, { redirect: "manual" });

      for (let i = 0; i < 10; i++) {
        await new Promise((r) => setTimeout(r, 500));
      }
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

    it("cannot establish TLS connection to mock server", () => {
      const requestsAfter = readCapturedRequests().length;
      expect(requestsAfter).toBe(requestsBefore);
    });
  });
});
