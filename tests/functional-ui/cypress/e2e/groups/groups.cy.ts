import { Context } from 'support/commands';
import { ensureGroupDoesntExist, ensureGroupExists } from '../../support/api/groups';

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
    beforeEach('ensure the group exists', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupExists(ctx.api, testGroupNameMembers);
        cy.visit('/groups');
      });
    });

    it('should show the members dialog', () => {
      cy.contains('tr', testGroupNameMembers).find('[data-e2e="group-members-button"]').click({ force: true });
      cy.get('cnsl-search-user-autocomplete');
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
