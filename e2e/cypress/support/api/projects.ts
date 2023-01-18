import { ZITADELTarget } from 'support/commands';
import { standardCreate, standardEnsureDoesntExist, standardEnsureExists, standardRemove, standardSearch } from './standard';

export function ensureProjectDoesntExist(target: ZITADELTarget, name: string) {
  return standardEnsureDoesntExist(ensureProjectExists(target, name), Cypress._.curry(remove)(target), () =>
    search(target, name),
  );
}

export function ensureProjectExists(target: ZITADELTarget, name: string) {
  return standardEnsureExists(create(target, name), () => search(target, name));
}

function create(target: ZITADELTarget, name: string) {
  return standardCreate<number>(
    target,
    `${target.mgmtBaseURL}/projects`,
    {
      name: name,
    },
    'id',
  );
}

function search(target: ZITADELTarget, name: string) {
  return standardSearch<number>(target, `${target.mgmtBaseURL}/projects/_search`, (entity) => entity.name == name, 'id');
}

function remove(target: ZITADELTarget, id: number) {
  return standardRemove(target, `${target.mgmtBaseURL}/projects/${id}`);
}
