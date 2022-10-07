import { ensureItemDoesntExist, ensureItemExists } from './ensure';
import { findFromList, searchSomething } from './search';
import { API } from './types';

export function ensureHumanIsNotOrgMember(api: API, username: string): Cypress.Chainable<number> {
  return ensureItemDoesntExist(
    api,
    `${api.mgntBaseURL}orgs/me/members/_search`,
    (member: any) => (<string>member.preferredLoginName).startsWith(username),
    (member) => `${api.mgntBaseURL}orgs/me/members/${member.userId}`,
  );
}

export function ensureHumanIsOrgMember(api: API, username: string, roles: string[]): Cypress.Chainable<number> {
  return searchSomething(
    api,
    `${api.mgntBaseURL}users/_search`,
    'POST',
    findFromList((user) => {
      return user.userName == username;
    }),
  ).then((user) => {
    return ensureItemExists(
      api,
      `${api.mgntBaseURL}orgs/me/members/_search`,
      (member: any) => member.userId == user.entity.id,
      `${api.mgntBaseURL}orgs/me/members`,
      {
        userId: user.entity.id,
        roles: roles,
      },
    );
  });
}

export function ensureHumanIsNotProjectMember(
  api: API,
  projectId: string,
  username: string,
  grantId?: number,
): Cypress.Chainable<number> {
  return ensureItemDoesntExist(
    api,
    `${api.mgntBaseURL}projects/${projectId}/${grantId ? `grants/${grantId}/` : ''}members/_search`,
    (member: any) => (<string>member.preferredLoginName).startsWith(username),
    (member) => `${api.mgntBaseURL}projects/${projectId}${grantId ? `grants/${grantId}/` : ''}/members/${member.userId}`,
  );
}

export function ensureHumanIsProjectMember(
  api: API,
  projectId: string,
  username: string,
  roles: string[],
  grantId?: number,
): Cypress.Chainable<number> {
  return searchSomething(
    api,
    `${api.mgntBaseURL}users/_search`,
    'POST',
    findFromList((user) => {
      return user.userName == username;
    }),
  ).then((user) => {
    return ensureItemExists(
      api,
      `${api.mgntBaseURL}projects/${projectId}/${grantId ? `grants/${grantId}/` : ''}members/_search`,
      (member: any) => member.userId == user.entity.id,
      `${api.mgntBaseURL}projects/${projectId}/${grantId ? `grants/${grantId}/` : ''}members`,
      {
        userId: user.entity.id,
        roles: roles,
      },
    );
  });
}
