import { defineConfig } from "cypress";

export default defineConfig({
  reporter: 'list',
  e2e: {
    specPattern:  'cypress/integration/**/*.cy.{js,jsx,ts,tsx}',
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
});
