import { loginByUsernamePassword } from '../login_ui';
import { createOrg, removeOrg } from '../org';
import {createMachine, User, addMachineSecret} from '../user';
import {clientCredentials, userinfo} from '../oidc';
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

  let credentials = (
    await Promise.all(
      machines.map((machine) => {
        return addMachineSecret(machine.userId, org, tokens.accessToken!);
      }),
    )
  ).map((credentials, i) => {
    return { userId: machines[i].userId, loginName: machines[i].loginName, password: credentials.clientSecret };
  });
  console.info(`setup: secrets added`);

  return { tokens, machines: credentials, org };
}

export default function (data: any) {
  clientCredentials(data.machines[__VU - 1].loginName, data.machines[__VU - 1].password)
    .then((token) => {
      userinfo(token.accessToken!)
    })
}

export function teardown(data: any) {
  removeOrg(data.org, data.tokens.accessToken);
  console.info('teardown: org removed');
}
