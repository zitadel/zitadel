import { ensureMachineDoesntExist, ensureMachineExists } from '../../support/api/users';
import { ensureDomainPolicy } from '../../support/api/policies';
import { newTarget } from 'support/api/target';
import { ZITADELTarget } from 'support/commands';
import { loginname } from 'support/login/login';

describe('machines', () => {
  const targetOrg = 'e2emachines';

  [
    { mustBeDomain: false, addName: 'e2emachineusernameaddGlobal', removeName: 'e2emachineusernameremoveGlobal' },
    { mustBeDomain: false, addName: 'e2emachineusernameadd@test.com', removeName: 'e2emachineusernameremove@test.com' },
    { mustBeDomain: true, addName: 'e2emachineusernameadd', removeName: 'e2emachineusernameremove' },
  ].forEach((machine) => {
    beforeEach(() => {
      newTarget(targetOrg)
        .as('target')
        .then((target) => {
          ensureDomainPolicy(target, machine.mustBeDomain, true, false);
        });
    });

    describe(`add "${machine.addName}" with domain setting "${machine.mustBeDomain}"`, () => {
      beforeEach(`ensure it doesn't exist already`, () => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          ensureMachineDoesntExist(target, machine.addName);
          navigateToMachines(target);
        });
      });

      it('should add a machine', () => {
        cy.get('[data-e2e="create-user-button"]').click();
        cy.url().should('contain', 'users/create-machine');
        //force needed due to the prefilled username prefix
        cy.get('[formcontrolname="userName"]').type(machine.addName);
        cy.get('[formcontrolname="name"]').type('e2emachinename');
        cy.get('[formcontrolname="description"]').type('e2emachinedescription');
        cy.get('[data-e2e="create-button"]').click();
        cy.shouldConfirmSuccess();
        let loginName = machine.addName;
        if (machine.mustBeDomain) {
          loginName = loginname(machine.addName, targetOrg);
        }
        cy.contains('[data-e2e="copy-loginname"]', loginName).click();
        cy.clipboardMatches(loginName);
      });
    });
  });
});

function navigateToMachines(target: ZITADELTarget) {
  // directly going to users is not working, atm
  cy.visit(`/org?org=${target.headers['x-zitadel-orgid']}`);
  cy.get('[data-e2e="users-nav"]').click();
  cy.get('[data-e2e="list-machines"] button').click();
}
