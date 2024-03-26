import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import crypto from 'k6/crypto';
import http from 'k6/http';

import url from './url.js';

export const Config = {
  host: __ENV.ZITADEL_HOST || 'http://localhost:8080',
  orgId: __ENV.ZITADEL_ORG_ID || '',
  codeVerifier: __ENV.CODE_VERIFIER || randomString(10)
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
