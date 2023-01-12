import { ZITADELTarget } from 'support/commands';
import { Entity, SearchResult } from './types';

export function searchSomething(
  target: ZITADELTarget,
  searchPath: string,
  method: string,
  mapResult: (body: any) => SearchResult,
  orgId?: number,
): Cypress.Chainable<SearchResult> {
  return cy
    .request({
      method: method,
      url: searchPath,
      headers: target.headers,
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
