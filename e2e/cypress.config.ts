import { defineConfig } from 'cypress';
import { readFileSync } from'fs';
import { Client } from "pg";

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
    BACKEND_URL: backendUrl(),
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
      on('task', {
        systemToken(): Promise<string> {
          const privateKey = readFileSync(process.env.CYPRESS_SYSTEM_USER_KEY_PATH || `${__dirname}/systemuser/cypress.pem.base64`, 'utf-8')
          let iat = Math.floor(Date.now() / 1000);
          let exp = iat + (24*60*60)
          return jwt.sign({
            "iss": "cypress",
            "sub": "cypress",
            "aud": backendUrl(),
            "iat": iat,
            "exp": exp
          }, Buffer.from(privateKey, 'base64').toString('ascii'), { algorithm: 'RS256' })
        }
      })
      on('task', {
        runSQL(statement: string){
          const client = new Client({
            connectionString: process.env.CYPRESS_DATABASE_CONNECTION_URL || 'postgresql://root@localhost:26257/zitadel'
          });

          return client.connect().then(() => {
            return client.query(statement).then((result) => {
              return client.end().then(() => {
                return result
              })
            })
          })
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
