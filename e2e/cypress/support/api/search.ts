import { API, Entity, SearchResult } from './types';

export function searchSomething(
  api: API,
  searchPath: string,
  method: string,
  mapResult: (body: any) => SearchResult,
  body?: any,
): Cypress.Chainable<SearchResult> {
  return cy
    .request({
      method: method,
      url: searchPath,
      headers: {
        Authorization: api.authHeader,
      },
      body: body,
    })
    .then((res) => {
      return mapResult(res.body);
    });
}

export function findFromList(find: (entity: Entity) => boolean): (body: any) => SearchResult {
  return (b) => {
    return {
      entity: b.result?.find(find),
      sequence: parseInt(<string>b.details.processedSequence),
    };
  };
}
