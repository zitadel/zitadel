import { loginByUsernamePassword } from '../login_ui.js';
import { createOrg } from '../setup.js';
import { createHuman, updateHuman, lockUser, deleteUser } from '../user.js';
import { removeOrg } from '../teardown.js';
import { Config } from '../config.js';
import { check, fail } from 'k6';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin);
  console.info("setup: admin signed in");

  const org = await createOrg(tokens.accessToken);
  console.info(`setup: org (${org.organizationId}) created`);

  return {tokens, org};
}

export default async function(data) {
    const human = await createHuman(`vu-${__VU}`, data.org, data.tokens.accessToken);
    console.log(human.userId);
    const updateRes = await updateHuman(
        {
            profile: {
                nickName: `${new Date(Date.now()).toISOString()}`
            }
        }, 
        human.userId,
        data.org,
        data.tokens.accessToken
    );
    check(updateRes, {
        "update user is status ok": (r) => r.status >= 200 && r.status < 300
    });

    const lockRes = await lockUser(human.userId, data.org, data.tokens.accessToken);
    check(lockRes, {
        "lock user is status ok": (r) => r.status >= 200 && r.status < 300
    });

    const deleteRes = await deleteUser(human.userId, data.org, data.tokens.accessToken);
    check(deleteRes, {
        "delete user is status ok": (r) => r.status >= 200 && r.status < 300
    });
}

export function teardown(data) {
    removeOrg(data.org, data.tokens.accessToken);
    console.info('teardown: org removed')
}
