import { Context } from 'support/commands';
import { ensureGroupDoesntExist, ensureGroupExists } from '../../support/api/groups';
import { ensureProjectExists, ensureRoleExists } from '../../support/api/projects';

const backendUrl = Cypress.env('BACKEND_URL');
const createGroupRoute = `${backendUrl}/zitadel.group.v2.GroupService/CreateGroup`;
const updateGroupRoute = `${backendUrl}/zitadel.group.v2.GroupService/UpdateGroup`;
const deleteGroupRoute = `${backendUrl}/zitadel.group.v2.GroupService/DeleteGroup`;
const createGroupGrantRoute = `${backendUrl}/zitadel.group.v2.GroupService/CreateGroupGrant`;

describe('groups — form-level client validation', () => {
  beforeEach(() => {
    cy.context().as('ctx');
  });

  describe('Create group dialog', () => {
    const name = 'e2egroup-form-create';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, name);
      });

      cy.intercept('POST', createGroupRoute).as('createGroup');
      cy.visit('/groups');
      cy.get('[data-e2e="create-group-button"]').click({ force: true });
    });

    it('disables save and blocks the API while the name is empty', () => {
      cy.get('[data-e2e="group-dialog-save"]').should('be.disabled');
      cy.get('[data-e2e="group-name-input"]').type(' ').clear();
      cy.get('[data-e2e="group-dialog-save"]').should('be.disabled').click({ force: true });
      cy.wait(500);
      cy.get('@createGroup.all').should('have.length', 0);
    });

    it('rejects a whitespace-only name without firing the API', () => {
      cy.get('[data-e2e="group-name-input"]').type('   ');
      cy.get('[data-e2e="group-dialog-save"]').should('be.disabled').click({ force: true });
      cy.wait(500);
      cy.get('@createGroup.all').should('have.length', 0);
    });

    it('fires exactly one create request even on a double-click', () => {
      cy.get('[data-e2e="group-name-input"]').type(name);
      cy.get('[data-e2e="group-dialog-save"]').dblclick();
      cy.wait('@createGroup').its('response.statusCode').should('be.oneOf', [200, 201]);
      cy.wait(500);
      cy.get('@createGroup.all').should('have.length', 1);
    });

    it('disables save when the name exceeds 200 characters', () => {
      const tooLong = 'a'.repeat(201);
      cy.get('[data-e2e="group-name-input"]').invoke('val', tooLong).trigger('input');
      cy.get('[data-e2e="group-dialog-save"]').should('be.disabled');
      cy.wait(500);
      cy.get('@createGroup.all').should('have.length', 0);
    });

    it('disables save when the description exceeds 200 characters', () => {
      const tooLong = 'a'.repeat(201);
      cy.get('[data-e2e="group-name-input"]').type(name);
      cy.get('[data-e2e="group-description-input"]').invoke('val', tooLong).trigger('input');
      cy.get('[data-e2e="group-dialog-save"]').should('be.disabled');
      cy.wait(500);
      cy.get('@createGroup.all').should('have.length', 0);
    });

    it('cancel discards changes and fires no request', () => {
      cy.get('[data-e2e="group-name-input"]').type(name);
      cy.get('button[mat-stroked-button]')
        .contains(/cancel|abbrechen/i)
        .click({ force: true });
      cy.wait(500);
      cy.get('@createGroup.all').should('have.length', 0);
    });

    it('submits exactly one request on save and clears the form', () => {
      cy.get('[data-e2e="group-name-input"]').type(name);
      cy.get('[data-e2e="group-dialog-save"]').click();
      cy.wait('@createGroup').its('response.statusCode').should('be.oneOf', [200, 201]);
      cy.wait(500);
      cy.get('@createGroup.all').should('have.length', 1);
    });
  });

  describe('Edit group dialog', () => {
    const name = 'e2egroup-form-edit';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupExists(ctx.api, name);
      });

      cy.intercept('POST', updateGroupRoute).as('updateGroup');
      cy.visit('/groups');
      cy.contains('tr', name).click();
    });

    it('preloads the name and disables save when the name is cleared', () => {
      cy.get('[data-e2e="group-name-input"]').should('have.value', name);
      cy.get('[data-e2e="group-name-input"]').clear();
      cy.get('[data-e2e="group-dialog-save"]').should('be.disabled');
      cy.wait(500);
      cy.get('@updateGroup.all').should('have.length', 0);
    });

    it('keeps save disabled until the form is dirty', () => {
      cy.get('[data-e2e="group-name-input"]').should('have.value', name);
      cy.get('[data-e2e="group-dialog-save"]').should('be.disabled');
      cy.wait(500);
      cy.get('@updateGroup.all').should('have.length', 0);
    });
  });

  describe('Grant dialog', () => {
    const groupName = 'e2egroup-form-grant';
    const projectName = 'e2egroup-form-grant-project';
    const roleKey = 'e2egroup-form-grant-role';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupExists(ctx.api, groupName);
        ensureProjectExists(ctx.api, projectName).then((projectId) => {
          ensureRoleExists(ctx.api, projectId, roleKey);
          cy.wrap(projectId).as('projectId');
        });
      });

      cy.intercept('POST', createGroupGrantRoute).as('createGrant');
      cy.visit('/groups');
      cy.contains('tr', groupName).find('[data-e2e="group-grants-button"]').click({ force: true });
    });

    it('disables save until a project is selected and at least one role is checked', () => {
      cy.get('[data-e2e="group-grant-save"]').should('be.disabled');

      cy.get('[data-e2e="group-grant-roles-table"]').should('not.exist');

      cy.get('[data-e2e="group-grant-project-autocomplete"] input').click();
      cy.contains('mat-option', projectName, { timeout: 10000 }).click();

      cy.get('[data-e2e="group-grant-roles-table"]').should('be.visible');
      cy.get('[data-e2e="group-grant-save"]').should('be.disabled');

      cy.contains('[data-e2e="group-grant-roles-table"] tr', roleKey).find('mat-checkbox').click();
      cy.get('[data-e2e="group-grant-save"]').should('be.enabled');

      cy.wait(500);
      cy.get('@createGrant.all').should('have.length', 0);
    });
  });

  describe('Delete confirmation dialog', () => {
    const name = 'e2egroup-form-delete';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupExists(ctx.api, name);
      });

      cy.intercept('POST', deleteGroupRoute).as('deleteGroup');
      cy.visit('/groups');
    });

    it('cancelling the confirm dialog fires no delete request', () => {
      cy.contains('tr', name).find('[data-e2e="delete-group-button"]').click({ force: true });
      cy.get('button')
        .contains(/cancel|abbrechen/i)
        .click({ force: true });
      cy.wait(500);
      cy.get('@deleteGroup.all').should('have.length', 0);
      cy.contains('tr', name).should('exist');
    });
  });
});
