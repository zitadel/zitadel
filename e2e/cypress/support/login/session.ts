import { login, loginname } from './login';

export enum User {
  OrgOwner = 'org_owner',
  OrgOwnerViewer = 'org_owner_viewer',
  OrgProjectCreator = 'org_project_creator',
  LoginPolicyUser = 'login_policy_user',
  PasswordComplexityUser = 'password_complexity_user',
  IAMAdminUser = 'zitadel-admin',
}

export function sessionAsPredefinedUser(user: User) {
  return session(loginname(<string>user, Cypress.env('ORGANIZATION')), null);
}

export function session(username: string, orgId: string): Cypress.Chainable<string> {
  // We want to have a clean session but miss cypresses sesssion cache
  return cy.session([username, orgId], () => login(username, orgId), {
    cacheAcrossSpecs: true,
  });
}
