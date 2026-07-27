import { check } from 'k6';

import { Config } from '../config';
import { loginByUsernamePassword } from '../login_ui';
import { createOrg, Org, removeOrg } from '../org';
import { createHuman, Human, listUsers, User } from '../user';

const userAmount = parseInt(__ENV.USER_AMOUNT) || 2500;

type SetupData = {
  tokens: { accessToken?: string };
  org: Org;
  targetLoginName: string;
};

export async function setup(): Promise<SetupData> {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.info('setup: admin signed in');

  const org = await createOrg(tokens.accessToken!);
  console.info(`setup: org (${org.organizationId}) created`);

  const users: Human[] = [];
  await Promise.all(
    Array.from({ length: userAmount }, async (_, i) => {
      const user = await createHuman(`zitizen-${i}`, org, tokens.accessToken!);
      users.push(user);
      if (i % 10 === 0) {
        console.log(`setup: ${i} of ${userAmount} users setup`);
      }
    }),
  );
  console.info(`setup: ${users.length} users created`);

  const targetLoginName = users[0].loginNames[0];
  console.info(`setup: target login name ${targetLoginName}`);

  return { tokens, org, targetLoginName };
}

export default async function (data: SetupData) {
  const result = await listUsers(
    {
      query: { limit: 2 },
      queries: [
        {
          loginNameQuery: {
            loginName: data.targetLoginName,
            method: 'TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE',
          },
        },
        {
          organizationIdQuery: {
            organizationId: data.org.organizationId,
          },
        },
      ],
    },
    data.tokens.accessToken!,
  );

  check(result, {
    'exact one user found': (res) => (res.result?.length ?? 0) === 1 || res.details.totalResult == 1,
  }) ||
    console.log(
      `unexpected list users result. expected 1 user but got result=${result.result?.length} totalResult=${result.details.totalResult}`,
    );
}

export function teardown(data: SetupData) {
  removeOrg(data.org, data.tokens.accessToken!);
  console.info('teardown: org removed');
}
