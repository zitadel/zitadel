import { defineConfig, devices } from "@playwright/test";
import * as dotenv from "dotenv";
import * as path from "path";

// Load .env.test so credentials and port are available when running on the host.
// On CI the same file is passed to the shell environment before invoking NX.
dotenv.config({ path: path.resolve(__dirname, ".env.test") });

export default defineConfig({
  testDir: "./tests",
  // Fail fast: one retry on flakiness, full trace on first retry for debugging.
  retries: 1,
  reporter: [["html", { open: "never" }], ["line"]],
  use: {
    // All navigation goes through Traefik on the published port.
    baseURL:
      process.env.PLAYWRIGHT_BASE_URL ??
      `http://localhost:${process.env.PROXY_HTTP_PUBLISHED_PORT ?? "8080"}`,
    trace: "on-first-retry",
    screenshot: "only-on-failure",
  },
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
  ],
});
