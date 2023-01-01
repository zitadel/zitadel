import { apiAuth, systemAuth } from 'support/api/apiauth';
import { instanceUnderTest } from 'support/api/instances';
import { addQuota, ensureQuotaIsAdded, ensureQuotaIsRemoved, removeQuota, Unit } from 'support/api/quota';
import { API, SystemAPI } from 'support/api/types';
import { Context } from 'support/commands';

beforeEach(() => {
  cy.context().as('ctx');
});

describe('quotas', () => {
  describe('management', () => {
    describe('add one quota', () => {
      it('should add a quota only once per unit', function () {
        cy.get<Context>('@ctx').then((ctx) => {
          addQuota(ctx).then(() => {
            cy.log('1');
          });
          addQuota(ctx, false).then((res) => {
            cy.log('2');
            expect(res.status).to.equal(409);
          });
        });
      });

      describe('add two quotas', () => {
        beforeEach(function () {
          cy.get<Context>('@ctx').then((ctx) => {
            ensureQuotaIsRemoved(ctx, Unit.ExecutionSeconds);
          });
        });
        it('should add a quota for each unit', function () {
          cy.get<Context>('@ctx').then((ctx) => {
            addQuota(ctx);
            addQuota(ctx, true, Unit.ExecutionSeconds);
          });
        });
      });
    });

    describe('edit', () => {
      describe('remove one quota', () => {
        beforeEach(function () {
          cy.get<Context>('@ctx').then((ctx) => {
            ensureQuotaIsAdded(ctx);
          });
        });
        it('should remove a quota only once per unit', function () {
          cy.get<Context>('@ctx').then((ctx) => {
            removeQuota(ctx);
          });
          cy.get<Context>('@ctx').then((ctx) => {
            removeQuota(ctx, false).then((res) => {
              expect(res.status).to.equal(404);
            });
          });
        });

        describe('remove two quotas', () => {
          beforeEach(function () {
            cy.get<Context>('@ctx').then((ctx) => {
              ensureQuotaIsAdded(ctx, Unit.ExecutionSeconds);
            });
          });
          it('should remove a quota for each unit', function () {
            cy.get<Context>('@ctx').then((ctx) => {
              removeQuota(ctx);
              removeQuota(ctx, true, Unit.ExecutionSeconds);
            });
          });
        });
      });
    });
  });

  describe('usage', () => {
    describe('authenticated requests', () => {
      describe('notifications', () => {
        it('authenticated requests are limited', function () {
          cy.get<Context>('@ctx').then((ctx) => {
            const urls = [
              `${ctx.api.authBaseURL}/users/me`,
              `${ctx.api.mgmtBaseURL}/iam`,
              `${ctx.api.adminBaseURL}/instances/me`,
              // `${api.assetsBaseURL}/instance/policy/label/icon`,
              `${ctx.api.oidcBaseURL}/userinfo`,
              `${ctx.api.oauthBaseURL}/keys`,
              `${ctx.api.samlBaseURL}/certificate`,
            ];
            ensureQuotaIsAdded(ctx, Unit.AuthenticatedRequests, urls.length, undefined, true);
            cy.task('runSQL', `TRUNCATE logstore.access;`);
            urls.forEach((url) => {
              cy.request({
                url: url,
                method: 'GET',
                auth: {
                  bearer: ctx.api.token,
                },
              });
            });
            cy.request({
              url: `${ctx.api.oidcBaseURL}/userinfo`,
              method: 'GET',
              auth: {
                bearer: ctx.api.token,
              },
              failOnStatusCode: false,
            }).then((res) => {
              expect(res.status).to.equal(429);
            });
          });
        });
      });
      describe('cleanup', () => {});
    });

    //    it('receives', () => {
    //      cy.task('receive').then((received) => {
    //        cy.log('receive returned', received);
    //      });
    //    });
  });
});
