import { defineConfig } from 'cypress';
import { Client } from "pg";
import { createServer } from 'http'
import { ZITADELWebhookEvent } from 'cypress/support/types';

const jwt = require('jsonwebtoken');

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAzi+FFSJL7f5yw4KTwzgMP34ePGycm/M+kT0M7V4Cgx5V3EaD
IvTQKTLfBaEB45zb9LtjIXzDw0rXRoS2hO6th+CYQCz3KCvh09C0IzxZiB2IS3H/
aT+5Bx9EFY+vnAkZjccbyG5YNRvmtOlnvIeIH7qZ0tEwkPfF5GEZNPJPtmy3UGV7
iofdVQS1xRj73+aMw5rvH4D8IdyiAC3VekIbpt0Vj0SUX3DwKtog337BzTiPk3aX
RF0sbFhQoqdJRI8NqgZjCwjq9yfI5tyxYswn+JGzHGdHvW3idODlmwEt5K2pasiR
IWK2OGfq+w0EcltQHabuqEPgZlmhCkRdNfixBwIDAQABAoIBAA9jNoBkRdxmH/R9
Wz+3gBqA9Aq4ZFuzJJk8QCm62V8ltWyyCnliYeKhPEm0QWrWOwghr/1AzW9Wt4g4
wVJcabD5TwODF5L0626eZcM3bsscwR44TMJzEgD5EWC2j3mKqFCPaoBj08tq4KXh
wW8tgjgz+eTk3cYD583qfTIZX1+SzSMBpetTBsssQtGhhOB/xPiuL7hi+fXmV2rh
8mc9X6+wJ5u3zepsyK0vBeEDmurD4ZUIXFrZ0WCB/wNkSW9VKyoH+RC1asQAgqTz
glJ/NPbDJSKGvSBQydoKkqoXx7MVJ8VObFddfgo4dtOoz6YCfUVBHt8qy+E5rz5y
CICjL/kCgYEA9MnHntVVKNXtEFZPo02xgCwS3eG27ZwjYgJ1ZkCHM5BuL4MS7qbr
743/POs1Ctaok0udHl1PFB4uAG0URnmkUnWzcoJYb6Plv03F0LRdsnfuhehfIxLP
nWvxSm5n21H4ytfxm0BWY09JkLDnJZtXrgTILbuqb9Wy6TmAvUaF2YUCgYEA16Ec
ywSaLVdqPaVpsTxi7XpRJAB2Isjp6RffNEecta4S0LL7s/IO3QXDH9SYpgmgCTah
3aXhpT4hIFlpg3eBjVfbOwgqub8DgirnSQyQt99edUtHIK+K8nMdGxz6X6pfTKzK
asSH7qPlt5tz1621vC0ocXSZR7zm99/FgwILwBsCgYBOsP8nJFV4By1qbxSy3qsN
FR4LjiAMSoFlZHzxHhVYkjmZtH1FkwuNuwwuPT6T+WW/1DLyK/Tb9se7A1XdQgV9
LLE/Qn/Dg+C7mvjYmuL0GHHpQkYzNDzh0m2DC/L/Il7kdn8I9anPyxFPHk9wW3vY
SVlAum+T/BLDvuSP9DfbMQKBgCc1j7PG8XYfOB1fj7l/volqPYjrYI/wssAE7Dxo
bTGIJrm2YhiVgmhkXNfT47IFfAlQ2twgBsjyZDmqqIoUWAVonV+9m29NMYkg3g+l
bkdRIa74ckWaRgzSK8+7VDfDFjMuFFyXwhP9z460gLsORkaie4Et75Vg3yrhkNvC
qnpTAoGBAMguDSWBbCewXnHlKGFpm+LH+OIvVKGEhtCSvfZojtNrg/JBeBebSL1n
mmT1cONO+0O5bz7uVaRd3JdnH2JFevY698zFfhVsjVCrm+fz31i5cxAgC39G2Lfl
YkTaa1AFLstnf348ZjuvBN3USUYZo3X3mxnS+uluVuRSGwIKsN0a
-----END RSA PRIVATE KEY-----`

let tokensCache = new Map<string,string>()

let webhookEvents = new Array<ZITADELWebhookEvent>()

export default defineConfig({
  reporter: 'mochawesome',

  reporterOptions: {
    reportDir: 'cypress/results',
    overwrite: false,
    html: true,
    json: true,
  },

  trashAssetsBeforeRuns: false,
  defaultCommandTimeout: 10000,

  env: {
    ORGANIZATION: process.env.CYPRESS_ORGANIZATION || 'zitadel',
    BACKEND_URL: backendUrl(),
    WEBHOOK_HANDLER_PORT: webhookHandlerPort(),
    WEBHOOK_HANDLER_HOST: process.env.CYPRESS_WEBHOOK_HANDLER_HOST || 'localhost',
  },

  e2e: {
    baseUrl: baseUrl(),
    experimentalRunAllSpecs: true,
    setupNodeEvents(on, config) {

      startWebhookEventHandler()

      on('task', {
        safetoken({key, token}) {
          tokensCache.set(key,token);
          return null
        },
        loadtoken({key}): string | null {
          return tokensCache.get(key) || null;
        },
        systemToken(): Promise<string> {
          let iat = Math.floor(Date.now() / 1000);
          let exp = iat + (999*12*30*24*60*60) // ~ 999 years
          return jwt.sign({
            "iss": "cypress",
            "sub": "cypress",
            "aud": backendUrl(),
            "iat": iat,
            "exp": exp
          }, Buffer.from(privateKey, 'ascii').toString('ascii'), { algorithm: 'RS256' })
        },
        async runSQL(statement: string) {
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
        },
        resetWebhookEvents() {
          webhookEvents = []
          return null
        },
        handledWebhookEvents(){
          return webhookEvents
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

function webhookHandlerPort() {
  return process.env.CYPRESS_WEBHOOK_HANDLER_PORT || '8900'
}

function startWebhookEventHandler() {
  const port = webhookHandlerPort()
  const server = createServer((req, res) => {
    const chunks = [];
    req.on("data", (chunk) => {
      chunks.push(chunk);
    });
    req.on("end", () => {
      webhookEvents.push(JSON.parse(Buffer.concat(chunks).toString()));
    });

    res.writeHead(200);
    res.end()
  });

  server.listen(port, () => {
    console.log(`Server is running on http://:${port}`);
  });
}
