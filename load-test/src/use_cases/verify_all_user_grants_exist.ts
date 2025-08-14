import { loginByUsernamePassword } from '../login_ui';
import { createOrg, removeOrg } from '../org';
import { createHuman, updateHuman, lockUser, deleteUser, User, createMachine } from '../user';
import { Config, MaxVUs } from '../config';
import { check } from 'k6';
import { createProject, Project } from '../project';
import { addUserGrant } from '../user_grant';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.info('setup: admin signed in');

  const org = await createOrg(tokens.accessToken!);
  console.info(`setup: org (${org.organizationId}) created`);

  const projects = await Promise.all(Array.from({ length: 600 }, (_, i) => {
    return createProject(`project-${i}`, org, tokens.accessToken!);
  }));
  console.log(`setup: ${projects.length} projects created`);

  let machines = (
    await Promise.all(
      Array.from({ length: MaxVUs() }, (_, i) => {
        return createMachine(`zitachine-${i}`, org, tokens.accessToken!);
      }),
    )
  ).map((machine) => {
    return { userId: machine.userId, loginName: machine.loginNames[0] };
  });

  return { tokens, org, machines, projects };
}

export default async function (data: any) {
  let userGrants = await Promise.all(
    data.projects.map((project: Project) => {
      return addUserGrant(data.org, data.machines[__VU - 1].userId, project, [], data.tokens.accessToken!);
    })
  );

  console.log(`${userGrants.length} user grants created`);

  return { userGrants };
}

export function teardown(data: any) {
  // removeOrg(data.org, data.tokens.accessToken);
  console.info('teardown: org removed');
}
