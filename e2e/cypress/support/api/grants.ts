import { Context } from 'support/commands';
import { ensureItemExists } from './ensure';
import { getOrgUnderTest } from './orgs';

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
