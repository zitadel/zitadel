import { loginByUsernamePassword } from '../login_ui';
import { createAPI, createAppKey } from '../app';
import { createProject } from '../project';
import { createOrg, removeOrg } from '../org';
import { introspect } from '../oidc';
import { Config, MaxVUs, Client } from '../config';
import { b64decode } from 'k6/encoding';
// @ts-ignore Import module
import zitadel from 'k6/x/zitadel';
import { User } from '../user';

export async function setup() {
  const adminTokens = loginByUsernamePassword(Config.admin as User);
  console.info('setup: admin signed in');

  const org = await createOrg(adminTokens.accessToken!);
  console.info(`setup: org (${org.organizationId}) created`);

  const projectPromises = Array.from({ length: MaxVUs() }, (_, i) => {
    return createProject(`project-${i}`, org, adminTokens.accessToken!);
  });
  const projects = await Promise.all(projectPromises);
  console.log(`setup: ${projects.length} projects created`);

  const apis = await Promise.all(
    projects.map((project, i) => {
      return createAPI(`api-${i}`, project.id, org, adminTokens.accessToken!);
    }),
  );
  console.info(`setup: ${apis.length} apis created`);

  const keys = await Promise.all(
    apis.map((api, i) => {
      return createAppKey(api.appId, projects[i].id, org, adminTokens.accessToken!);
    }),
  );
  console.info(`setup: ${keys.length} keys created`);

  const tokens = keys.map((key) => {
    return zitadel.jwtFromKey(b64decode(key.keyDetails, 'url', 's'), Config.host);
  });
  console.info(`setup: ${tokens.length} tokens generated`);

  const client = { ...Client() };
  client.scope = [client.scope, ...projects.map((p) => `urn:zitadel:iam:org:project:id:${p.id}:aud`)].join(',');
  console.info('setup: login user with scope %s', client.scope);

  const userTokens = loginByUsernamePassword(Config.admin as User, client);
  console.info('setup: user signed in');

  return { adminTokens, userTokens, tokens, org };
}

export default function (data: any) {
  introspect(data.tokens[__VU - 1], data.userTokens.accessToken);
}

export function teardown(data: any) {
  removeOrg(data.org, data.adminTokens.accessToken);
  console.info('teardown: org removed');
}
