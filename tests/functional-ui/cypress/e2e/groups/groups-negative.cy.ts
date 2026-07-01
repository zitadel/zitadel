import { Context } from 'support/commands';
import {
  addUsersToGroup,
  createGroup,
  createGroupGrant,
  defaultOrgId,
  deleteGroup,
  deleteGroupGrant,
  ensureGroupDoesntExist,
  ensureGroupExists,
  removeUsersFromGroup,
  updateGroup,
  updateGroupGrant,
} from '../../support/api/groups';
import { ensureHumanUserExists } from '../../support/api/users';
import { ensureProjectExists, ensureRoleExists } from '../../support/api/projects';
import { assertNoGroupEventsAppendedSince, listGroupEventTypes, listGrantEventTypes } from '../../support/api/events';

describe('groups — negative paths (no event must be appended)', () => {
  beforeEach(() => {
    cy.context().as('ctx');
  });

  describe('CreateGroup', () => {
    const name = 'e2egroup-neg-duplicate';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, name);
      });
    });

    it('rejects a duplicate name in the same org and emits no extra event', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        defaultOrgId(ctx.api).then((orgId) => {
          createGroup(ctx.api, { organizationId: orgId, name }).then((first) => {
            const groupId = first.body.id;
            listGroupEventTypes(ctx.api, groupId).then((before) => {
              createGroup(ctx.api, { organizationId: orgId, name }, { failOnStatusCode: false }).then((second) => {
                expect(second.status, 'duplicate name rejected').to.be.gte(400);
                assertNoGroupEventsAppendedSince(ctx.api, groupId, before.length);
              });
            });
          });
        });
      });
    });
  });

  describe('UpdateGroup', () => {
    const nameA = 'e2egroup-neg-rename-a';
    const nameB = 'e2egroup-neg-rename-b';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, nameA);
        ensureGroupDoesntExist(ctx.api, nameB);
      });
    });

    it('rejects renaming to a name already taken by another group in the org', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupExists(ctx.api, nameA).then((groupA) => {
          ensureGroupExists(ctx.api, nameB).then((groupB) => {
            listGroupEventTypes(ctx.api, groupB).then((before) => {
              updateGroup(ctx.api, { id: groupB, name: nameA }, { failOnStatusCode: false }).then((res) => {
                expect(res.status, 'rename collision rejected').to.be.gte(400);
                assertNoGroupEventsAppendedSince(ctx.api, groupB, before.length);
              });
            });
          });
        });
      });
    });
  });

  describe('DeleteGroup', () => {
    it('returns success for an unknown group id (idempotent) and emits no event', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        const fakeId = '999999999999999999';
        deleteGroup(ctx.api, fakeId, { failOnStatusCode: false }).then((res) => {
          expect(res.status, 'idempotent delete returns 200').to.equal(200);
          listGroupEventTypes(ctx.api, fakeId).should((types) => {
            expect(types, 'no events on fake aggregate').to.deep.equal([]);
          });
        });
      });
    });
  });

  describe('AddUsersToGroup', () => {
    const groupName = 'e2egroup-neg-adduser';
    const username = 'e2egroup-neg-user';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
        ensureHumanUserExists(ctx.api, username);
      });
    });

    it('rejects when any user id is unknown — all-or-nothing — and emits no event', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureHumanUserExists(ctx.api, username).then((knownUserId) => {
          ensureGroupExists(ctx.api, groupName).then((groupId) => {
            listGroupEventTypes(ctx.api, groupId).then((before) => {
              addUsersToGroup(ctx.api, groupId, [knownUserId, '888888888888888888'], {
                failOnStatusCode: false,
              }).then((res) => {
                expect(res.status, 'mixed batch with unknown user rejected').to.be.gte(400);
                assertNoGroupEventsAppendedSince(ctx.api, groupId, before.length);
              });
            });
          });
        });
      });
    });

    it('is idempotent when re-adding the same user — appends no second group.user.added', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureHumanUserExists(ctx.api, username).then((userId) => {
          ensureGroupExists(ctx.api, groupName).then((groupId) => {
            addUsersToGroup(ctx.api, groupId, [userId]).then(() => {
              listGroupEventTypes(ctx.api, groupId).then((after1) => {
                addUsersToGroup(ctx.api, groupId, [userId]).then(() => {
                  assertNoGroupEventsAppendedSince(ctx.api, groupId, after1.length);
                });
              });
            });
          });
        });
      });
    });
  });

  describe('RemoveUsersFromGroup', () => {
    const groupName = 'e2egroup-neg-removeuser';
    const username = 'e2egroup-neg-removeuser-name';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
      });
    });

    it('is idempotent when removing a non-member and appends no event', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureHumanUserExists(ctx.api, username).then((userId) => {
          ensureGroupExists(ctx.api, groupName).then((groupId) => {
            listGroupEventTypes(ctx.api, groupId).then((before) => {
              removeUsersFromGroup(ctx.api, groupId, [userId]).then(() => {
                assertNoGroupEventsAppendedSince(ctx.api, groupId, before.length);
              });
            });
          });
        });
      });
    });
  });

  describe('CreateGroupGrant', () => {
    const groupName = 'e2egroup-neg-grant';
    const projectName = 'e2egroup-neg-grant-project';
    const roleKey = 'e2egroup-neg-grant-role';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
        ensureProjectExists(ctx.api, projectName).then((projectId) => {
          ensureRoleExists(ctx.api, projectId, roleKey);
          cy.wrap(projectId).as('projectId');
        });
      });
    });

    it('rejects a role the project does not have and emits no grant event on the group', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        cy.get<string>('@projectId').then((projectId) => {
          ensureGroupExists(ctx.api, groupName).then((groupId) => {
            listGroupEventTypes(ctx.api, groupId).then((before) => {
              createGroupGrant(
                ctx.api,
                { groupId, projectId, roleKeys: ['ROLE_THAT_DOES_NOT_EXIST'] },
                { failOnStatusCode: false },
              ).then((res) => {
                expect(res.status, 'bad role rejected').to.be.gte(400);
                assertNoGroupEventsAppendedSince(ctx.api, groupId, before.length);
              });
            });
          });
        });
      });
    });

    it('rejects against an unknown project id', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupExists(ctx.api, groupName).then((groupId) => {
          listGroupEventTypes(ctx.api, groupId).then((before) => {
            createGroupGrant(
              ctx.api,
              { groupId, projectId: '777777777777777777', roleKeys: [roleKey] },
              { failOnStatusCode: false },
            ).then((res) => {
              expect(res.status, 'unknown project rejected').to.be.gte(400);
              assertNoGroupEventsAppendedSince(ctx.api, groupId, before.length);
            });
          });
        });
      });
    });
  });

  describe('UpdateGroupGrant — removed grant', () => {
    const groupName = 'e2egroup-neg-update-removed';
    const projectName = 'e2egroup-neg-update-removed-project';
    const roleKey = 'e2egroup-neg-update-removed-role';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
        ensureProjectExists(ctx.api, projectName).then((projectId) => {
          ensureRoleExists(ctx.api, projectId, roleKey);
          cy.wrap(projectId).as('projectId');
        });
      });
    });

    it('rejects updating a grant that has already been deleted', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        cy.get<string>('@projectId').then((projectId) => {
          ensureGroupExists(ctx.api, groupName).then((groupId) => {
            createGroupGrant(ctx.api, { groupId, projectId, roleKeys: [roleKey] }).then((res) => {
              const grantId = res.body.id;
              deleteGroupGrant(ctx.api, grantId).then(() => {
                listGrantEventTypes(ctx.api, grantId).then((before) => {
                  updateGroupGrant(ctx.api, grantId, [roleKey], { failOnStatusCode: false }).then((res2) => {
                    expect(res2.status, 'update on removed grant rejected').to.be.gte(400);
                    cy.wait(1000).then(() => {
                      listGrantEventTypes(ctx.api, grantId).should((after) => {
                        expect(after.length, 'no new events on removed grant').to.equal(before.length);
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
});
