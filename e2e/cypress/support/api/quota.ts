import { Context } from 'support/commands';
import { SystemAPI } from './types';

export enum Unit {
  Unimplemented,
  AuthenticatedRequests,
  ExecutionSeconds,
}

interface notification {
  percent: 100;
  repeat?: boolean;
  callUrl: string;
}

export function addQuota(
  ctx: Context,
  failOnStatusCode = true,
  unit: Unit = Unit.AuthenticatedRequests,
  amount: number = 25000,
  intervalSeconds: string = `${30 * 24 * 60 * 60}s`,
  limit: boolean = true,
  notifications?: Array<notification>,
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
      interval: intervalSeconds,
      limit: limit,
      notifications: notifications,
    },
    failOnStatusCode: failOnStatusCode,
  });
}

export function ensureQuotaIsAdded(
  ctx: Context,
  unit?: Unit,
  amount?: number,
  intervalSeconds?: string,
  limit?: boolean,
  notifications?: Array<notification>,
): Cypress.Chainable<null> {
  return addQuota(ctx, false, unit, amount, intervalSeconds, limit, notifications).then((res) => {
    if (!res.isOkStatusCode) {
      expect(res.status).to.equal(409);
    }
    return null;
  });
}

export function removeQuota(
  ctx: Context,
  failOnStatusCode = true,
  unit: Unit = Unit.AuthenticatedRequests,
): Cypress.Chainable<Cypress.Response<any>> {
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
  return removeQuota(ctx, false, unit).then((res) => {
    if (!res.isOkStatusCode) {
      expect(res.status).to.equal(404);
    }
    return null;
  });
}
