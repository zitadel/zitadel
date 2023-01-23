import { ensureOrgExists } from 'support/api/orgs';
import { apiAuth } from '../../support/api/apiauth';
import { v4 as uuidv4 } from 'uuid';

const orgPath = `/org`;

const orgNameOnCreation = 'e2eorgrename';
const testOrgNameChange = uuidv4();

describe('organizations', () => {
  describe('rename', () => {
    beforeEach(() => {
      apiAuth()
        .as('api')
        .then((api) => {
          ensureOrgExists(api, orgNameOnCreation)
          .as('newOrgId')
          .then((newOrgId) => {
            cy.visit(`${orgPath}?org=${newOrgId}`).as('orgsite');
          });
      });
    });

    it('should rename the organization', () => {
      cy.get('[data-e2e="actions"]').click();
      cy.get('[data-e2e="rename"]', { timeout: 1000 }).should('be.visible').click();

      cy.get('[data-e2e="name"]').focus().clear().type(testOrgNameChange);
      cy.get('[data-e2e="dialog-submit"]').click();
      cy.shouldConfirmSuccess();
      cy.visit(orgPath);
      cy.get('[data-e2e="top-view-title"').should('contain', testOrgNameChange);
    });
  });

  it('should add an organization with the personal account as org owner');
  describe('changing the current organization', () => {
    it('should update displayed organization details');
  });
});
