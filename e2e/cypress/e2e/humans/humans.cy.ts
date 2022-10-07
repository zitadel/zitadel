import { apiAuth } from '../../support/api/apiauth';
import { ensureHumanUserExists, ensureUserDoesntExist } from '../../support/api/users';
import { loginname } from '../../support/login/users';

describe('humans', () => {
  const humansPath = `/users?type=human`;
  const testHumanUserNameAdd = 'e2ehumanusernameadd';
  const testHumanUserNameRemove = 'e2ehumanusernameremove';

  beforeEach(() => {
    apiAuth().as('api');
  });

  describe('add', () => {
    beforeEach(`ensure it doesn't exist already`, function () {
      ensureUserDoesntExist(this.api, loginname(testHumanUserNameAdd, Cypress.env('ORGANIZATION')));
      cy.visit(humansPath);
    });

    it('should add a user', () => {
      cy.get('[data-e2e="create-user-button"]').click();
      cy.url().should('contain', 'users/create');
      cy.get('[formcontrolname="email"]').type('dummy@dummy.com');
      //force needed due to the prefilled username prefix
      cy.get('[formcontrolname="userName"]').type(loginname(testHumanUserNameAdd, Cypress.env('ORGANIZATION')));
      cy.get('[formcontrolname="firstName"]').type('e2ehumanfirstname');
      cy.get('[formcontrolname="lastName"]').type('e2ehumanlastname');
      cy.get('[formcontrolname="phone"]').type('+41 123456789');
      cy.get('[data-e2e="create-button"]').click();
      cy.get('.data-e2e-success');
      const loginName = loginname(testHumanUserNameAdd, Cypress.env('ORGANIZATION'));
      cy.contains('[data-e2e="copy-loginname"]', loginName).click();
      cy.clipboardMatches(loginName);
      cy.shouldNotExist({ selector: '.data-e2e-failure' });
    });
  });

  describe('remove', () => {
    beforeEach('ensure it exists', function () {
      ensureHumanUserExists(this.api, loginname(testHumanUserNameRemove, Cypress.env('ORGANIZATION')));
      cy.visit(humansPath);
    });

    it('should delete a human user', () => {
      const rowSelector = `tr:contains(${testHumanUserNameRemove})`;
      cy.get(rowSelector).find('[data-e2e="enabled-delete-button"]').click({ force: true });
      cy.get('[data-e2e="confirm-dialog-input"]')
        .focus()
        .type(loginname(testHumanUserNameRemove, Cypress.env('ORGANIZATION')));
      cy.get('[data-e2e="confirm-dialog-button"]').click();
      cy.get('.data-e2e-success');
      cy.shouldNotExist({ selector: '.data-e2e-failure' });
      cy.shouldNotExist({ selector: rowSelector, timeout: 2000 });
    });
  });
});
