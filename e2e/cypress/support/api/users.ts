import { values } from 'cypress/types/lodash';
import { ensureItemDoesntExist, ensureItemExists } from './ensure';
import { API } from './types';

export function ensureHumanUserExists(api: API, username: string): Cypress.Chainable<number> {
  return ensureItemExists(
    api,
    `${api.mgmtBaseURL}/users/_search`,
    (user: any) => user.userName === username,
    `${api.mgmtBaseURL}/users/human`,
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
    undefined,
    'userId',
  );
}

export function ensureMachineUserExists(api: API, username: string): Cypress.Chainable<number> {
  return ensureItemExists(
    api,
    `${api.mgmtBaseURL}/users/_search`,
    (user: any) => user.userName === username,
    `${api.mgmtBaseURL}/users/machine`,
    {
      user_name: username,
      name: 'e2emachinename',
      description: 'e2emachinedescription',
    },
    undefined,
    'userId',
  );
}

export function ensureUserDoesntExist(api: API, username: string): Cypress.Chainable<null> {
  return ensureItemDoesntExist(
    api,
    `${api.mgmtBaseURL}/users/_search`,
    (user: any) => user.userName === username,
    (user) => `${api.mgmtBaseURL}/users/${user.id}`,
  );
}

export function ensureUserMetadataExists(api: API, userId: string, key: string): Cypress.Chainable<number> {
  return ensureItemExists(
    api,
    `${api.mgmtBaseURL}/users/${userId}/metadata/_search`,
    (metadata: any) => metadata.key === key,
    `${api.mgmtBaseURL}/users/${userId}/metadata/${key}`,
    {
      key: key,
      value: 'testvalue',
    },
    undefined,
  );
}

export function ensureUserMetadataDoesntExist(api: API, userId: string, key: string): Cypress.Chainable<null> {
  return ensureItemDoesntExist(
    api,
    `${api.mgmtBaseURL}/users/${userId}/metadata/_search`,
    (metadata: any) => metadata.key === key,
    (user) => `${api.mgmtBaseURL}/users/${userId}/metadata/${key}`,
  );
}
