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
      return instances.find(instance => api.baseURL.indexOf(instance.domain) > -1).id
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

export function createInstance(api: SystemAPI, name: string, domain: string, failOnStatusCode = true) {
  return cy.request({
    method: 'POST',
    url: `${api.baseURL}/instances/_create`,
    auth: {
      bearer: api.token,
    },
    body: {
      instanceName: name,
      custom_domain: domain,
      human: {
        userName: "zitadel-admin@zitadel.localhost",
        email: {
          email: "zitadel-admin@zitadel.localhost",
          isEmailVerified: true
        },
        password: {
          password: "Password1!",
          passwordChangeRequrired: false
        },
        profile: {
          firstName: "ZITADEL",
          lastName: "Admin"
        }
      }
    },
    failOnStatusCode: failOnStatusCode,
  });
}

export function deleteInstance(api: SystemAPI, instanceId: string, failOnStatusCode = true) {
  return cy.request({
    method: 'DELETE',
    url: `${api.baseURL}/instances/${instanceId}`,
    auth: {
      bearer: api.token,
    },
    failOnStatusCode: failOnStatusCode,
  });
}
