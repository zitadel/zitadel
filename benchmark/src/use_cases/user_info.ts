import { loginByUsernamePassword } from '../login_ui';
import { userinfo } from '../oidc';
import { Config } from '../config';
import { User, createHuman } from '../user';
import { createOrg, removeOrg } from '../org';

export async function setup() {
  const adminTokens = loginByUsernamePassword(Config.admin as User);
  console.info('setup: admin signed in');

  const org = await createOrg(adminTokens.accessToken!);
  console.info(`setup: org (${org.organizationId}) created`);

  const user = await createHuman('gigi', org, adminTokens.accessToken!);
  console.info(`setup: user (${user.userId}) created`);

  // The measured call needs the user's token, but only the admin may remove the org. Tearing
  // down with the user's token failed with `membership not found (AUTHZ-cdgFk)` and leaked the
  // org on every run.
  return {
    org,
    adminTokens,
    tokens: loginByUsernamePassword({ loginName: user.loginNames[0], password: 'Password1!' } as User),
  };
}

export default function (data: any) {
  userinfo(data.tokens.accessToken);
}

export function teardown(data: any) {
  removeOrg(data.org, data.adminTokens.accessToken);
  console.info('teardown: org removed');
}
