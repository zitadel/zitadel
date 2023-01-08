import { API } from './types';

export function ensureActionDoesntExist(api: API, name: string) {
  return getAction(api, name).then((action) => {
    if (action) {
      return removeAction(api, name);
    }
  });
}

export function ensureActionExists(api: API, name: string, script: string): Cypress.Chainable<number> {
  return getAction(api, name).then((action) => {
    if (!action) {
      return createAction(api, name, script);
    }

    if (action.script != script) {
      updateAction(api, action.id, name, script);
    }
    return cy.wrap(<number>action.id);
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

function getAction(api: API, name: string): Cypress.Chainable<any> {
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

function createAction(api: API, name: string, script: string): Cypress.Chainable<number> {
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

function updateAction(api: API, id: string, name: string, script: string) {
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

function removeAction(api: API, id: string) {
  return cy.request({
    method: 'DELETE',
    url: `${api.mgmtBaseURL}/actions/${id}`,
    ...auth(api),
  });
}

function auth(api: API) {
  return { auth: { bearer: api.token } };
}
