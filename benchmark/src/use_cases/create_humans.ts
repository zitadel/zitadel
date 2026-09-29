import { Config } from '../config';
import { loginByUsernamePassword } from '../login_ui';
import { createOrg, Org } from '../org';
import { createHuman, User } from '../user';

type SetupData = {
  tokens: { accessToken?: string };
  org: Org;
};

export async function setup(): Promise<SetupData> {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.info('setup: admin signed in');

  const org = await createOrg(tokens.accessToken!);
  console.info(`setup: org (${org.organizationId}) created`);

  return { tokens, org };
}

export default async function (data: SetupData) {
  await createHuman(`zitizen3-${ __VU }-${ __ITER }`, data.org, data.tokens.accessToken!);
}

