import { last } from 'cypress/types/lodash';
import { ensureItemDoesntExist, ensureItemExists } from './ensure';
import { API } from './types';

export function ensureHumanUserExists(
  api: API,
  username: string,
  password: string,
  emailIsVerified = false,
): Cypress.Chainable<number> {
  return ensureItemExists(
    api,
    `${api.mgmtBaseURL}/users/_search`,
    (user: any) => user.userName === username,
    `${api.mgmtBaseURL}/users/human/_import`,
    {
      user_name: username,
      profile: {
        first_name: 'e2efirstName',
        last_name: 'e2elastName',
      },
      email: {
        email: 'e2e@email.ch',
        isEmailVerified: emailIsVerified,
      },
      phone: {
        phone: '+41 123456789',
      },
      password: password,
      password_change_required: false,
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

export function setMetadata(api: API, userId: number, key: string, value: string) {
  return cy.request({
    method: 'POST',
    url: `${api.mgmtBaseURL}/users/${userId}/metadata/${key}`,
    body: { value: value },
    auth: { bearer: api.token },
  });
}
