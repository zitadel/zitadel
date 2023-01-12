import { ZITADELTarget } from 'support/commands';
import { ensureItemExists } from './ensure';

export function ensureProjectGrantExists(
  api: ZITADELTarget,
  projectId: number,
  grantOrgId: number,
): Cypress.Chainable<number> {
  return ensureItemExists(
    api,
    `${api.mgmtBaseURL}/projectgrants/_search`,
    (grant: any) => grant.grantedOrgId == api.headers['x-zitadel-orgid'] && grant.projectId == projectId,
    `${api.mgmtBaseURL}/projects/${projectId}/grants`,
    { granted_org_id: api.headers['x-zitadel-orgid'] },
    grantOrgId,
    'grantId',
    'grantId',
  );
}
