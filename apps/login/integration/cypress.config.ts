import { defineConfig } from "cypress";

export default defineConfig({
  reporter: "list",

  e2e: {
    baseUrl: process.env.LOGIN_BASE_URL || "http://localhost:3000",
    specPattern: "integration/**/*.cy.{js,jsx,ts,tsx}",
    supportFile: "support/e2e.{js,jsx,ts,tsx}",
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
});
