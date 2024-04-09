import { loginByUsernamePassword } from '../login_ui.js';
import { createOrg } from '../setup.js';
import { createHuman, updateHuman, lockUser, deleteUser } from '../user.js';
import { removeOrg } from '../teardown.js';
import { Config } from '../config.js';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin);
  console.info("setup: admin signed in");

  const org = await createOrg(tokens.accessToken);
  console.info(`setup: org (${org.organizationId}) created`);

  return {tokens, org};
}

export default async function(data) {
    let human = await createHuman(`vu-${__VU}`, data.org, data.tokens.accessToken);
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
    const lockRes = await lockUser(human.userId, data.org, data.tokens.accessToken);
    const deleteRes = await deleteUser(human.userId, data.org, data.tokens.accessToken);
}

export function teardown(data) {
    removeOrg(data.org, data.tokens.accessToken);
    console.info('teardown: org removed')
  }