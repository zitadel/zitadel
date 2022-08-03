import { defineConfig } from 'cypress';

let tokensCache = new Map<string,string>()
let initmfaandpwrequired = true

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
      on('task', {
        initmfaandpwrequired(){
          if (config.env.noInitMFAAndPWRequired == 'true'){
            initmfaandpwrequired = false
          }
          if (initmfaandpwrequired){
            initmfaandpwrequired = false
            return true
          }
          return false
        }
      })
    },
  },
});
