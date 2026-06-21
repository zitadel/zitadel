import { Context } from 'support/commands';
import {
  addUsersToGroup,
  createGroup,
  createGroupGrant,
  defaultOrgId,
  ensureGroupDoesntExist,
  ensureGroupExists,
} from '../../support/api/groups';
import { ensureOrgExists } from '../../support/api/orgs';
import { ensureHumanUserExists, ensureHumanUserExistsInOrg } from '../../support/api/users';
import { ensureProjectExists, ensureRoleExists } from '../../support/api/projects';
import {
  assertNoGroupEventsAppendedSince,
  GroupEventTypes,
  listGroupEventTypes,
} from '../../support/api/events';

const foreignOrgName = 'e2eforeignorg-groups';

describe('groups — cross-org invariants', () => {
  beforeEach(() => {
    cy.context().as('ctx');
    cy.get<Context>('@ctx').then((ctx) => {
      ensureOrgExists(ctx, foreignOrgName).as('foreignOrgId');
    });
  });

  describe('CreateGroup — name uniqueness is scoped to org', () => {
    const sharedName = 'e2egroup-xorg-shared';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, sharedName);
      });
    });

    it('allows the same group name to be created in two distinct orgs', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        cy.get<string>('@foreignOrgId').then((foreignOrgId) => {
          defaultOrgId(ctx.api).then((homeOrgId) => {
            createGroup(ctx.api, { organizationId: homeOrgId, name: sharedName }).then((homeRes) => {
              expect(homeRes.status, 'home org create').to.be.oneOf([200, 201]);
              const homeId = homeRes.body.id;

              createGroup(ctx.api, { organizationId: foreignOrgId, name: sharedName }).then((foreignRes) => {
                expect(foreignRes.status, 'foreign org create with same name').to.be.oneOf([200, 201]);
                expect(foreignRes.body.id, 'distinct aggregate ID').to.not.equal(homeId);
              });
            });
          });
        });
      });
    });
  });

  describe('AddUsersToGroup — rejects users from a different org', () => {
    const groupName = 'e2egroup-xorg-add-foreign-user';
    const foreignUsername = 'e2egroup-xorg-foreign-user';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
      });
    });

    it('rejects adding a user that lives in another org and appends no event', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        cy.get<string>('@foreignOrgId').then((foreignOrgId) => {
          ensureHumanUserExistsInOrg(ctx.api, foreignOrgId, foreignUsername).then((foreignUserId) => {
            ensureGroupExists(ctx.api, groupName).then((groupId) => {
              listGroupEventTypes(ctx.api, groupId).then((before) => {
                addUsersToGroup(ctx.api, groupId, [foreignUserId], { failOnStatusCode: false }).then((res) => {
                  expect(res.status, 'cross-org user rejected').to.be.gte(400);
                  assertNoGroupEventsAppendedSince(ctx.api, groupId, before.length);
                });
              });
            });
          });
        });
      });
    });
  });

  describe('CreateGroupGrant — rejects projects the group org does not own', () => {
    const groupName = 'e2egroup-xorg-grant-foreign-project';
    const foreignProjectName = 'e2egroup-xorg-foreign-project';
    const foreignRoleKey = 'e2egroup-xorg-foreign-role';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
        cy.get<string>('@foreignOrgId').then((foreignOrgId) => {
          ensureProjectExists(ctx.api, foreignProjectName, foreignOrgId).then((foreignProjectId) => {
            ensureRoleExists(ctx.api, foreignProjectId, foreignRoleKey);
            cy.wrap(foreignProjectId).as('foreignProjectId');
          });
        });
      });
    });

    it('rejects a grant on a project that is neither owned by nor granted to the group org', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        cy.get<number | string>('@foreignProjectId').then((foreignProjectId) => {
          ensureGroupExists(ctx.api, groupName).then((groupId) => {
            listGroupEventTypes(ctx.api, groupId).then((before) => {
              createGroupGrant(
                ctx.api,
                { groupId, projectId: `${foreignProjectId}`, roleKeys: [foreignRoleKey] },
                { failOnStatusCode: false },
              ).then((res) => {
                expect(res.status, 'cross-org project grant rejected').to.be.gte(400);
                assertNoGroupEventsAppendedSince(ctx.api, groupId, before.length);
              });
            });
          });
        });
      });
    });
  });
});
