import { ensureProjectGrantExists } from 'support/api/grants';
import {
  ensureHumanIsOrgMember,
  ensureHumanIsNotOrgMember,
  ensureHumanIsNotProjectMember,
  ensureHumanIsProjectMember,
} from 'support/api/members';
import { ensureOrgExists } from 'support/api/orgs';
import { ensureHumanUserExists, ensureUserDoesntExist } from 'support/api/users';
import { Context } from 'support/commands';
import { loginname } from 'support/login/users';
import { apiAuth } from '../../support/api/apiauth';
import { ensureProjectExists, ensureProjectResourceDoesntExist, Roles } from '../../support/api/projects';

describe('permissions', () => {
  beforeEach(() => {
    cy.context().as('ctx');
  });

  describe('management', () => {
    const testManagerLoginname = loginname('e2ehumanmanager', Cypress.env('ORGANIZATION'));
    function testAuthorizations(
      roles: string[],
      beforeCreate: (ctx: Context) => void,
      beforeMutate: (ctx: Context) => void,
      navigate: () => void,
    ) {
      beforeEach(function () {
        cy.get<Context>('@ctx').then((ctx) => {
          ensureUserDoesntExist(ctx.api, testManagerLoginname);
          ensureHumanUserExists(ctx.api, testManagerLoginname);
        });
      });

      describe('create authorization', () => {
        beforeEach(function () {
          cy.get<Context>('@ctx').then((ctx) => {
            beforeCreate(ctx);
            navigate();
          });
        });

        it('should add a manager', () => {
          cy.get('[data-e2e="add-member-button"]').click();
          cy.get('[data-e2e="add-member-input"]').type(testManagerLoginname);
          cy.get('[data-e2e="user-option"]').click();
          cy.contains('[data-e2e="role-checkbox"]', roles[0]).click();
          cy.get('[data-e2e="confirm-add-member-button"]').click();
          cy.get('.data-e2e-success');
          cy.contains('[data-e2e="member-avatar"]', 'ee');
          cy.shouldNotExist({ selector: '.data-e2e-failure' });
        });
      });

      describe('mutate authorization', () => {
        const rowSelector = `tr:contains(${testManagerLoginname})`;

        beforeEach(() => {
          cy.get<Context>('@ctx').then((ctx) => {
            beforeMutate(ctx);
            navigate();
            cy.contains('[data-e2e="member-avatar"]', 'ee').click();
            cy.get(rowSelector).as('managerRow');
          });
        });

        it('should remove a manager', () => {
          cy.get('@managerRow').find('[data-e2e="remove-member-button"]').click({ force: true });
          cy.get('[data-e2e="confirm-dialog-button"]').click();
          cy.get('.data-e2e-success');
          cy.shouldNotExist({ selector: rowSelector, timeout: 2000 });
          cy.shouldNotExist({ selector: '.data-e2e-failure' });
        });

        it('should remove a managers authorization', () => {
          cy.get('@managerRow').find('[data-e2e="role"]').should('have.length', roles.length);
          cy.get('@managerRow')
            .contains('[data-e2e="role"]', roles[0])
            .find('[data-e2e="remove-role-button"]')
            .click({ force: true }); // TODO: Is this a bug?
          cy.get('[data-e2e="confirm-dialog-button"]').click();
          cy.get('.data-e2e-success');
          cy.get('@managerRow')
            .find('[data-e2e="remove-role-button"]')
            .should('have.length', roles.length - 1);
          cy.shouldNotExist({ selector: '.data-e2e-failure' });
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
        function (ctx: Context) {
          ensureHumanIsNotOrgMember(ctx.api, testManagerLoginname);
        },
        function (ctx: Context) {
          ensureHumanIsOrgMember(
            ctx.api,
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
        beforeEach(() => {
          cy.get<Context>('@ctx').then((ctx) => {
            ensureProjectExists(ctx.api, 'e2eprojectpermission').as('projectId');
          });
        });

        const visitOwnedProject = function () {
          cy.get<number>('@projectId').then((projectId) => {
            cy.visit(`/projects/${projectId}`);
          });
        };

        describe('authorizations', () => {
          const roles = [
            { internal: 'PROJECT_OWNER_GLOBAL', display: 'Project Owner Global' },
            { internal: 'PROJECT_OWNER_VIEWER_GLOBAL', display: 'Project Owner Viewer Global' },
          ];

          testAuthorizations(
            roles.map((role) => role.display),
            function (ctx) {
              cy.get<number>('@projectId').then((projectId) => {
                ensureHumanIsNotProjectMember(ctx.api, projectId, testManagerLoginname);
              });
            },
            function (ctx) {
              cy.get<number>('@projectId').then((projectId) => {
                ensureHumanIsProjectMember(
                  ctx.api,
                  projectId,
                  testManagerLoginname,
                  roles.map((role) => role.internal),
                );
              });
            },
            visitOwnedProject,
          );
        });

        describe('roles', () => {
          const testRoleName = 'e2eroleundertestname';

          beforeEach(function () {
            cy.get<Context>('@ctx').then((ctx) => {
              cy.get<number>('@projectId').then((projectId) => {
                ensureProjectResourceDoesntExist(ctx.api, projectId, Roles, testRoleName);
                visitOwnedProject();
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
            cy.get('.data-e2e-success');
            cy.contains('tr', testRoleName);
            cy.shouldNotExist({ selector: '.data-e2e-failure' });
          });
          it('should remove a role');
        });
      });

      describe('granted projects', () => {
        beforeEach(function () {
          cy.get<Context>('@ctx').then((ctx) => {
            ensureOrgExists(ctx.api, 'e2eforeignorg').then((foreignOrgId) => {
              ensureProjectExists(ctx.api, 'e2eprojectgrants', foreignOrgId)
                .as('foreignProjectId')
                .then((foreignProjectId) => {
                  ensureProjectGrantExists(ctx.api, foreignOrgId, foreignProjectId).as('grantId');
                });
            });
          });
        });

        function visitGrantedProject() {
          cy.get<number>('@foreignProjectId').then((foreignProjectId) => {
            cy.get<number>('@grantId').then((grantId) => {
              cy.visit(`/granted-projects/${foreignProjectId}/grant/${grantId}`);
            });
          });
        }

        describe('authorizations', () => {
          const roles = [
            { internal: 'PROJECT_GRANT_OWNER', display: 'Project Grant Owner' },
            { internal: 'PROJECT_GRANT_OWNER_VIEWER', display: 'Project Grant Owner Viewer' },
          ];

          testAuthorizations(
            roles.map((role) => role.display),
            function (ctx: Context) {
              cy.get<number>('@foreignProjectId').then((foreignProjectId) => {
                cy.get<number>('@grantId').then((grantId) => {
                  ensureHumanIsNotProjectMember(ctx.api, foreignProjectId, testManagerLoginname, grantId);
                });
              });
            },
            function (ctx: Context) {
              cy.get<number>('@foreignProjectId').then((foreignProjectId) => {
                cy.get<number>('@grantId').then((grantId) => {
                  ensureHumanIsProjectMember(
                    ctx.api,
                    foreignProjectId,
                    testManagerLoginname,
                    roles.map((role) => role.internal),
                    grantId,
                  );
                });
              });
            },
            visitGrantedProject,
          );
        });
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
