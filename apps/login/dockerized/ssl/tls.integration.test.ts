import { describe, it, expect, beforeAll, afterAll } from "vitest";
import { GenericContainer, type StartedTestContainer, Wait } from "testcontainers";
import * as https from "node:https";
import * as fs from "fs";
import * as path from "path";
import { fileURLToPath } from "node:url";
import { generateCertificates } from "../ca/utils/tls.ts";

const TEST_DIR = path.dirname(fileURLToPath(import.meta.url));
const OUTPUT_DIR = path.join(TEST_DIR, "output");
const CERTS_DIR = path.join(OUTPUT_DIR, "certs");
const LOGIN_APP_DIR = path.join(TEST_DIR, "../..");

/**
 * Integration tests for TLS support in the login application.
 *
 * These tests verify that the login container can serve traffic over HTTPS
 * when configured with ZITADEL_TLS_ENABLED, ZITADEL_TLS_CERTPATH, and
 * ZITADEL_TLS_KEYPATH environment variables. This is essential for deployments
 * where the login application must terminate TLS itself rather than relying on
 * a reverse proxy.
 *
 * The test suite builds the login container image and runs it in two modes:
 * HTTP (default) and HTTPS (with TLS enabled). For HTTPS mode, a self-signed
 * certificate is generated and mounted into the container. The tests verify
 * that the application starts correctly and responds to health checks in both
 * modes.
 */
describe("Login Container TLS Support", () => {
  let loginImage: GenericContainer;

  beforeAll(async () => {
    if (fs.existsSync(OUTPUT_DIR)) {
      fs.rmSync(OUTPUT_DIR, { recursive: true });
    }
    fs.mkdirSync(CERTS_DIR, { recursive: true });

    const certs = generateCertificates({
      serverCN: "localhost",
      dns: ["localhost"],
      ips: ["127.0.0.1"],
    });

    fs.writeFileSync(path.join(CERTS_DIR, "server.crt"), certs.server.cert);
    fs.writeFileSync(path.join(CERTS_DIR, "server.key"), certs.server.key);

    loginImage = await GenericContainer.fromDockerfile(LOGIN_APP_DIR).build();
  });

  describe("when TLS is not enabled", () => {
    let container: StartedTestContainer;
    let appUrl: string;

    beforeAll(async () => {
      container = await loginImage
        .withExposedPorts(3000)
        .withWaitStrategy(Wait.forHttp("/ui/v2/login/healthy", 3000))
        .start();

      appUrl = `http://${container.getHost()}:${container.getMappedPort(3000)}`;
    });

    afterAll(async () => {
      await container?.stop();
    });

    it("returns 200 from healthy endpoint", async () => {
      const response = await fetch(`${appUrl}/ui/v2/login/healthy`);
      expect(response.status).toBe(200);
    });
  });

  describe("when TLS is enabled", () => {
    let container: StartedTestContainer;
    let appUrl: string;

    beforeAll(async () => {
      container = await loginImage
        .withExposedPorts(3000)
        .withEnvironment({
          ZITADEL_TLS_ENABLED: "true",
          ZITADEL_TLS_CERTPATH: "/certs/server.crt",
          ZITADEL_TLS_KEYPATH: "/certs/server.key",
        })
        .withCopyFilesToContainer([
          {
            source: path.join(CERTS_DIR, "server.crt"),
            target: "/certs/server.crt",
          },
          {
            source: path.join(CERTS_DIR, "server.key"),
            target: "/certs/server.key",
          },
        ])
        .withWaitStrategy(Wait.forLogMessage(/Ready on https/))
        .start();

      appUrl = `https://${container.getHost()}:${container.getMappedPort(3000)}`;
    });

    afterAll(async () => {
      await container?.stop();
    });

    it("returns 200 from healthy endpoint over HTTPS", async () => {
      const response = await new Promise<https.IncomingMessage>((resolve, reject) => {
        https
          .get(`${appUrl}/ui/v2/login/healthy`, { rejectUnauthorized: false }, (res) => {
            res.resume();
            resolve(res);
          })
          .on("error", reject);
      });
      expect(response.statusCode).toBe(200);
    });
  });
});
