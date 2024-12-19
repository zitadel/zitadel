import { loginByUsernamePassword } from '../login_ui';
import { createOrg, removeOrg } from '../org';
import {createMachine, User, addMachineKey} from '../user';
import {JWTProfileRequest, token, userinfo} from '../oidc';
import { Config, MaxVUs } from '../config';
import encoding from 'k6/encoding';

const publicKey = encoding.b64encode(open('../.keys/key.pem.pub'));

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
  
  let keys = (
    await Promise.all(
      machines.map((machine) => {
        return addMachineKey(
          machine.userId, 
          org, 
          tokens.accessToken!,
          publicKey,        
        );
      }),
    )
  ).map((key, i) => {
    return { userId: machines[i].userId, keyId: key.keyId };
  });
  console.info(`setup: ${keys.length} keys added`);

  return { tokens, machines: keys, org };
}

export default function (data: any) {
  token(new JWTProfileRequest(data.machines[__VU - 1].userId, data.machines[__VU - 1].keyId))
}

export function teardown(data: any) {
  removeOrg(data.org, data.tokens.accessToken);
  console.info('teardown: org removed');
}
