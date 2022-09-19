import { ensureHumanIsMember, ensureHumanIsNotMember } from 'support/api/members';
import { ensureHumanUserExists } from 'support/api/users';
import { apiAuth } from '../../support/api/apiauth';
import { ensureProjectExists, ensureProjectResourceDoesntExist, Roles } from '../../support/api/projects';

describe('permissions', () => {
  const testProjectName = 'e2eprojectpermission';
  const testAppName = 'e2eapppermission';
  const testRoleName = 'e2eroleundertestname';
  const testRoleDisplay = 'e2eroleundertestdisplay';
  const testRoleGroup = 'e2eroleundertestgroup';
  const testGrantName = 'e2egrantundertest';

  beforeEach(() => {
    apiAuth().as('api');
  });

  describe('management', () => {
    describe('organizations', () => {
      const testManagerName = 'e2ehumanmanager';

      beforeEach(function () {
        ensureHumanUserExists(this.api, testManagerName);
      });

      describe('add manager', () => {
        beforeEach(function () {
          ensureHumanIsNotMember(this.api, testManagerName);
          cy.visit('/orgs');
        });

        it('should add an organization manager', () => {
          cy.contains('tr', Cypress.env('ORGANIZATION')).click();
          cy.get('[data-e2e="add-member-button"]').click();
          cy.get('[data-e2e="add-member-input"]').type(testManagerName);
          cy.get('[data-e2e="user-option"]').click();
          cy.contains('[data-e2e="role-checkbox"]', 'Org Owner').click();
          cy.get('[data-e2e="confirm-add-member-button"]').click();
          cy.get('.data-e2e-success');
          cy.contains('[data-e2e="member-avatar"]', 'ee');
          cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist');
        });
      });

      describe('edit authorizations', () => {
        beforeEach(function () {
          ensureHumanIsMember(this.api, testManagerName, ['ORG_OWNER', 'ORG_OWNER_VIEWER']);
          cy.visit('/orgs');
          cy.contains('tr', Cypress.env('ORGANIZATION')).click();
          cy.contains('[data-e2e="member-avatar"]', 'ee').click();
          cy.contains('tr', testManagerName).as('managerRow');
        });

        it('should remove a manager', () => {
          cy.get('@managerRow').find('[data-e2e="remove-member-button"]').click({ force: true });
          cy.get('[data-e2e="confirm-dialog-button"]').click();
          cy.get('.data-e2e-success');
          // https://github.com/NoriSte/cypress-wait-until/issues/75#issuecomment-572685623
          cy.waitUntil(() => Cypress.$(`tr:contains('${testManagerName}')`).length === 0);
          cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist');
        });

        describe('roles', () => {
          it('should remove a managers authorization', () => {
            cy.get('@managerRow').find('[data-e2e="role"]').should('have.length', 2);
            cy.get('@managerRow')
              .contains('[data-e2e="role"]', 'Org Owner Viewer')
              .find('[data-e2e="remove-role-button"]')
              .click();
            cy.get('[data-e2e="confirm-dialog-button"]').click();
            cy.get('.data-e2e-success');
            cy.get('@managerRow').find('[data-e2e="remove-role-button"]').should('have.length', 1);
            cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist');
          });
        });
      });
    });

    describe('projects', () => {
      const testProjectName = 'e2eprojectpermission';
      const testRoleName = 'e2eroleundertestname';
      const testRoleDisplay = 'e2eroleundertestdisplay';
      const testRoleGroup = 'e2eroleundertestgroup';

      beforeEach(function () {
        ensureProjectExists(this.api, testProjectName).as('projectId');
      });

      describe('managers', () => {
        it('should add a project manager');
        it('should remove a project manager');
      });

      describe('authorizations', () => {
        it('should add an authorization');
        it('should remove an authorization');
      });

      describe('owned projects', () => {
        describe('roles', () => {
          beforeEach(function () {
            ensureProjectResourceDoesntExist(this.api, this.projectId, Roles, testRoleName);
            cy.visit(`/projects/${this.projectId}?id=roles`);
          });

          it('should add a role', () => {
            cy.get('[data-e2e="add-new-role"]').click();
            cy.get('[formcontrolname="key"]').type(testRoleName);
            cy.get('[formcontrolname="displayName"]').type(testRoleDisplay);
            cy.get('[formcontrolname="group"]').type(testRoleGroup);
            cy.get('[data-e2e="save-button"]').click();
            cy.get('.data-e2e-success');
            cy.contains('tr', testRoleName);
            cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist');
          });
          it('should remove a role');
        });

        describe('grants', () => {
          it('should add a grant');
          it('should remove a grant');
        });
      });
    });
  });

  describe('validations', () => {
    describe('owned projects', () => {
      describe('no ownership', () => {
        it('a user without project global ownership can ...');
        it('a user without project global ownership can not ...');
      });
      describe('project owner viewer global', () => {
        it('a project owner viewer global additionally can ...');
        it('a project owner viewer global still can not ...');
      });
      describe('project owner global', () => {
        it('a project owner global additionally can ...');
        it('a project owner global still can not ...');
      });
    });

    describe('granted projects', () => {
      describe('no ownership', () => {
        it('a user without project grant ownership can ...');
        it('a user without project grant ownership can not ...');
      });
      describe('project grant owner viewer', () => {
        it('a project grant owner viewer additionally can ...');
        it('a project grant owner viewer still can not ...');
      });
      describe('project grant owner', () => {
        it('a project grant owner additionally can ...');
        it('a project grant owner still can not ...');
      });
    });
    describe('organization', () => {
      describe('org owner', () => {
        it('a project owner global can ...');
        it('a project owner global can not ...');
      });
    });
  });
});
