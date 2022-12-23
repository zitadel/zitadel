import { apiAuth } from '../../support/api/apiauth';
import { ensureMachineUserExists, ensureUserDoesntExist } from '../../support/api/users';
import { loginname } from '../../support/login/users';
import { ensureDomainPolicy } from '../../support/api/policies';

describe('machines', () => {
  const machinesPath = `/users?type=machine`;

  beforeEach(() => {
    apiAuth().as('api');
  });

  [
    { mustBeDomain: false, addName: 'e2emachineusernameaddGlobal', removeName: 'e2emachineusernameremoveGlobal' },
    { mustBeDomain: false, addName: 'e2emachineusernameadd@test.com', removeName: 'e2emachineusernameremove@test.com' },
    { mustBeDomain: true, addName: 'e2emachineusernameadd', removeName: 'e2emachineusernameremove' },
  ].forEach((machine) => {
    describe(`add "${machine.addName}" with domain setting "${machine.mustBeDomain}"`, () => {
      beforeEach(`ensure it doesn't exist already`, function () {
        ensureDomainPolicy(this.api, machine.mustBeDomain, false, false);
        ensureUserDoesntExist(this.api, machine.addName);
        cy.visit(machinesPath);
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
          loginName = loginname(machine.addName, Cypress.env('ORGANIZATION'));
        }
        cy.contains('[data-e2e="copy-loginname"]', loginName).click();
        cy.clipboardMatches(loginName);
      });
    });

    describe(`remove "${machine.removeName}" with domain setting "${machine.mustBeDomain}"`, () => {
      beforeEach('ensure it exists', function () {
        ensureMachineUserExists(this.api, machine.removeName);
        cy.visit(machinesPath);
      });

      let loginName = machine.removeName;
      if (machine.mustBeDomain) {
        loginName = loginname(machine.removeName, Cypress.env('ORGANIZATION'));
      }
      it('should delete a machine', () => {
        const rowSelector = `tr:contains(${machine.removeName})`;
        cy.get(rowSelector).find('[data-e2e="enabled-delete-button"]').click({ force: true });
        cy.get('[data-e2e="confirm-dialog-input"]').focus().type(loginName);
        cy.get('[data-e2e="confirm-dialog-button"]').click();
        cy.shouldConfirmSuccess();
        cy.shouldNotExist({ selector: rowSelector, timeout: 2000 });
      });

      it('should create a personal access token');
    });
  });
});
