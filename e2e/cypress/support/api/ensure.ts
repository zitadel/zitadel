import { findFromList as mapFromList, searchSomething } from './search';
import { API, Entity, SearchResult } from './types';

export function ensureItemExists(
  api: API,
  searchPath: string,
  findInList: (entity: Entity) => boolean,
  createPath: string,
  body: Entity,
  orgId?: number,
  idField?: string,
): Cypress.Chainable<number> {
  return ensureSomething(
    api,
    () => searchSomething(api, searchPath, 'POST', mapFromList(findInList)),
    () => createPath,
    'POST',
    body,
    (entity) => !!entity,
    orgId,
    idField,
  );
}

export function ensureItemDoesntExist(
  api: API,
  searchPath: string,
  findInList: (entity: Entity) => boolean,
  deletePath: (entity: Entity) => string,
): Cypress.Chainable<null> {
  return ensureSomething(
    api,
    () => searchSomething(api, searchPath, 'POST', mapFromList(findInList)),
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
): Cypress.Chainable<number> {
  return ensureSomething(
    api,
    () => searchSomething(api, path, 'GET', mapResult),
    () => createPath,
    'PUT',
    body,
    (entity) => !!entity,
  );
}

function awaitDesired(
  trials: number,
  expectEntity: (entity: Entity) => boolean,
  initialSequence: number,
  search: () => Cypress.Chainable<SearchResult>,
) {
  search().then((resp) => {
    const foundExpectedEntity = expectEntity(resp.entity);
    const foundExpectedSequence = resp.sequence >= initialSequence;

    if (!foundExpectedEntity || !foundExpectedSequence) {
      expect(trials, `trying ${trials} more times`).to.be.greaterThan(0);
      cy.wait(1000);
      awaitDesired(trials - 1, expectEntity, initialSequence, search);
    }
  });
}

export function ensureSomething(
  api: API,
  search: () => Cypress.Chainable<SearchResult>,
  apiPath: (entity: Entity) => string,
  ensureMethod: string,
  body: Entity,
  expectEntity: (entity: any) => boolean,
  orgId?: number,
  idFieldname: string = 'id',
): Cypress.Chainable<number> {
  return search()
    .then((sRes) => {
      if (expectEntity(sRes.entity)) {
        return cy.wrap({
          id: ensureMethod == 'DELETE' ? NaN : sRes.entity[idFieldname],
          expectSequenceFrom: sRes.sequence,
        });
      }

      const req = {
        method: ensureMethod,
        url: apiPath(sRes.entity),
        headers: {
          Authorization: api.authHeader,
        },
        body: body,
        failOnStatusCode: false,
        followRedirect: false,
      };

      if (orgId) {
        req.headers['x-zitadel-orgid'] = orgId;
      }

      return cy.request(req).then((cRes) => {
        expect(cRes.status).to.equal(200);
        return {
          id: cRes.body[idFieldname],
          expectSequenceFrom: sRes.sequence,
        };
      });
    })
    .then((data) => {
      awaitDesired(90, expectEntity, data.expectSequenceFrom, search);
      return cy.wrap<number>(data.id);
    });
}
