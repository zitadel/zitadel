import { Context } from 'support/commands';
import {
  addUsersToGroup,
  ensureGroupDoesntExist,
  ensureGroupExists,
  listGroupUsers,
  removeUsersFromGroup,
} from '../../support/api/groups';
import { ensureHumanUserExists, ensureUserDoesntExist } from '../../support/api/users';
import {
  assertGroupEvents,
  assertGroupEventsContain,
  GroupEventTypes,
} from '../../support/api/events';

describe('groups — members (API + event log)', () => {
  beforeEach(() => {
    cy.context().as('ctx');
  });

  describe('AddUsersToGroup — bulk', () => {
    const groupName = 'e2egroup-bulk-add';
    const userA = 'e2egroup-bulkuser-a';
    const userB = 'e2egroup-bulkuser-b';
    const userC = 'e2egroup-bulkuser-c';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
        ensureHumanUserExists(ctx.api, userA).then((idA) => cy.wrap(idA).as('idA'));
        ensureHumanUserExists(ctx.api, userB).then((idB) => cy.wrap(idB).as('idB'));
        ensureHumanUserExists(ctx.api, userC).then((idC) => cy.wrap(idC).as('idC'));
      });
    });

    it('emits one group.user.added per user when adding many in a single call', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        cy.get<string>('@idA').then((idA) => {
          cy.get<string>('@idB').then((idB) => {
            cy.get<string>('@idC').then((idC) => {
              ensureGroupExists(ctx.api, groupName).then((groupId) => {
                addUsersToGroup(ctx.api, groupId, [idA, idB, idC]).then(() => {
                  assertGroupEvents(ctx.api, groupId, [
                    GroupEventTypes.Added,
                    GroupEventTypes.UserAdded,
                    GroupEventTypes.UserAdded,
                    GroupEventTypes.UserAdded,
                  ]);

                  listGroupUsers(ctx.api, groupId).should((users) => {
                    expect(users, 'all three users present').to.have.length(3);
                    const ids = users.map((u: any) => u.userId);
                    expect(ids).to.include.members([idA, idB, idC]);
                  });
                });
              });
            });
          });
        });
      });
    });

    it('ignores duplicates inside a single bulk add — emits one event per fresh user only', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        cy.get<string>('@idA').then((idA) => {
          cy.get<string>('@idB').then((idB) => {
            ensureGroupExists(ctx.api, groupName).then((groupId) => {
              addUsersToGroup(ctx.api, groupId, [idA]).then(() => {
                addUsersToGroup(ctx.api, groupId, [idA, idB]).then(() => {
                  assertGroupEvents(ctx.api, groupId, [
                    GroupEventTypes.Added,
                    GroupEventTypes.UserAdded,
                    GroupEventTypes.UserAdded,
                  ]);
                });
              });
            });
          });
        });
      });
    });
  });

  describe('RemoveUsersFromGroup — bulk', () => {
    const groupName = 'e2egroup-bulk-remove';
    const userA = 'e2egroup-rmuser-a';
    const userB = 'e2egroup-rmuser-b';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
        ensureHumanUserExists(ctx.api, userA).then((idA) => cy.wrap(idA).as('idA'));
        ensureHumanUserExists(ctx.api, userB).then((idB) => cy.wrap(idB).as('idB'));
      });
    });

    it('emits group.user.removed for each user supplied in a single call', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        cy.get<string>('@idA').then((idA) => {
          cy.get<string>('@idB').then((idB) => {
            ensureGroupExists(ctx.api, groupName).then((groupId) => {
              addUsersToGroup(ctx.api, groupId, [idA, idB]).then(() => {
                removeUsersFromGroup(ctx.api, groupId, [idA, idB]).then(() => {
                  assertGroupEvents(ctx.api, groupId, [
                    GroupEventTypes.Added,
                    GroupEventTypes.UserAdded,
                    GroupEventTypes.UserAdded,
                    GroupEventTypes.UserRemoved,
                    GroupEventTypes.UserRemoved,
                  ]);

                  listGroupUsers(ctx.api, groupId).should((users) => {
                    expect(users, 'no users left').to.have.length(0);
                  });
                });
              });
            });
          });
        });
      });
    });

    it('is idempotent — removing a non-member does not append an event', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        cy.get<string>('@idA').then((idA) => {
          cy.get<string>('@idB').then((idB) => {
            ensureGroupExists(ctx.api, groupName).then((groupId) => {
              addUsersToGroup(ctx.api, groupId, [idA]).then(() => {
                removeUsersFromGroup(ctx.api, groupId, [idB]).then(() => {
                  cy.wait(1000).then(() => {
                    assertGroupEvents(ctx.api, groupId, [
                      GroupEventTypes.Added,
                      GroupEventTypes.UserAdded,
                    ]);
                  });
                });
              });
            });
          });
        });
      });
    });
  });

  describe('User deletion propagates to group membership', () => {
    const groupName = 'e2egroup-userdel-prop';
    const username = 'e2egroup-userdel';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
        ensureUserDoesntExist(ctx.api, username);
      });
    });

    it('appends group.user.removed on every group containing the user when the user is deleted', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureHumanUserExists(ctx.api, username).then((userId) => {
          ensureGroupExists(ctx.api, groupName).then((groupId) => {
            addUsersToGroup(ctx.api, groupId, [userId]).then(() => {
              assertGroupEventsContain(ctx.api, groupId, [GroupEventTypes.UserAdded]);

              ensureUserDoesntExist(ctx.api, username).then(() => {
                assertGroupEvents(ctx.api, groupId, [
                  GroupEventTypes.Added,
                  GroupEventTypes.UserAdded,
                  GroupEventTypes.UserRemoved,
                ]);
              });
            });
          });
        });
      });
    });
  });
});
