import { defineConfig } from "cypress";
import { startMockContainer, stopMockContainer } from "./integration/mock-container";

export default defineConfig({
  reporter: "list",
  video: true,
  retries: {
    runMode: 2
  },
  e2e: {
    baseUrl: `http://localhost:3001${process.env.NEXT_PUBLIC_BASE_PATH || ""}`,
    specPattern: "integration/integration/**/*.cy.{js,jsx,ts,tsx}",
    supportFile: "integration/support/e2e.{js,jsx,ts,tsx}",
    pageLoadTimeout: 120_0000,
    async setupNodeEvents(on, config) {
      // Start the grpc-mock container via testcontainers.
      // This replaces the old Nx `serve` target + `wait-on` approach.
      // The container binds to fixed ports (22220 for stubs, 22222 for mock)
      // matching the login app's ZITADEL_API_URL configuration.
      const mock = await startMockContainer();
      config.env.API_MOCK_STUBS_URL = mock.stubsUrl;

      on("after:run", async () => {
        await stopMockContainer();
      });

      return config;
    },
    env: {
      API_MOCK_STUBS_URL: process.env.API_MOCK_STUBS_URL || "http://localhost:22220/v1/stubs"
    }
  },
});
