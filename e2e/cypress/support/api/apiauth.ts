import { login, User } from 'support/login/users';
import { API, SystemAPI, Token } from './types';

const authHeaderKey = 'Authorization',
  orgIdHeaderKey = 'x-zitadel-orgid';

export function apiAuth(): Cypress.Chainable<API> {
  return login(User.IAMAdminUser, 'Password1!', false, true).then((token) => {
    return <API>{
      token: token,
      mgmtBaseURL: `${Cypress.env('BACKEND_URL')}/management/v1`,
      adminBaseURL: `${Cypress.env('BACKEND_URL')}/admin/v1`,
    };
  });
}

export function systemAuth(): Cypress.Chainable<SystemAPI> {
  return cy.task('systemToken').then((tokem) => {
    return <SystemAPI>{
      token: tokem,
      baseURL: `${Cypress.env('BACKEND_URL')}/system/v1`,
    };
  });
}

export function requestHeaders(token: Token, orgId?: number): object {
  const headers = { [authHeaderKey]: `Bearer ${token}` };
  if (orgId) {
    headers[orgIdHeaderKey] = orgId;
  }
  return headers;
}
