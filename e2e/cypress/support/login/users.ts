export enum User {
  OrgOwner = 'org_owner',
  OrgOwnerViewer = 'org_owner_viewer',
  OrgProjectCreator = 'org_project_creator',
  LoginPolicyUser = 'login_policy_user',
  PasswordComplexityUser = 'password_complexity_user',
  IAMAdminUser = 'zitadel-admin',
}

export function loginAsPredefinedUser(user: User) {
  return login(loginname(<string>user, Cypress.env('ORGANIZATION')));
}

export function login(
  username: string,
  pw?: string,
  onPasswordScreen?: () => void,
): Cypress.Chainable<string> {
    // We want to have a clean session but miss cypresses sesssion cache
    return cy.session(Math.random().toString(), plainLogin(username, pw, onPasswordScreen)).then(() => {
    return cy.task('loadtoken', { key: username });
  });
}

function plainLogin(username: string, pw = 'Password1!', onPasswordScreen?: () => void): () => void {
  const loginUrl: string = '/ui/login';
  const issuerUrl: string = '/oauth/v2';

  return () => {
    cy.intercept({
      method: 'POST',
      url: `${issuerUrl}/token*`,
    }).as('token');

    cy.intercept({
      method: 'POST',
      url: `${loginUrl}/password*`,
    }).as('password');

    cy.intercept(
      {
        method: 'GET',
        url: `${issuerUrl}/authorize*`,
      },
      (req) => {
        req.query['login_hint'] = username;
        req.query['prompt'] = 'login';
        req.continue();
      },
    ).as('loginAuthReq');

    cy.visit('/users/me');

    cy.wait('@loginAuthReq');

    cy.contains(username);

    cy.get('#password').type(pw);
    cy.get('#submit-button').click();
    onPasswordScreen ? onPasswordScreen() : null;

    cy.wait('@password').then((interception) => {
      if (interception.response.body.indexOf('/ui/login/mfa/prompt') === -1) {
        return;
      }
      cy.contains('button', 'skip').click();
    });

    cy.wait('@token').then((interception) => {
      cy.task('safetoken', { key: username, token: interception.response.body.access_token });
    });

    cy.get('[data-e2e="top-view-title"]');
  };
}

export function loginname(withoutDomain: string, org?: string): string {
  return `${withoutDomain}@${org}.${host(Cypress.config('baseUrl'))}`;
}

export function host(url: string): string {
  return stripPort(stripProtocol(url));
}

function stripPort(s: string): string {
  const idx = s.indexOf(':');
  return idx === -1 ? s : s.substring(0, idx);
}

function stripProtocol(url: string): string {
  return url.replace('http://', '').replace('https://', '');
}
