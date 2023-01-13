import { ZITADELTarget } from 'support/commands';
import { standardCreate, standardEnsureExists, standardSearch } from './standard';

export function ensureProjectGrantExists(target: ZITADELTarget, projectId: number, grantOrgId): Cypress.Chainable<any> {
  return standardEnsureExists(create(target, projectId, grantOrgId), () => search(target, projectId, grantOrgId));
}

function create(target: ZITADELTarget, projectId: number, grantedOrgId: number): Cypress.Chainable<number> {
  return standardCreate(
    target,
    `${target.mgmtBaseURL}/projects/${projectId}/grants`,
    { grantedOrgId: grantedOrgId.toString() },
    'grantId',
  );
}

function search(target: ZITADELTarget, projectId: number, grantedOrgId: number): Cypress.Chainable<number> {
  return standardSearch(
    target,
    `${target.mgmtBaseURL}/projects/${projectId}/grants/_search`,
    (entity) => entity.projectId == projectId && entity.grantedOrgId == grantedOrgId,
    'grantId',
  );
}
