import { ZITADELTarget } from 'support/commands';

export function ensureRoleExists(api: ZITADELTarget, projectId: number, roleKey: string): Cypress.Chainable<null> {
  return cy
    .request({
      method: 'POST',
      url: `${api.mgmtBaseURL}/projects/${projectId}/roles`,
      body: {
        roleKey: roleKey,
        displayName: roleKey,
      },
      headers: api.headers,
      failOnStatusCode: false,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(409);
      }
      return null;
    });
}

export function ensureRoleDoesntExist(api: ZITADELTarget, projectId: number, roleKey: string): Cypress.Chainable<null> {
  return cy
    .request({
      method: 'DELETE',
      url: `${api.mgmtBaseURL}/projects/${projectId}/roles/${roleKey}`,
      headers: api.headers,
      failOnStatusCode: false,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(404);
      }
      return null;
    });
}
