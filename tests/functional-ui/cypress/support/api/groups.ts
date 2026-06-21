import { requestHeaders } from './apiauth';
import { API } from './types';

const backendUrl = Cypress.env('BACKEND_URL');

function groupService(method: string): string {
  return `${backendUrl}/zitadel.group.v2.GroupService/${method}`;
}

export function defaultOrgId(api: API): Cypress.Chainable<string> {
  return cy
    .request({
      method: 'GET',
      url: `${api.mgmtBaseURL}/orgs/me`,
      headers: requestHeaders(api),
    })
    .then((res) => res.body.org.id);
}

function searchGroup(api: API, name: string): Cypress.Chainable<any | null> {
  return cy
    .request({
      method: 'POST',
      url: groupService('ListGroups'),
      headers: requestHeaders(api),
      body: {
        filters: [
          {
            nameFilter: {
              name: name,
              method: 'TEXT_FILTER_METHOD_EQUALS',
            },
          },
        ],
      },
    })
    .then((res) => {
      const groups = res.body.groups || [];
      return groups.length ? groups[0] : null;
    });
}

function awaitGroupSearch(api: API, name: string, expectFound: boolean, trials = 20): Cypress.Chainable<any | null> {
  return searchGroup(api, name).then((group) => {
    if (!!group === expectFound) {
      return cy.wrap(group);
    }
    expect(trials, `group ${name} ${expectFound ? 'visible' : 'gone'} in projection`).to.be.greaterThan(0);
    return cy.wait(500).then(() => awaitGroupSearch(api, name, expectFound, trials - 1));
  });
}

export function ensureGroupExists(api: API, name: string, orgId?: string): Cypress.Chainable<string> {
  return searchGroup(api, name).then((group) => {
    if (group) {
      return cy.wrap(group.id);
    }
    const resolveOrg = orgId ? cy.wrap(orgId) : defaultOrgId(api);
    return resolveOrg.then((resolvedOrgId) =>
      createGroup(api, { organizationId: resolvedOrgId, name }).then((res) => {
        const id = res.body.id;
        return awaitGroupSearch(api, name, true).then(() => id);
      }),
    );
  });
}

export function ensureGroupDoesntExist(api: API, name: string): Cypress.Chainable<null> {
  return searchGroup(api, name).then((group) => {
    if (!group) {
      return cy.wrap(null);
    }
    return deleteGroup(api, group.id)
      .then(() => awaitGroupSearch(api, name, false))
      .then(() => null);
  });
}

export interface CreateGroupBody {
  organizationId: string;
  name: string;
  description?: string;
  id?: string;
}

export function createGroup(
  api: API,
  body: CreateGroupBody,
  options: { failOnStatusCode?: boolean } = {},
): Cypress.Chainable<Cypress.Response<any>> {
  return cy.request({
    method: 'POST',
    url: groupService('CreateGroup'),
    headers: requestHeaders(api),
    body,
    failOnStatusCode: options.failOnStatusCode ?? true,
  });
}

export interface UpdateGroupBody {
  id: string;
  name?: string;
  description?: string;
}

export function updateGroup(
  api: API,
  body: UpdateGroupBody,
  options: { failOnStatusCode?: boolean } = {},
): Cypress.Chainable<Cypress.Response<any>> {
  return cy.request({
    method: 'POST',
    url: groupService('UpdateGroup'),
    headers: requestHeaders(api),
    body,
    failOnStatusCode: options.failOnStatusCode ?? true,
  });
}

export function deleteGroup(
  api: API,
  groupId: string,
  options: { failOnStatusCode?: boolean } = {},
): Cypress.Chainable<Cypress.Response<any>> {
  return cy.request({
    method: 'POST',
    url: groupService('DeleteGroup'),
    headers: requestHeaders(api),
    body: { id: groupId },
    failOnStatusCode: options.failOnStatusCode ?? true,
  });
}

export function getGroup(api: API, groupId: string): Cypress.Chainable<any | null> {
  return cy
    .request({
      method: 'POST',
      url: groupService('GetGroup'),
      headers: requestHeaders(api),
      body: { id: groupId },
      failOnStatusCode: false,
    })
    .then((res) => (res.status === 200 ? res.body.group : null));
}

export function addUsersToGroup(
  api: API,
  groupId: string,
  userIds: string[],
  options: { failOnStatusCode?: boolean } = {},
): Cypress.Chainable<Cypress.Response<any>> {
  return cy.request({
    method: 'POST',
    url: groupService('AddUsersToGroup'),
    headers: requestHeaders(api),
    body: { id: groupId, userIds },
    failOnStatusCode: options.failOnStatusCode ?? true,
  });
}

