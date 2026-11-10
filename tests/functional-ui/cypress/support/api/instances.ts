import { SystemAPI } from './types';

export function instanceUnderTest(api: SystemAPI): Cypress.Chainable<string> {
  return cy
    .request({
      method: 'POST',
      url: `${api.baseURL}/instances/_search`,
      auth: {
        bearer: api.token,
      },
    })
    .then((res) => {
      const instances = <Array<any>>res.body.result;
      expect(instances.length).to.equal(
        1,
        'instanceUnderTest just supports running against an API with exactly one instance, yet',
      );
      return instances[0].id;
    });
}

export function getInstance(api: SystemAPI, instanceId: string, failOnStatusCode = true) {
  return cy.request({
    method: 'GET',
    url: `${api.baseURL}/instances/${instanceId}`,
    auth: {
      bearer: api.token,
    },
    failOnStatusCode: failOnStatusCode,
  });
}
