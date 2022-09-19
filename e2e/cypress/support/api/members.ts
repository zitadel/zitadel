import { apiCallProperties } from './apiauth';
import { ensureSomethingDoesntExist, ensureSomethingExists, searchSomething } from './ensure';

export function ensureHumanIsNotOrgMember(api: apiCallProperties, username: string): Cypress.Chainable<number> {
  return ensureSomethingDoesntExist(
    api,
    'orgs/me/members/_search',
    (member: any) => (<string>member.preferredLoginName).startsWith(username),
    (member) => `orgs/me/members/${member.userId}`,
  );
}

export function ensureHumanIsOrgMember(
  api: apiCallProperties,
  username: string,
  roles: string[],
): Cypress.Chainable<number> {
  return searchSomething(api, 'users/_search', (user) => {
    return user.userName == username;
  }).then((user) => {
    return ensureSomethingExists(
      api,
      'orgs/me/members/_search',
      (member: any) => member.userId == user.entity.id,
      'orgs/me/members',
      {
        userId: user.entity.id,
        roles: roles,
      },
    );
  });
}

export function ensureHumanIsNotProjectMember(
  api: apiCallProperties,
  projectId: string,
  username: string,
): Cypress.Chainable<number> {
  return ensureSomethingDoesntExist(
    api,
    `projects/${projectId}/members/_search`,
    (member: any) => (<string>member.preferredLoginName).startsWith(username),
    (member) => `projects/${projectId}/members/${member.userId}`,
  );
}

export function ensureHumanIsProjectMember(
  api: apiCallProperties,
  projectId: string,
  username: string,
  roles: string[],
): Cypress.Chainable<number> {
  return searchSomething(api, 'users/_search', (user) => {
    return user.userName == username;
  }).then((user) => {
    return ensureSomethingExists(
      api,
      `projects/${projectId}/members/_search`,
      (member: any) => member.userId == user.entity.id,
      `projects/${projectId}/members`,
      {
        userId: user.entity.id,
        roles: roles,
      },
    );
  });
}
