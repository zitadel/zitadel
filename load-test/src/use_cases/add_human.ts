import { loginByUsernamePassword } from '../login_ui';
import { createOrg, removeOrg } from '../org';
import { User, createHuman } from '../user';
import { Trend } from 'k6/metrics';
import { Config, MaxVUs } from '../config';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.log('setup: admin signed in');

  const orgPromises = Array.from({ length: MaxVUs() }, (_, i) => {
    return createOrg(tokens.accessToken!, `${Date.now().toString()}-project-${i}`);
  });

  return { tokens, orgs: await Promise.all(orgPromises) };
}

const addHumanTrend = new Trend('add_human_duration', true);
export default async function (data: any) {  
  const start = new Date();
  await createHuman(`zitizen-${__VU}-${__ITER}`, data.orgs[__VU - 1], data.tokens.accessToken);

  addHumanTrend.add(new Date().getTime() - start.getTime());
}
  
export function teardown(data: any) {
  data.orgs.forEach((org: any) => {
    removeOrg(org, data.tokens.accessToken);
  });
}
