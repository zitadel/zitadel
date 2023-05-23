export enum User {
  OrgOwner = 'org_owner',
  OrgOwnerViewer = 'org_owner_viewer',
  OrgProjectCreator = 'org_project_creator',
  LoginPolicyUser = 'login_policy_user',
  PasswordComplexityUser = 'password_complexity_user',
  IAMAdminUser = 'zitadel-admin',
}

export function login(
  user: User,
  pw?: string,
  force?: boolean,
  skipMFAChangePW?: boolean,
  onUsernameScreen?: () => void,
  onPasswordScreen?: () => void,
  onAuthenticated?: () => void,
): Cypress.Chainable<string> {
  let creds = credentials(user, pw);

  const loginUrl: string = '/ui/login';
  const issuerUrl: string = '/oauth/v2';

  return cy
    .session(
      creds.username,
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

        cy.visit(Cypress.config('baseUrl'), { retryOnNetworkFailure: true });

        onUsernameScreen ? onUsernameScreen() : null;
        cy.get('#loginName').type(creds.username);
        cy.get('#submit-button').click();

        onPasswordScreen ? onPasswordScreen() : null;
        cy.get('#password').type(creds.password);
        cy.get('#submit-button').click();

        cy.wait('@password').then((interception) => {
          if (interception.response.body.indexOf(`${loginUrl}/mfa/prompt`) === -1) {
            return;
          }

          cy.contains('button', 'skip').click();
        });

        cy.wait('@token').then(() => {
          cy.task('safetoken', { key: creds.username, token: userToken });
        });

        onAuthenticated ? onAuthenticated() : null;

        cy.visit('/');

        cy.get('[data-e2e=authenticated-welcome]', {
          timeout: 50_000,
        });
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
      return cy.task('loadtoken', { key: creds.username });
    });
}

export function loginname(withoutDomain: string, org?: string): string {
  return `${withoutDomain}@${org}.${host(Cypress.config('baseUrl'))}`;
}

function credentials(user: User, pw?: string) {
  // TODO: ugly
  const woDomain = user == User.IAMAdminUser ? User.IAMAdminUser : `${user}_user_name`;
  const org = Cypress.env('ORGANIZATION') ? Cypress.env('ORGANIZATION') : 'zitadel';

  return {
    username: loginname(woDomain, org),
    password: pw ? pw : Cypress.env(`${user}_password`),
  };
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
