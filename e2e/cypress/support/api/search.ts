import { API, Entity, SearchResult } from './types';

export function searchSomething(
  api: API,
  searchPath: string,
  method: string,
  mapResult: (body: any) => SearchResult,
  orgId?: number,
): Cypress.Chainable<SearchResult> {
  const req = {
    method: method,
    url: searchPath,
    headers: {
      Authorization: api.authHeader,
    },
    failOnStatusCode: method == 'POST',
  };

  if (orgId) {
    req.headers['x-zitadel-orgid'] = orgId;
  }

  return cy.request(req).then((res) => {
    return mapResult(res.body);
  });
}

export function findFromList(find: (entity: Entity) => boolean, idField: string = 'id'): (body: any) => SearchResult {
  return (b) => {
    const entity = b.result?.find(find);
    return {
      entity: entity,
      sequence: parseInt(<string>b.details.processedSequence),
      id: entity?.[idField],
    };
  };
}
