import { Context } from 'support/commands';
import { ensureItemDoesntExist, ensureItemExists } from './ensure';
import { getOrgUnderTest } from './orgs';
import { API, Entity } from './types';

export function ensureProjectGrantExists(ctx: Context, foreignOrgId: string, foreignProjectId: string) {
  return getOrgUnderTest(ctx).then((orgUnderTest) => {
    return ensureItemExists(
      ctx.api,
      `${ctx.api.mgmtBaseURL}/projectgrants/_search`,
      (grant: any) => grant.grantedOrgId == orgUnderTest && grant.projectId == foreignProjectId,
      `${ctx.api.mgmtBaseURL}/projects/${foreignProjectId}/grants`,
      { granted_org_id: orgUnderTest },
      foreignOrgId,
      'grantId',
      'grantId',
    );
  });
}

export function ensureProjectGrantDoesntExist(ctx: Context, projectId: number, foreignOrgId: string) {
  return getOrgUnderTest(ctx).then((orgUnderTest) => {
    console.log('removing grant to foreignOrgId', foreignOrgId, 'in orgUnderTest', orgUnderTest, 'projectId', projectId);
    return ensureItemDoesntExist(
      ctx.api,
      `${ctx.api.mgmtBaseURL}/projectgrants/_search`,
      (grant: any) => grant.grantedOrgId == foreignOrgId && grant.projectId == projectId,
      (grant: any) => `${ctx.api.mgmtBaseURL}/projects/${projectId}/grants/${grant.grantId}`,
      orgUnderTest.toString(),
    );
  });
}
