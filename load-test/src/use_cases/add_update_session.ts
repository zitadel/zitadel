import { loginByUsernamePassword } from '../login_ui';
import { createOrg, removeOrg } from '../org';
import { User } from '../user';
import { Trend } from 'k6/metrics';
import { Config } from '../config';
import { createEmptySession, updateSession } from '../session';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.log('setup: admin signed in');

  const org = await createOrg(tokens.accessToken!);
  console.log(`setup: org (${org.organizationId}) created`);

  return { tokens, org };
}

const addAndUpdateSessionTrend = new Trend('add_and_update_session_duration', true);
export default async function (data: any) {  
  const start = new Date();
  const session = await createEmptySession(data.org, data.tokens.accessToken!);;
  await updateSession(session, data.org, data.tokens.accessToken);

  addAndUpdateSessionTrend.add(new Date().getTime() - start.getTime());
}
  
export function teardown(data: any) {
  removeOrg(data.org, data.tokens.accessToken);
}
