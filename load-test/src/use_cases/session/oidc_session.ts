import { loginByUsernamePassword, initLogin } from '../../login_ui';
import { createOrg, removeOrg } from '../../org';
import { User, addMachineKey, createMachine } from '../../user';
import { Trend } from 'k6/metrics';
import { Config } from '../../config';
import { check } from 'k6';
import { finalizeAuthRequest, JWTProfileRequest, token } from '../../oidc';
import { createSession } from '../../session';
import encoding from 'k6/encoding';
import { addIAMMember } from '../../membership';

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

  return { tokens, user: loginUser, key: loginUserKey, org, adminTokens };
}

// implements the flow described in
// https://zitadel.com/docs/guides/integrate/login-ui/oidc-standard
const addSessionTrend = new Trend('oidc_session_duration', true);
export default async function (data: any) {
  const start = new Date();
  const authorizeResponse = initLogin(data.tokens.info.client_id);

  const authRequestId = new URLSearchParams(authorizeResponse.headers['Location']).values().next().value;
  check(authRequestId, {
    'auth request id returned': (s) => s !== '',
  });

  const session = await createSession(data.org, data.tokens.accessToken, {
    user: {
      userId: data.user.userId,
    },
  });
  await finalizeAuthRequest(authRequestId!, session, data.tokens!);

  addSessionTrend.add(new Date().getTime() - start.getTime());
}

export function teardown(data: any) {
  removeOrg(data.org, data.adminTokens.accessToken);
}
