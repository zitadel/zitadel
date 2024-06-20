import { loginByUsernamePassword } from '../login_ui';
import { createOrg, removeOrg, Org } from '../org';
import { User, getMyUser } from '../user';
import { Trend } from 'k6/metrics';
import { Config, MaxVUs } from '../config';
import { createProject, createProjectGrant } from '../project';

export const options = {
  scenarios: {
    contacts: {
      executor: 'per-vu-iterations',
      vus: 100,
      iterations: 100,
    },
  },
};

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.log('setup: admin signed in');

  const myUser = await getMyUser(tokens.accessToken!);
  const myOrg: Org = {organizationId: myUser.details.resourceOwner}

  const orgPromises = Array.from({ length: MaxVUs() }, (_, i) => {
    return createOrg(tokens.accessToken!, `${Date.now().toString()}-grant-${i}`);
  });
  
  const projectPromises = Array.from({ length: MaxVUs() }, (_, i) => {
    return createProject(`${Date.now().toString()}-project-${i}`, myOrg, tokens.accessToken!);
  });
    
  const orgs = await Promise.all(orgPromises);
  console.log(`setup: ${orgs.length} orgs created`);
    
  const projects = await Promise.all(projectPromises);
  console.log(`setup: ${projects.length} projects created`);

  return { tokens, orgs, projects };
}

const addProjectGrantTrend = new Trend('add_project_grant_duration', true);
export default async function (data: any) {  
  const start = new Date();
  await createProjectGrant(data.projects[__ITER], data.orgs[__VU - 1], [], data.tokens.accessToken);

  addProjectGrantTrend.add(new Date().getTime() - start.getTime());
}
  
export function teardown(data: any) {
  data.orgs.forEach((org: any) => {
    removeOrg(org, data.tokens.accessToken);
  });
}
