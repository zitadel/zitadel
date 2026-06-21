import { Context } from 'support/commands';
import {
  createGroupGrant,
  deleteGroupGrant,
  ensureGroupDoesntExist,
  ensureGroupExists,
  listGroupGrants,
  updateGroupGrant,
} from '../../support/api/groups';
import { ensureProjectExists, ensureRoleExists } from '../../support/api/projects';
import { assertGrantEvents, assertGroupEvents, GroupEventTypes } from '../../support/api/events';

describe('groups — grants (API + event log)', () => {
  beforeEach(() => {
    cy.context().as('ctx');
  });

  describe('UpdateGroupGrant — change roles', () => {
    const groupName = 'e2egroup-grant-change';
    const projectName = 'e2egroup-grant-change-project';
    const roleA = 'e2egroup-grant-roleA';
    const roleB = 'e2egroup-grant-roleB';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
        ensureProjectExists(ctx.api, projectName).then((projectId) => {
          ensureRoleExists(ctx.api, projectId, roleA);
          ensureRoleExists(ctx.api, projectId, roleB);
          cy.wrap(projectId).as('projectId');
        });
      });
    });

    it('emits group.grant.changed when roles are replaced', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        cy.get<string>('@projectId').then((projectId) => {
          ensureGroupExists(ctx.api, groupName).then((groupId) => {
            createGroupGrant(ctx.api, { groupId, projectId, roleKeys: [roleA] }).then((res) => {
              const grantId = res.body.id;

              updateGroupGrant(ctx.api, grantId, [roleA, roleB]).then(() => {
                assertGrantEvents(ctx.api, grantId, [
                  GroupEventTypes.GrantAdded,
                  GroupEventTypes.GrantChanged,
                ]);

                listGroupGrants(ctx.api, groupId).should((grants) => {
                  const match = grants.find((g) => g.id === grantId);
                  expect(match, 'grant present').to.not.be.undefined;
                  expect(match.roleKeys, 'roles replaced').to.have.members([roleA, roleB]);
                });
              });
            });
          });
        });
      });
    });

    it('is idempotent — updating with the same roles emits no further event', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        cy.get<string>('@projectId').then((projectId) => {
          ensureGroupExists(ctx.api, groupName).then((groupId) => {
            createGroupGrant(ctx.api, { groupId, projectId, roleKeys: [roleA] }).then((res) => {
              const grantId = res.body.id;

              updateGroupGrant(ctx.api, grantId, [roleA]).then(() => {
                cy.wait(1000).then(() => {
                  assertGrantEvents(ctx.api, grantId, [GroupEventTypes.GrantAdded]);
                });
              });
            });
          });
        });
      });
    });
  });

  describe('DeleteGroupGrant — explicit revoke', () => {
    const groupName = 'e2egroup-grant-revoke';
    const projectName = 'e2egroup-grant-revoke-project';
    const roleKey = 'e2egroup-grant-revoke-role';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
        ensureProjectExists(ctx.api, projectName).then((projectId) => {
          ensureRoleExists(ctx.api, projectId, roleKey);
          cy.wrap(projectId).as('projectId');
        });
      });
    });

    it('emits group.grant.removed (non-cascade) when the grant is deleted directly', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        cy.get<string>('@projectId').then((projectId) => {
          ensureGroupExists(ctx.api, groupName).then((groupId) => {
            createGroupGrant(ctx.api, { groupId, projectId, roleKeys: [roleKey] }).then((res) => {
              const grantId = res.body.id;

              deleteGroupGrant(ctx.api, grantId).then(() => {
                assertGrantEvents(ctx.api, grantId, [
                  GroupEventTypes.GrantAdded,
                  GroupEventTypes.GrantRemoved,
                ]);

                assertGroupEvents(ctx.api, groupId, [GroupEventTypes.Added]);

                listGroupGrants(ctx.api, groupId).should((grants) => {
                  expect(grants.find((g) => g.id === grantId), 'grant gone').to.be.undefined;
                });
              });
            });
          });
        });
      });
    });
  });
});
