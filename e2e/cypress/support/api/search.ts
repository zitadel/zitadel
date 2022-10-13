import { requestHeaders } from './apiauth';
import { API, Entity, SearchResult } from './types';

export function searchSomething(
  api: API,
  searchPath: string,
  method: string,
  mapResult: (body: any) => SearchResult,
  orgId?: number,
): Cypress.Chainable<SearchResult> {
  return cy
    .request({
      method: method,
      url: searchPath,
      headers: requestHeaders(api, orgId),
      failOnStatusCode: method == 'POST',
    })
    .then((res) => {
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
