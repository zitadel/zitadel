import { ensureSomething } from './ensure';
import { searchSomething } from './search';
import { API } from './types';
import { host } from '../login/users';
import { requestHeaders } from './apiauth';

export function ensureOrgExists(api: Context, name: string) {
  return ensureSomething(
    api,
    () =>
      searchSomething(
        api,
        encodeURI(`${api.mgmtBaseURL}/global/orgs/_by_domain?domain=${name}.${host(Cypress.config('baseUrl'))}`),
        'GET',
        (res) => {
          return { entity: res.org, id: res.org?.id, sequence: parseInt(<string>res.org?.details?.sequence) };
        },
      ),
    () => `${api.mgmtBaseURL}/orgs`,
    'POST',
    { name: name },
    (org) => org?.name === name,
    (res) => res.id,
  );
}

export function isDefaultOrg(api: API, orgId: number): Cypress.Chainable<boolean> {
  console.log('huhu', orgId);
  return cy
    .request({
      method: 'GET',
      url: encodeURI(`${api.mgmtBaseURL}/iam`),
      headers: requestHeaders(api, orgId),
    })
    .then((res) => {
      const { defaultOrgId } = res.body;
      expect(defaultOrgId).to.equal(orgId);
      return defaultOrgId === orgId;
    });
}

export function ensureOrgIsDefault(api: API, orgId: number): Cypress.Chainable<boolean> {
  return cy
    .request({
      method: 'GET',
      url: encodeURI(`${api.mgmtBaseURL}/iam`),
      headers: requestHeaders(api, orgId),
    })
    .then((res) => {
      return res.body;
    })
    .then(({ defaultOrgId }) => {
      if (defaultOrgId === orgId) {
        return true;
      } else {
        return cy
          .request({
            method: 'PUT',
            url: `${api.adminBaseURL}/orgs/default/${orgId}`,
            headers: requestHeaders(api, orgId),
            failOnStatusCode: true,
            followRedirect: false,
          })
          .then((cRes) => {
            expect(cRes.status).to.equal(200);
            return !!cRes.body;
          });
      }
    });
}

export function getOrgUnderTest(api: API): Cypress.Chainable<number> {
  return searchSomething(api, `${api.mgmtBaseURL}/orgs/me`, 'GET', (res) => {
    return { entity: res.org, id: res.org.id, sequence: parseInt(<string>res.org.details.sequence) };
  }).then((res) => res.entity.id);
}
