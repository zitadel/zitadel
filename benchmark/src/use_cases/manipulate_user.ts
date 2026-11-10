import { loginByUsernamePassword } from '../login_ui';
import { createOrg, removeOrg } from '../org';
import { createHuman, updateHuman, lockUser, deleteUser, User } from '../user';
import { Config } from '../config';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.info('setup: admin signed in');

  const org = await createOrg(tokens.accessToken!);
  console.info(`setup: org (${org.organizationId}) created`);

  return { tokens, org };
}

export default async function (data: any) {
  const human = await createHuman(`vu-${__VU}-${new Date(Date.now()).getTime()}`, data.org, data.tokens.accessToken);
  // updateHuman/lockUser/deleteUser each register their own check; repeating them
  // here would double count, and k6 aggregates checks by name.
  await updateHuman(
    {
      human: {
        profile: {
          nickName: `${new Date(Date.now()).toISOString()}`,
        },
      },
    },
    human.userId,
    data.org,
    data.tokens.accessToken,
  );
  await lockUser(human.userId, data.org, data.tokens.accessToken);
  await deleteUser(human.userId, data.org, data.tokens.accessToken);
}

export function teardown(data: any) {
  removeOrg(data.org, data.tokens.accessToken);
  console.info('teardown: org removed');
}
