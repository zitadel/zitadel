import { Context } from 'support/commands';
import { ensureGroupDoesntExist, ensureGroupExists, setGroupManagerRoles } from '../../support/api/groups';
import {
  assertGroupEvents,
  assertGroupEventsContain,
  assertNoGroupEventsAppendedSince,
  GroupEventTypes,
  listGroupEventTypes,
} from '../../support/api/events';

const VALID_ROLE = 'ORG_OWNER_VIEWER';
const VALID_ROLE_2 = 'ORG_USER_MANAGER';

describe('groups — manager roles (API + event log)', () => {
  beforeEach(() => {
    cy.context().as('ctx');
  });

  describe('SetGroupManagerRoles — set initial', () => {
    const groupName = 'e2egroup-mgr-set';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
      });
    });

    it('emits group.manager.roles.set when roles are first applied', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupExists(ctx.api, groupName).then((groupId) => {
          setGroupManagerRoles(ctx.api, groupId, [VALID_ROLE]).then(() => {
            assertGroupEvents(ctx.api, groupId, [GroupEventTypes.Added, GroupEventTypes.ManagerRolesSet]);
          });
        });
      });
    });
  });

  describe('SetGroupManagerRoles — replace', () => {
    const groupName = 'e2egroup-mgr-replace';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
      });
    });

    it('emits a second group.manager.roles.set when the role list changes', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupExists(ctx.api, groupName).then((groupId) => {
          setGroupManagerRoles(ctx.api, groupId, [VALID_ROLE]).then(() => {
            setGroupManagerRoles(ctx.api, groupId, [VALID_ROLE, VALID_ROLE_2]).then(() => {
              assertGroupEvents(ctx.api, groupId, [
                GroupEventTypes.Added,
                GroupEventTypes.ManagerRolesSet,
                GroupEventTypes.ManagerRolesSet,
              ]);
            });
          });
        });
      });
    });

    it('is idempotent — re-applying the same set order-independent emits nothing extra', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupExists(ctx.api, groupName).then((groupId) => {
          setGroupManagerRoles(ctx.api, groupId, [VALID_ROLE, VALID_ROLE_2]).then(() => {
            assertGroupEventsContain(ctx.api, groupId, [GroupEventTypes.ManagerRolesSet]);

            setGroupManagerRoles(ctx.api, groupId, [VALID_ROLE_2, VALID_ROLE]).then(() => {
              cy.wait(1000).then(() => {
                listGroupEventTypes(ctx.api, groupId).then((types) => {
                  const setCount = types.filter((t) => t === GroupEventTypes.ManagerRolesSet).length;
                  expect(setCount, 'one set event despite order-only re-apply').to.equal(1);
                });
              });
            });
          });
        });
      });
    });
  });

  describe('SetGroupManagerRoles — clear', () => {
    const groupName = 'e2egroup-mgr-clear';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
      });
    });

    it('emits group.manager.roles.set with an empty list when roles are cleared', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupExists(ctx.api, groupName).then((groupId) => {
          setGroupManagerRoles(ctx.api, groupId, [VALID_ROLE]).then(() => {
            setGroupManagerRoles(ctx.api, groupId, []).then(() => {
              assertGroupEvents(ctx.api, groupId, [
                GroupEventTypes.Added,
                GroupEventTypes.ManagerRolesSet,
                GroupEventTypes.ManagerRolesSet,
              ]);
            });
          });
        });
      });
    });

    it('emits no event when clearing an already-empty role set', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupExists(ctx.api, groupName).then((groupId) => {
          setGroupManagerRoles(ctx.api, groupId, []).then(() => {
            cy.wait(1000).then(() => {
              assertGroupEvents(ctx.api, groupId, [GroupEventTypes.Added]);
            });
          });
        });
      });
    });
  });

  describe('SetGroupManagerRoles — invalid roles', () => {
    const groupName = 'e2egroup-mgr-invalid';

    beforeEach(() => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupDoesntExist(ctx.api, groupName);
      });
    });

    it('rejects a role without the ORG prefix and appends no event', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupExists(ctx.api, groupName).then((groupId) => {
          listGroupEventTypes(ctx.api, groupId).then((before) => {
            const baseline = before.length;
            setGroupManagerRoles(ctx.api, groupId, ['NOT_ORG_PREFIXED'], { failOnStatusCode: false }).then((res) => {
              expect(res.status, 'API rejected the call').to.be.gte(400);
              assertNoGroupEventsAppendedSince(ctx.api, groupId, baseline);
            });
          });
        });
      });
    });

    it('rejects an unknown role even with the ORG prefix and appends no event', () => {
      cy.get<Context>('@ctx').then((ctx) => {
        ensureGroupExists(ctx.api, groupName).then((groupId) => {
          listGroupEventTypes(ctx.api, groupId).then((before) => {
            const baseline = before.length;
            setGroupManagerRoles(ctx.api, groupId, ['ORG_NONEXISTENT_ROLE'], { failOnStatusCode: false }).then((res) => {
              expect(res.status, 'API rejected the call').to.be.gte(400);
              assertNoGroupEventsAppendedSince(ctx.api, groupId, baseline);
            });
          });
        });
      });
    });
  });
});
