import { ensureProjectGrantExists } from 'support/api/grants';
import { ensureHumanIsNotOrgMember, ensureHumanIsNotProjectMember, ensureHumanIsOrgMember, ensureHumanIsProjectMember } from 'support/api/members';
import { ensureOrgExists } from 'support/api/orgs';
import { ensureProjectExists } from 'support/api/projects';
import { ensureRoleDoesntExist } from 'support/api/roles';
import { newTarget } from 'support/api/target';
import { ensureHumanDoesntExist, ensureHumanExists } from 'support/api/users';
import { ZITADELTarget } from 'support/commands';
import { loginname } from 'support/login/login';

describe('permissions', () => {
  const targetOrg = 'e2epermissionsmyorg';

  beforeEach(() => {
    newTarget(targetOrg).as('target');
  });

  describe('management', () => {
    const testManagerLoginname = loginname('e2ehumanmanager', targetOrg);
    function testAuthorizations(
      roles: string[],
      beforeCreate: (target: ZITADELTarget, userId: number) => void,
      beforeMutate: (target: ZITADELTarget, userId: number) => void,
      navigate: (target: ZITADELTarget, userId: number) => void,
    ) {
      beforeEach(() => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          ensureHumanDoesntExist(target, testManagerLoginname);
          ensureHumanExists(target, testManagerLoginname).as('userId');
        });
      });

      describe('create authorization', () => {
        beforeEach(() => {
          cy.get<ZITADELTarget>('@target').then((target) => {
            cy.get<number>('@userId').then(userId => {
              beforeCreate(target, userId);
              navigate(target, userId);
            })
          });
        });

        it('should add a manager', () => {
          cy.get('[data-e2e="add-member-button"]').should("be.visible").click();
          cy.get('[data-e2e="add-member-input"]').should("be.visible").type(testManagerLoginname);
          cy.get('[data-e2e="user-option"]').should("be.visible").click();
          cy.contains('[data-e2e="role-checkbox"]', roles[0]).should("be.visible").click();
          cy.get('[data-e2e="confirm-add-member-button"]').should("be.visible").click();
          cy.shouldConfirmSuccess();
          cy.contains('[data-e2e="member-avatar"]', 'ee');
        });
      });

      describe('mutate authorization', () => {
        const rowSelector = `tr:contains(${testManagerLoginname})`;

        beforeEach(() => {
          cy.get<ZITADELTarget>('@target').then((target) => {
            cy.get<number>('@userId').then(userId => {
              beforeMutate(target,userId);
            navigate(target, userId);
          });
          cy.contains('[data-e2e="member-avatar"]', 'ee').should("be.visible").click();
            cy.get(rowSelector).as('managerRow');
          });
        });

        it('should remove a manager', () => {
          cy.get('@managerRow').find('[data-e2e="remove-member-button"]').should("be.visible").click({ force: true });
          cy.get('[data-e2e="confirm-dialog-button"]').should("be.visible").click();
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
            .should("be.visible").click({ force: true }); // TODO: Is this a bug?
          cy.get('[data-e2e="confirm-dialog-button"]').should("be.visible").click();
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
        (target: ZITADELTarget, userId:number) => {
          ensureHumanIsNotOrgMember(target, userId);
        },
        (target: ZITADELTarget, userId:number) => {
          ensureHumanIsNotOrgMember(target, userId);
          ensureHumanIsOrgMember(
            target,
            userId,
            roles.map((role) => role.internal),
          );
        },
        (target: ZITADELTarget) => {
          cy.visit(`/org?org=${target.headers['x-zitadel-orgid']}`);
        },
      );
    });

    describe('projects', () => {
      describe('owned projects', () => {

        function visitOwnedProject(target: ZITADELTarget) {
          cy.get<number>('@projectId').then((projectId) => {
            cy.visit(`/projects/${projectId}?org=${target.headers['x-zitadel-orgid']}`);
          });
        }

        beforeEach(()=>{
          cy.get<ZITADELTarget>('@target').then(target => {
            ensureProjectExists(target, 'e2ecreateauthorization')
            .as('projectId')
          })
        })

        describe('authorizations', () => {
          const roles = [
            { internal: 'PROJECT_OWNER_GLOBAL', display: 'Project Owner Global' },
            { internal: 'PROJECT_OWNER_VIEWER_GLOBAL', display: 'Project Owner Viewer Global' },
          ];

          testAuthorizations(
            roles.map((role) => role.display),
            (target: ZITADELTarget, userId:number) => {
              cy.get<number>('@projectId').then((projectId) => {
                  ensureHumanIsNotProjectMember(target, projectId, userId);
                });
            },
            (target: ZITADELTarget, userId:number) => {
              ensureProjectExists(target, 'e2emutateauthorization')
                .as('projectId')
                .then((projectId) => {
                  ensureHumanIsNotProjectMember(target, projectId, userId);
                  ensureHumanIsProjectMember(
                    target,
                    projectId,
                    userId,
                    roles.map((role) => role.internal),
                  );
                });
            },
            visitOwnedProject,
          );
        });

        describe('roles', () => {
          const testRoleName = 'e2eroleundertestname';

          beforeEach(() => {
            cy.get<ZITADELTarget>('@target').then((target) => {
              cy.get<number>('@projectId').then((projectId) => {
                ensureRoleDoesntExist(target, projectId, testRoleName);
                visitOwnedProject(target);
              });
            });
          });

          it('should add a role', () => {
            cy.get('[data-e2e="sidenav-element-roles"]').should("be.visible").click();
            cy.get('[data-e2e="add-new-role"]').should("be.visible").click();
            cy.get('[formcontrolname="key"]').should("be.visible").type(testRoleName);
            cy.get('[formcontrolname="displayName"]').should("be.visible").type('e2eroleundertestdisplay');
            cy.get('[formcontrolname="group"]').should("be.visible").type('e2eroleundertestgroup');
            cy.get('[data-e2e="save-button"]').should("be.visible").click();
            cy.shouldConfirmSuccess();
            cy.contains('tr', testRoleName);
          });
          it('should remove a role');
        });
      });

      describe('granted projects', () => {
        describe('authorizations', () => {
          const roles = [
            { internal: 'PROJECT_GRANT_OWNER', display: 'Project Grant Owner' },
            { internal: 'PROJECT_GRANT_OWNER_VIEWER', display: 'Project Grant Owner Viewer' },
          ];

          testAuthorizations(
            roles.map((role) => role.display),
            (target: ZITADELTarget, userId:number) => {
              ensureOrgExists(target, 'e2eforeignorg').then((foreignOrgTarget) => {
                ensureProjectExists(foreignOrgTarget, 'e2eprojectgrants')
                  .as('projectId')
                  .then((foreignProjectId) => {
                    ensureProjectGrantExists(foreignOrgTarget, foreignProjectId, parseInt(target.headers['x-zitadel-orgid']))
                      .as('grantId')
                      .then((grantId) => {
                        ensureHumanIsNotProjectMember(target, foreignProjectId, userId, grantId);
                      });
                  });
              });
            },
            (target: ZITADELTarget, userId:number) => {
              ensureOrgExists(target, 'e2eforeignorg').then((foreignOrgTarget) => {
                ensureProjectExists(foreignOrgTarget, 'e2eprojectgrants')
                  .as('projectId')
                  .then((foreignProjectId) => {
                    ensureProjectGrantExists(foreignOrgTarget, foreignProjectId, parseInt(target.headers['x-zitadel-orgid']))
                      .as('grantId')
                      .then((grantId) => {
                        ensureHumanIsNotProjectMember(target, foreignProjectId, userId, grantId);
                        ensureHumanIsProjectMember(
                          target,
                          foreignProjectId,
                          userId,
                          roles.map((role) => role.internal),
                          grantId,
                        );
                      });
                  });
              });
            },
            (target: ZITADELTarget) => {
              cy.get<number>('@projectId').then((projectId) => {
                cy.get<number>('@grantId').then((grantId) => {
                  cy.visit(`/granted-projects/${projectId}/grant/${grantId}?org=${target.headers['x-zitadel-orgid']}`);
                });
              });
            },
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
