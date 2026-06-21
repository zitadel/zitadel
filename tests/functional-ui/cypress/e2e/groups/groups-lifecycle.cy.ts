import { Context } from 'support/commands';
import {
  createGroup,
  deleteGroup,
  ensureGroupDoesntExist,
  ensureGroupExists,
  getGroup,
  listGroupGrants,
  createGroupGrant,
  updateGroup,
  defaultOrgId,
} from '../../support/api/groups';
import { ensureProjectExists, ensureRoleExists } from '../../support/api/projects';
import {
  assertGrantEvents,
  assertGroupEvents,
  GroupEventTypes,
} from '../../support/api/events';

describe('groups — lifecycle (API + event log)', () => {
  beforeEach(() => {
    cy.context().as('ctx');
  });

  describe('UpdateGroup — description only', () => {
    const name = 'e2egroup-desc-only';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, name);
      });
    });

    it('emits group.added then group.changed when description is updated without renaming', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        defaultOrgId(ctx.api).then((orgId) => {
          createGroup(ctx.api, { organizationId: orgId, name, description: 'initial' }).then((res) => {
            const groupId = res.body.id;

            updateGroup(ctx.api, { id: groupId, description: 'updated description' }).then(() => {
              assertGroupEvents(ctx.api, groupId, [GroupEventTypes.Added, GroupEventTypes.Changed]);

              getGroup(ctx.api, groupId).should((g) => {
                expect(g, 'group still present').to.not.be.null;
                expect(g.name, 'name unchanged').to.equal(name);
                expect(g.description, 'description updated').to.equal('updated description');
              });
            });
          });
        });
      });
    });

    it('emits no group.changed when the desired state already matches (idempotent update)', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        defaultOrgId(ctx.api).then((orgId) => {
          createGroup(ctx.api, { organizationId: orgId, name, description: 'same' }).then((res) => {
            const groupId = res.body.id;

            updateGroup(ctx.api, { id: groupId, description: 'same' }).then(() => {
              cy.wait(1000).then(() => {
                assertGroupEvents(ctx.api, groupId, [GroupEventTypes.Added]);
              });
            });
          });
        });
      });
    });
  });

  describe('DeleteGroup — cascade group grants', () => {
    const name = 'e2egroup-cascade-delete';
    const projectAName = 'e2egroup-cascade-projectA';
    const projectBName = 'e2egroup-cascade-projectB';
    const roleKey = 'e2egroup-cascade-role';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, name);
        ensureProjectExists(ctx.api, projectAName).then((projectA) => {
          ensureRoleExists(ctx.api, projectA, roleKey);
          cy.wrap(projectA).as('projectA');
        });
        ensureProjectExists(ctx.api, projectBName).then((projectB) => {
          ensureRoleExists(ctx.api, projectB, roleKey);
          cy.wrap(projectB).as('projectB');
        });
      });
    });

    it('cascades a group.grant.cascade.removed for each grant when the group is deleted', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        cy.get<string>('@projectA').then((projectA) => {
          cy.get<string>('@projectB').then((projectB) => {
            ensureGroupExists(ctx.api, name).then((groupId) => {
              createGroupGrant(ctx.api, { groupId, projectId: projectA, roleKeys: [roleKey] }).then((resA) => {
                const grantA = resA.body.id;
                createGroupGrant(ctx.api, { groupId, projectId: projectB, roleKeys: [roleKey] }).then((resB) => {
                  const grantB = resB.body.id;

                  deleteGroup(ctx.api, groupId).then(() => {
                    assertGroupEvents(ctx.api, groupId, [
                      GroupEventTypes.Added,
                      GroupEventTypes.Removed,
                    ]);

                    assertGrantEvents(ctx.api, grantA, [
                      'group.grant.added',
                      GroupEventTypes.GrantCascadeRemoved,
                    ]);
                    assertGrantEvents(ctx.api, grantB, [
                      'group.grant.added',
                      GroupEventTypes.GrantCascadeRemoved,
                    ]);

                    listGroupGrants(ctx.api, groupId).should((grants) => {
                      expect(grants, 'no grants remain after cascade').to.have.length(0);
                    });
                  });
                });
              });
            });
          });
        });
      });
    });
  });
});
