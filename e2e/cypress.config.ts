import { defineConfig } from 'cypress';

let tokensCache = new Map<string,string>()

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

  env: {
    ORGANIZATION: process.env.CYPRESS_ORGANIZATION || 'zitadel'
  },

  e2e: {
    baseUrl: process.env.CYPRESS_BASE_URL || 'http://localhost:8080',
    experimentalSessionAndOrigin: true,
    setupNodeEvents(on, config) {

      on('task', {
        safetoken({key, token}) {
          tokensCache.set(key,token);
          return null
        }
      })
      on('task', {
        loadtoken({key}): string | null {
          return tokensCache.get(key) || null;
        }
      })
    },
  },
});
