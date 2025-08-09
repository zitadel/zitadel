import { defineConfig } from "cypress";

export default defineConfig({
  reporter: "list",
  video: true,
  pageLoadTimeout: 2 * 60_000,
  e2e: {
    baseUrl: process.env.LOGIN_BASE_URL || "http://localhost:3001/ui/v2/login",
    specPattern: "integration/integration/**/*.cy.{js,jsx,ts,tsx}",
    supportFile: "integration/support/e2e.{js,jsx,ts,tsx}",
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
});
