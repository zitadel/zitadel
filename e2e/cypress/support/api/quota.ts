import { SystemAPI } from './types';

export function addQuota(
  api: SystemAPI,
  instanceId: string,
  unit: number,
  failOnStatusCode = true,
): Cypress.Chainable<Cypress.Response<any>> {
  return cy.request({
    method: 'POST',
    url: `${api.baseURL}/instances/${instanceId}/quotas`,
    auth: {
      bearer: api.token,
    },
    body: {
      unit: unit,
    },
    failOnStatusCode: failOnStatusCode,
  });
}

export function removeQuota(
  api: SystemAPI,
  instanceId: string,
  unit: number,
  failOnStatusCode = true,
): Cypress.Chainable<Cypress.Response<any>> {
  return cy.request({
    method: 'DELETE',
    url: `${api.baseURL}/instances/${instanceId}/quotas/${unit}`,
    auth: {
      bearer: api.token,
    },
    failOnStatusCode: failOnStatusCode,
  });
}
