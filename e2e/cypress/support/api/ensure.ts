import { toNumber } from 'cypress/types/lodash';
import { apiCallProperties } from './apiauth';

// Entity is an object but not a function
type Entity = { [k: string]: any } & ({ bind?: never } | { call?: never });

export function ensureSomething(
  api: apiCallProperties,
  searchPath: string,
  find: (entity: Entity) => boolean,
  apiPath: (entity: Entity) => string,
  method: string,
  body: Entity,
  expectEntity: (entity: Entity) => boolean,
): Cypress.Chainable<number> {
  return searchSomething(api, searchPath, find)
    .then((sRes) => {
      if (expectEntity(sRes.entity)) {
        return cy.wrap({
          id: sRes?.entity?.id,
          expectSequenceFrom: sRes.sequence,
        });
      }

      const req = {
        method: method,
        url: `${api.mgntBaseURL}${apiPath(sRes.entity)}`,
        headers: {
          Authorization: api.authHeader,
        },
        body: body,
        failOnStatusCode: false,
        followRedirect: false,
      };

      return cy.request(req).then((cRes) => {
        expect(cRes.status).to.equal(200);
        return {
          id: cRes.body.id,
          expectSequenceFrom: sRes.sequence,
        };
      });
    })
    .then((data) => {
      awaitDesired(90, expectEntity, data.expectSequenceFrom, api, searchPath, find);
      return cy.wrap<number>(data.id);
    });
}

export function ensureSomethingExists(
  api: apiCallProperties,
  searchPath: string,
  find: (entity: Entity) => boolean,
  createPath: string,
  body: Entity,
): Cypress.Chainable<number> {
  return ensureSomething(
    api,
    searchPath,
    find,
    () => createPath,
    'POST',
    body,
    (entity) => !!entity,
  );
}

export function ensureSomethingIsSet(
  api: apiCallProperties,
  path: string,
  find: (entity: any) => SearchResult,
  createPath: string,
  body: any,
): Cypress.Chainable<number> {
  return getSomething(api, path, find)
    .then((sRes) => {
      if (sRes.entity) {
        return cy.wrap({
          id: sRes.entity.id,
          initialSequence: 0,
        });
      }
      return cy
        .request({
          method: 'PUT',
          url: createPath,
          headers: {
            Authorization: api.authHeader,
          },
          body: body,
          failOnStatusCode: false,
          followRedirect: false,
        })
        .then((cRes) => {
          expect(cRes.status).to.equal(200);
          return {
            id: cRes.body.id,
            initialSequence: sRes.sequence,
          };
        });
    })
    .then((data) => {
      awaitDesiredById(90, (entity) => !!entity, data.initialSequence, api, path, find);
      return cy.wrap<number>(data.id);
    });
}

export function ensureSomethingDoesntExist(
  api: apiCallProperties,
  searchPath: string,
  find: (entity: Entity) => boolean,
  deletePath: (entity: Entity) => string,
): Cypress.Chainable<null> {
  return ensureSomething(api, searchPath, find, deletePath, 'DELETE', null, (entity) => !entity).then(() => {
    return null;
  });
}

type SearchResult = {
  entity: Entity;
  sequence: number;
};

export function searchSomething(
  api: apiCallProperties,
  searchPath: string,
  find: (entity: Entity) => boolean,
): Cypress.Chainable<SearchResult> {
  return cy
    .request({
      method: 'POST',
      url: `${api.mgntBaseURL}${searchPath}`,
      headers: {
        Authorization: api.authHeader,
      },
    })
    .then((res) => {
      return {
        entity: res.body.result?.find(find) || null,
        sequence: parseInt(<string>res.body.details.processedSequence),
      };
    });
}

function getSomething(
  api: apiCallProperties,
  searchPath: string,
  find: (entity: any) => SearchResult,
): Cypress.Chainable<SearchResult> {
  return cy
    .request({
      method: 'GET',
      url: searchPath,
      headers: {
        Authorization: api.authHeader,
      },
    })
    .then((res) => {
      return find(res.body);
    });
}

function awaitDesired(
  trials: number,
  expectEntity: (entity: Entity) => boolean,
  initialSequence: number,
  api: apiCallProperties,
  searchPath: string,
  find: (entity: Entity) => boolean,
) {
  searchSomething(api, searchPath, find).then((resp) => {
    const foundExpectedEntity = expectEntity(resp.entity);
    const foundExpectedSequence = resp.sequence >= initialSequence;

    if (!foundExpectedEntity || !foundExpectedSequence) {
      expect(trials, `trying ${trials} more times`).to.be.greaterThan(0);
      cy.wait(1000);
      awaitDesired(trials - 1, expectEntity, initialSequence, api, searchPath, find);
    }
  });
}

function awaitDesiredById(
  trials: number,
  expectEntity: (entity: any) => boolean,
  initialSequence: number,
  api: apiCallProperties,
  path: string,
  find: (entity: any) => SearchResult,
) {
  getSomething(api, path, find).then((resp) => {
    const foundExpectedEntity = expectEntity(resp.entity);
    const foundExpectedSequence = resp.sequence > initialSequence;

    if (!foundExpectedEntity || !foundExpectedSequence) {
      expect(trials, `trying ${trials} more times`).to.be.greaterThan(0);
      cy.wait(1000);
      awaitDesiredById(trials - 1, expectEntity, initialSequence, api, path, find);
    }
  });
}
