import { ensureOrgExists, ensureOrgIsDefault, isDefaultOrg } from 'support/api/orgs';
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

  const orgOverviewPath = `/orgs`;
  const initialDefaultOrg = 'e2eorgolddefault';
  const orgNameForNewDefault = 'e2eorgnewdefault';

  describe('set default org', () => {
    beforeEach(() => {
      apiAuth()
        .as('api')
        .then((api) => {
          ensureOrgExists(api, orgNameForNewDefault)
            .as('newDefaultOrgId')
            .then(() => {
              ensureOrgExists(api, initialDefaultOrg)
                .as('defaultOrg')
                .then((id) => {
                  ensureOrgIsDefault(api, id)
                    .as('orgWasDefault')
                    .then(() => {
                      cy.visit(`${orgOverviewPath}`).as('orgsite');
                    });
                });
            });
        });
    });

    it('should rename the organization', function () {
      const rowSelector = `tr:contains(${orgNameForNewDefault})`;
      cy.get(rowSelector).find('[data-e2e="table-actions-button"]').click({ force: true });
      cy.get('[data-e2e="set-default-button"]', { timeout: 1000 }).should('be.visible').click();
      cy.shouldConfirmSuccess();
      isDefaultOrg(this.api, this.newDefaultOrgId);
    });
  });

  it('should add an organization with the personal account as org owner');
  describe('changing the current organization', () => {
    it('should update displayed organization details');
  });
});
