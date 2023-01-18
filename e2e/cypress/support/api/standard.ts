import { ZITADELTarget } from 'support/commands';

type IDType = number | string;

export function standardEnsureExists<IDType>(
  create: Cypress.Chainable<IDType>,
  search: () => Cypress.Chainable<IDType>,
  update: (id: IDType) => Cypress.Chainable<any> = () => cy.wrap(null),
): Cypress.Chainable<IDType> {
  return create.then((id) => {
    if (id) {
      return cy.wrap(id);
    }
    return search()
      .should((id) => id)
      .then((id) => {
        return update(id).wrap(id);
      });
  });
}

export function standardEnsureDoesntExist<IDType>(
  ensureExists: Cypress.Chainable<IDType>,
  remove: (id: IDType) => any,
  search: () => Cypress.Chainable<IDType>,
): Cypress.Chainable<null> {
  return ensureExists.then(remove).then(() => {
    search().should((id) => !id);
    return null;
  });
}

export function standardCreate<IDType>(
  target: ZITADELTarget,
  url: string,
  body: any,
  idField: string,
): Cypress.Chainable<IDType> {
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

export function standardSearch<IDType>(
  target: ZITADELTarget,
  url: string,
  find: (entity: any) => boolean,
  idField: string,
): Cypress.Chainable<IDType> {
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
