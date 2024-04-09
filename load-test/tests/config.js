import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import crypto from 'k6/crypto';
import http from 'k6/http';

import url from './url.js';
import execution from 'k6/execution';

export const Config = {
  host: __ENV.ZITADEL_HOST || 'http://localhost:8080',
  orgId: __ENV.ZITADEL_ORG_ID || '',
  codeVerifier: __ENV.CODE_VERIFIER || randomString(10),
  admin: {
    loginName: __ENV.ADMIN_LOGIN_NAME || 'zitadel-admin@zitadel.localhost',
    password: __ENV.ADMIN_PASSWORD || 'Password1!'
  }
};

const client = {
    'response_type': 'code',
    'scope': 'openid email profile urn:zitadel:iam:org:project:id:zitadel:aud',
    'prompt': 'login',
    'code_challenge_method': 'S256',
    'code_challenge': crypto.sha256(Config.codeVerifier, "base64rawurl"),
    'client_id': __ENV.CLIENT_ID || '',
    'redirect_uri': url('/ui/console/auth/callback')
};

export function Client() {
  if (client.client_id) {
    return client;
  }
  client.client_id = http.get(url('/ui/console/assets/environment.json')).json().clientid;
  return client
}

let maxVUs;
export function MaxVUs() {
  if (maxVUs != undefined) {
    return maxVUs;
  }

  let max = 1;
  if (execution.test.options.stages != undefined) {
    max = execution.test.options.stages.reduce((acc, value) => {
      if (acc <= value.target) {
        return;
      }
      acc = value.target;
    });
  }

  new Map(Object.entries(execution.test.options.scenarios)).forEach((value) => {
    if (max < value.vus) {
      max = value.vus;
    }
  })

  maxVUs = max;
  return maxVUs;
}