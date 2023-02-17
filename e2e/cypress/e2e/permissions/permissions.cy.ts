import { ensureProjectGrantExists } from 'support/api/grants';
import {
  ensureHumanIsOrgMember,
  ensureHumanIsNotOrgMember,
  ensureHumanIsNotProjectMember,
  ensureHumanIsProjectMember,
} from 'support/api/members';
import { ensureOrgExists } from 'support/api/orgs';
import { ensureHumanUserExists, ensureUserDoesntExist } from 'support/api/users';
import { loginname } from 'support/login/users';
import { apiAuth } from '../../support/api/apiauth';
import { ensureProjectExists, ensureProjectResourceDoesntExist, Roles } from '../../support/api/projects';

describe('permissions', () => {
  beforeEach(() => {
    apiAuth().as('api');
  });

  describe('management', () => {
    const testManagerLoginname = loginname('e2ehumanmanager', Cypress.env('ORGANIZATION'));
    function testAuthorizations(
      roles: string[],
      beforeCreate: Mocha.HookFunction,
      beforeMutate: Mocha.HookFunction,
      navigate: Mocha.HookFunction,
    ) {
      beforeEach(function () {
        ensureUserDoesntExist(this.api, testManagerLoginname);
        ensureHumanUserExists(this.api, testManagerLoginname);
      });

      describe('create authorization', () => {
        beforeEach(beforeCreate);
        beforeEach(navigate);

        it('should add a manager', () => {
          cy.get('[data-e2e="add-member-button"]').click();
          cy.get('[data-e2e="add-member-input"]').type(testManagerLoginname);
          cy.get('[data-e2e="user-option"]').first().click();
          cy.contains('[data-e2e="role-checkbox"]', roles[0]).click();
          cy.get('[data-e2e="confirm-add-member-button"]').click();
          cy.shouldConfirmSuccess();
          cy.contains('[data-e2e="member-avatar"]', 'ee');
        });
      });

      describe('mutate authorization', () => {
        const rowSelector = `tr:contains(${testManagerLoginname})`;

        beforeEach(beforeMutate);
        beforeEach(navigate);

        beforeEach(() => {
          cy.contains('[data-e2e="member-avatar"]', 'ee').click();
          cy.get(rowSelector).as('managerRow');
        });

        it('should remove a manager', () => {
          cy.get('@managerRow').find('[data-e2e="remove-member-button"]').click({ force: true });
          cy.get('[data-e2e="confirm-dialog-button"]').click();
          cy.shouldConfirmSuccess();
          cy.shouldNotExist({
            selector: rowSelector,
            timeout: { ms: 2000, errMessage: 'timed out before manager disappeared from the table' },
          });
        });

        it('should remove a managers authorization', () => {
          cy.get('@managerRow').find('[data-e2e="role"]').should('have.length', roles.length);
          cy.get('@managerRow')
            .contains('[data-e2e="role"]', roles[0])
            .find('[data-e2e="remove-role-button"]')
            .click({ force: true }); // TODO: Is this a bug?
          cy.get('[data-e2e="confirm-dialog-button"]').click();
          cy.shouldConfirmSuccess();
          cy.get('@managerRow')
            .find('[data-e2e="remove-role-button"]')
            .should('have.length', roles.length - 1);
        });
      });
    }

    describe('organizations', () => {
      const roles = [
        { internal: 'ORG_OWNER', display: 'Org Owner' },
        { internal: 'ORG_OWNER_VIEWER', display: 'Org Owner Viewer' },
      ];

      testAuthorizations(
        roles.map((role) => role.display),
        function () {
          ensureHumanIsNotOrgMember(this.api, testManagerLoginname);
        },
        function () {
          ensureHumanIsNotOrgMember(this.api, testManagerLoginname);
          ensureHumanIsOrgMember(
            this.api,
            testManagerLoginname,
            roles.map((role) => role.internal),
          );
        },
        () => {
          cy.visit('/orgs');
          cy.contains('tr', Cypress.env('ORGANIZATION')).click();
        },
      );
    });

    describe('projects', () => {
      describe('owned projects', () => {
        beforeEach(function () {
          ensureProjectExists(this.api, 'e2eprojectpermission').as('projectId');
        });

        const visitOwnedProject: Mocha.HookFunction = function () {
          cy.visit(`/projects/${this.projectId}`);
        };

        describe('authorizations', () => {
          const roles = [
            { internal: 'PROJECT_OWNER_GLOBAL', display: 'Project Owner Global' },
            { internal: 'PROJECT_OWNER_VIEWER_GLOBAL', display: 'Project Owner Viewer Global' },
          ];

          testAuthorizations(
            roles.map((role) => role.display),
            function () {
              ensureHumanIsNotProjectMember(this.api, this.projectId, testManagerLoginname);
            },
            function () {
              ensureHumanIsNotProjectMember(this.api, this.projectId, testManagerLoginname);
              ensureHumanIsProjectMember(
                this.api,
                this.projectId,
                testManagerLoginname,
                roles.map((role) => role.internal),
              );
            },
            visitOwnedProject,
          );
        });

        describe('roles', () => {
          const testRoleName = 'e2eroleundertestname';

          beforeEach(function () {
            ensureProjectResourceDoesntExist(this.api, this.projectId, Roles, testRoleName);
          });

          beforeEach(visitOwnedProject);

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
          it('should remove a role');
        });
      });

      describe('granted projects', () => {
        beforeEach(function () {
          ensureOrgExists(this.api, 'e2eforeignorg')
            .as('foreignOrgId')
            .then((foreignOrgId) => {
              ensureProjectExists(this.api, 'e2eprojectgrants', foreignOrgId)
                .as('projectId')
                .then((projectId) => {
                  ensureProjectGrantExists(this.api, foreignOrgId, projectId).as('grantId');
                });
            });
        });

        const visitGrantedProject: Mocha.HookFunction = function () {
          cy.visit(`/granted-projects/${this.projectId}/grant/${this.grantId}`);
        };

        describe('authorizations', () => {
          const roles = [
            { internal: 'PROJECT_GRANT_OWNER', display: 'Project Grant Owner' },
            { internal: 'PROJECT_GRANT_OWNER_VIEWER', display: 'Project Grant Owner Viewer' },
          ];

          testAuthorizations(
            roles.map((role) => role.display),
            function () {
              ensureHumanIsNotProjectMember(this.api, this.projectId, testManagerLoginname, this.grantId);
            },
            function () {
              ensureHumanIsNotProjectMember(this.api, this.projectId, testManagerLoginname, this.grantId);
              ensureHumanIsProjectMember(
                this.api,
                this.projectId,
                testManagerLoginname,
                roles.map((role) => role.internal),
                this.grantId,
              );
            },
            visitGrantedProject,
          );
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
