import { defineConfig } from "cypress";

export default defineConfig({
  reporter: "list",
  video: true,
  retries: {
    runMode: 2
  },
  e2e: {
    baseUrl: process.env.LOGIN_BASE_URL || "http://localhost:3000/ui/v2/login",
    specPattern: "integration/integration/**/*.cy.{js,jsx,ts,tsx}",
    supportFile: "integration/support/e2e.{js,jsx,ts,tsx}",
    env: {
      CORE_MOCK_STUBS_URL: process.env.CORE_MOCK_STUBS_URL || "http://localhost:22220/v1/stubs"
    },    
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
});
