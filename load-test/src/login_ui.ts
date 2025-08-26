import { JSONObject, check, fail } from 'k6';
import http, { Response } from 'k6/http';
// @ts-ignore Import module
import { URL } from 'https://jslib.k6.io/url/1.0.0/index.js';
import { Trend } from 'k6/metrics';

import { Config, Client } from './config';
import url from './url';
import { User } from './user';
import { Tokens } from './oidc';

export function loginByUsernamePassword(user: User) {
  check(user, {
    'user defined': (u) => u !== undefined || fail(`user is undefined`),
  });

  const loginUI = initLogin();
  const loginNameResponse = enterLoginName(loginUI, user);
  const passwordResponse = enterPassword(loginNameResponse, user);
  return token(new URL(passwordResponse.url).searchParams.get('code'));
}

const initLoginTrend = new Trend('login_ui_init_login_duration', true);
export function initLogin(clientId?: string): Response {
  let params = {};
  let expectedStatus = 200;
  if (clientId) {
    params = {
      headers: {
        'x-zitadel-login-client': clientId,
      },
      redirects: 0,
    };
    expectedStatus = 302;
  }

  const response = http.get(
    url('/oauth/v2/authorize', {
      searchParams: Client(),
    }),
    params,
  );
  check(response, {
    'authorize status ok': (r) => r.status == expectedStatus || fail(`init login failed: ${JSON.stringify(r)}`),
  });
  initLoginTrend.add(response.timings.duration);
  return response;
}

const enterLoginNameTrend = new Trend('login_ui_enter_login_name_duration', true);
function enterLoginName(page: Response, user: User): Response {
  const response = page.submitForm({
    formSelector: 'form',
    fields: {
      loginName: user.loginName,
    },
  });

  check(response, {
    'login name status ok': (r) => (r && r.status >= 200 && r.status < 300) || fail('enter login name failed'),
    'login shows password page': (r) => r && r.body !== null && r.body.toString().includes('password'),
    // 'login has no error': (r) => r && r.body != null && r.body.toString().includes('error') || fail(`error in enter login name ${r.body}`)
  });

  enterLoginNameTrend.add(response.timings.duration);

  return response;
}

const enterPasswordTrend = new Trend('login_ui_enter_password_duration', true);
function enterPassword(page: Response, user: User): Response {
  let response = page.submitForm({
    formSelector: 'form',
    fields: {
      password: user.password,
    },
  });
  enterPasswordTrend.add(response.timings.duration);

  // skip 2fa init
  if (response.url.endsWith('/password')) {
    response = response.submitForm({
      formSelector: 'form',
      submitSelector: '[name="skip"]',
    });
  }

  check(response, {
    'password status ok': (r) => r.status >= 200 && r.status < 300 || fail('enter password failed'),
    'password callback': (r) =>
      r.url.startsWith(url('/ui/console/auth/callback?code=')) || fail(`wrong password callback: ${r.url}`),
  });

  return response;
}

const tokenTrend = new Trend('login_ui_token_duration', true);
function token(code = '') {
  check(code, {
    'code set': (c) => (c !== undefined && c !== null) || fail('code was not set'),
  });
  const response = http.post(
    url('/oauth/v2/token'),
    {
      grant_type: 'authorization_code',
      code: code,
      redirect_uri: Client().redirect_uri,
      code_verifier: Config.codeVerifier,
      client_id: Client().client_id,
    },
    {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    },
  );

  tokenTrend.add(response.timings.duration);
  check(response, {
    'token status ok': (r) => r.status >= 200 && r.status < 300 || fail(`invalid token response status: ${r.status} body: ${r.body}`),
  });
  const token = new Tokens(response.json() as JSONObject);
  check(token, {
    'access token created': (t) => t.accessToken !== undefined,
    'id token created': (t) => t.idToken !== undefined,
    'info created': (t) => t.info !== undefined,
  });

  return token;
}
