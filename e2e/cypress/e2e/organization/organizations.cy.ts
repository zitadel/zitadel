import { ensureOrgExists } from 'support/api/orgs';
import { v4 as uuidv4 } from 'uuid';
import { Context } from 'support/commands';

const orgPath = `/org`;

const orgNameOnCreation = 'e2eorgrename';
const testOrgNameChange = uuidv4();

beforeEach(() => {
  cy.context().as('ctx');
});

describe('organizations', () => {
  describe('rename', () => {
    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureOrgExists(ctx.api, orgNameOnCreation).then((newOrgId) => {
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

    it('should add an organization with the personal account as org owner');
    describe('changing the current organization', () => {
      it('should update displayed organization details');
    });
  });
});
