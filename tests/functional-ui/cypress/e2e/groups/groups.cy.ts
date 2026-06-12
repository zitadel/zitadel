import { Context } from 'support/commands';

describe('groups', () => {
  beforeEach(() => {
    cy.context().as('ctx');
  });

  const testGroupName = 'e2egroupundertest';
  const renamedGroupName = 'e2egrouprenamed';

  // delete the group through the UI if it is present, so the suite is self-cleaning
  function deleteGroupIfExists(name: string) {
    cy.visit('/groups');
    cy.get('body').then(($body) => {
      if ($body.text().includes(name)) {
        cy.contains('tr', name).find('[data-e2e="delete-group-button"]').click({ force: true });
        cy.get('[data-e2e="confirm-dialog-button"]').click();
        cy.shouldConfirmSuccess();
      }
    });
  }

  describe('add group', () => {
    beforeEach(`ensure it doesn't exist already`, () => {
      cy.get<Context>('@ctx').then(() => {
        deleteGroupIfExists(testGroupName);
        deleteGroupIfExists(renamedGroupName);
        cy.visit('/groups');
      });
    });

    it('should add a group', () => {
      cy.get('[data-e2e="create-group-button"]').click({ force: true });
      cy.get('[data-e2e="group-name-input"]').should('be.enabled').type(testGroupName);
      cy.get('[data-e2e="group-description-input"]').should('be.enabled').type('e2e group description');
      cy.get('[data-e2e="group-dialog-save"]').click();
      cy.shouldConfirmSuccess();
      cy.contains('tr', testGroupName);
    });
  });

  describe('edit group', () => {
    beforeEach('ensure it exists', () => {
      cy.get<Context>('@ctx').then(() => {
        deleteGroupIfExists(renamedGroupName);
        cy.visit('/groups');
        cy.get('body').then(($body) => {
          if (!$body.text().includes(testGroupName)) {
            cy.get('[data-e2e="create-group-button"]').click({ force: true });
            cy.get('[data-e2e="group-name-input"]').should('be.enabled').type(testGroupName);
            cy.get('[data-e2e="group-dialog-save"]').click();
            cy.shouldConfirmSuccess();
          }
        });
      });
    });

    it('should rename the group', () => {
      cy.contains('tr', testGroupName).click();
      cy.get('[data-e2e="group-name-input"]').should('be.enabled').clear().type(renamedGroupName);
      cy.get('[data-e2e="group-dialog-save"]').click();
      cy.shouldConfirmSuccess();
      cy.contains('tr', renamedGroupName);
    });
  });

  describe('manage members', () => {
    beforeEach('ensure the group exists', () => {
      cy.visit('/groups');
      cy.get('body').then(($body) => {
        if (!$body.text().includes(testGroupName)) {
          cy.get('[data-e2e="create-group-button"]').click({ force: true });
          cy.get('[data-e2e="group-name-input"]').should('be.enabled').type(testGroupName);
          cy.get('[data-e2e="group-dialog-save"]').click();
          cy.shouldConfirmSuccess();
        }
      });
    });

    it('should show the members dialog', () => {
      cy.contains('tr', testGroupName).find('[data-e2e="group-members-button"]').click({ force: true });
      cy.get('cnsl-search-user-autocomplete');
    });
  });

  describe('remove group', () => {
    beforeEach('ensure it exists', () => {
      cy.visit('/groups');
      cy.get('body').then(($body) => {
        if (!$body.text().includes(testGroupName)) {
          cy.get('[data-e2e="create-group-button"]').click({ force: true });
          cy.get('[data-e2e="group-name-input"]').should('be.enabled').type(testGroupName);
          cy.get('[data-e2e="group-dialog-save"]').click();
          cy.shouldConfirmSuccess();
        }
      });
    });

    it('should delete the group', () => {
      cy.contains('tr', testGroupName).find('[data-e2e="delete-group-button"]').click({ force: true });
      cy.get('[data-e2e="confirm-dialog-button"]').click();
      cy.shouldConfirmSuccess();
      cy.contains('tr', testGroupName).should('not.exist');
    });
  });
});
