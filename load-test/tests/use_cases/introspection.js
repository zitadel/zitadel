import { loginByUsernamePassword } from '../login_ui.js';
import { createOrg, createProject, createAPI, createAppKey } from '../setup.js';
import { removeOrg } from '../teardown.js';
import { introspect } from '../oidc.js';
import { Trend } from 'k6/metrics';
import { Config, MaxVUs } from '../config.js';
import { b64decode } from 'k6/encoding';
import zitadel from 'k6/x/zitadel';

export async function setup() {
  const adminTokens = loginByUsernamePassword(Config.admin);
  console.info("setup: admin signed in");
  
  const org = await createOrg(adminTokens.accessToken);
  console.info(`setup: org (${org.organizationId}) created`);

  let projects = Array.from({length: MaxVUs()}, (_, i) => {
     return createProject(`project-${i}`, org, adminTokens.accessToken);
  });
  projects = await Promise.all(projects);
  console.log(`setup: ${projects.length} projects created`)

  let apps = projects.map((project, i) => {
    return createAPI(`api-${i}`, project.id, org, adminTokens.accessToken);
  });
  apps = await Promise.all(apps);
  console.info(`setup: ${apps.length} apis created`);

  let keys = apps.map((app, i) => {
    return createAppKey(app.appId, projects[i].id, org, adminTokens.accessToken);
  });
  keys = await Promise.all(keys);
  console.info(`setup: ${keys.length} keys created`);
  
  let tokens = keys.map((key) => {
    return zitadel.jwtFromKey(b64decode(key.keyDetails, 'url', 's'), Config.host);
  });
  console.info(`setup: ${tokens.length} tokens generated`);

  return {adminTokens, tokens, org};
}

const humanPasswordLoginTrend = new Trend('machine_pat_login_duration', true);
export default function (data) {
  const token = introspect(data.tokens[__VU-1], data.adminTokens.accessToken);
}

export function teardown(data) {
  removeOrg(data.org, data.adminTokens.accessToken);
  console.info('teardown: org removed')
}