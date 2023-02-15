import { login, User } from 'support/login/users';
import { API, SystemAPI, Token } from './types';

const authHeaderKey = 'Authorization',
  orgIdHeaderKey = 'x-zitadel-orgid',
  backendUrl = Cypress.env('BACKEND_URL');

export function apiAuth(): Cypress.Chainable<API> {
  return login(User.IAMAdminUser, 'Password1!', false, true).then((token) => {
    return <API>{
      token: token,
      mgmtBaseURL: `${backendUrl}/management/v1`,
      adminBaseURL: `${backendUrl}/admin/v1`,
      authBaseURL: `${backendUrl}/auth/v1`,
      assetsBaseURL: `${backendUrl}/assets/v1`,
      oauthBaseURL: `${backendUrl}/oauth/v2`,
      oidcBaseURL: `${backendUrl}/oidc/v1`,
      samlBaseURL: `${backendUrl}/saml/v2`,
    };
  });
}

export function systemAuth(): Cypress.Chainable<SystemAPI> {
  return cy.task('systemToken').then((token) => {
    return <SystemAPI>{
      token: token,
      baseURL: `${backendUrl}/system/v1`,
    };
  });
}

export function requestHeaders(token: Token, orgId?: string): object {
  const headers = { [authHeaderKey]: `Bearer ${token.token}` };
  if (orgId) {
    headers[orgIdHeaderKey] = orgId;
  }
  return headers;
}
