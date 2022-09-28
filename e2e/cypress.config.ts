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
    ORGANIZATION: process.env.CYPRESS_ORGANIZATION || 'zitadel',
    BACKEND_URL: process.env.CYPRESS_BACKEND_URL || baseUrl().replace("/ui/console", "")
  },

  e2e: {
    baseUrl: baseUrl(),
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

function baseUrl(){
  return process.env.CYPRESS_BASE_URL || 'http://localhost:8080/ui/console'
}
