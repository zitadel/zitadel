import { loginByUsernamePassword } from '../login_ui';
import { createOrg } from '../org';
import { User, createMachine } from '../user';
import { Config, MaxVUs } from '../config';
import { createProject, Project } from '../project';
import { addUserGrant } from '../user_grant';

export async function setup() { 
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.info('setup: admin signed in');

  const org = await createOrg(tokens.accessToken!);
  console.info(`setup: org (${org.organizationId}) created`);

  const projects = await Promise.all(
    Array.from({ length: 50 }, (_, i) => {
      return createProject(`project-${i}`, org, tokens.accessToken!);
    }),
  );
  console.log(`setup: ${projects.length} projects created`);

  let machines = (
    await Promise.all(
      Array.from({ length: MaxVUs() }, async (_, i) => {
        return await createMachine(`zitachine-${i}`, org, tokens.accessToken!);
      }),
    )
  ).map((machine) => {
    return { userId: machine.userId, loginName: machine.loginNames[0] };
  });
  console.log(`setup: ${machines.length} machines created`);

  return { tokens, org, machines, projects };
}

export default async function (data: any) {
  const machine = await createMachine(`zitachine-${__VU}-${__ITER}`, data.org, data.tokens.accessToken!);
  let userGrants = await Promise.all(
    data.projects.map((project: Project) => {
      return addUserGrant(data.org, machine.userId, project, [], data.tokens.accessToken!);
    }),
  );

  return { userGrants };
}

export function teardown(data: any) {
  // removeOrg(data.org, data.tokens.accessToken);
  console.info('teardown: org removed');
}