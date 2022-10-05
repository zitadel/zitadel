import { login, User } from 'support/login/users';
import { API } from './types';

export function apiAuth(): Cypress.Chainable<API> {
  return login(User.IAMAdminUser, 'Password1!', false, true).then((token) => {
    return <API>{
      authHeader: `Bearer ${token}`,
      mgntBaseURL: `${Cypress.env('BACKEND_URL')}/management/v1/`,
      adminBaseURL: `${Cypress.env('BACKEND_URL')}/admin/v1/`,
    };
  });
}