export function removeUsersFromGroup(
  api: API,
  groupId: string,
  userIds: string[],
  options: { failOnStatusCode?: boolean } = {},
): Cypress.Chainable<Cypress.Response<any>> {
  return cy.request({
    method: 'POST',
    url: groupService('RemoveUsersFromGroup'),
    headers: requestHeaders(api),
    body: { id: groupId, userIds },
    failOnStatusCode: options.failOnStatusCode ?? true,
  });
}

export function listGroupUsers(api: API, groupId: string): Cypress.Chainable<any[]> {
  return cy
    .request({
      method: 'POST',
      url: groupService('ListGroupUsers'),
      headers: requestHeaders(api),
      body: {
        filters: [{ groupIdFilter: { groupId } }],
      },
    })
    .then((res) => res.body.groupUsers ?? []);
}

export interface CreateGroupGrantBody {
  groupId: string;
  projectId: string;
  projectGrantId?: string;
  roleKeys: string[];
}

export function createGroupGrant(
  api: API,
  body: CreateGroupGrantBody,
  options: { failOnStatusCode?: boolean } = {},
): Cypress.Chainable<Cypress.Response<any>> {
  return cy.request({
    method: 'POST',
    url: groupService('CreateGroupGrant'),
    headers: requestHeaders(api),
    body,
    failOnStatusCode: options.failOnStatusCode ?? true,
  });
}

export function updateGroupGrant(
  api: API,
  grantId: string,
  roleKeys: string[],
  options: { failOnStatusCode?: boolean } = {},
): Cypress.Chainable<Cypress.Response<any>> {
  return cy.request({
    method: 'POST',
    url: groupService('UpdateGroupGrant'),
    headers: requestHeaders(api),
    body: { id: grantId, roleKeys },
    failOnStatusCode: options.failOnStatusCode ?? true,
  });
}

export function deleteGroupGrant(
  api: API,
  grantId: string,
  options: { failOnStatusCode?: boolean } = {},
): Cypress.Chainable<Cypress.Response<any>> {
  return cy.request({
    method: 'POST',
    url: groupService('DeleteGroupGrant'),
    headers: requestHeaders(api),
    body: { id: grantId },
    failOnStatusCode: options.failOnStatusCode ?? true,
  });
}

export function listGroupGrants(api: API, groupId: string): Cypress.Chainable<any[]> {
  return cy
    .request({
      method: 'POST',
      url: groupService('ListGroupGrants'),
      headers: requestHeaders(api),
      body: {
        filters: [{ groupIdFilter: { groupId } }],
      },
    })
    .then((res) => res.body.groupGrants ?? []);
}

export function setGroupManagerRoles(
  api: API,
  groupId: string,
  roles: string[],
  options: { failOnStatusCode?: boolean } = {},
): Cypress.Chainable<Cypress.Response<any>> {
  return cy.request({
    method: 'POST',
    url: groupService('SetGroupManagerRoles'),
    headers: requestHeaders(api),
    body: { groupId, roles },
    failOnStatusCode: options.failOnStatusCode ?? true,
  });
}

export function ensureGroupGrantDoesntExist(
  api: API,
  groupId: string,
  projectId: string,
): Cypress.Chainable<null> {
  return listGroupGrants(api, groupId).then((grants) => {
    const match = grants.find((g) => g.projectId === projectId);
    if (!match) {
      return cy.wrap(null);
    }
    return deleteGroupGrant(api, match.id).then(() => null);
  });
}

export function assertGroupAggregateScopedToOrg(
  api: API,
  groupId: string,
  expectedOrgId: string,
): Cypress.Chainable<any> {
  return getGroup(api, groupId).then((g) => {
    expect(g, `group ${groupId} resolvable`).to.not.be.null;
    expect(g.organizationId, `group ${groupId} scoped to org`).to.equal(expectedOrgId);
    return cy.wrap(g);
  });
}

export function assertNoOrphanGrants(api: API, groupId: string): Cypress.Chainable<any[]> {
  return listGroupGrants(api, groupId).then((grants) => {
    expect(grants, `no grants remain for group ${groupId}`).to.have.length(0);
    return cy.wrap(grants);
  });
}

export function assertMembershipCount(
  api: API,
  groupId: string,
  expected: number,
): Cypress.Chainable<any[]> {
  return listGroupUsers(api, groupId).then((users) => {
    expect(users.length, `group ${groupId} membership count`).to.equal(expected);
    return cy.wrap(users);
  });
}
