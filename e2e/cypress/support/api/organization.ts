import { apiCallProperties } from './apiauth';
import { ensureKeyIsSet, getSomething } from './ensure';

export function ensureOrganizationMetadataExists(
  api: apiCallProperties,
  key: string,
  value: string,
): Cypress.Chainable<number> {
  console.log('ensure key is set');
  return ensureKeyIsSet(
    api,
    `${api.mgntBaseURL}metadata/${key}`,
    (metadata: any) => {
      let entity = null;
      if (metadata.metadata?.key === key) {
        entity = metadata.metadata;
      }
      return {
        entity: entity,
        sequence: metadata.details?.sequence,
      };
    },
    `${api.mgntBaseURL}metadata/${key}`,
    {
      key: key,
      value: value,
    },
    'POST',
  );
}

export function ensureOrganizationMetadataDoesntExist(api: apiCallProperties, key: string): Cypress.Chainable<boolean> {
  return getSomething(api, `/metadata/${key}`, (metadata: any) => metadata.key === key).then((sRes) => {
    return !!!sRes.entity;
  });
}
