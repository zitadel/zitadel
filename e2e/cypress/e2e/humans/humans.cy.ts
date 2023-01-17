import { ensureDomainPolicy } from '../../support/api/policies';
import { ZITADELTarget } from 'support/commands';
import { newTarget } from 'support/api/target';
import { ensureHumanDoesntExist, ensureHumanExists } from 'support/api/users';
import { loginname } from 'support/login/login';

describe('humans', () => {
  const targetOrg = 'e2ehumans';

  [
    { mustBeDomain: false, addName: 'e2ehumanusernameaddGlobal', removeName: 'e2ehumanusernameremoveGlobal' },
    { mustBeDomain: false, addName: 'e2ehumanusernameadd@test.com', removeName: 'e2ehumanusernameremove@test.com' },
    { mustBeDomain: true, addName: 'e2ehumanusernameaddSimple', removeName: 'e2ehumanusernameremoveSimple' },
  ].forEach((user) => {
    beforeEach(() => {
      newTarget(targetOrg)
        .as('target')
        .then((target) => {
          ensureDomainPolicy(target, user.mustBeDomain, true, false);
        });
    });

    describe(`add "${user.addName}" with domain setting "${user.mustBeDomain}"`, () => {
      beforeEach(`ensure it doesn't exist already`, () => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          ensureHumanDoesntExist(target, user.addName);
          navigateToUsers(target);
        });
      });

      it('should add a user', () => {
        cy.contains('tr', user.addName).should('not.exist')
        cy.get('[data-e2e="create-user-button"]').should('be.visible').click();
        cy.url().should('contain', 'users/create');
        cy.get('[formcontrolname="email"]').should('be.visible').type('dummy@dummy.com');
        //force needed due to the prefilled username prefix
        cy.get('[formcontrolname="userName"]').should('be.visible').type(user.addName);
        cy.get('[formcontrolname="firstName"]').should('be.visible').type('e2ehumanfirstname');
        cy.get('[formcontrolname="lastName"]').should('be.visible').type('e2ehumanlastname');
        cy.get('[formcontrolname="phone"]').should('be.visible').type('+41 123456789');
        cy.get('[data-e2e="create-button"]').should('be.visible').click();
        cy.shouldConfirmSuccess();
        let loginName = user.addName;
        if (user.mustBeDomain) {
          loginName = loginname(user.addName, targetOrg);
        }
        cy.contains('[data-e2e="copy-loginname"]', loginName).should('be.visible').click();
        cy.clipboardMatches(loginName);
      });
    });

    describe(`remove "${user.removeName}" with domain setting "${user.mustBeDomain}"`, () => {
      beforeEach('ensure it exists', () => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          ensureHumanExists(target, user.removeName);
          navigateToUsers(target);
        });
      });

      // TODO: fix exact username matching (same for machines)
      // TODO: fix confirm-dialog username (same for machines)
      it.skip('should delete a human user', () => {
        Cypress.$.expr[':'].textEquals = Cypress.$.expr.createPseudo((arg) => {
          return (elem) => {
            return (
              Cypress.$(elem)
                .text()
                .trim()
                .match('^' + arg + '$').length === 1
            );
          };
        });

        let loginName = user.removeName;
        if (user.mustBeDomain) {
          loginName = loginname(user.removeName, targetOrg);
        }
        const rowSelector = `tr:contains(${user.removeName})`;
        // TODO: Is there a way to make the button visible?
        cy.get(rowSelector).find('[data-e2e="enabled-delete-button"]').click({ force: true });
        cy.get('[data-e2e="confirm-dialog-input"]').focus().should('be.visible').type(loginName);
        cy.get('[data-e2e="confirm-dialog-button"]').should('be.visible').click();
        cy.shouldConfirmSuccess();
        cy.shouldNotExist({
          selector: rowSelector,
          timeout: { ms: 2000, errMessage: 'timed out before human disappeared from the table' },
        });
      });
    });
  });
});

function navigateToUsers(target: ZITADELTarget) {
  // directly going to users is not working, atm
  cy.visit(`/org?org=${target.headers['x-zitadel-orgid']}`);
  cy.get('[data-e2e="users-nav"]').should('be.visible').click();
  cy.get('[data-e2e="list-humans"] button').should('be.visible').click();
}
