import { ensureOrgExists, renameOrg } from 'support/api/orgs';
import { apiAuth } from '../../support/api/apiauth';

const orgListPath = `/orgs`;
const orgPath = `/org`;

const orgNameOnCreation = 'e2eneworg';
const testOrgNameChange = 'e2erenametest';

describe('organizations', () => {
  describe('rename', () => {
    beforeEach(() => {
      apiAuth()
        .as('api')
        .then((api) => {
          ensureOrgExists(api, orgNameOnCreation)
            .as('newOrgId')
            .then((newOrgId) => {
              cy.visit(orgListPath);
            });
        });
    });

    it('should rename the organization', () => {
      const rowSelector = `tr:contains(${orgNameOnCreation})`;
      cy.get(rowSelector).children('.mat-cell').first().click({ force: true });

      cy.get('[data-e2e="actions"]').click();
      cy.get('[data-e2e="rename"]', { timeout: 1000 }).should('be.visible').click();

      cy.get('[data-e2e="name"]').focus().clear().type(testOrgNameChange);
      cy.get('[data-e2e="dialog-submit"]').click();
      cy.get('.data-e2e-success');
      cy.shouldNotExist({ selector: '.data-e2e-failure' });
    });
  });

  it('should add an organization with the personal account as org owner');
  describe('changing the current organization', () => {
    it('should update displayed organization details');
  });
});
