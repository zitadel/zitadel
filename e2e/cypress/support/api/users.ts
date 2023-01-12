import { ZITADELTarget } from 'support/commands';
import { loginname } from 'support/login/login';

export function ensureHumanExists(target: ZITADELTarget, username: string): Cypress.Chainable<number> {
  return ensureUserExists(target, username, () => createHuman(target, username));
}

export function ensureMachineExists(target: ZITADELTarget, username: string): Cypress.Chainable<number> {
  return ensureUserExists(target, username, () => createMachine(target, username));
}

function ensureUserExists(
  target: ZITADELTarget,
  username: string,
  create: () => Cypress.Chainable<number>,
): Cypress.Chainable<number> {
  return create().then((id) => {
    if (id) {
      return cy.wrap(id);
    }
    return search(target, username).then((id) => {
      if (id) {
        return cy.wrap(id);
      }
      sleep(6_000);
      cy.log('retrying');
      return search(target, username).then((id) => {
        if (id) {
          return cy.wrap(id);
        }
        sleep(6_000);
        cy.log('retrying');
        debugger;
        return search(target, username);
      });
    });
  });
}

export function ensureHumanDoesntExist(api: ZITADELTarget, username: string) {
  return ensureUserDoesntExist(api, username, ensureHumanExists);
}

export function ensureMachineDoesntExist(api: ZITADELTarget, username: string) {
  return ensureUserDoesntExist(api, username, ensureMachineExists);
}

function ensureUserDoesntExist(
  target: ZITADELTarget,
  username: string,
  ensureExists: (target: ZITADELTarget, username: string) => Cypress.Chainable<number>,
) {
  function ensure(loginname: string) {
    ensureExists(target, loginname).then((userId) => {
      remove(target, userId);
    });
  }

  ensure(username);
  ensure(loginname(username, target.org));
}

function search(target: ZITADELTarget, username: string): Cypress.Chainable<number> {
  return cy
    .request({
      method: 'POST',
      url: `${target.mgmtBaseURL}/users/_search`,
      headers: target.headers,
    })
    .then((res) => {
      return res.body?.result?.find((entity) => entity.userName == username)?.id || cy.wrap(null);
    });
}

function createHuman(target: ZITADELTarget, username: string): Cypress.Chainable<number> {
  return cy
    .request({
      method: 'POST',
      url: `${target.mgmtBaseURL}/users/human/_import`,
      body: {
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
        passwordVhangeRequired: false,
      },
      failOnStatusCode: false,
      headers: target.headers,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(409);
        return null;
      }
      return res.body?.userId || null;
    });
}

function createMachine(target: ZITADELTarget, username: string): Cypress.Chainable<any> {
  return cy
    .request({
      method: 'POST',
      url: `${target.mgmtBaseURL}/users/machine`,
      body: {
        userName: username,
        name: 'e2emachinename',
        description: 'e2emachinedescription',
      },
      failOnStatusCode: false,
      headers: target.headers,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(409);
        return null;
      }
      return res.body;
    });
}

function remove(target: ZITADELTarget, id: number) {
  return cy
    .request({
      method: 'DELETE',
      url: `${target.mgmtBaseURL}/users/${id}`,
      failOnStatusCode: false,
      headers: target.headers,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(404);
      }
    });
}

function sleep(ms: number) {
  (async () => {
    await new Promise((f) => setTimeout(f, ms));
  })();
}
