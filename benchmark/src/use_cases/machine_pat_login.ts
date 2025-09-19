import { loginByUsernamePassword } from '../login_ui';
import { createOrg, removeOrg } from '../org';
import { createMachine, addMachinePat, User } from '../user';
import { userinfo } from '../oidc';
import { Config, MaxVUs } from '../config';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.info('setup: admin signed in');

  const org = await createOrg(tokens.accessToken!);
  console.info(`setup: org (${org.organizationId}) created`);

  let machines = (
    await Promise.all(
      Array.from({ length: MaxVUs() }, (_, i) => {
        return createMachine(`zitachine-${i}`, org, tokens.accessToken!);
      }),
    )
  ).map((machine) => {
    return { userId: machine.userId, loginName: machine.loginNames[0] };
  });
  console.info(`setup: ${machines.length} machines created`);

  let pats = (
    await Promise.all(
      machines.map((machine) => {
        return addMachinePat(machine.userId, org, tokens.accessToken!);
      }),
    )
  ).map((pat, i) => {
    return { userId: machines[i].userId, loginName: machines[i].loginName, pat: pat.token };
  });
  console.info(`setup: Pats added`);

  return { tokens, machines: pats, org };
}

export default function (data: any) {
  userinfo(data.machines[__VU - 1].pat);
}

export function teardown(data: any) {
  removeOrg(data.org, data.tokens.accessToken);
  console.info('teardown: org removed');
}
