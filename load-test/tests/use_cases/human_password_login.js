import { loginByUsernamePassword } from '../login_ui.js';
import { createOrg, createHuman } from '../setup.js';
import { removeOrg } from '../teardown.js';
import { userinfo } from '../oidc.js';
import { Trend } from 'k6/metrics';
import { Config } from '../config.js';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin);
  console.log("admin signed in");
  
  const org = await createOrg(tokens.accessToken);
  console.log(`org (${org.organizationId}) created`);

  let humans = Array.from({length: __ENV.MAX_VUS || 1}, (_, i) => {
    return createHuman(`zitizen-${i}`, org, tokens.accessToken);
  });
  humans = await Promise.all(humans);
  humans = humans.map((user, i) => {
    return {userId: user.userId, loginName: user.loginNames[0], password: 'Password1!'};
  })
  console.log(`${humans.length} users created`);
  return {tokens, users: humans, org};
}

const humanPasswordLoginTrend = new Trend('human_password_login_duration', true);
export default function (data) {
  const start = new Date();
  const token = loginByUsernamePassword(data.users[__VU-1]);
  userinfo(token.accessToken);

  humanPasswordLoginTrend.add(new Date() - start);
}

export function teardown(data) {
  removeOrg(data.org, data.tokens.accessToken);
}

