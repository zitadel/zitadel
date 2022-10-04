import { ensureOrganizationMetadataDoesntExist, ensureOrganizationMetadataExists } from 'support/api/organization';
import { apiAuth } from '../../support/api/apiauth';

describe('organization', () => {
  const testMetadataKeyAdd = 'testkey';
  const testMetadataValueAdd = 'testvalue';

  const testMetadataKeyRemove = 'testkey';
  const testMetadataValueRemove = 'testvalue';

  describe('add org metadata', () => {
    beforeEach(`ensure it doesn't exist already`, () => {
      apiAuth().then((api) => {
        ensureOrganizationMetadataDoesntExist(api, testMetadataKeyAdd);
      });
      cy.visit(`/org`);
    });

    it('should add a metadata entry', () => {
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

  describe('remove org metadata', () => {
    beforeEach('ensure it exists', () => {
      apiAuth().then((api) => {
        cy.visit(`/org`);
        ensureOrganizationMetadataExists(api, testMetadataKeyRemove, testMetadataValueRemove);
      });
    });

    it('removes the metadata', () => {
      cy.contains('tr', testMetadataKeyRemove, { timeout: 1000 })
        .get('[data-e2e="edit-metadata-button"]')
        .click({ force: true });
      cy.get('[data-e2e="key-input-0"]').should('have.value', testMetadataKeyRemove);
      cy.get('[data-e2e="value-input-0"]').type(testMetadataValueAdd);
      cy.get(`[data-e2e="metadata-remove-button-0"]`).click();
      cy.get('.data-e2e-success');
      cy.wait(200);
      cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist');
    });
  });
});
