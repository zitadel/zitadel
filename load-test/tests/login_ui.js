import { check, fail } from 'k6';
import http from 'k6/http';
import { URL } from 'https://jslib.k6.io/url/1.0.0/index.js';
import { Trend } from 'k6/metrics';
import encoding from 'k6/encoding';

import { Config, Client } from './config.js';
import url from './url.js';

export function loginByUsernamePassword(user) {
  check(user, {
    'user defined': (u) => u !== undefined || fail(`user is undefined`)
  });

  const loginUI = initLogin();
  const loginNameResponse = enterLoginName(loginUI, user);
  const passwordResponse = enterPassword(loginNameResponse, user);
  return token(new URL(passwordResponse.url).searchParams.get('code'));
}

const initLoginTrend = new Trend('login_ui_init_login_duration', true);
function initLogin() {
  const response = http.get(url('/oauth/v2/authorize', {searchParams: Client()}));
  check(response, {
    'authorize status ok': (r) => r.status == 200 || fail('init login failed', r)
  });

  initLoginTrend.add(response.timings.duration);

  return response;
}

const enterLoginNameTrend = new Trend('login_ui_enter_login_name_duration', true);
function enterLoginName(page, user) {
  const response = page.submitForm({
    formSelector: 'form',
    fields: {
      'loginName': user.loginName
    }
  });

  check(response, {
    'login name status ok': (r) => r.status == 200 || fail('enter login name failed'),
    'login shows password page': (r) => r.body.includes('password'),
    'login has no error': (r) => !r.body.includes('error') || fail(`error in enter login name ${r.body}`)
  });

  enterLoginNameTrend.add(response.timings.duration);

  return response;
}

const enterPasswordTrend = new Trend('login_ui_enter_password_duration', true);
function enterPassword(page, user) {
  let response = page.submitForm({
    formSelector: 'form',
    fields: {
      'password': user.password
    }
  });
  enterPasswordTrend.add(response.timings.duration);

  // skip 2fa init
  if (response.url.endsWith('/password')) {
    response = response.submitForm({
      formSelector: 'form',
      submitSelector: '[name="skip"]'
    });;
  }

  check(response, {
    'password status ok': (r) => r.status == 200 || fail('enter password failed'),
    'password callback': (r) => r.url.startsWith(url('/ui/console/auth/callback?code='))  || fail(`wrong password callback: ${r.url}`)
  });

  return response;
}

const tokenTrend = new Trend('login_ui_token_duration', true);
function token(code = '') {
  check(code, {
    'code set': (c) => (code !== undefined && code !== null) || fail('code was not set')
  });
  const response = http.post(url('/oauth/v2/token'), {
      grant_type: 'authorization_code',
      code: code,
      redirect_uri: Client().redirect_uri,
      code_verifier: Config.codeVerifier,
      client_id: Client().client_id
  }, {
      headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
      },
  });

  check(response, {
    'token status ok': (r) => r.status == 200 || fail(`invalid token response status: ${r.status} body: ${r.body}`),
    'access token created': (r) => r.json().access_token !== undefined
  });
  
  tokenTrend.add(response.timings.duration);
  return {
    accessToken: response.json().access_token,
    idToken: response.json().id_token,
    info: JSON.parse(encoding.b64decode(response.json().id_token.split('.')[1].toString(), 'rawstd', 's'))
  }
}
