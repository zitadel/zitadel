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
  console.info('teardown: org is not removed to verify correctness of projections, do not forget to remove the org afterwards');
}

/** 
 * To verify the correctness of the projections you can use the following statements:
 * 
 * set the owner of the events:
 * 
 * set my.owner = '<org id of the created org>';
 * 
 * check if the amount of events is the same as amount of objects
 * 
 *  select * from (
    select 'projections.user_grants5', count(*) from projections.user_grants5 where resource_owner = (select current_setting('my.owner'))
    union all
    select 'projections.users14', count(*) from projections.users14 where resource_owner = (select current_setting('my.owner'))
    union all
    select 'projections.sessions8', count(*) from projections.sessions8 where user_resource_owner = (select current_setting('my.owner'))
    union all
    select aggregate_type, count(*) from eventstore.events2
    where
        aggregate_type in ('user', 'usergrant', 'session')
        and event_type in ('user.machine.added', 'user.human.added', 'user.grant.added', 'session.user.checked')
        and (owner = (select current_setting('my.owner'))
            OR payload->>'userResourceOwner' = (select current_setting('my.owner'))
        )
        group by aggregate_type
) order by 2;
 */