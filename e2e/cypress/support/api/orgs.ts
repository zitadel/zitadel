import { ZITADELTarget } from 'support/commands';
import { newOrgTarget } from './target';

export function ensureOrgExists(target: ZITADELTarget, name: string): Cypress.Chainable<ZITADELTarget> {
  return createOrg(target, name).then((id) => {
    if (id) {
      return cy.wrap(newOrgTarget(target, id));
    }
    return search(target, name).then((id) => {
      if (id) {
        return cy.wrap(newOrgTarget(target, id));
      }
      sleep(6_000);
      cy.log('retrying');
      return search(target, name).then((id) => {
        if (id) {
          return cy.wrap(newOrgTarget(target, id));
        }
        sleep(6_000);
        cy.log('retrying');
        debugger;
        return search(target, name).then((id) => cy.wrap(newOrgTarget(target, id)));
      });
    });
  });
}

function search(target: ZITADELTarget, name: string): Cypress.Chainable<number> {
  return cy
    .request({
      method: 'POST',
      url: `${target.adminBaseURL}/orgs/_search`,
      headers: target.headers,
    })
    .then((res) => {
      return res.body?.result?.find((entity) => entity.name == name)?.id || cy.wrap(null);
    });
}

function createOrg(target: ZITADELTarget, name: string): Cypress.Chainable<number> {
  return cy
    .request({
      method: 'POST',
      url: `${target.mgmtBaseURL}/orgs`,
      body: { name: name },
      headers: target.headers,
      failOnStatusCode: false,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(409);
        return null;
      }
      return res.body.id;
    });
}

export function removeOrg(target: ZITADELTarget): Cypress.Chainable<null> {
  return cy
    .request({
      method: 'DELETE',
      url: `${target.mgmtBaseURL}/orgs/me`,
      headers: target.headers,
      failOnStatusCode: false,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(404);
      }
      return null;
    });
}

function sleep(ms: number) {
  (async () => {
    await new Promise((f) => setTimeout(f, ms));
  })();
}
