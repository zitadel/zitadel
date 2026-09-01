import { check } from 'k6';

import { Config } from '../config';
import { loginByUsernamePassword } from '../login_ui';
import { createOrg, Org, removeOrg } from '../org';
import { mapPool } from '../pool';
import { createHuman, listUsers, User } from '../user';

const userAmount = parseInt(__ENV.USER_AMOUNT) || 2500;
// Unbounded Promise.all over large USER_AMOUNT exhausts ephemeral ports
// (dial tcp ... can't assign requested address). Keep a modest in-flight cap.
const setupConcurrency = parseInt(__ENV.SETUP_CONCURRENCY) || 50;

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

  const progressEvery = Math.max(1, Math.floor(userAmount / 100));
  const users = await mapPool(
    Array.from({ length: userAmount }, (_, i) => i),
    setupConcurrency,
    async (i) => {
      const user = await createHuman(`zitizen-${i}`, org, tokens.accessToken!);
      if (i % progressEvery === 0 || i === userAmount - 1) {
        console.log(`setup: ${i + 1} of ${userAmount} users setup`);
      }
      return user;
    },
  );
  console.info(`setup: ${users.length} users created (concurrency=${setupConcurrency})`);

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
