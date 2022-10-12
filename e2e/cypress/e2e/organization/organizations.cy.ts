import { renameOrg } from 'support/api/orgs';
import { apiAuth } from '../../support/api/apiauth';

const orgPath = `/org`;
const testOrgNameChange = 'e2erenametest';

describe('organizations', () => {
  describe('rename', () => {
    beforeEach(() => {
      apiAuth().as('api');
      cy.visit(orgPath);
    });

    afterEach(() => {
      apiAuth()
        .as('api')
        .then((api) => {
          renameOrg(api, 'ZITADEL');
        });
    });

    it.only('should rename the organization', () => {
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
