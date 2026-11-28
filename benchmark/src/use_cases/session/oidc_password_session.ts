import { loginByUsernamePassword, initLogin, exchangeAuthorizationCode } from '../../login_ui';
import { createOrg, removeOrg } from '../../org';
import { User, addMachineKey, createMachine, createHuman } from '../../user';
import { Trend } from 'k6/metrics';
import { Config, MaxVUs } from '../../config';
import { check, fail } from 'k6';
import { finalizeAuthRequest, JWTProfileRequest, token } from '../../oidc';
import { createSession, setSession } from '../../session';
import encoding from 'k6/encoding';
import { addIAMMember } from '../../membership';
// @ts-ignore Import module
import { URL } from 'https://jslib.k6.io/url/1.0.0/index.js';

const publicKey = encoding.b64encode(open('../.keys/key.pem.pub'));

export async function setup() {
  const adminTokens = loginByUsernamePassword(Config.admin as User);
  console.log('setup: admin signed in');

  const org = await createOrg(adminTokens.accessToken!);
  console.log(`setup: org (${org.organizationId}) created`);

  const loginUser = await createMachine('load-test', org, adminTokens.accessToken!);
  const loginUserKey = await addMachineKey(loginUser.userId, org, adminTokens.accessToken!, publicKey);
  await addIAMMember(loginUser.userId, ['IAM_LOGIN_CLIENT'], adminTokens.accessToken!);
  const tokens = await token(new JWTProfileRequest(loginUser.userId, loginUserKey.keyId));

  const humanPromises = Array.from({ length: MaxVUs() }, (_, i) => {
    return createHuman(`zitizen-${i}`, org, adminTokens.accessToken!);
  });

  const users = (await Promise.all(humanPromises)).map((user) => {
    return { userId: user.userId, loginName: user.loginNames[0], password: 'Password1!' };
  });
  console.log(`setup: ${users.length} users created`);

  return { tokens, org, adminTokens, users };
}

// implements the username/password + OIDC flow described in
// https://zitadel.com/docs/guides/integrate/login-ui/username-password
// https://zitadel.com/docs/guides/integrate/login-ui/oidc-standard
const oidcPasswordSessionTrend = new Trend('oidc_password_session_duration', true);
export default async function (data: any) {
  const start = new Date();
  const user = data.users[__VU - 1];

  const authorizeResponse = initLogin(data.tokens.info.client_id);
  const authRequestId = new URLSearchParams(authorizeResponse.headers['Location']).values().next().value;
  check(authRequestId, {
    'auth request id returned': (s) => s !== '',
  });

  let session = await createSession(data.org, data.tokens.accessToken, {
    user: {
      loginName: user.loginName,
    },
  });
  const sessionId = (session as any).sessionId;

  session = await setSession(sessionId, session, data.tokens.accessToken, null, {
    password: {
      password: user.password,
    },
  });

  const finalizeResponse = await finalizeAuthRequest(
    authRequestId!,
    { sessionId, sessionToken: (session as any).sessionToken },
    data.tokens!,
  );
  const callbackUrl = finalizeResponse.json('callbackUrl') as string;
  check(callbackUrl, {
    'callback url returned': (u) => (u !== undefined && u !== null && u !== '') || fail('callback url missing'),
  });

  const code = new URL(callbackUrl).searchParams.get('code');
  const tokens = exchangeAuthorizationCode(code);
  check(tokens, {
    'oidc access token returned': (t) => t.accessToken !== undefined,
    'oidc id token returned': (t) => t.idToken !== undefined,
  });

  oidcPasswordSessionTrend.add(new Date().getTime() - start.getTime());
}

export function teardown(data: any) {
  removeOrg(data.org, data.adminTokens.accessToken);
}
