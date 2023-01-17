import { ZITADELTarget } from 'support/commands';

export function standardEnsureExists(
  create: Cypress.Chainable<number>,
  search: () => Cypress.Chainable<number>,
  update: (id: number) => Cypress.Chainable<any> = () => cy.wrap(null),
): Cypress.Chainable<number> {
  return create.then((id) => {
    if (id) {
      return cy.wrap(id);
    }
    return search().should((id) => id).then(id => {
      return update(id).wrap(id)
    });
  });
}

export function standardEnsureDoesntExist(ensureExists: Cypress.Chainable<number>, remove: (id: number) => any, search: () => Cypress.Chainable<number>) {
  ensureExists.then(remove).then(()=> {
    search().should((id) => !id)
  });
}

export function standardCreate(target: ZITADELTarget, url: string, body: any, idField: string): Cypress.Chainable<number> {
  return cy
    .request({
      method: 'POST',
      url: url,
      body: body,
      failOnStatusCode: false,
      headers: target.headers,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(409);
        return null;
      }
      return res.body?.[idField] || null;
    });
}

export function standardSearch(
  target: ZITADELTarget,
  url: string,
  find: (entity: any) => boolean,
  idField: string,
): Cypress.Chainable<number> {
  return cy
    .request({
      method: 'POST',
      url: url,
      headers: target.headers,
    })
    .then((res) => {
      const found = res.body?.result?.find(find);
      if (!found) {
        cy.log("couldn't find entity");
      }
      return found?.[idField] || null;
    });
}

export function standardUpdate(target: ZITADELTarget, url: string, body: any): Cypress.Chainable<Cypress.Response<any>> {
  return cy
    .request({
      method: 'PUT',
      url: url,
      body: body,
      failOnStatusCode: false,
      headers: target.headers,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(400);
        expect(res.body.message).to.contain('No changes');
      }
    });
}

export function standardRemove(target: ZITADELTarget, url: string) {
  return cy
    .request({
      method: 'DELETE',
      url: url,
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
