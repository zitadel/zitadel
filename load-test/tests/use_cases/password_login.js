import { loginByUsernamePassword } from '../login_ui.js';
import { createOrg, createUser } from '../setup.js';
import { removeOrg } from '../teardown.js';
import { userinfo } from '../oidc.js';
import { Trend } from 'k6/metrics';

const admin = JSON.parse(open('../data/admin.json'));

export async function setup() {
  const tokens = loginByUsernamePassword(admin);
  console.log("admin signed in");
  
  const org = await createOrg(tokens.accessToken);
  console.log(`org (${org.organizationId}) created`);

  let users = Array.from({length: __ENV.MAX_VUS || 1}, (_, i) => {
    return createUser(`zitizen-${i}`, org, tokens.accessToken);
  });
  users = await Promise.all(users);
  users = users.map((user, i) => {
    return {userId: user.userId, loginName: user.loginNames[0], password: 'Password1!'};
  })
  console.log(`${users.length} users created`);
  return {tokens, users: users, org};
}

const passwordLoginTrend = new Trend('login_password_login_duration', true);
export default function (data) {
  const start = new Date();
  const token = loginByUsernamePassword(data.users[__VU-1]);
  userinfo(token.accessToken);

  passwordLoginTrend.add(new Date() - start);
}

export function teardown(data) {
  removeOrg(data.org, data.tokens.accessToken);
}

