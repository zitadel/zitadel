import { requestHeaders } from './apiauth';
import { API } from './types';

export interface ZitadelEvent {
  type: { type: string };
  aggregate: { id: string; type: { type: string }; resourceOwner: string };
  editor: { userId: string; displayName: string; service: string };
  sequence: string;
  creationDate: string;
}

export interface ListEventsRequest {
  aggregateId?: string;
  aggregateTypes?: string[];
  eventTypes?: string[];
  editorUserId?: string;
  resourceOwner?: string;
  limit?: number;
  asc?: boolean;
}

export function listGrantEvents(api: API, grantId: string): Cypress.Chainable<ZitadelEvent[]> {
  return listEvents(api, {
    aggregateId: grantId,
    aggregateTypes: ['groupgrant'],
    asc: true,
  });
}

export function listGrantEventTypes(api: API, grantId: string): Cypress.Chainable<string[]> {
  return listGrantEvents(api, grantId).then((events) => events.map((e) => e.type.type));
}

export function awaitGrantEventTypes(
  api: API,
  grantId: string,
  predicate: (types: string[]) => boolean,
  trials = 20,
): Cypress.Chainable<string[]> {
  return listGrantEventTypes(api, grantId).then((types) => {
    if (predicate(types)) {
      return cy.wrap(types);
    }
    expect(trials, `event projection on grant ${grantId}`).to.be.greaterThan(0);
    return cy.wait(500).then(() => awaitGrantEventTypes(api, grantId, predicate, trials - 1));
  });
}

export function assertGrantEvents(api: API, grantId: string, expected: string[]): Cypress.Chainable<string[]> {
  return awaitGrantEventTypes(api, grantId, (types) => arraysEqual(types, expected)).then((types) => {
    expect(types, `event sequence for grant ${grantId}`).to.deep.equal(expected);
    return cy.wrap(types);
  });
}

export function listEvents(api: API, req: ListEventsRequest): Cypress.Chainable<ZitadelEvent[]> {
  return cy
    .request({
      method: 'POST',
      url: `${api.adminBaseURL}/events/_search`,
      headers: requestHeaders(api),
      body: {
        aggregateId: req.aggregateId,
        aggregateTypes: req.aggregateTypes,
        eventTypes: req.eventTypes,
        editorUserId: req.editorUserId,
        resourceOwner: req.resourceOwner,
        limit: req.limit ?? 200,
        asc: req.asc ?? true,
      },
    })
    .then((res) => (res.body.events ?? []) as ZitadelEvent[]);
}

export function listGroupEvents(api: API, groupId: string): Cypress.Chainable<ZitadelEvent[]> {
  return listEvents(api, {
    aggregateId: groupId,
    aggregateTypes: ['group', 'groupgrant'],
    asc: true,
  });
}

export function listGroupEventTypes(api: API, groupId: string): Cypress.Chainable<string[]> {
  return listGroupEvents(api, groupId).then((events) => events.map((e) => e.type.type));
}

export function awaitGroupEventTypes(
  api: API,
  groupId: string,
  predicate: (types: string[]) => boolean,
  trials = 20,
): Cypress.Chainable<string[]> {
  return listGroupEventTypes(api, groupId).then((types) => {
    if (predicate(types)) {
      return cy.wrap(types);
    }
    expect(trials, `event projection on ${groupId}`).to.be.greaterThan(0);
    return cy.wait(500).then(() => awaitGroupEventTypes(api, groupId, predicate, trials - 1));
  });
}

export function assertGroupEvents(api: API, groupId: string, expected: string[]): Cypress.Chainable<string[]> {
  return awaitGroupEventTypes(api, groupId, (types) => arraysEqual(types, expected)).then((types) => {
    expect(types, `event sequence for group ${groupId}`).to.deep.equal(expected);
    return cy.wrap(types);
  });
}

export function assertGroupEventsContain(api: API, groupId: string, expected: string[]): Cypress.Chainable<string[]> {
  return awaitGroupEventTypes(api, groupId, (types) => expected.every((t) => types.includes(t))).then((types) => {
    for (const t of expected) {
      expect(types, `expected group ${groupId} to contain event ${t}`).to.include(t);
    }
    return cy.wrap(types);
  });
}

export function assertLastGroupEvent(api: API, groupId: string, expectedType: string): Cypress.Chainable<string> {
  return awaitGroupEventTypes(api, groupId, (types) => types[types.length - 1] === expectedType).then((types) => {
    expect(types[types.length - 1], `last event for group ${groupId}`).to.equal(expectedType);
    return cy.wrap(types[types.length - 1]);
  });
}

export function assertNoGroupEventsAppendedSince(
  api: API,
  groupId: string,
  baselineCount: number,
  waitMs = 1500,
): Cypress.Chainable<string[]> {
  return cy.wait(waitMs).then(() =>
    listGroupEventTypes(api, groupId).then((types) => {
      expect(types.length, `no new events appended for group ${groupId}`).to.equal(baselineCount);
      return cy.wrap(types);
    }),
  );
}

export function assertEventOwner(api: API, groupId: string, expectedOrgId: string): Cypress.Chainable<ZitadelEvent[]> {
  return listGroupEvents(api, groupId).then((events) => {
    for (const e of events) {
      expect(e.aggregate.resourceOwner, `event ${e.type.type} on group ${groupId}`).to.equal(expectedOrgId);
    }
    return cy.wrap(events);
  });
}

function arraysEqual(a: string[], b: string[]): boolean {
  if (a.length !== b.length) return false;
  for (let i = 0; i < a.length; i++) {
    if (a[i] !== b[i]) return false;
  }
  return true;
}

export const GroupEventTypes = {
  Added: 'group.added',
  Changed: 'group.changed',
  Removed: 'group.removed',
  UserAdded: 'group.user.added',
  UserRemoved: 'group.user.removed',
  ManagerRolesSet: 'group.manager.roles.set',
  GrantAdded: 'group.grant.added',
  GrantChanged: 'group.grant.changed',
  GrantRemoved: 'group.grant.removed',
  GrantCascadeRemoved: 'group.grant.cascade.removed',
} as const;
