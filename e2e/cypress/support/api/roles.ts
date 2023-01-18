import { ZITADELTarget } from 'support/commands';
import { standardCreate, standardEnsureDoesntExist, standardEnsureExists, standardRemove, standardSearch } from './standard';

export function ensureRoleExists(target: ZITADELTarget, projectId: number, roleKey: string) {
  return standardEnsureExists(create(target, projectId, roleKey), () => search(target, projectId, roleKey));
}

export function ensureRoleDoesntExist(target: ZITADELTarget, projectId: number, roleKey: string) {
  return standardEnsureDoesntExist(
    ensureRoleExists(target, projectId, roleKey),
    Cypress._.curry(remove)(target, projectId),
    () => search(target, projectId, roleKey),
  );
}

function create(target: ZITADELTarget, projectId: number, roleKey: string) {
  return standardCreate<string>(
    target,
    `${target.mgmtBaseURL}/projects/${projectId}/roles`,
    {
      roleKey: roleKey,
      displayName: roleKey,
    },
    'key',
  );
}

function search(target: ZITADELTarget, projectId: number, roleKey: string) {
  return standardSearch<string>(
    target,
    `${target.mgmtBaseURL}/projects/${projectId}/roles/_search`,
    (entity) => entity.key === roleKey,
    'key',
  );
}

function remove(target: ZITADELTarget, projectId: number, roleKey: string) {
  return standardRemove(target, `${target.mgmtBaseURL}/projects/${projectId}/roles/${roleKey}`);
}
