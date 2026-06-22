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
 * custom Certificate Authorities via SSL_CERT_FILE, SSL_CERT_DIR, and the
 * load-ssl-cert-dir.cjs preload script.
 *
 * OpenSSL (used by Node with --use-openssl-ca) only reads certs from
 * SSL_CERT_DIR when they have hashed filenames (e.g. "9d66eef0.0").  Go's
 * crypto/x509 reads any file.  The preload script bridges this gap so that
 * operators get the same behaviour with both the Go backend and the login
 * container.
 *
 * We use Smocker as a TLS server whose certificate is signed by a test CA
 * that is NOT in the system trust store.  Each test block configures the
 * login container differently and asserts whether the TLS handshake to
 * Smocker succeeds or fails — verifying actual trust-store behaviour, not
 * just environment variable values.
 */
const LOGIN_IMAGE_TAG = `zitadel-login-test:${process.env.GITHUB_SHA || "local"}`;

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

  describe("when custom CA is provided via SSL_CERT_FILE", () => {
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
  });

  describe("when custom CA is provided via SSL_CERT_DIR (non-hashed filename)", () => {
    let container: StartedTestContainer;
    let customCaDir: string;

    beforeAll(async () => {
      customCaDir = fs.mkdtempSync(path.join(os.tmpdir(), "ca-dir-test-"));
      fs.writeFileSync(path.join(customCaDir, "corp-root.crt"), fs.readFileSync(path.join(certsDir, "ca.crt")));

      container = await new GenericContainer(LOGIN_IMAGE_TAG)
        .withNetwork(network)
        .withExposedPorts(3000)
        .withEnvironment({
          SSL_CERT_DIR: "/custom-certs",
        })
        .withCopyFilesToContainer([
          {
            source: path.join(customCaDir, "corp-root.crt"),
            target: "/custom-certs/corp-root.crt",
          },
        ])
        .withStartupTimeout(180_000)
        .withWaitStrategy(Wait.forHttp("/ui/v2/login/healthy", 3000))
        .start();
    }, 180_000);

    afterAll(async () => {
      await container?.stop().catch(() => {});
    });

    it("can establish TLS connection with non-hashed cert in SSL_CERT_DIR", async () => {
      const result = await testTlsConnection(container, "https://mock-zitadel/");
      expect(result.exitCode).toBe(0);
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
  });
});
