import { Context } from 'support/commands';
import { ensureProjectDoesntExist, ensureProjectExists } from '../../support/api/projects';
import { ensureOrgExists } from 'support/api/orgs';

describe('projects', () => {
  beforeEach(() => {
    cy.context().as('ctx');
  });

  const defaultOrg = 'e2eorgnewdefault';
  const testProjectNameCreate = 'e2eprojectcreate';
  const testProjectNameDelete = 'e2eprojectdelete';

  describe('add project', () => {
    beforeEach(`ensure it doesn't exist already`, () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureOrgExists(ctx, defaultOrg).then(() => {
          ensureProjectDoesntExist(ctx.api, testProjectNameCreate);
          cy.visit(`/projects`);
        });
      });
    });

    it('should add a project', () => {
      cy.get('.add-project-button').click({ force: true });
      cy.get('input').type(testProjectNameCreate);
      cy.get('[data-e2e="continue-button"]').click();
      cy.shouldConfirmSuccess();
    });

    it('should configure a project to assert roles on authentication');
  });

  describe('create project grant', () => {
    const testRoleName = 'e2eroleundertestname';

    beforeEach('ensure it exists', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureProjectExists(ctx.api, testProjectNameCreate).as('projectId');
        cy.get<number>('@projectId').then((projectId) => {
          cy.visit(`/projects/${projectId}`);
        });
      });
    });

    it('should add a role', () => {
      cy.get('[data-e2e="sidenav-element-roles"]').click();
      cy.get('[data-e2e="add-new-role"]').click();
      cy.get('[formcontrolname="key"]').type(testRoleName);
      cy.get('[formcontrolname="displayName"]').type('e2eroleundertestdisplay');
      cy.get('[formcontrolname="group"]').type('e2eroleundertestgroup');
      cy.get('[data-e2e="save-button"]').click();
      cy.shouldConfirmSuccess();
      cy.contains('tr', testRoleName);
    });

    it('should add a project grant', () => {
      const rowSelector = `tr:contains(${testRoleName})`;

      cy.get('[data-e2e="sidenav-element-projectgrants"]').click();
      cy.get('[data-e2e="create-project-grant-button"]').click();
      cy.get('[data-e2e="add-org-input"]').type(defaultOrg);
      cy.get('mat-option').contains(defaultOrg).click();
      cy.get('button').should('be.enabled');
      cy.get('[data-e2e="project-grant-continue"]').first().click();
      cy.get(rowSelector).find('input').click({ force: true });
      cy.get('[data-e2e="save-project-grant-button"]').click();
      cy.contains('tr', defaultOrg);
      cy.contains('tr', testRoleName);
    });
  });

  describe('edit project', () => {
    beforeEach('ensure it exists', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureProjectExists(ctx.api, testProjectNameDelete).as('projectId');
        cy.visit(`/projects`);
      });
    });

    describe('remove project', () => {
      it('removes the project from list view', () => {
        const rowSelector = `tr:contains(${testProjectNameDelete})`;
        cy.get('[data-e2e="toggle-grid"]').click();
        cy.get('[data-e2e="timestamp"]');
        cy.get(rowSelector).find('[data-e2e="delete-project-button"]').click({ force: true });
        cy.get('[data-e2e="confirm-dialog-input"]').focus().type(testProjectNameDelete);
        cy.get('[data-e2e="confirm-dialog-button"]').click();
        cy.shouldConfirmSuccess();
        cy.shouldNotExist({
          selector: rowSelector,
          timeout: { ms: 2000, errMessage: 'timed out before project disappeared from the table' },
        });
      });

      it('removes the project from grid view', () => {
        const cardSelector = `[data-e2e="grid-card"]:contains(${testProjectNameDelete})`;
        cy.get(cardSelector).find('[data-e2e="delete-project-button"]').click({ force: true });
        cy.get('[data-e2e="confirm-dialog-input"]').focus().type(testProjectNameDelete);
        cy.get('[data-e2e="confirm-dialog-button"]').click();
        cy.shouldConfirmSuccess();
        cy.shouldNotExist({
          selector: cardSelector,
          timeout: { ms: 2000, errMessage: 'timed out before project disappeared from the grid' },
        });
      });
    });

    it('should add a project manager');
    it('should remove a project manager');
  });
});
