import { requestHeaders } from './apiauth';
import { API } from './types';

const backendUrl = Cypress.env('BACKEND_URL');

function groupService(method: string): string {
  return `${backendUrl}/zitadel.group.v2.GroupService/${method}`;
}

function defaultOrgId(api: API): Cypress.Chainable<string> {
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

export function ensureGroupExists(api: API, name: string): Cypress.Chainable<string> {
  return searchGroup(api, name).then((group) => {
    if (group) {
      return cy.wrap(group.id);
    }
    return defaultOrgId(api).then((orgId) =>
      cy
        .request({
          method: 'POST',
          url: groupService('CreateGroup'),
          headers: requestHeaders(api),
          body: {
            organizationId: orgId,
            name: name,
          },
        })
        .then((res) => {
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
    return cy
      .request({
        method: 'POST',
        url: groupService('DeleteGroup'),
        headers: requestHeaders(api),
        body: { id: group.id },
      })
      .then(() => awaitGroupSearch(api, name, false))
      .then(() => null);
  });
}
