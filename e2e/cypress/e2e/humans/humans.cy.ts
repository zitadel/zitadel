import { ensureDomainPolicy } from '../../support/api/policies';
import { ZITADELTarget } from 'support/commands';
import { newTarget } from 'support/api/target';
import { ensureHumanDoesntExist, ensureHumanExists } from 'support/api/users';
import { loginname } from 'support/login/login';

describe('humans', () => {
  const targetOrg = 'e2ehumans';
  beforeEach(() => {
    newTarget(targetOrg).as('target');
  });

  [
    { mustBeDomain: false, addName: 'e2ehumanusernameaddGlobal', removeName: 'e2ehumanusernameremoveGlobal' },
    { mustBeDomain: false, addName: 'e2ehumanusernameadd@test.com', removeName: 'e2ehumanusernameremove@test.com' },
    { mustBeDomain: true, addName: 'e2ehumanusernameaddSimple', removeName: 'e2ehumanusernameremoveSimple' },
  ].forEach((user) => {
    describe(`must ${user.mustBeDomain ? '' : 'not '}be domain`, () => {
      beforeEach(() => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          ensureDomainPolicy(target, user.mustBeDomain, true, false);
        });
      });

      describe(`add ${user.addName}`, () => {
        beforeEach(`ensure it doesn't exist already`, () => {
          cy.get<ZITADELTarget>('@target').then((target) => {
            ensureHumanDoesntExist(target, user.addName);
            navigateToUsers(target);
          });
        });

        it('should add a user', () => {
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
          cy.contains('[data-e2e="copy-loginname"]', user.addName).should('be.visible').click();
          cy.clipboardMatches(user.addName);
        });
      });

      describe(`remove ${user.removeName}`, () => {
        beforeEach('ensure it exists', () => {
          cy.get<ZITADELTarget>('@target').then((target) => {
            ensureHumanExists(target, user.removeName);
            navigateToUsers(target);
          });
        });

        it('should delete a human user', () => {
          getUsernameCell(user.removeName)
            .parents('tr')
            .find('[data-e2e="enabled-delete-button"]')
          // TODO: Is there a way to make the button visible?
          .click({ force: true });
          cy.get('[data-e2e="confirm-dialog-input"]').focus().should('be.visible').type(user.removeName);
          cy.get('[data-e2e="confirm-dialog-button"]').should('be.visible').click();
          cy.shouldConfirmSuccess();
          usernameCellDoesntExist(user.removeName);
        });
      });
    });
  });
});

function navigateToUsers(target: ZITADELTarget) {
  // directly going to users is not working, atm
  cy.visit(`/org?org=${target.orgId}`);
  cy.get('[data-e2e="users-nav"]').should('be.visible').click();
  cy.get('[data-e2e="list-humans"] button').should('be.visible').click();
}


function usernameCellDoesntExist(username: string) {
  cy.waitUntil(() => {
    return getUsernameCell(username).then(($el) => $el.length === 0);
  });
}

function getUsernameCell(username: string) {
  return cy.get('[data-e2e="username-cell"]').containsExactly(username);
}
