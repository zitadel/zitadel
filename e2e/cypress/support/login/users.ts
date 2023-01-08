export enum User {
  OrgOwner = 'org_owner',
  OrgOwnerViewer = 'org_owner_viewer',
  OrgProjectCreator = 'org_project_creator',
  LoginPolicyUser = 'login_policy_user',
  PasswordComplexityUser = 'password_complexity_user',
  IAMAdminUser = 'zitadel-admin',
}

export function loginAsPredefinedUser(user: User) {
  return login(loginname(<string>user, Cypress.env('ORGANIZATION')), undefined, false);
}

export function login(
  username: string,
  pw = 'Password1!',
  force?: boolean,
  onUsernameScreen?: () => void,
  onPasswordScreen?: () => void,
  onAuthenticated?: () => void,
): Cypress.Chainable<string> {
  const loginUrl: string = '/ui/login';
  const issuerUrl: string = '/oauth/v2';

  return cy
    .session(
      username,
      () => {
        const cookies = new Map<string, string>();

        cy.intercept(
          {
            times: 6,
          },
          (req) => {
            req.headers['cookie'] = requestCookies(cookies);
            req.continue((res) => {
              updateCookies(res.headers['set-cookie'] as string[], cookies);
            });
          },
        );

        let userToken: string;
        cy.intercept(
          {
            method: 'POST',
            url: `${issuerUrl}/token`,
          },
          (req) => {
            req.continue((res) => {
              userToken = res.body['access_token'];
            });
          },
        ).as('token');

        cy.intercept({
          method: 'POST',
          url: `${loginUrl}/password*`,
          times: 1,
        }).as('password');

        cy.visit(loginUrl, { retryOnNetworkFailure: true });

        onUsernameScreen ? onUsernameScreen() : null;
        cy.get('#loginName').type(username);
        cy.get('#submit-button').click();

        onPasswordScreen ? onPasswordScreen() : null;
        cy.get('#password').type(pw);
        cy.get('#submit-button').click();

        cy.wait('@password').then((interception) => {
          if (interception.response.body.indexOf('/ui/login/mfa/prompt') === -1) {
            return;
          }

          cy.contains('button', 'skip').click();
        });

        cy.wait('@token').then(() => {
          cy.task('safetoken', { key: username, token: userToken });
        });

        onAuthenticated ? onAuthenticated() : null;
      },
      {
        validate: () => {
          if (force) {
            throw new Error('clear session');
          }
        },
      },
    )
    .then(() => {
      return cy.task('loadtoken', { key: username });
    });
}

export function loginname(withoutDomain: string, org?: string): string {
  return `${withoutDomain}@${org}.${host(Cypress.config('baseUrl'))}`;
}

function updateCookies(newCookies: string[] | undefined, currentCookies: Map<string, string>) {
  if (newCookies === undefined) {
    return;
  }
  newCookies.forEach((cs) => {
    cs.split('; ').forEach((cookie) => {
      const idx = cookie.indexOf('=');
      currentCookies.set(cookie.substring(0, idx), cookie.substring(idx + 1));
    });
  });
}

function requestCookies(currentCookies: Map<string, string>): string[] {
  let list: Array<string> = [];
  currentCookies.forEach((val, key) => {
    list.push(key + '=' + val);
  });
  return list;
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
