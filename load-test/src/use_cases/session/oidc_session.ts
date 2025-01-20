import { loginByUsernamePassword, initLogin } from '../../login_ui';
import { createOrg, removeOrg } from '../../org';
import { User, createHuman } from '../../user';
import { Trend } from 'k6/metrics';
import { Config, MaxVUs } from '../../config';
import { check } from 'k6';
import { authRequestByID, finalizeAuthRequest } from '../../oidc';
import http from 'k6/http';
import { createSession } from '../../session';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.log('setup: admin signed in');

  const org = await createOrg(tokens.accessToken!);
  console.log(`setup: org (${org.organizationId}) created`);

  const humanPromises = Array.from({ length: MaxVUs() }, (_, i) => {
    return createHuman(`zitizen-${i}`, org, tokens.accessToken!);
  });

  const humans = (await Promise.all(humanPromises)).map((user) => {
    return { userId: user.userId, loginName: user.loginNames[0], password: 'Password1!' };
  });
  console.log(`setup: ${humans.length} users created`);
  return { tokens, users: humans, org };
}

// implements the flow described in
// https://zitadel.com/docs/guides/integrate/login-ui/oidc-standard
const addSessionTrend = new Trend('oidc_session_duration', true);
export default async function (data: any) {
  const start = new Date();
  const authorizeResponse = initLogin(true);
  check(authorizeResponse, {
    'authorize is status ok': (s) => s.status === 302,
  });

  const authRequestId = new URLSearchParams(authorizeResponse.headers['Location']).values().next().value;
  check(authRequestId, {
    'auth request id returned': (s) => s !== '',
  });

  const authRequest = await authRequestByID(authRequestId!, data.tokens!);
  check(authRequest, {
    'auth request is status ok': (s) => s.status === 200,
  });

  const session = await createSession(data.users[__VU -1], data.org, data.tokens.accessToken);
  
  const finalizedAuthRequest = await finalizeAuthRequest(authRequestId!, session, data.tokens!);
  console.log(`finalizedAuthRequest: ${JSON.stringify(finalizedAuthRequest)}`);
  addSessionTrend.add(new Date().getTime() - start.getTime());
}

export function teardown(data: any) {
  removeOrg(data.org, data.tokens.accessToken);
}
