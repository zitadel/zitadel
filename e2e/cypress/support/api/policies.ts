import { requestHeaders } from './apiauth';
import { ensureSetting } from './ensure';
import { API } from './types';

export enum Policy {
  Label = 'label',
}

export function resetPolicy(api: API, policy: Policy) {
  cy.request({
    method: 'DELETE',
    url: `${api.mgmtBaseURL}/policies/${policy}`,
    headers: requestHeaders(api),
  }).then((res) => {
    expect(res.status).to.equal(200);
    return null;
  });
}

export function ensureLoginPolicy(api: API, policy: any) {
  ensureSetting(
    api,
    `${api.mgmtBaseURL}/policies/login`,
    (body: any) => {
      return {
        sequence: body.policy?.details?.sequence,
        id: body.policy.id,
        entity: Cypress._.includes(body.policy, policy) ? body.policy : null,
      };
    },
    '/policies/login',
    policy,
  );
}
