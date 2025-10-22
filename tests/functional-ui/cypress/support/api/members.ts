import { ensureItemDoesntExist, ensureItemExists } from './ensure';
import { findFromList, searchSomething } from './search';
import { API } from './types';

export function ensureHumanIsNotOrgMember(api: API, username: string) {
  return ensureItemDoesntExist(
    api,
    `${api.mgmtBaseURL}/orgs/me/members/_search`,
    (member: any) => (<string>member.preferredLoginName).startsWith(username),
    (member) => `${api.mgmtBaseURL}/orgs/me/members/${member.userId}`,
  );
}

export function ensureHumanIsOrgMember(api: API, username: string, roles: string[]) {
  return searchSomething(
    api,
    `${api.mgmtBaseURL}/users/_search`,
    'POST',
    findFromList((user) => {
      return user.userName == username;
    }),
  ).then((user) => {
    return ensureItemExists(
      api,
      `${api.mgmtBaseURL}/orgs/me/members/_search`,
      (member: any) => member.userId == user.entity.id,
      `${api.mgmtBaseURL}/orgs/me/members`,
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
  grantId?: string,
): Cypress.Chainable<string> {
  return ensureItemDoesntExist(
    api,
    `${api.mgmtBaseURL}/projects/${projectId}/${grantId ? `grants/${grantId}/` : ''}members/_search`,
    (member: any) => (<string>member.preferredLoginName).startsWith(username),
    (member) => `${api.mgmtBaseURL}/projects/${projectId}/${grantId ? `grants/${grantId}/` : ''}members/${member.userId}`,
  );
}

export function ensureHumanIsProjectMember(
  api: API,
  projectId: string,
  username: string,
  roles: string[],
  grantId?: string,
): Cypress.Chainable<string> {
  return searchSomething(
    api,
    `${api.mgmtBaseURL}/users/_search`,
    'POST',
    findFromList((user) => {
      return user.userName == username;
    }),
  ).then((user) => {
    return ensureItemExists(
      api,
      `${api.mgmtBaseURL}/projects/${projectId}/${grantId ? `grants/${grantId}/` : ''}members/_search`,
      (member: any) => member.userId == user.entity.id,
      `${api.mgmtBaseURL}/projects/${projectId}/${grantId ? `grants/${grantId}/` : ''}members`,
      {
        userId: user.entity.id,
        roles: roles,
      },
    );
  });
}
