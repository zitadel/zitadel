import { requestHeaders } from './apiauth';
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
