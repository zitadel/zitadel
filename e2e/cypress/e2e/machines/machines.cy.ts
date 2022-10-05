import { apiAuth } from '../../support/api/apiauth';
import { ensureMachineUserExists, ensureUserDoesntExist } from '../../support/api/users';
import { loginname } from '../../support/login/users';

describe('machines', () => {
  beforeEach(() => {
    apiAuth().as('api');
  });

  const machinesPath = `/users?type=machine`;
  const testMachineUserNameAdd = 'e2emachineusernameadd';
  const testMachineUserNameRemove = 'e2emachineusernameremove';

  describe('add', () => {
    beforeEach(`ensure it doesn't exist already`, function () {
      ensureUserDoesntExist(this.api, testMachineUserNameAdd);
      cy.visit(machinesPath);
    });

    it('should add a machine', () => {
      cy.get('[data-e2e="create-user-button"]').click();
      cy.url().should('contain', 'users/create-machine');
      //force needed due to the prefilled username prefix
      cy.get('[formcontrolname="userName"]').type(testMachineUserNameAdd);
      cy.get('[formcontrolname="name"]').type('e2emachinename');
      cy.get('[formcontrolname="description"]').type('e2emachinedescription');
      cy.get('[data-e2e="create-button"]').click();
      cy.get('.data-e2e-success');
      cy.contains('[data-e2e="copy-loginname"]', testMachineUserNameAdd).click();
      cy.clipboardMatches(testMachineUserNameAdd);
      cy.shouldNotExist({ selector: '.data-e2e-failure' });
    });
  });

  describe('edit', () => {
    beforeEach('ensure it exists', function () {
      ensureMachineUserExists(this.api, testMachineUserNameRemove);
      cy.visit(machinesPath);
    });

    it('should delete a machine', () => {
      cy.contains('tr', testMachineUserNameRemove)
        .as('machineUserRow')
        .find('[data-e2e="enabled-delete-button"]')
        .click({ force: true });
      cy.get('[data-e2e="confirm-dialog-input"]').focus().type(testMachineUserNameRemove);
      cy.get('[data-e2e="confirm-dialog-button"]').click();
      cy.get('.data-e2e-success');
      cy.shouldNotExist({ selector: '.data-e2e-failure' });
      cy.get('@machineUserRow').shouldNotExist({ timeout: 2000 });
    });

    it('should create a personal access token');
  });
});
