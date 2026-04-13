// @ts-ignore Import module
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import crypto from 'k6/crypto';
import http from 'k6/http';
import execution from 'k6/execution';
import { Stage } from 'k6/options';

import url from './url';

export const Config = {
  host: __ENV.ZITADEL_HOST || 'http://localhost:8080',
  orgId: '',
  codeVerifier: __ENV.CODE_VERIFIER || randomString(10),
  admin: {
    loginName: __ENV.ADMIN_LOGIN_NAME || 'zitadel-admin@zitadel.localhost',
    password: __ENV.ADMIN_PASSWORD || 'Password1!',
  },
};

const client = {
  response_type: 'code',
  scope: 'openid email profile urn:zitadel:iam:org:project:id:zitadel:aud',
  prompt: 'login',
  code_challenge_method: 'S256',
  code_challenge: crypto.sha256(Config.codeVerifier, 'base64rawurl'),
  client_id: __ENV.CLIENT_ID || '',
  redirect_uri: url('/ui/console/auth/callback'),
};

export function Client() {
  if (client.client_id) {
    return client;
  }
  const env = http.get(url('/ui/console/assets/environment.json'));

  client.client_id = env.json('clientid') ? env.json('clientid')?.toString()! : '';

  return client;
}

let maxVUs: number;
export function MaxVUs() {
  if (maxVUs != undefined) {
    return maxVUs;
  }

  let max: number = execution.test.options.stages
    ? execution.test.options.stages
        .map((value: Stage): number => value.target)
        .reduce((acc: number, value: number): number => {
          return acc <= value ? acc : value;
        })
    : 1;

  if (execution.test.options.scenarios) {
    new Map(Object.entries(execution.test.options.scenarios)).forEach((value) => {
      if ('vus' in value) {
        max = value.vus && max < value.vus ? value.vus : max;
      } else if ('maxVUs' in value) {
        max = value.maxVUs && max < value.maxVUs ? value.maxVUs : max;
      }
    });
  }

  maxVUs = max;
  return maxVUs;
}
