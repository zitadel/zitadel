export function login(username: string, pw = 'Password1!', orgId?: string, onPasswordScreen?: () => void, passwordless = false): void {
  cy.clearAllSessionStorage();

  cy.intercept({
    method: 'POST',
    url: `/ui/login/password*`,
  }).as('password');

  cy.intercept(
    {
      method: 'GET',
      url: `/oauth/v2/authorize*`,
    },
    (req) => {
      req.query['login_hint'] = username;
      req.query['prompt'] = 'login';
      if (orgId) {
        req.query['scope'] = `${req.query['scope']} urn:zitadel:iam:org:id:${orgId}`;
      }
      req.continue();
    },
  ).as('loginAuthReq');

  cy.visit(`${Cypress.config('baseUrl')}/users/me`);

  cy.wait('@loginAuthReq');

  cy.url().should('contain', '/login');
  cy.contains(username);

  if (passwordless) {
    cy.get('#btn-login').should('be.visible').click()
  } else {
    cy.get('#password').should('be.visible').type(pw);
    cy.get('#submit-button').should('be.visible').click();

    cy.wait('@password').then((interception) => {
      if (interception.response.body.indexOf('/ui/login/mfa/prompt') === -1) {
        return;
      }
      cy.contains('button', 'skip').click();
    });

    onPasswordScreen ? onPasswordScreen() : null;
  }

  cy.contains('[data-e2e="top-view-subtitle"]', username).then(($el) => {
    expect($el.text().trim()).to.equal(username);
  });
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
