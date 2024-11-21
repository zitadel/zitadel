import { loginByUsernamePassword } from '../login_ui';
import { createOrg, removeOrg } from '../org';
import {createMachine, User, addMachineKey} from '../user';
import {JWTProfileRequest, token, userinfo} from '../oidc';
import { Config } from '../config';
import encoding from 'k6/encoding';

const publicKey = encoding.b64encode(open('../.keys/key.pem.pub'));

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.info('setup: admin signed in');
  
  const org = await createOrg(tokens.accessToken!);
  console.info(`setup: org (${org.organizationId}) created`);
  
  const machine = await createMachine(`zitachine`, org, tokens.accessToken!);
  console.info(`setup: machine ${machine.userId} created`);
  const key = await addMachineKey(machine.userId, org, tokens.accessToken!, publicKey);
  console.info(`setup: key ${key.keyId} added`);

  return { tokens, machine: {userId: machine.userId, keyId: key.keyId}, org };
}

export default function (data: any) {
  token(new JWTProfileRequest(data.machine.userId, data.machine.keyId))
}

export function teardown(data: any) {
  removeOrg(data.org, data.tokens.accessToken);
  console.info('teardown: org removed');
}
