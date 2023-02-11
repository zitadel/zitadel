import { ensureItemExists } from './ensure';
import { getOrgUnderTest } from './orgs';
import { API } from './types';

export function ensureProjectGrantExists(api: API, foreignOrgId: string, foreignProjectId: string) {
  return getOrgUnderTest(api).then((orgUnderTest) => {
    return ensureItemExists(
      api,
      `${api.mgmtBaseURL}/projectgrants/_search`,
      (grant: any) => grant.grantedOrgId == orgUnderTest && grant.projectId == foreignProjectId,
      `${api.mgmtBaseURL}/projects/${foreignProjectId}/grants`,
      { granted_org_id: orgUnderTest },
      foreignOrgId,
      'grantId',
      'grantId',
    );
  });
}
