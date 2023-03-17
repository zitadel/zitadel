import { ensureSomething } from './ensure';
import { searchSomething } from './search';
import { API } from './types';
import { host } from '../login/users';
import { requestHeaders } from './apiauth';
import { Context } from 'support/commands';

export function ensureOrgExists(ctx: Context, name: string) {
  return ensureSomething(
    ctx.api,
    () =>
      searchSomething(
        ctx.api,
        encodeURI(`${ctx.api.mgmtBaseURL}/global/orgs/_by_domain?domain=${name}.${host(Cypress.config('baseUrl'))}`),
        'GET',
        (res) => {
          return { entity: res.org, id: res.org?.id, sequence: parseInt(<string>res.org?.details?.sequence) };
        },
      ),
    () => `${ctx.api.mgmtBaseURL}/orgs`,
    'POST',
    { name: name },
    (org) => org?.name === name,
    (res) => res.id,
  );
}

export function isDefaultOrg(ctx: Context, orgId: string): Cypress.Chainable<boolean> {
  console.log('huhu', orgId);
  return cy
    .request({
      method: 'GET',
      url: encodeURI(`${ctx.api.mgmtBaseURL}/iam`),
      headers: requestHeaders(ctx.api, orgId),
    })
    .then((res) => {
      const { defaultOrgId } = res.body;
      expect(defaultOrgId).to.equal(orgId);
      return defaultOrgId === orgId;
    });
}

export function ensureOrgIsDefault(ctx: Context, orgId: string): Cypress.Chainable<boolean> {
  return cy
    .request({
      method: 'GET',
      url: encodeURI(`${ctx.api.mgmtBaseURL}/iam`),
      headers: requestHeaders(ctx.api, orgId),
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
            url: `${ctx.api.adminBaseURL}/orgs/default/${orgId}`,
            headers: requestHeaders(ctx.api, orgId),
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

export function getOrgUnderTest(ctx: Context): Cypress.Chainable<number> {
  return searchSomething(ctx.api, `${ctx.api.mgmtBaseURL}/orgs/me`, 'GET', (res) => {
    return { entity: res.org, id: res.org.id, sequence: parseInt(<string>res.org.details.sequence) };
  }).then((res) => res.entity.id);
}
