import { defineConfig } from "cypress";
const watchApp = require("cypress-app-watcher-preprocessor");

export default defineConfig({
  reporter: "list",
  e2e: {
    baseUrl: "http://localhost:3000",
    specPattern: "cypress/integration/**/*.cy.{js,jsx,ts,tsx}",
    setupNodeEvents(on, config) {
        on("file:preprocessor", watchApp());
    },
  },
});
