import { ZITADELTarget } from 'support/commands';
import {
  standardCreate,
  standardEnsureDoesntExist,
  standardEnsureExists,
  standardRemove,
  standardSearch,
  standardUpdate,
} from './standard';

export function ensureActionDoesntExist(target: ZITADELTarget, name: string) {
  return standardEnsureDoesntExist(ensureActionExists(target, name, ''), Cypress._.curry(remove)(target));
}

export function ensureActionExists(target: ZITADELTarget, name: string, script: string): Cypress.Chainable<number> {
  return standardEnsureExists(
    create(target, name, script),
    () => search(target, name),
    Cypress._.curry(update)(target, name, script),
  );
}

function create(target: ZITADELTarget, name: string, script: string): Cypress.Chainable<any> {
  return standardCreate(
    target,
    `${target.mgmtBaseURL}/actions`,
    {
      name: name,
      script: script,
      allowedToFail: false,
      timeout: '10s',
    },
    'id',
  );
}

function search(target: ZITADELTarget, name: string): Cypress.Chainable<number> {
  return standardSearch(target, `${target.mgmtBaseURL}/actions/_search`, (entity) => entity.name == name, 'id');
}

function update(target: ZITADELTarget, name: string, script: string, id: number) {
  return standardUpdate(target, `${target.mgmtBaseURL}/actions/${id}`, { name: name, script: script });
}

function remove(target: ZITADELTarget, id: number) {
  return standardRemove(target, `${target.mgmtBaseURL}/actions/${id}`);
}

export function setTriggerTypes(target: ZITADELTarget, flowType: number, triggerType: number, actionIds: Array<number>) {
  return cy
    .request({
      method: 'POST',
      url: `${target.mgmtBaseURL}/flows/${flowType}/trigger/${triggerType}`,
      body: {
        actionIds: actionIds,
      },
      failOnStatusCode: false,
      headers: target.headers,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.body.message).to.contain('No Changes');
      }
    });
}
