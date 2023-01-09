import { API } from './types';

export function ensureUserDoesntExist(api: API, username: string) {
  return search(api, username).then((entity) => {
    if (entity) {
      return remove(api, entity.id);
    }
  });
}

export function ensureHumanUserExists(api: API, username: string, emailIsVerified = false): Cypress.Chainable<number> {
  return search(api, username).then((entity) => {
    if (!entity) {
      return createHuman(api, username, emailIsVerified);
    }

    return cy.wrap(<number>entity.id);
  });
}

export function ensureMachineUserExists(api: API, username: string): Cypress.Chainable<number> {
  return search(api, username).then((entity) => {
    if (!entity) {
      return createMachine(api, username);
    }

    return cy.wrap(<number>entity.id);
  });
}
/*
export function legacyEnsureHumanUserExists(
  api: API,
  username: string,
  password?: string,
  otpCode?: string,
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
      otp_code: otpCode,
    },
    undefined,
    'userId',
  );
}

export function legacyEnsureMachineUserExists(api: API, username: string): Cypress.Chainable<number> {
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

export function legacyensureUserDoesntExist(api: API, username: string): Cypress.Chainable<null> {
  return ensureItemDoesntExist(
    api,
    `${api.mgmtBaseURL}/users/_search`,
    (user: any) => user.userName === username,
    (user) => `${api.mgmtBaseURL}/users/${user.id}`,
  );
}
*/

function search(api: API, username: string): Cypress.Chainable<any> {
  return cy
    .request({
      method: 'POST',
      url: `${api.mgmtBaseURL}/users/_search`,
      ...auth(api),
    })
    .then((res) => {
      return res.body?.result?.find((entity) => entity.userName == username) || cy.wrap(null);
    });
}

function createHuman(api: API, username: string, emailIsVerified = false): Cypress.Chainable<number> {
  return cy
    .request({
      method: 'POST',
      url: `${api.mgmtBaseURL}/users/human/_import`,
      body: {
        userName: username,
        profile: {
          firstName: 'e2efirstName',
          lastName: 'e2elastName',
        },
        email: {
          email: 'e2e@email.ch',
          isEmailVerified: emailIsVerified,
        },
        phone: {
          phone: '+41 123456789',
        },
        password: 'Password1!',
        passwordVhangeRequired: false,
      },
      ...auth(api),
    })
    .its('body.userId');
}

function createMachine(api: API, username: string): Cypress.Chainable<number> {
  return cy
    .request({
      method: 'POST',
      url: `${api.mgmtBaseURL}/users/machine`,
      body: {
        userName: username,
        name: 'e2emachinename',
        description: 'e2emachinedescription',
      },
      ...auth(api),
    })
    .its('body.userId');
}

function remove(api: API, id: string) {
  return cy.request({
    method: 'DELETE',
    url: `${api.mgmtBaseURL}/users/${id}`,
    ...auth(api),
  });
}

function auth(api: API) {
  return { auth: { bearer: api.token } };
}
