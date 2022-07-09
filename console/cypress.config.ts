import { defineConfig } from 'cypress';
import { readFileSync } from 'fs';

export default defineConfig({
  reporter: 'mochawesome',

  reporterOptions: {
    reportDir: 'cypress/results',
    overwrite: false,
    html: true,
    json: true,
  },

  chromeWebSecurity: false,
  trashAssetsBeforeRuns: false,
  defaultCommandTimeout: 10000,

  e2e: {
    experimentalSessionAndOrigin: true,
    setupNodeEvents(on, config) {

      require('cypress-terminal-report/src/installLogsPrinter')(on);

      config.defaultCommandTimeout = 10_000
      config.env.parsedServiceAccountKey = config.env.serviceAccountKey
      if (config.env.serviceAccountKeyPath) {
        config.env.parsedServiceAccountKey = JSON.parse(readFileSync(config.env.serviceAccountKeyPath, 'utf-8'))
      }
      return config
    },
  },
});
