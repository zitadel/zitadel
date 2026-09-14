import { loginByUsernamePassword } from '../login_ui';
import { createOrg, removeOrg } from '../org';
import { mapPool } from '../pool';
import { createHuman, User, createMachine, setUserMetadata, listUsers } from '../user';
import { Config } from '../config';
import { check } from 'k6';
import encoding from 'k6/encoding';

const userAmount = parseInt(__ENV.USER_AMOUNT) || 2500;
const setupConcurrency = parseInt(__ENV.SETUP_CONCURRENCY) || 50;

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.info('setup: admin signed in');

  const org = await createOrg(tokens.accessToken!);
  console.info(`setup: org (${org.organizationId}) created`);

  const progressEvery = Math.max(1, Math.floor(userAmount / 100));
  const users = await mapPool(
    Array.from({ length: userAmount }, (_, i) => i),
    setupConcurrency,
    async (i) => {
      let user: User;
      let type: 'human' | 'machine';
      if (i % 2 === 0) {
        user = await createHuman(`zitizen-${i}`, org, tokens.accessToken!);
        type = 'human';
      } else {
        user = await createMachine(`zitachine-${i}`, org, tokens.accessToken!);
        type = 'machine';
      }
      await setUserMetadata(
        [
          { key: 'type', value: encoding.b64encode(type, 'rawurl') },
          { key: 'org', value: encoding.b64encode(org.organizationId, 'rawurl') },
          { key: 'id', value: encoding.b64encode(user.userId, 'rawurl') },
        ],
        user.userId,
        tokens.accessToken!,
      );
      if (i % progressEvery === 0 || i === userAmount - 1) {
        console.log(`setup: ${i + 1} of ${userAmount} users setup`);
      }
      return user;
    },
  );
  console.info(`setup: ${users.length} users created (concurrency=${setupConcurrency})`);

  return { tokens, org, users };
}

export default async function (data: any) {
  const result = await listUsers(
    {
      queries: [{ metadataKeyFilter: { key: 'org', method: 'TEXT_FILTER_METHOD_EQUALS' } }],
    },
    data.tokens.accessToken!,
  );
  check(result, {
    'total result length': (res) => res.details.totalResult == userAmount,
  }) || console.log(`unexpected amount of users. expected ${userAmount} but got ${result.details.totalResult}`);
}

export function teardown(data: any) {
  removeOrg(data.org, data.tokens.accessToken);
  console.info('teardown: org removed');
}
