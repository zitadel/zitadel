import { login, User } from 'support/login/users';

export interface apiCallProperties {
  authHeader: string;
  mgntBaseURL: string;
  adminBaseURL: string;
}

export function apiAuth(): Cypress.Chainable<apiCallProperties> {
  return login(User.IAMAdminUser, 'Password1!', false, true).then((token) => {
    return <apiCallProperties>{
      authHeader: `Bearer ${token}`,
      mgntBaseURL: `${Cypress.env('BACKEND_URL')}/management/v1/`,
      adminBaseURL: `${Cypress.env('BACKEND_URL')}/admin/v1/`,
    };
  });
}
