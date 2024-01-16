import { defineConfig } from "cypress";

export default defineConfig({
  reporter: "list",
  e2e: {
    baseUrl: "http://localhost:3000",
    specPattern: "cypress/integration/**/*.cy.{js,jsx,ts,tsx}",
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
});
