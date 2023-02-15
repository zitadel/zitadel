import { Context } from 'support/commands';

export enum Unit {
  Unimplemented,
  AuthenticatedRequests,
  ExecutionSeconds,
}

interface notification {
  percent: number;
  repeat?: boolean;
  callUrl: string;
}

export function addQuota(
  ctx: Context,
  unit: Unit = Unit.AuthenticatedRequests,
  limit: boolean,
  amount: number,
  notifications?: Array<notification>,
  from: Date = (() => {
    const date = new Date();
    date.setMonth(0, 1);
    date.setMinutes(0, 0, 0);
    // default to start of current year
    return date;
  })(),
  intervalSeconds: string = `${315_576_000_000}s`, // proto max duration is 1000 years
  failOnStatusCode = true,
): Cypress.Chainable<Cypress.Response<any>> {
  return cy.request({
    method: 'POST',
    url: `${ctx.system.baseURL}/instances/${ctx.instanceId}/quotas`,
    auth: {
      bearer: ctx.system.token,
    },
    body: {
      unit: unit,
      amount: amount,
      resetInterval: intervalSeconds,
      limit: limit,
      from: from,
      notifications: notifications,
    },
    failOnStatusCode: failOnStatusCode,
  });
}

export function ensureQuotaIsAdded(
  ctx: Context,
  unit: Unit,
  limit: boolean,
  amount?: number,
  notifications?: Array<notification>,
  from?: Date,
  intervalSeconds?: string,
): Cypress.Chainable<null> {
  return addQuota(ctx, unit, limit, amount, notifications, from, intervalSeconds, false).then((res) => {
    if (!res.isOkStatusCode) {
      expect(res.status).to.equal(409);
    }
    return null;
  });
}

export function removeQuota(ctx: Context, unit: Unit, failOnStatusCode = true): Cypress.Chainable<Cypress.Response<any>> {
  return cy.request({
    method: 'DELETE',
    url: `${ctx.system.baseURL}/instances/${ctx.instanceId}/quotas/${unit}`,
    auth: {
      bearer: ctx.system.token,
    },
    failOnStatusCode: failOnStatusCode,
  });
}

export function ensureQuotaIsRemoved(ctx: Context, unit?: Unit): Cypress.Chainable<null> {
  return removeQuota(ctx, unit, false).then((res) => {
    if (!res.isOkStatusCode) {
      expect(res.status).to.equal(404);
    }
    return null;
  });
}
