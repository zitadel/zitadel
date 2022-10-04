import { ensureOrganizationMetadataDoesntExist } from 'support/api/organization';
import { apiAuth } from '../../support/api/apiauth';
import { ensureHumanUserExists, ensureUserDoesntExist, ensureUserMetadataDoesntExist } from '../../support/api/users';
import { credentials, loginname } from '../../support/login/users';
import { User } from '../../support/login/users';

describe('humans', () => {
  const humansPath = `/users?type=human`;
  const testHumanUserNameAdd = 'e2ehumanusernameadd';
  const testHumanUserNameRemove = 'e2ehumanusernameremove';

  describe('add', () => {
    before(`ensure it doesn't exist already`, () => {
      apiAuth().then((apiCallProperties) => {
        ensureUserDoesntExist(apiCallProperties, testHumanUserNameAdd).then(() => {
          cy.visit(humansPath);
        });
      });
    });

    it('should add a user', () => {
      cy.get('[data-e2e="create-user-button"]').click();
      cy.url().should('contain', 'users/create');
      cy.get('[formcontrolname="email"]').type(loginname('e2ehuman', Cypress.env('ORGANIZATION')));
      //force needed due to the prefilled username prefix
      cy.get('[formcontrolname="userName"]').type(loginname(testHumanUserNameAdd, Cypress.env('ORGANIZATION')));
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
        ensureHumanUserExists(api, loginname(testHumanUserNameRemove, Cypress.env('ORGANIZATION'))).then(() => {
          cy.visit(humansPath);
        });
      });
    });

    it('should delete a human user', () => {
      cy.contains('tr', testHumanUserNameRemove)
        // doesn't work, need to force click.
        // .trigger('mouseover')
        .find('[data-e2e="enabled-delete-button"]')
        .click({ force: true });
      cy.get('[data-e2e="confirm-dialog-input"]')
        .focus()
        .type(loginname(testHumanUserNameRemove, Cypress.env('ORGANIZATION')));
      cy.get('[data-e2e="confirm-dialog-button"]').click();
      cy.get('.data-e2e-success');
      cy.wait(200);
      cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist');
    });
  });

  const testMetadataKeyAdd = 'testkey';
  const testMetadataValueAdd = 'testvalue';
  const sidenavId = 'metadata';

  describe('add user metadata', () => {
    beforeEach(`ensure it doesn't exist already`, () => {
      apiAuth().then((api) => {
        cy.visit(`/users/me?id=${sidenavId}`);
        const userIdText = cy.get('[data-e2e="user-id"]').invoke('text');
        userIdText.then((userId) => {
          ensureUserMetadataDoesntExist(api, userId, testMetadataKeyAdd);
        });
      });
    });

    it('should add a user metadata entry', () => {
      cy.get(`[data-e2e="sidenav-${sidenavId}"]`).click({ force: true });
      cy.get('[data-e2e="edit-metadata-button"]').click({ force: true });
      cy.get('[data-e2e="add-key-value"]').click({ force: true });
      cy.get('[data-e2e="key-input-0"]').type(testMetadataKeyAdd);
      cy.get('[data-e2e="value-input-0"]').type(testMetadataValueAdd);
      cy.get(`[data-e2e="metadata-save-button-${testMetadataKeyAdd}"]`).click();
      cy.get('.data-e2e-success');
      cy.wait(200);
      cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist');
    });
  });
});
