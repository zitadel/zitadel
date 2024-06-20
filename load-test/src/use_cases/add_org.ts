import { loginByUsernamePassword } from '../login_ui';
import { createOrg, removeOrg } from '../org';
import { User } from '../user';
import { Trend } from 'k6/metrics';
import { Config } from '../config';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.log('setup: admin signed in');

  return { tokens, now: Date.now().toString() };
}

const addOrgTrend = new Trend('add_org_duration', true);
export default async function (data: any) {  
  const start = new Date();
  await createOrg(data.tokens.accessToken!, `${data.now}_${__VU}_${__ITER}`);

  addOrgTrend.add(new Date().getTime() - start.getTime());
}
  
export function teardown(data: any) {
  removeOrg(data.org, data.tokens.accessToken);
}
