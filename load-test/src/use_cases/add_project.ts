import { loginByUsernamePassword } from '../login_ui';
import { createOrg, removeOrg } from '../org';
import { User } from '../user';
import { Trend } from 'k6/metrics';
import { Config, MaxVUs } from '../config';
import { createProject } from '../project';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.log('setup: admin signed in');

  const orgPromises = Array.from({ length: MaxVUs() }, (_, i) => {
    return createOrg(tokens.accessToken!, `${Date.now().toString()}-project-${i}`);
  });

  return { tokens, orgs: await Promise.all(orgPromises) };
}

const addProjectTrend = new Trend('add_project_duration', true);
export default async function (data: any) {  
  const start = new Date();
  await createProject(`${__VU}_${__ITER}`, data.orgs[__VU - 1], data.tokens.accessToken);

  addProjectTrend.add(new Date().getTime() - start.getTime());
}
  
export function teardown(data: any) {
  data.orgs.forEach((org: any) => {
    removeOrg(org, data.tokens.accessToken);
  });
}
