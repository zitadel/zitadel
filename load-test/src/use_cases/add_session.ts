import { loginByUsernamePassword } from '../login_ui';
import { createOrg, removeOrg } from '../org';
import { User, createHuman } from '../user';
import { Trend } from 'k6/metrics';
import { Config } from '../config';
import { createSession } from '../session';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.log('setup: admin signed in');

  const org = await createOrg(tokens.accessToken!);
  console.log(`setup: org (${org.organizationId}) created`);

  // console.log(`setup: ${humans.length} users created`);
  return { tokens, user: await createHuman(`zitizen`, org, tokens.accessToken!), org };
}

const addSessionTrend = new Trend('add_session_duration', true);
export default async function (data: any) {  
  const start = new Date();
  await createSession(data.user, data.org, data.tokens.accessToken);

  addSessionTrend.add(new Date().getTime() - start.getTime());
}
  
export function teardown(data: any) {
  removeOrg(data.org, data.tokens.accessToken);
}
