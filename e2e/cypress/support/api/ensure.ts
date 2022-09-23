import { apiCallProperties } from './apiauth';

export function ensureSomethingExists(
  api: apiCallProperties,
  searchPath: string,
  find: (entity: any) => boolean,
  createPath: string,
  body: any,
): Cypress.Chainable<number> {
  return searchSomething(api, searchPath, find)
    .then((sRes) => {
      if (sRes.entity) {
        return cy.wrap({
          id: sRes.entity.id,
          initialSequence: 0,
        });
      }
      return cy
        .request({
          method: 'POST',
          url: `${api.mgntBaseURL}${createPath}`,
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
      awaitDesired(90, (entity) => !!entity, data.initialSequence, api, searchPath, find);
      return cy.wrap<number>(data.id);
    });
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
  find: (entity: any) => boolean,
  deletePath: (entity: any) => string,
): Cypress.Chainable<null> {
  return searchSomething(api, searchPath, find)
    .then((sRes) => {
      if (!sRes.entity) {
        return cy.wrap(0);
      }
      return cy
        .request({
          method: 'DELETE',
          url: `${api.mgntBaseURL}${deletePath(sRes.entity)}`,
          headers: {
            Authorization: api.authHeader,
          },
          failOnStatusCode: false,
        })
        .then((dRes) => {
          expect(dRes.status).to.equal(200);
          return sRes.sequence;
        });
    })
    .then((initialSequence) => {
      awaitDesired(90, (entity) => !entity, initialSequence, api, searchPath, find);
      return null;
    });
}

type SearchResult = {
  entity: any;
  sequence: number;
};

function searchSomething(
  api: apiCallProperties,
  searchPath: string,
  find: (entity: any) => boolean,
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
        sequence: res.body.details.processedSequence,
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
  expectEntity: (entity: any) => boolean,
  initialSequence: number,
  api: apiCallProperties,
  searchPath: string,
  find: (entity: any) => boolean,
) {
  searchSomething(api, searchPath, find).then((resp) => {
    const foundExpectedEntity = expectEntity(resp.entity);
    const foundExpectedSequence = resp.sequence > initialSequence;

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
