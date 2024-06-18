import { loginByUsernamePassword } from '../login_ui';
import { createOrg, removeOrg } from '../org';
import { User } from '../user';
import { Trend } from 'k6/metrics';
import { Config, MaxVUs } from '../config';
import { createEmptySession, updateSession } from '../session';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.log('setup: admin signed in');

  const org = await createOrg(tokens.accessToken!);
  console.log(`setup: org (${org.organizationId}) created`);

  const sessions = await Promise.all(Array.from({ length: MaxVUs() }, () => {
    return createEmptySession(org, tokens.accessToken!);
  }));
  console.log(`setup: ${sessions.length} sessions created`);

  return { tokens, sessions, org };
}

const updateSessionTrend = new Trend('update_session_duration', true);
export default async function (data: any) {  
  const start = new Date();
  await updateSession(data.sessions[__VU - 1], data.org, data.tokens.accessToken);

  updateSessionTrend.add(new Date().getTime() - start.getTime());
}
  
export function teardown(data: any) {
  removeOrg(data.org, data.tokens.accessToken);
}
