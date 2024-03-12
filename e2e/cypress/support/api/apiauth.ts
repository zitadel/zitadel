import { login, User } from 'support/login/users';
import { API, SystemAPI, Token } from './types';

const authHeaderKey = 'Authorization',
  orgIdHeaderKey = 'x-zitadel-orgid',
  backendUrl = Cypress.env('BACKEND_URL');

export function apiAuth(instanceDomain?: string): Cypress.Chainable<API> {
  return login(User.IAMAdminUser, 'Password1!', false, true, undefined, undefined, undefined, instanceDomain).then((token) => {
    return <API>{
      token: token,
      mgmtBaseURL: `${instanceDomain ? instanceDomain : backendUrl}/management/v1`,
      adminBaseURL: `${instanceDomain ? instanceDomain : backendUrl}/admin/v1`,
      authBaseURL: `${instanceDomain ? instanceDomain : backendUrl}/auth/v1`,
      assetsBaseURL: `${instanceDomain ? instanceDomain : backendUrl}/assets/v1`,
      oauthBaseURL: `${instanceDomain ? instanceDomain : backendUrl}/oauth/v2`,
      oidcBaseURL: `${instanceDomain ? instanceDomain : backendUrl}/oidc/v1`,
      samlBaseURL: `${instanceDomain ? instanceDomain : backendUrl}/saml/v2`,
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
