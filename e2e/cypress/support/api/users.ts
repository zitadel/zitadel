import { ZITADELTarget } from 'support/commands';
import { standardCreate, standardEnsureDoesntExist, standardEnsureExists, standardRemove, standardSearch } from './standard';

export function ensureHumanExists(target: ZITADELTarget, username: string) {
  return standardEnsureExists(createHuman(target, username), () => search(target, username));
}

export function ensureMachineExists(target: ZITADELTarget, username: string) {
  return standardEnsureExists(createMachine(target, username), () => search(target, username));
}

export function ensureHumanDoesntExist(target: ZITADELTarget, username: string) {
  return standardEnsureDoesntExist(ensureHumanExists(target, username), Cypress._.curry(remove)(target), () =>
    search(target, username),
  );
}

export function ensureMachineDoesntExist(target: ZITADELTarget, username: string) {
  return standardEnsureDoesntExist(ensureMachineExists(target, username), Cypress._.curry(remove)(target), () =>
    search(target, username),
  );
}

function search(target: ZITADELTarget, username: string) {
  return standardSearch<number>(
    target,
    `${target.mgmtBaseURL}/users/_search`,
    (entity) => entity.userName == username,
    'id',
  );
}

function createHuman(target: ZITADELTarget, username: string) {
  return standardCreate<number>(
    target,
    `${target.mgmtBaseURL}/users/human/_import`,
    {
      userName: username,
      profile: {
        firstName: 'e2efirstName',
        lastName: 'e2elastName',
      },
      email: {
        email: 'e2e@email.ch',
        isEmailVerified: true,
      },
      phone: {
        phone: '+41 123456789',
      },
      password: 'Password1!',
      passwordChangeRequired: false,
    },
    'userId',
  );
}

function createMachine(target: ZITADELTarget, username: string) {
  return standardCreate<number>(
    target,
    `${target.mgmtBaseURL}/users/machine`,
    {
      userName: username,
      name: 'e2emachinename',
      description: 'e2emachinedescription',
    },
    'userId',
  );
}

function remove(target: ZITADELTarget, id: number) {
  return standardRemove(target, `${target.mgmtBaseURL}/users/${id}`);
}
