import { Context } from '../../support/api/target';
import { ensureHumanExists } from '../../support/api/users';
import { User } from '../../support/login/session';

describe('login policy', () => {
  const orgPath = `/org`;

  [User.OrgOwner].forEach((user) => {
    describe(`as user "${user}"`, () => {
      beforeEach(() => {
        // login(user);
        cy.visit(orgPath);
        // TODO: Why force?
        cy.contains('[data-e2e="policy-card"]', 'Login Policy').contains('button', 'Modify').click({ force: true }); // TODO: select data-e2e
        ctx().then((api) => {
          ensureHumanExists(api, User.LoginPolicyUser);
        });
      });

      // TODO: verify email

      it(`username and password disallowed`);
      it(`registering is allowed`);
      it(`registering is disallowed`);
      it(`login by an external IDP is allowed`);
      it(`login by an external IDP is disallowed`);
      it(`MFA is forced`);
      it(`MFA is not forced`);
      it(`the password reset option is hidden`);
      it(`the password reset option is shown`);
      it(`passwordless login is allowed`);
      it(`passwordless login is disallowed`);
      describe('identity providers', () => {});
    });
  });
});
