import { requestHeaders } from './apiauth';
import { findFromList as mapFromList, searchSomething } from './search';
import { API, Entity, SearchResult, Token } from './types';

export function ensureItemExists(
  token: Token,
  searchPath: string,
  findInList: (entity: Entity) => boolean,
  createPath: string,
  body: Entity,
  orgId?: string,
  newItemIdField: string = 'id',
  searchItemIdField?: string,
) {
  return ensureSomething(
    token,
    () => searchSomething(token, searchPath, 'POST', mapFromList(findInList, searchItemIdField), orgId),
    () => createPath,
    'POST',
    body,
    (entity) => !!entity,
    (body) => body[newItemIdField],
    orgId,
  );
}

export function ensureItemDoesntExist(
  token: Token,
  searchPath: string,
  findInList: (entity: Entity) => boolean,
  deletePath: (entity: Entity) => string,
  orgId?: string,
): Cypress.Chainable<null> {
  return ensureSomething(
    token,
    () => searchSomething(token, searchPath, 'POST', mapFromList(findInList), orgId),
    deletePath,
    'DELETE',
    null,
    (entity) => !entity,
  ).then(() => null);
}

export function ensureSetting(
  api: API,
  path: string,
  mapResult: (entity: any) => SearchResult,
  createPath: string,
  body: any,
  orgId?: string,
): Cypress.Chainable<string> {
  return ensureSomething(
    api,
    () => searchSomething(api, path, 'GET', mapResult, orgId),
    () => createPath,
    'PUT',
    body,
    (entity) => !!entity,
  );
}

function awaitDesired(trials: number, eventTimestamp: number, search: () => Cypress.Chainable<SearchResult>) {
  return search().then((resp) => {
    if (resp.viewTimeStamp < eventTimestamp) {
      expect(trials, `trying ${trials} more times`).to.be.greaterThan(0);
      cy.wait(1000);
      return awaitDesired(trials - 1, eventTimestamp, search);
    }
  });
}

export function ensureSomething(
  token: Token,
  search: () => Cypress.Chainable<SearchResult>,
  apiPath: (entity: Entity) => string,
  ensureMethod: string,
  body: Entity,
  expectEntity: (entity: Entity) => boolean,
  mapId?: (body: any) => string,
  orgId?: string,
): Cypress.Chainable<string> {
  return search().then((sRes) => {
    if (expectEntity(sRes.entity)) {
      return cy.wrap(sRes.id);
    }
    return cy
      .request({
        method: ensureMethod,
        url: apiPath(sRes.entity),
        headers: requestHeaders(token, orgId),
        body: body,
        followRedirect: false,
      })
      .then((cRes) => {
        expect(cRes.status).to.equal(200);
        const id = mapId ? mapId(cRes.body) : undefined;
        return awaitDesired(90, cRes.body.details.changeDate || cRes.body.details.creationDate, search).then(() => {
          return cy.wrap(id);
        });
      });
  });
}
