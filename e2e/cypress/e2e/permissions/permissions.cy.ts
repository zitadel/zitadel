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

// TODO: .as is asynchronous. Stop using Mocha.Context with `this.grantId and this.projectId`
describe('permissions', () => {
  beforeEach(() => {
    cy.context().as('ctx');
    cy.get<Context>('@ctx').then(ctx=> {
      ensureProjectExists(ctx.api, 'e2eprojectpermission').as('projectId')
    })
  });

  describe('management', () => {
    const testManagerLoginname = loginname('e2ehumanmanager', Cypress.env('ORGANIZATION'));
    function testAuthorizations(
      roles: string[],
      beforeCreate: (ctx: Context, projectId: number) => void,
      beforeMutate: (ctx: Context, projectId: number) => void,
      navigate: (projectId: number) => void,
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
            cy.get<number>('@projectId').then(projectId=>{
              beforeCreate(ctx, projectId);
              navigate(projectId);
            });
          });
        });

        it('should add a manager', () => {
          cy.get('[data-e2e="add-member-button"]').click();
          cy.get('[data-e2e="add-member-input"]').type(testManagerLoginname);
          cy.get('[data-e2e="user-option"]').click();
          cy.contains('[data-e2e="role-checkbox"]', roles[0]).click();
          cy.get('[data-e2e="confirm-add-member-button"]').click();
          cy.get('.data-e2e-success');
          cy.get('[data-e2e="member-avatar"]');
          cy.shouldNotExist({ selector: '.data-e2e-failure' });
        });
      });

      describe('mutate authorization', () => {
        const rowSelector = `tr:contains(${testManagerLoginname})`;

        beforeEach(() => {
          cy.get<Context>('@ctx').then((ctx) => {
            beforeMutate(ctx, this);
            navigate(this);
          });
        });

        beforeEach(() => {
          debugger
          cy.get('[data-e2e="member-avatar"]').click();
          cy.get(rowSelector).as('managerRow');
        });

        it.only('should remove a manager', () => {
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
          return () => ensureHumanIsNotOrgMember(ctx.api, testManagerLoginname);
        },
        function (ctx: Context) {
          return () => {
            ensureHumanIsNotOrgMember(ctx.api, testManagerLoginname);
            ensureHumanIsOrgMember(
              ctx.api,
              testManagerLoginname,
              roles.map((role) => role.internal),
            );
          };
        },
        () => {
          cy.visit('/orgs');
          cy.contains('tr', Cypress.env('ORGANIZATION')).click();
        },
      );
    });

    describe('projects', () => {
      describe('owned projects', () => {

        const visitOwnedProject = function (projectId: number) {
          cy.visit(`/projects/${projectId}`);
        };

        describe.only('authorizations', () => {
          const roles = [
            { internal: 'PROJECT_OWNER_GLOBAL', display: 'Project Owner Global' },
            { internal: 'PROJECT_OWNER_VIEWER_GLOBAL', display: 'Project Owner Viewer Global' },
          ];

          testAuthorizations(
            roles.map((role) => role.display),
            function (ctx, projectId) {
              return () => ensureHumanIsNotProjectMember(ctx.api, projectId, testManagerLoginname);
            },
            function (ctx, projectId) {
              return () => {
                ensureHumanIsNotProjectMember(ctx.api, projectId, testManagerLoginname);
                ensureHumanIsProjectMember(
                  ctx.api,
                  projectId,
                  testManagerLoginname,
                  roles.map((role) => role.internal),
                );
              };
            },
            visitOwnedProject,
          );
        });

        describe('roles', () => {
          const testRoleName = 'e2eroleundertestname';

          beforeEach(function () {
            cy.get<Context>('@ctx').then((ctx) => {
              ensureProjectResourceDoesntExist(ctx.api, this.projectId, Roles, testRoleName);
            });
            cy.get<number>('@projectId').then(projectId=>{
              visitOwnedProject(projectId);
            })
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
                .then((projectId) => {
                  ensureProjectGrantExists(ctx.api, foreignOrgId, projectId).as('grantId');
                });
            });
          });
        });

        function visitGrantedProject(projectId: number) {
          cy.get('@grantId').then((grantId)=>{
            cy.visit(`/granted-projects/${projectId}/grant/${grantId}`);

          })
        }

        describe('authorizations', () => {
          const roles = [
            { internal: 'PROJECT_GRANT_OWNER', display: 'Project Grant Owner' },
            { internal: 'PROJECT_GRANT_OWNER_VIEWER', display: 'Project Grant Owner Viewer' },
          ];

          testAuthorizations(
            roles.map((role) => role.display),
            function (ctx: Context, projectId: number) {
              return () => {
                cy.get<number>('@grantId').then((grantId)=>{
                  ensureHumanIsNotProjectMember(ctx.api, projectId, testManagerLoginname, grantId);
                })
              }
            },
            function (ctx: Context, projectId: number) {
              return () => {
                cy.get<number>('@grantId').then((grantId)=>{
               ensureHumanIsNotProjectMember(ctx.api, projectId, testManagerLoginname, grantId);
                ensureHumanIsProjectMember(
                  ctx.api,
                  projectId,
                  testManagerLoginname,
                  roles.map((role) => role.internal),
                  grantId,
                );
                })
              };
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
