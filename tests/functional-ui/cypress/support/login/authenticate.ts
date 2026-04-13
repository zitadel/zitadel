export function authenticate(
  loginUrl: string,
  creds: {
    password: any;
    username: string;
  },
  onUsernameScreen: () => void,
  onPasswordScreen: () => void,
  onAuthenticated: () => void,
) {
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

  cy.wait('@token');

  onAuthenticated ? onAuthenticated() : null;

  cy.wait(1000);
}
