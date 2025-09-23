import { defineConfig } from "cypress";

export default defineConfig({
  reporter: "list",
  video: true,
  retries: {
    runMode: 2
  },
  e2e: {
    baseUrl: `http://localhost:3000${process.env.NEXT_PUBLIC_BASE_PATH || ""}`,
    specPattern: "integration/integration/**/*.cy.{js,jsx,ts,tsx}",
    supportFile: "integration/support/e2e.{js,jsx,ts,tsx}",
    pageLoadTimeout: 120_0000,
    env: {
      API_MOCK_STUBS_URL: process.env.API_MOCK_STUBS_URL || "http://localhost:22220/v1/stubs"
    },
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
});
