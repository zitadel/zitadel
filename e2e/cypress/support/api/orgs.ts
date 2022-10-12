import { ensureSomething } from './ensure';
import { searchSomething } from './search';
import { API } from './types';
import { host } from '../login/users';
import { requestHeaders } from './apiauth';

export function ensureOrgExists(api: API, name: string): Cypress.Chainable<number> {
  return ensureSomething(
    api,
    () =>
      searchSomething(
        api,
        encodeURI(`${api.mgmtBaseURL}/global/orgs/_by_domain?domain=${name}.${host(Cypress.config('baseUrl'))}`),
        'GET',
        (res) => {
          return { entity: res.org, id: res.org?.id, sequence: res.org?.details?.sequence };
        },
      ),
    () => `${api.mgmtBaseURL}/orgs`,
    'POST',
    { name: name },
    (org: any) => org?.name === name,
    (res) => res.id,
  );
}

export function getOrgUnderTest(api: API): Cypress.Chainable<number> {
  return searchSomething(api, `${api.mgmtBaseURL}/orgs/me`, 'GET', (res) => {
    return { entity: res.org, id: res.org.id, sequence: res.org.details.sequence };
  }).then((res) => res.entity.id);
}

export function renameOrg(api: API, name): Cypress.Chainable {
  return cy.request({
    method: 'PUT',
    url: `${api.mgmtBaseURL}/orgs/me`,
    body: {
      name: name,
    },
    headers: requestHeaders(api),
  });
}
