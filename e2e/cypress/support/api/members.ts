import { ZITADELTarget } from 'support/commands';
import { standardCreate, standardEnsureDoesntExist, standardEnsureExists, standardRemove, standardSearch } from './standard';

export function ensureHumanIsOrgMember(target: ZITADELTarget, userId: number, roles: string[]) {
  return standardEnsureExists(addOrgMember(target, userId, roles), () => searchOrgMembers(target, userId));
}

export function ensureHumanIsNotOrgMember(target: ZITADELTarget, userId: number, anyExistingRole: string) {
  return standardEnsureDoesntExist(
    ensureHumanIsOrgMember(target, userId, [anyExistingRole]),
    Cypress._.curry(removeOrgMember)(target),
    () => searchOrgMembers(target, userId),
  );
}

function addOrgMember(target: ZITADELTarget, userId: number, roles: string[]) {
  return standardCreate<number>(
    target,
    `${target.mgmtBaseURL}/orgs/me/members`,
    {
      userId: userId,
      roles: roles,
    },
    'userId',
  );
}

function searchOrgMembers(target: ZITADELTarget, userId: number) {
  return standardSearch<number>(
    target,
    `${target.mgmtBaseURL}/orgs/me/members/_search`,
    (entity) => entity.userId === userId,
    'userId',
  );
}

function removeOrgMember(target: ZITADELTarget, userId: number) {
  return standardRemove(target, `${target.mgmtBaseURL}/orgs/me/members/${userId}`);
}

export function ensureHumanIsProjectMember(target: ZITADELTarget, projectId: number, userId: number, roles: string[]) {
  return standardEnsureExists(addProjectMember(target, projectId, userId, roles), () =>
    searchProjectMembers(target, projectId, userId),
  );
}

export function ensureHumanIsNotProjectMember(
  target: ZITADELTarget,
  projectId: number,
  userId: number,
  anyExistingRole: string,
) {
  return standardEnsureDoesntExist(
    ensureHumanIsProjectMember(target, projectId, userId, [anyExistingRole]),
    Cypress._.curry(removeProjectMember)(target, projectId),
    () => searchProjectMembers(target, projectId, userId),
  );
}

function addProjectMember(target: ZITADELTarget, projectId: number, userId: number, roles: string[]) {
  return standardCreate<number>(
    target,
    `${target.mgmtBaseURL}/projects/${projectId}/members`,
    {
      userId: userId,
      roles: roles,
    },
    'userId',
  );
}

function searchProjectMembers(target: ZITADELTarget, projectId: number, userId: number) {
  return standardSearch<number>(
    target,
    `${target.mgmtBaseURL}/projects/${projectId}/members/_search`,
    (entity) => entity.userId === userId,
    'userId',
  );
}

function removeProjectMember(target: ZITADELTarget, projectId: number, userId: number) {
  return standardRemove(target, `${target.mgmtBaseURL}/projects/${projectId}/members/${userId}`);
}

export function ensureHumanIsGrantedProjectMember(
  target: ZITADELTarget,
  projectId: number,
  grantId: number,
  userId: number,
  roles: string[],
) {
  return standardEnsureExists(addGrantedProjectMember(target, projectId, grantId, userId, roles), () =>
    searchGrantedProjectMembers(target, projectId, grantId, userId),
  );
}

export function ensureHumanIsNotGrantedProjectMember(
  target: ZITADELTarget,
  projectId: number,
  grantId: number,
  userId: number,
  anyExistingRole: string,
) {
  return standardEnsureDoesntExist(
    ensureHumanIsGrantedProjectMember(target, projectId, grantId, userId, [anyExistingRole]),
    Cypress._.curry(removeGrantedProjectMember)(target, projectId, grantId),
    () => searchGrantedProjectMembers(target, projectId, grantId, userId),
  );
}

function addGrantedProjectMember(
  target: ZITADELTarget,
  projectId: number,
  grantId: number,
  userId: number,
  roles: string[],
) {
  return standardCreate<number>(
    target,
    `${target.mgmtBaseURL}/projects/${projectId}/grants/${grantId}/members`,
    {
      userId: userId,
      roles: roles,
    },
    'userId',
  );
}

function searchGrantedProjectMembers(target: ZITADELTarget, projectId: number, grantId: number, userId: number) {
  return standardSearch<number>(
    target,
    `${target.mgmtBaseURL}/projects/${projectId}/grants/${grantId}/members/_search`,
    (entity) => entity.userId === userId,
    'userId',
  );
}

function removeGrantedProjectMember(target: ZITADELTarget, projectId: number, grantId: number, userId: number) {
  return standardRemove(target, `${target.mgmtBaseURL}/projects/${projectId}/grants/${grantId}/members/${userId}`);
}
