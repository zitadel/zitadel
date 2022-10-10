import { ensureItemExists } from './ensure';
import { getOrgUnderTest } from './orgs';
import { API } from './types';

export function ensureProjectGrantExists(
  api: API,
  foreignOrgId: number,
  foreignProjectId: number,
): Cypress.Chainable<number> {
  return getOrgUnderTest(api).then((orgUnderTest) => {
    return ensureItemExists(
      api,
      `${api.mgntBaseURL}projectgrants/_search`,
      (grant: any) => grant.grantedOrgId == orgUnderTest && grant.projectId == foreignProjectId,
      `${api.mgntBaseURL}projects/${foreignProjectId}/grants`,
      { granted_org_id: orgUnderTest },
      foreignOrgId,
      'grantId',
      'grantId',
    );
  });
}
