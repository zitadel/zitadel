import { Context } from 'support/commands';
import { ensureProjectDoesntExist, ensureProjectExists } from '../../support/api/projects';

describe('projects', () => {
  beforeEach(() => {
    cy.context().as('ctx');
  });

  const testProjectNameCreate = 'e2eprojectcreate';
  const testProjectNameDelete = 'e2eprojectdelete';

  describe('add project', () => {
    beforeEach(`ensure it doesn't exist already`, () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureProjectDoesntExist(ctx.api, testProjectNameCreate);
        cy.visit(`/projects`);
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

  describe('edit project', () => {
    beforeEach('ensure it exists', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureProjectExists(ctx.api, testProjectNameDelete);
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
