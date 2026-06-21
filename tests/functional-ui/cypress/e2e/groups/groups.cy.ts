import { Context } from 'support/commands';
import { ensureGroupDoesntExist, ensureGroupExists } from '../../support/api/groups';
import { ensureProjectExists, ensureRoleExists } from '../../support/api/projects';
import { ensureHumanUserExists } from '../../support/api/users';

describe('groups', () => {
  beforeEach(() => {
    cy.context().as('ctx');
  });

  const testGroupNameCreate = 'e2egroupcreate';
  const testGroupNameRename = 'e2egrouprename';
  const testGroupNameRenamed = 'e2egrouprenamed';
  const testGroupNameMembers = 'e2egroupmembers';
  const testGroupNameDelete = 'e2egroupdelete';

  describe('add group', () => {
    beforeEach(`ensure it doesn't exist already`, () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, testGroupNameCreate);
        cy.visit('/groups');
      });
    });

    it('should show the navigation on a direct load', () => {
      cy.get('#mainnav').should('be.visible');
    });

    it('should add a group', () => {
      cy.get('[data-e2e="create-group-button"]').click({ force: true });
      cy.get('[data-e2e="group-name-input"]').should('be.enabled').type(testGroupNameCreate);
      cy.get('[data-e2e="group-description-input"]').should('be.enabled').type('e2e group description');
      cy.get('[data-e2e="group-dialog-save"]').click();
      cy.shouldConfirmSuccess();
      cy.contains('tr', testGroupNameCreate);
    });
  });

  describe('edit group', () => {
    beforeEach('ensure it exists', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, testGroupNameRenamed);
        ensureGroupExists(ctx.api, testGroupNameRename);
        cy.visit('/groups');
      });
    });

    it('should rename the group', () => {
      cy.contains('tr', testGroupNameRename).click();
      cy.get('[data-e2e="group-name-input"]').should('be.enabled').clear().type(testGroupNameRenamed);
      cy.get('[data-e2e="group-dialog-save"]').click();
      cy.shouldConfirmSuccess();
      cy.contains('tr', testGroupNameRenamed);
    });
  });

  describe('manage members', () => {
    const testMemberUsername = 'e2egroupmemberuser';

    beforeEach('ensure the group and the user exist', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, testGroupNameMembers);
        ensureGroupExists(ctx.api, testGroupNameMembers);
        ensureHumanUserExists(ctx.api, testMemberUsername);
        cy.visit('/groups');
      });
    });

    it('should show the members dialog', () => {
      cy.contains('tr', testGroupNameMembers).find('[data-e2e="group-members-button"]').click({ force: true });
      cy.get('cnsl-search-user-autocomplete');
    });

    it('should add a member through the autocomplete and update the user count', () => {
      cy.contains('tr', testGroupNameMembers).find('[data-e2e="group-members-button"]').click({ force: true });
      cy.get('[data-e2e="add-member-input"]').type(testMemberUsername);
      cy.get('[data-e2e="user-option"]').first().click();
      cy.get('[data-e2e="group-members-add"]').click();
      cy.shouldConfirmSuccess();
      cy.contains('.member-row', testMemberUsername);

      cy.get('[mat-dialog-actions] button').first().click();
      cy.contains('tr', testGroupNameMembers).should('contain', '1');
    });

    it('should remove a member', () => {
      cy.contains('tr', testGroupNameMembers).find('[data-e2e="group-members-button"]').click({ force: true });
      cy.get('[data-e2e="add-member-input"]').type(testMemberUsername);
      cy.get('[data-e2e="user-option"]').first().click();
      cy.get('[data-e2e="group-members-add"]').click();
      cy.shouldConfirmSuccess();
      cy.contains('.member-row', testMemberUsername).find('[data-e2e="group-member-remove"]').click();
      cy.shouldConfirmSuccess();
      cy.contains('.member-row', testMemberUsername).should('not.exist');
    });
  });

  describe('manage grants', () => {
    const testGroupNameGrants = 'e2egroupgrants';
    const testProjectName = 'e2egroupgrantproject';
    const testRoleKey = 'e2egroupgrantrole';

    beforeEach('ensure the group, project, and role exist', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, testGroupNameGrants);
        ensureGroupExists(ctx.api, testGroupNameGrants);
        ensureProjectExists(ctx.api, testProjectName).then((projectId) => {
          ensureRoleExists(ctx.api, projectId, testRoleKey);
          cy.wrap(projectId).as('projectId');
        });
        cy.visit('/groups');
      });
    });

    it('should create and revoke a group grant', () => {
      cy.get<number>('@projectId').then((projectId) => {
        cy.contains('tr', testGroupNameGrants).find('[data-e2e="group-grants-button"]').click({ force: true });

        cy.get('[data-e2e="group-grant-project-autocomplete"] input').click();
        cy.contains('mat-option', testProjectName, { timeout: 10000 }).click();
        cy.contains('[data-e2e="group-grant-roles-table"] tr', testRoleKey)
          .find('mat-checkbox')
          .click();
        cy.get('[data-e2e="group-grant-save"]').click();
        cy.shouldConfirmSuccess();
        cy.contains('.grant-row', `${projectId}`).should('contain', testRoleKey);

        cy.contains('.grant-row', `${projectId}`).find('[data-e2e="group-grant-remove"]').click();
        cy.shouldConfirmSuccess();
        cy.contains('.grant-row', `${projectId}`).should('not.exist');
      });
    });
  });

  describe('remove group', () => {
    beforeEach('ensure it exists', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupExists(ctx.api, testGroupNameDelete);
        cy.visit('/groups');
      });
    });

    it('should delete the group', () => {
      cy.contains('tr', testGroupNameDelete).find('[data-e2e="delete-group-button"]').click({ force: true });
      cy.get('[data-e2e="confirm-dialog-button"]').click();
      cy.shouldConfirmSuccess();
      cy.contains('tr', testGroupNameDelete).should('not.exist');
    });
  });
});
