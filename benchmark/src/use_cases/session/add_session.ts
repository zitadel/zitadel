import { loginByUsernamePassword } from '../../login_ui';
import { createOrg, removeOrg } from '../../org';
import { User, createHuman, createMachine } from '../../user';
import { Trend } from 'k6/metrics';
import { Config, MaxVUs } from '../../config';
import { createSession } from '../../session';
import { check } from 'k6';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.log('setup: admin signed in');

  const org = await createOrg(tokens.accessToken!);
  console.log(`setup: org (${org.organizationId}) created`);

  const humanPromises = Array.from({ length: MaxVUs() }, (_, i) => {
    return createHuman(`zitizen-${i}`, org, tokens.accessToken!);
  });

  const humans = (await Promise.all(humanPromises)).map((user) => {
    return { userId: user.userId, loginName: user.loginNames[0], password: 'Password1!' };
  });
  console.log(`setup: ${humans.length} users created`);
  return { tokens, users: humans, org };
}

const addSessionTrend = new Trend('add_session_duration', true);
export default async function (data: any) {
  const start = new Date();
  const session = await createSession(data.org, data.tokens.accessToken, {
    user: {
      userId: data.users[__VU - 1].userId,
    },
  });

  check(session, {
    'add session is status ok': (s) => s.id !== '',
  });

  addSessionTrend.add(new Date().getTime() - start.getTime());
}

export function teardown(data: any) {
  removeOrg(data.org, data.tokens.accessToken);
}
