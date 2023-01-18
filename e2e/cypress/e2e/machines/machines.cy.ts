import { ensureMachineDoesntExist, ensureMachineExists } from '../../support/api/users';
import { ensureDomainPolicy } from '../../support/api/policies';
import { newTarget } from 'support/api/target';
import { ZITADELTarget } from 'support/commands';
import { loginname } from 'support/login/login';

describe('machines', () => {
  const targetOrg = 'e2emachines';

  beforeEach(() => {
    newTarget(targetOrg).as('target');
  });

  [
    { mustBeDomain: false, addName: 'e2emachineusernameaddGlobal', removeName: 'e2emachineusernameremoveGlobal' },
    { mustBeDomain: false, addName: 'e2emachineusernameadd@test.com', removeName: 'e2emachineusernameremove@test.com' },
    { mustBeDomain: true, addName: 'e2emachineusernameaddSimple', removeName: 'e2emachineusernameremoveSimple' },
  ].forEach((machine) => {
    describe(`${machine.mustBeDomain ? 'must' : 'can'} be domain`, () => {
      beforeEach(() => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          ensureDomainPolicy(target, machine.mustBeDomain, true, false);
        });
      });
      describe(`add ${machine.addName}`, () => {
        beforeEach(`ensure it doesn't exist already`, () => {
          cy.get<ZITADELTarget>('@target').then((target) => {
            ensureMachineDoesntExist(target, machine.addName);
            navigateToMachines(target);
          });
        });

        it('should add a machine', () => {
          usernameCellDoesntExist(machine.addName);
          cy.get('[data-e2e="create-user-button"]').should('be.visible').click();
          cy.url().should('contain', 'users/create-machine');
          //force needed due to the prefilled username prefix
          cy.get('[formcontrolname="userName"]').should('be.visible').type(machine.addName);
          cy.get('[formcontrolname="name"]').should('be.visible').type('e2emachinename');
          cy.get('[formcontrolname="description"]').should('be.visible').type('e2emachinedescription');
          cy.get('[data-e2e="create-button"]').should('be.visible').click();
          cy.shouldConfirmSuccess();
          let loginName = machine.addName;
          if (machine.mustBeDomain) {
            loginName = loginname(machine.addName, targetOrg);
          }
          // TODO: Should contain loginname, not username
          cy.contains('[data-e2e="copy-loginname"]', machine.addName).should('be.visible').click();
          cy.clipboardMatches(machine.addName);
          cy.get<ZITADELTarget>('@target').then((target) => {
            navigateToMachines(target);
          });
          usernameCellExists(machine.addName);
        });
      });

      describe(`remove ${machine.removeName}`, () => {
        beforeEach('ensure it exists', () => {
          cy.get<ZITADELTarget>('@target').then((target) => {
            ensureMachineExists(target, machine.removeName);
            navigateToMachines(target);
          });
        });

        let test: Mocha.TestFunction | Mocha.PendingTestFunction = it;
        if (machine.mustBeDomain) {
          // This is flaky
          test = it.skip;
        }
        test('should delete a machine', () => {
          usernameCellExists(machine.removeName)
            .parents('tr')
            .find('[data-e2e="enabled-delete-button"]')
            // TODO: Is there a way to make the button visible?
            .click({ force: true });
          cy.get('[data-e2e="confirm-dialog-input"]').focus().type(machine.removeName);
          cy.get('[data-e2e="confirm-dialog-button"]').click();
          cy.shouldConfirmSuccess();
          usernameCellDoesntExist(machine.removeName);
        });

        it('should create a personal access token');
      });
    });
  });

  function navigateToMachines(target: ZITADELTarget) {
    // directly going to users is not working, atm
    cy.visit(`/org?org=${target.orgId}`);
    cy.get('[data-e2e="users-nav"]').should('be.visible').click();
    cy.get('[data-e2e="list-machines"] button').should('be.visible').click();
  }
});

function usernameCellDoesntExist(username: string) {
  expect(Cypress.$('[data-e2e="username-cell"]')).to.satisfy(($el: JQuery<HTMLElement>) => {
    return $el.length == 0 || cy.wrap($el).getContainingExactText(username).should('not.exist');
  });
}

function usernameCellExists(username: string) {
  return cy.get('[data-e2e="username-cell"]').getContainingExactText(username).should('exist');
}
