import { ensureOrgExists } from 'support/api/orgs';
import { newTarget } from 'support/api/target';
import { ZITADELTarget } from 'support/commands';
import { v4 as uuidv4 } from 'uuid';

const orgPath = `/org`;

const orgNameOnCreation = 'e2eorgrename';
const testOrgNameChange = uuidv4();

describe('organizations', () => {
  beforeEach(() => {
    newTarget('e2eorgs').as('target');
  });

  describe('routing', () => {
    // TODO: Fix console bug
    it.skip('routing works', () => {
      cy.get<ZITADELTarget>('@target').then((target) => {
        cy.visit(`/users?type=human&org=${target.headers['x-zitadel-orgid']}`);
        cy.contains('cnsl-nav', 'Users');
        cy.get('tr:contains(ZITADEL Admin)', { timeout: 0 }).should('not.exist');
      });
    });
  });

  describe('rename', () => {
    beforeEach(() => {
      cy.get<ZITADELTarget>('@target').then((target) => {
        ensureOrgExists(target, orgNameOnCreation)
          .as('newOrgId')
          .then((newOrgId) => {
            cy.visit(`${orgPath}?org=${newOrgId}`).as('orgsite');
          });
      });
    });

    it('should rename the organization', () => {
      cy.get('[data-e2e="actions"]').should('be.visible').click();
      cy.get('[data-e2e="rename"]', { timeout: 1000 }).should('be.visible').click();

      cy.get('[data-e2e="name"]').focus().clear().should('be.visible').type(testOrgNameChange);
      cy.get('[data-e2e="dialog-submit"]').should('be.visible').click();
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
