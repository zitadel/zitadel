import { apiAuth } from '../../support/api/apiauth';
import { ensureHumanUserExists, ensureUserDoesntExist } from '../../support/api/users';
import { loginname } from '../../support/login/users';
import { ensureDomainPolicy } from '../../support/api/policies';
import { Context } from 'support/commands';

describe('humans', () => {
  const humansPath = `/users?type=human`;

  beforeEach(() => {
    cy.context().as('ctx');
  });

  [
    { mustBeDomain: false, addName: 'e2ehumanusernameaddGlobal', removeName: 'e2ehumanusernameremoveGlobal' },
    { mustBeDomain: false, addName: 'e2ehumanusernameadd@test.com', removeName: 'e2ehumanusernameremove@test.com' },
    { mustBeDomain: true, addName: 'e2ehumanusernameadd', removeName: 'e2ehumanusernameremove' },
  ].forEach((user) => {
    describe(`add "${user.addName}" with domain setting "${user.mustBeDomain}"`, () => {
      beforeEach(`ensure it doesn't exist already`, function () {
        cy.get<Context>('@ctx').then((ctx) => {
          ensureDomainPolicy(ctx.api, user.mustBeDomain, true, false);
          ensureUserDoesntExist(ctx.api, user.addName);
          cy.visit(humansPath);
        });
      });

      it('should add a user', () => {
        cy.get('[data-e2e="create-user-button"]').click();
        cy.url().should('contain', 'users/create');
        cy.get('[formcontrolname="email"]').type('dummy@dummy.com');
        //force needed due to the prefilled username prefix
        cy.get('[formcontrolname="userName"]').type(user.addName);
        cy.get('[formcontrolname="firstName"]').type('e2ehumanfirstname');
        cy.get('[formcontrolname="lastName"]').type('e2ehumanlastname');
        cy.get('[formcontrolname="phone"]').type('+41 123456789');
        cy.get('[data-e2e="create-button"]').click();
        cy.get('.data-e2e-success');
        let loginName = user.addName;
        if (user.mustBeDomain) {
          loginName = loginname(user.addName, Cypress.env('ORGANIZATION'));
        }
        cy.contains('[data-e2e="copy-loginname"]', loginName).click();
        cy.clipboardMatches(loginName);
        cy.shouldNotExist({ selector: '.data-e2e-failure' });
      });
    });

    describe(`remove "${user.removeName}" with domain setting "${user.mustBeDomain}"`, () => {
      beforeEach('ensure it exists', function () {
        cy.get<Context>('@ctx').then((ctx) => {
          ensureHumanUserExists(ctx.api, user.removeName);
        });
        cy.visit(humansPath);
      });

      let loginName = user.removeName;
      if (user.mustBeDomain) {
        loginName = loginname(user.removeName, Cypress.env('ORGANIZATION'));
      }
      it('should delete a human user', () => {
        const rowSelector = `tr:contains(${user.removeName})`;
        cy.get(rowSelector).find('[data-e2e="enabled-delete-button"]').click({ force: true });
        cy.get('[data-e2e="confirm-dialog-input"]').focus().type(loginName);
        cy.get('[data-e2e="confirm-dialog-button"]').click();
        cy.get('.data-e2e-success');
        cy.shouldNotExist({ selector: rowSelector, timeout: 2000 });
        cy.shouldNotExist({ selector: '.data-e2e-failure' });
      });
    });
  });
});
