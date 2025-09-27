import { login, User } from '../../support/login/users';

describe('password complexity', () => {
  const orgPath = `/org`;
  const testProjectName = 'e2eproject';

  [User.OrgOwner].forEach((user) => {
    describe(`as user "${user}"`, () => {
      beforeEach(() => {
        login(user);
        cy.visit(orgPath);
        // TODO: Why force?
        cy.contains('[data-e2e="policy-card"]', 'Password Complexity').contains('button', 'Modify').click({ force: true }); // TODO: select data-e2e
      });

      // TODO: fix saving password complexity policy bug

      it(`should restrict passwords that don't have the minimal length`);
      it(`should require passwords to contain a number if option is switched on`);
      it(`should not require passwords to contain a number if option is switched off`);
      it(`should require passwords to contain a symbol if option is switched on`);
      it(`should not require passwords to contain a symbol if option is switched off`);
      it(`should require passwords to contain a lowercase letter if option is switched on`);
      it(`should not require passwords to contain a lowercase letter if option is switched off`);
      it(`should require passwords to contain an uppercase letter if option is switched on`);
      it(`should not require passwords to contain an uppercase letter if option is switched off`);
    });
  });
});
