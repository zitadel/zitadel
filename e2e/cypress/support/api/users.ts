import { ensureItemDoesntExist, ensureItemExists } from './ensure';
import { API } from './types';

export function ensureHumanUserExists(api: API, username: string): Cypress.Chainable<number> {
  return ensureItemExists(
    api,
    `${api.mgntBaseURL}users/_search`,
    (user: any) => user.userName === username,
    `${api.mgntBaseURL}users/human`,
    {
      user_name: username,
      profile: {
        first_name: 'e2efirstName',
        last_name: 'e2elastName',
      },
      email: {
        email: 'e2e@email.ch',
      },
      phone: {
        phone: '+41 123456789',
      },
    },
  );
}

export function ensureMachineUserExists(api: API, username: string): Cypress.Chainable<number> {
  return ensureItemExists(
    api,
    `${api.mgntBaseURL}users/_search`,
    (user: any) => user.userName === username,
    `${api.mgntBaseURL}users/machine`,
    {
      user_name: username,
      name: 'e2emachinename',
      description: 'e2emachinedescription',
    },
  );
}

export function ensureUserDoesntExist(api: API, username: string): Cypress.Chainable<null> {
  return ensureItemDoesntExist(
    api,
    `${api.mgntBaseURL}users/_search`,
    (user: any) => user.userName === username,
    (user) => `${api.mgntBaseURL}users/${user.id}`,
  );
}
