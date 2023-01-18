import { ZITADELTarget } from 'support/commands';
import { standardCreate, standardEnsureExists, standardRemove, standardSearch } from './standard';
import { newOrgTarget } from './target';

export function ensureOrgExists(target: ZITADELTarget, name: string): Cypress.Chainable<ZITADELTarget> {
  return standardEnsureExists(create(target, name), () => search(target, name)).then((id) => {
    return newOrgTarget(target, id, name);
  });
}

function search(target: ZITADELTarget, name: string) {
  return standardSearch<number>(target, `${target.adminBaseURL}/orgs/_search`, (entity) => entity.name == name, 'id');
}

function create(target: ZITADELTarget, name: string) {
  return standardCreate<number>(target, `${target.mgmtBaseURL}/orgs`, { name: name }, 'id');
}

export function remove(target: ZITADELTarget) {
  return standardRemove(target, `${target.mgmtBaseURL}/orgs/me`);
}
