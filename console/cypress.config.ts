import { defineConfig } from 'cypress';

export default defineConfig({
  reporter: 'mochawesome',

  reporterOptions: {
    reportDir: 'cypress/results',
    overwrite: false,
    html: true,
    json: true,
  },

  chromeWebSecurity: false,
  //   experimentalSessionSupport: true,
  trashAssetsBeforeRuns: false,
  defaultCommandTimeout: 10000,

  e2e: {
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
});
