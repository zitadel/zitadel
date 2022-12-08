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
    (body) => body?.settings?.id,
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

    const check = (!foundExpectedEntity || !foundExpectedSequence) && trials > 0;
    if (check) {
      cy.log(`trying ${trials} more times`);
      cy.wait(1000);
      return awaitDesired(trials - 1, expectEntity, search, initialSequence);
    } else {
      return;
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

// export function ensureKeyIsSet(
//   api: API,
//   path: string,
//   find: (entity: any) => SearchResult,
//   createPath: string,
//   body: any,
//   createMethod: string = 'PUT',
// ): Cypress.Chainable<number> {
//   return getSomething(api, path, find)
//     .then((sRes) => {
//       if (sRes.entity) {
//         return cy.wrap({
//           key: sRes.entity.key,
//           initialSequence: 0,
//         });
//       }
//       return cy
//         .request({
//           method: createMethod,
//           url: createPath,
//           headers: requestHeaders(api),
//           body: body,
//           failOnStatusCode: false,
//           followRedirect: false,
//         })
//         .then((cRes) => {
//           expect(cRes.status).to.equal(200);
//           return {
//             key: cRes.body.key,
//             initialSequence: sRes.sequence,
//           };
//         });
//     })
//     .then((data) => {
//       awaitDesiredByKey(90, (entity) => !!entity, data.initialSequence, api, path, find);
//       return cy.wrap<number>(data.key);
//     });
// }

// export function ensureSomethingDoesntExist(
//   api: API,
//   searchPath: string,
//   find: (entity: any) => boolean,
//   deletePath: (entity: any) => string,
// ): Cypress.Chainable<null> {
//   return searchSomething(api, searchPath, find)
//     .then((sRes) => {
//       if (!sRes.entity) {
//         return cy.wrap(0);
//       }
//       return cy
//         .request({
//           method: 'DELETE',
//           url: `${api.mgmtBaseURL}${deletePath(sRes.entity)}`,
//           headers: requestHeaders(api),
//           failOnStatusCode: false,
//         })
//         .then((dRes) => {
//           expect(dRes.status).to.equal(200);
//           return sRes.sequence;
//         });
//     })
//     .then((initialSequence) => {
//       awaitDesired(90, (entity) => !entity, initialSequence, api, searchPath, find);
//       return null;
//     });
// }

// type SearchResult = {
//   entity: any;
//   sequence: number;
// };

// function searchSomething(api: API, searchPath: string, find: (entity: any) => boolean): Cypress.Chainable<SearchResult> {
//   return cy
//     .request({
//       method: 'POST',
//       url: `${api.mgmtBaseURL}${searchPath}`,
//       headers: requestHeaders(api),
//     })
//     .then((res) => {
//       return {
//         entity: res.body.result?.find(find) || null,
//         sequence: res.body.details.processedSequence,
//       };
//     });
// }

// export function getSomething(api: API, searchPath: string, find: (entity: any) => any): Cypress.Chainable<any> {
//   return cy
//     .request({
//       method: 'GET',
//       url: searchPath,
//       headers: requestHeaders(api),
//     })
//     .then((res) => {
//       return find(res.body);
//     });
// }

// function awaitDesired(
//   trials: number,
//   expectEntity: (entity: any) => boolean,
//   initialSequence: number,
//   api: API,
//   searchPath: string,
//   find: (entity: any) => boolean,
// ) {
//   searchSomething(api, searchPath, find).then((resp) => {
//     const foundExpectedEntity = expectEntity(resp.entity);
//     const foundExpectedSequence = resp.sequence > initialSequence;

//     if (!foundExpectedEntity || !foundExpectedSequence) {
//       expect(trials, `trying ${trials} more times`).to.be.greaterThan(0);
//       cy.wait(1000);
//       awaitDesired(trials - 1, expectEntity, initialSequence, api, searchPath, find);
//     }
//   });
// }

// function awaitDesiredById(
//   trials: number,
//   expectEntity: (entity: any) => boolean,
//   initialSequence: number,
//   api: API,
//   path: string,
//   find: (entity: any) => SearchResult,
// ) {
//   getSomething(api, path, find).then((resp) => {
//     const foundExpectedEntity = expectEntity(resp.entity);
//     const foundExpectedSequence = resp.sequence > initialSequence;

//     if (!foundExpectedEntity || !foundExpectedSequence) {
//       expect(trials, `trying ${trials} more times`).to.be.greaterThan(0);
//       cy.wait(1000);
//       awaitDesiredById(trials - 1, expectEntity, initialSequence, api, path, find);
//     }
//   });
// }

// function awaitDesiredByKey(
//   trials: number,
//   expectEntity: (entity: any) => boolean,
//   initialSequence: number,
//   api: API,
//   path: string,
//   find: (entity: any) => SearchResult,
// ) {
//   getSomething(api, path, find).then((resp) => {
//     const foundExpectedKey = expectEntity(resp);

//     if (!foundExpectedKey) {
//       expect(trials, `trying ${trials} more times`).to.be.greaterThan(0);
//       cy.wait(1000);
//       awaitDesiredById(trials - 1, expectEntity, initialSequence, api, path, find);
//     }
//   });
// }
