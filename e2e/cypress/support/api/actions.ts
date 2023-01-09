import { API } from './types';

export function ensureActionDoesntExist(api: API, name: string) {
  return search(api, name).then((entity) => {
    if (entity) {
      return remove(api, entity.id);
    }
  });
}

export function ensureActionExists(api: API, name: string, script: string): Cypress.Chainable<number> {
  return search(api, name).then((entity) => {
    if (!entity) {
      return create(api, name, script);
    }

    if (entity.script != script) {
      update(api, entity.id, name, script);
    }
    return cy.wrap(<number>entity.id);
  });
}

export function setTriggerTypes(api: API, flowType: number, triggerType: number, actionIds: Array<number>) {
  return cy
    .request({
      method: 'POST',
      url: `${api.mgmtBaseURL}/flows/${flowType}/trigger/${triggerType}`,
      body: {
        actionIds: actionIds,
      },
      ...auth(api),
      failOnStatusCode: false,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.body.message).to.contain('No Changes');
      }
    });
}

function search(api: API, name: string): Cypress.Chainable<any> {
  return cy
    .request({
      method: 'POST',
      url: `${api.mgmtBaseURL}/actions/_search`,
      ...auth(api),
    })
    .then((res) => {
      return res.body?.result?.find((action) => action.name == name) || cy.wrap(null);
    });
}

function create(api: API, name: string, script: string): Cypress.Chainable<number> {
  return cy
    .request({
      method: 'POST',
      url: `${api.mgmtBaseURL}/actions`,
      body: {
        name: name,
        script: script,
        allowedToFail: false,
        timeout: '10s',
      },
      ...auth(api),
    })
    .its('body.id');
}

function update(api: API, id: string, name: string, script: string) {
  return cy.request({
    method: 'PUT',
    url: `${api.mgmtBaseURL}/actions/${id}`,
    body: {
      name: name,
      script: script,
      allowedToFail: false,
      timeout: '10s',
    },
    ...auth(api),
  });
}

function remove(api: API, id: string) {
  return cy.request({
    method: 'DELETE',
    url: `${api.mgmtBaseURL}/actions/${id}`,
    ...auth(api),
  });
}

function auth(api: API) {
  return { auth: { bearer: api.token } };
}
