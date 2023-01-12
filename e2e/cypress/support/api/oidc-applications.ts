import { ZITADELTarget } from 'support/commands';
import {
  standardCreate,
  standardEnsureDoesntExist,
  standardEnsureExists,
  standardRemove,
  standardSearch,
  standardUpdate,
} from './standard';

export function ensureOIDCAppDoesntExist(target: ZITADELTarget, projectId: number, name: string) {
  return standardEnsureDoesntExist(ensureOIDCAppExists(target, projectId, name), Cypress._.curry(remove)(target, projectId));
}

export function ensureOIDCAppExists(target: ZITADELTarget, projectId: number, name: string): Cypress.Chainable<number> {
  return standardEnsureExists(
    create(target, projectId, name),
     () => search(target, projectId, name),
    );
}

function create(target: ZITADELTarget, projectId: number, name: string): Cypress.Chainable<any> {
  return standardCreate(
    target,
    `${target.mgmtBaseURL}/projects/${projectId}/apps/oidc`,
    {
      name: name,
      redirectUris: ['https://e2eredirecturl.org'],
      responseTypes: ['OIDC_RESPONSE_TYPE_CODE'],
      grantTypes: ['OIDC_GRANT_TYPE_AUTHORIZATION_CODE'],
      authMethodType: 'OIDC_AUTH_METHOD_TYPE_NONE',
      postLogoutRedirectUris: ['https://e2elogoutredirecturl.org'],
    },
    'id',
  );
}

function search(target: ZITADELTarget, projectId: number, name: string): Cypress.Chainable<number> {
  return standardSearch(
    target,
    `${target.mgmtBaseURL}/projects/${projectId}/apps/_search`,
    (entity) => entity.name == name,
    'id',
  );
}

function remove(target: ZITADELTarget, projectId: number, id: number) {
  return standardRemove(target, `${target.mgmtBaseURL}/projects/${projectId}/apps/${id}`);
}
