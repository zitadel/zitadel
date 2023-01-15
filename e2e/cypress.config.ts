import { defineConfig } from 'cypress';
import * as CRD from 'chrome-remote-interface'

let tokensCache = new Map<string,string>()
let crdPort: number
let crdClient: Promise<CRD.Client> = null

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
    setupNodeEvents(on, config) {
      on("before:browser:launch", (browser, browserCfg) => {
        crdPort = parseInt(browserCfg.args.find(arg => arg.startsWith('--remote-debugging-port'))?.split('=')[1]) || parseInt(process.env.CYPRESS_REMOTE_DEBUGGING_PORT) || 4201
      }),
      on('task', {
        safetoken({key, token}) {
          tokensCache.set(key,token);
          return null
        },
        loadtoken({key}): string | null {
          return tokensCache.get(key) || null;
        },
        generateOTP: require("cypress-otp"),
        resetCRDInterface: async () => {
          if (crdClient) {
            (await crdClient).close()
            crdClient = null
          }
          return null
        },
        remoteDebuggerCommand: async (args) => {
          crdClient = crdClient || CRD({port: crdPort});
          return (await crdClient).send(args.event, args.params)
        }
      })
    },
  },
});

function baseUrl(){
  return process.env.CYPRESS_BASE_URL || 'http://localhost:8080/ui/console'
}
