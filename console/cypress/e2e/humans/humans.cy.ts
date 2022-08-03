import { apiAuth } from '../../support/api/apiauth';
import { ensureHumanUserExists, ensureUserDoesntExist } from '../../support/api/users';
import { login, User, loginname } from '../../support/login/users';

describe('humans', () => {
  const humansPath = `/ui/console/users?type=human`;
  const testHumanUserNameAdd = 'e2ehumanusernameadd';
  const testHumanUserNameRemove = 'e2ehumanusernameremove';

  describe('add', () => {
    before(`ensure it doesn't exist already`, () => {
      apiAuth().then((apiCallProperties) => {
        ensureUserDoesntExist(apiCallProperties, testHumanUserNameAdd).then(()=>{
          cy.visit(humansPath);
        });
      });
    });

    it('should add a user', () => {
      cy.get('[data-e2e="action-key-add"]').parents('[data-e2e="create-user-button"]').click();
      cy.url().should('contain', 'users/create');
      cy.get('[formcontrolname="email"]').type(loginname('e2ehuman'));
      //force needed due to the prefilled username prefix
      cy.get('[formcontrolname="userName"]').type(testHumanUserNameAdd, { force: true });
      cy.get('[formcontrolname="firstName"]').type('e2ehumanfirstname');
      cy.get('[formcontrolname="lastName"]').type('e2ehumanlastname');
      cy.get('[formcontrolname="phone"]').type('+41 123456789');
      cy.get('[data-e2e="create-button"]').click();
      cy.get('.data-e2e-success');
      cy.wait(200);
      cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist');
    });
  });

  describe('remove', () => {
    before('ensure it exists', () => {
      apiAuth().then((api) => {
        ensureHumanUserExists(api, testHumanUserNameRemove).then(()=>{
          cy.visit(humansPath);
        });
      });
    });

    it('should delete a human user', () => {
      cy.contains('tr', testHumanUserNameRemove)
        // doesn't work, need to force click.
        // .trigger('mouseover')
        .find('[e2e-data="enabled-delete-button"]')
        .click({force: true});
      cy.get('[e2e-data="confirm-dialog-input"]').click().type(loginname(testHumanUserNameRemove, Cypress.env('org')));
      cy.get('[e2e-data="confirm-dialog-button"]').click();
      cy.get('.data-e2e-success');
      cy.wait(200);
      cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist');
    });
  });
});