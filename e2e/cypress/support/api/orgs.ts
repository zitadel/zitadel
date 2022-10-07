import { ensureSomething } from './ensure';
import { searchSomething } from './search';
import { API } from './types';
import { host } from '../login/users';

export function ensureOrgExists(api: API, name: string): Cypress.Chainable<number> {
  return ensureSomething(
    api,
    () =>
      searchSomething(
        api,
        encodeURI(`${api.mgntBaseURL}global/orgs/_by_domain?domain=${name}.${host(Cypress.config('baseUrl'))}`),
        'GET',
        (res) => {
          return { entity: res.org, sequence: res.org.details?.sequence };
        },
      ),
    () => `${api.mgntBaseURL}orgs`,
    'POST',
    {
      name: name,
    },
    (org: any) => org.name === name,
  );
}

export function getOrgUnderTest(api: API): Cypress.Chainable<number> {
  return searchSomething(api, `${api.mgntBaseURL}orgs/me`, 'GET', (res) => {
    return { entity: res.org, sequence: res.org.details.sequence };
  }).then((res) => res.entity.id);
}
