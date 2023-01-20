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
  return standardEnsureDoesntExist(ensureActionExists(target, name, ''), Cypress._.curry(remove)(target), () =>
    search(target, name),
  );
}

export function ensureActionExists(target: ZITADELTarget, name: string, script: string, allowedToFail = false) {
  return standardEnsureExists(
    create(target, name, script, allowedToFail),
    () => search(target, name),
    Cypress._.curry(update)(target, name, script, allowedToFail),
  );
}

function create(target: ZITADELTarget, name: string, script: string, allowedToFail: boolean) {
  return standardCreate<number>(
    target,
    `${target.mgmtBaseURL}/actions`,
    {
      name: name,
      script: script,
      allowedToFail: allowedToFail,
      timeout: '10s',
    },
    'id',
  );
}

function search(target: ZITADELTarget, name: string) {
  return standardSearch<number>(target, `${target.mgmtBaseURL}/actions/_search`, (entity) => entity.name == name, 'id');
}

function update(target: ZITADELTarget, name: string, script: string, allowedToFail: boolean, id: number) {
  return standardUpdate(target, `${target.mgmtBaseURL}/actions/${id}`, {
    name: name,
    script: script,
    allowedToFail: allowedToFail,
  });
}

function remove(target: ZITADELTarget, id: number) {
  return standardRemove(target, `${target.mgmtBaseURL}/actions/${id}`);
}

export function triggerActions(target: ZITADELTarget, flowType: number, triggerType: number, actionIds: Array<number>) {
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

export function resetAllTriggers(target: ZITADELTarget) {
  [
    { flow: 3, trigger: 1 },
    { flow: 3, trigger: 2 },
    { flow: 3, trigger: 3 },
  ].forEach((combo) => triggerActions(target, combo.flow, combo.trigger, []));
}
