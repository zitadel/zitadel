import { ensureOrgExists, ensureOrgIsDefault, isDefaultOrg } from 'support/api/orgs';
import { v4 as uuidv4 } from 'uuid';
import { Context } from 'support/commands';

const orgPath = `/org`;
const orgsPath = `/orgs`;
const orgsPathCreate = `/orgs/create`;

const orgNameOnCreation = 'e2eorgrename';
const testOrgNameChange = uuidv4();
const newOrg = uuidv4();

beforeEach(() => {
  cy.context().as('ctx');
});

describe('organizations', () => {
  describe('add and delete org', () => {
    it('should create an org', () => {
      cy.visit(orgsPathCreate);
      cy.get('[data-e2e="org-name-input"]').focus().clear().type(newOrg);
      cy.get('[data-e2e="create-org-button"]').click();
      cy.contains('tr', newOrg);
    });

    it('should delete an org', () => {
      cy.visit(orgsPath);
      cy.contains('tr', newOrg).click();
      cy.get('[data-e2e="actions"]').click();
      cy.get('[data-e2e="delete"]', { timeout: 3000 }).should('be.visible').click();
      cy.get('[data-e2e="confirm-dialog-input"]').focus().clear().type(newOrg);
      cy.get('[data-e2e="confirm-dialog-button"]').click();
      cy.shouldConfirmSuccess();
      cy.contains('tr', newOrg).should('not.exist');
    });
  });

  describe('rename', () => {
    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureOrgExists(ctx, orgNameOnCreation).then((newOrgId) => {
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

    const orgOverviewPath = `/orgs`;
    const initialDefaultOrg = 'e2eorgolddefault';
    const orgNameForNewDefault = 'e2eorgnewdefault';

    describe('set default org', () => {
      beforeEach(() => {
        cy.get<Context>('@ctx').then((ctx) => {
          ensureOrgExists(ctx, orgNameForNewDefault)
            .as('newDefaultOrgId')
            .then(() => {
              ensureOrgExists(ctx, initialDefaultOrg)
                .as('defaultOrg')
                .then((id) => {
                  ensureOrgIsDefault(ctx, id)
                    .as('orgWasDefault')
                    .then(() => {
                      cy.visit(`${orgOverviewPath}`).as('orgsite');
                    });
                });
            });
        });
      });

      it('should rename the organization', () => {
        cy.get<Context>('@ctx').then((ctx) => {
          const rowSelector = `tr:contains(${orgNameForNewDefault})`;
          cy.get(rowSelector).find('[data-e2e="table-actions-button"]').click({ force: true });
          cy.get('[data-e2e="set-default-button"]', { timeout: 1000 }).should('be.visible').click();
          cy.shouldConfirmSuccess();
          cy.get<string>('@newDefaultOrgId').then((newDefaultOrgId) => {
            isDefaultOrg(ctx, newDefaultOrgId);
          });
        });
      });
    });

    it('should add an organization with the personal account as org owner');
    describe('changing the current organization', () => {
      it('should update displayed organization details');
    });
  });
});
