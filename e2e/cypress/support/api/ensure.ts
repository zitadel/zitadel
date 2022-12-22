import { requestHeaders } from './apiauth';
import { findFromList as mapFromList, searchSomething } from './search';
import { API, Entity, SearchResult } from './types';

export function ensureItemExists(
  api: API,
  searchPath: string,
  findInList: (entity: Entity) => boolean,
  createPath: string,
  body: Entity,
  orgId?: number,
  newItemIdField: string = 'id',
  searchItemIdField?: string,
): Cypress.Chainable<number> {
  return ensureSomething(
    api,
    () => searchSomething(api, searchPath, 'POST', mapFromList(findInList, searchItemIdField), orgId),
    () => createPath,
    'POST',
    body,
    (entity) => !!entity,
    (body) => body[newItemIdField],
    orgId,
  );
}

export function ensureItemDoesntExist(
  api: API,
  searchPath: string,
  findInList: (entity: Entity) => boolean,
  deletePath: (entity: Entity) => string,
  orgId?: number,
): Cypress.Chainable<null> {
  return ensureSomething(
    api,
    () => searchSomething(api, searchPath, 'POST', mapFromList(findInList), orgId),
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
  orgId?: number,
): Cypress.Chainable<number> {
  return ensureSomething(
    api,
    () => searchSomething(api, path, 'GET', mapResult, orgId),
    () => createPath,
    'PUT',
    body,
    (entity) => !!entity,
  );
}

function awaitDesired(
  trials: number,
  expectEntity: (entity: Entity) => boolean,
  search: () => Cypress.Chainable<SearchResult>,
  initialSequence?: number,
) {
  return search().then((resp) => {
    const foundExpectedEntity = expectEntity(resp.entity);
    const foundExpectedSequence = !initialSequence || resp.sequence >= initialSequence;

    const check = !foundExpectedEntity || !foundExpectedSequence;
    if (check) {
      expect(trials, `trying ${trials} more times`).to.be.greaterThan(0);
      cy.wait(1000);
      return awaitDesired(trials - 1, expectEntity, search, initialSequence);
    }
  });
}

interface EnsuredResult {
  id: number;
  sequence: number;
}

export function ensureSomething(
  api: API,
  search: () => Cypress.Chainable<SearchResult>,
  apiPath: (entity: Entity) => string,
  ensureMethod: string,
  body: Entity,
  expectEntity: (entity: Entity) => boolean,
  mapId?: (body: any) => number,
  orgId?: number,
): Cypress.Chainable<number> {
  return search()
    .then<EnsuredResult>((sRes) => {
      if (expectEntity(sRes.entity)) {
        return cy.wrap({ id: sRes.id, sequence: sRes.sequence });
      }

      return cy
        .request({
          method: ensureMethod,
          url: apiPath(sRes.entity),
          headers: requestHeaders(api, orgId),
          body: body,
          failOnStatusCode: false,
          followRedirect: false,
        })
        .then((cRes) => {
          expect(cRes.status).to.equal(200);
          return {
            id: mapId ? mapId(cRes.body) : undefined,
            sequence: sRes.sequence,
          };
        });
    })
    .then((data) => {
      return awaitDesired(90, expectEntity, search, data.sequence).then(() => {
        return cy.wrap<number>(data.id);
      });
    });
}
