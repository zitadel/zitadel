import { requestHeaders } from './apiauth';
import { API } from './types';

export function createJWTIDP(api: API, name: string): Cypress.Chainable<string> {
  return cy
    .request({
      method: 'POST',
      url: `${api.adminBaseURL}/idps/jwt`,
      headers: requestHeaders(api),
      body: {
        name,
        jwtEndpoint: `https://${name}.example.invalid/jwt`,
        issuer: `https://${name}.example.invalid`,
        keysEndpoint: `https://${name}.example.invalid/keys`,
        headerName: 'authorization',
      },
    })
    .then((response) => {
      expect(response.status).to.equal(200);
      return response.body.idpId as string;
    });
}

export function deleteIDP(api: API, idpId: string): Cypress.Chainable<null> {
  return cy
    .request({
      method: 'DELETE',
      url: `${api.adminBaseURL}/idps/templates/${idpId}`,
      headers: requestHeaders(api),
      failOnStatusCode: false,
    })
    .then((response) => {
      expect([200, 404]).to.include(response.status);
      return null;
    });
}

export function addUserIDPLink(
  api: API,
  userId: string,
  idpId: string,
  linkedUserId: string,
  userName: string,
): Cypress.Chainable<null> {
  const userBaseURL = api.mgmtBaseURL.replace('/management/v1', '/v2');

  return cy
    .request({
      method: 'POST',
      url: `${userBaseURL}/users/${userId}/links`,
      headers: requestHeaders(api),
      body: {
        idpLink: {
          idpId,
          userId: linkedUserId,
          userName,
        },
      },
    })
    .then((response) => {
      expect(response.status).to.equal(200);
      return null;
    });
}

export function waitForLinkedIDPCount(api: API, userId: string, expectedCount: number): Cypress.Chainable<null> {
  return cy
    .waitUntil(
      () =>
        cy
          .request({
            method: 'POST',
            url: `${api.mgmtBaseURL}/users/${userId}/idps/_search`,
            headers: requestHeaders(api),
            body: {
              query: {
                limit: 20,
              },
            },
          })
          .then((response) => {
            return (response.body.result?.length ?? 0) === expectedCount;
          }),
      {
        timeout: 90_000,
        interval: 1_000,
        errorMsg: `timed out waiting for ${expectedCount} linked idps on user ${userId}`,
      },
    )
    .then(() => null);
}
