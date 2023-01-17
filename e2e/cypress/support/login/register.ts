export function register(email: string, orgId: number): Cypress.Chainable<string> {
  cy.clearAllSessionStorage();

  const pw = 'Password1!';
  const userFirstname = 'e2efirstname';
  const userLastname = 'e2elastname';

  cy.intercept(
    {
      method: 'GET',
      url: `/oauth/v2/authorize*`,
    },
    (req) => {
      req.query['prompt'] = 'create';
      req.query['login_hint'] = email;
      req.query['scope'] = `${req.query['scope']} urn:zitadel:iam:org:id:${orgId}`;
      req.continue();
    },
  ).as('regAuthReq');
  cy.visit(`${Cypress.config('baseUrl')}/users/me`);
  cy.wait('@regAuthReq');
  cy.get('#firstname').type(userFirstname);
  cy.get('#lastname').type(userLastname);
  cy.get('#register-password').type(pw);
  cy.get('#register-password-confirmation').type(pw);
  cy.get('#register-term-confirmation').check({ force: true });
  cy.get('#register-term-confirmation-privacy').check({ force: true });
  cy.get('form').submit();
  cy.get('#password').type(pw);

  cy.intercept({
    method: 'POST',
    url: `/ui/login/password*`,
    times: 1,
  }).as('password');
  cy.get('form').submit();

  cy.wait('@password').then((interception) => {
    if (interception.response.body.indexOf('/ui/login/mfa/prompt') === -1) {
      return;
    }
    cy.contains('button', 'skip').click();
  });

  cy.contains('[data-e2e="top-view-subtitle"]', email).then(($el) => {
    expect($el.text().trim()).to.equal(email);
  });

  return cy.get('[data-e2e="user-id"]').then(($el) => {
    return $el.text().trim();
  });
}
