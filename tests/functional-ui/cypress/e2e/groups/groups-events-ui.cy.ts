import { Context } from 'support/commands';
import {
  defaultOrgId,
  createGroup,
  updateGroup,
  deleteGroup,
  ensureGroupDoesntExist,
} from '../../support/api/groups';
import { assertGroupEvents, GroupEventTypes } from '../../support/api/events';

describe('groups — admin events log surfaces group actions', () => {
  beforeEach(() => {
    cy.context().as('ctx');
  });

  const name = 'e2egroup-events-ui';
  const renamed = 'e2egroup-events-ui-renamed';

  beforeEach(() => {
    cy.get<Context>('@ctx').then((ctx) => {
      ensureGroupDoesntExist(ctx.api, name);
      ensureGroupDoesntExist(ctx.api, renamed);
    });
  });

  it('shows group.added, group.changed and group.removed rows when filtered by aggregate id', () => {
    cy.get<Context>('@ctx').then((ctx) => {
      defaultOrgId(ctx.api).then((orgId) => {
        createGroup(ctx.api, { organizationId: orgId, name }).then((res) => {
          const groupId = res.body.id;

          updateGroup(ctx.api, { id: groupId, name: renamed }).then(() => {
            deleteGroup(ctx.api, groupId).then(() => {
              assertGroupEvents(ctx.api, groupId, [
                GroupEventTypes.Added,
                GroupEventTypes.Changed,
                GroupEventTypes.Removed,
              ]);

              cy.visit('/instance?id=events');
              cy.get('[data-e2e="open-filter-button"]').click();
              cy.get('mat-checkbox#aggregateFilterSet').click();
              cy.get('input#aggregateId').type(groupId);
              cy.get('[data-e2e="filter-finish-button"]').click();

              cy.get('[data-e2e="event-type-cell"]').should('have.length.gte', 3);
            });
          });
        });
      });
    });
  });
});
