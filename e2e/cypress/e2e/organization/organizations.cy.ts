import { ensureOrgExists, renameOrg } from 'support/api/orgs';
import { apiAuth } from '../../support/api/apiauth';

const orgListPath = `/orgs`;
const orgPath = `/org`;

const orgNameOnCreation = '1';
const testOrgNameChange = '2';

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

    afterEach(() => {
      this.orgsite.then(() => {
        apiAuth()
          .as('api')
          .then((api) => {
            renameOrg(api, orgNameOnCreation);
          });
      });
    });

    it('should rename the organization', () => {
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
