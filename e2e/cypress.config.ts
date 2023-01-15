import { defineConfig } from 'cypress';
import { readFileSync } from'fs';
const jwt = require('jsonwebtoken');

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
    BACKEND_URL: backendUrl()
  },

  e2e: {
    baseUrl: baseUrl(),
    setupNodeEvents(on, config) {

      on('task', {
        safetoken({key, token}) {
          tokensCache.set(key,token);
          return null
        },
        loadtoken({key}): string | null {
          return tokensCache.get(key) || null;
        },
        generateOTP: require("cypress-otp")
      })
      on('task', {
        systemToken(): Promise<string> {
          const privateKey = readFileSync(process.env.CYPRESS_SYSTEM_USER_KEY_PATH || `${__dirname}/systemuser/cypress.pem`, 'utf-8')
          console.log("pk",  privateKey)
          let iat = Math.floor(Date.now() / 1000);
          let exp = iat + (24*60*60)
          return jwt.sign({
            "iss": "cypress",
            "sub": "cypress",
            "aud": "http://localhost:8080",
            "iat": iat,
            "exp": exp
          }, privateKey, { algorithm: 'RS256' })
        }
      })
    },
  },
});

function baseUrl(){
  return process.env.CYPRESS_BASE_URL || 'http://localhost:8080/ui/console'
}

function backendUrl(){
  return process.env.CYPRESS_BACKEND_URL || baseUrl().replace("/ui/console", "")
}