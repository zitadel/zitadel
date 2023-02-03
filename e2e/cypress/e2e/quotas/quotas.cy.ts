import { addQuota, ensureQuotaIsAdded, removeQuota, Unit } from 'support/api/quota';
import { Context } from 'support/commands';
import { ZITADELWebhookEvent } from 'support/types';

beforeEach(() => {
  cy.context().as('ctx');
});

describe('quotas', () => {
  describe('management', () => {
    describe('add one quota', () => {
      it('should add a quota only once per unit', () => {
        cy.get<Context>('@ctx').then((ctx) => {
          addQuota(ctx, Unit.AuthenticatedRequests, true, 1);
          addQuota(ctx, Unit.AuthenticatedRequests, true, 1, undefined, undefined, undefined, false).then((res) => {
            expect(res.status).to.equal(409);
          });
        });
      });

      describe('add two quotas', () => {
        it('should add a quota for each unit', () => {
          cy.get<Context>('@ctx').then((ctx) => {
            addQuota(ctx, Unit.AuthenticatedRequests, true, 1);
            addQuota(ctx, Unit.ExecutionSeconds, true, 1);
          });
        });
      });
    });

    describe('edit', () => {
      describe('remove one quota', () => {
        beforeEach(() => {
          cy.get<Context>('@ctx').then((ctx) => {
            ensureQuotaIsAdded(ctx, Unit.AuthenticatedRequests, true, 1);
          });
        });
        it('should remove a quota only once per unit', () => {
          cy.get<Context>('@ctx').then((ctx) => {
            removeQuota(ctx, Unit.AuthenticatedRequests);
          });
          cy.get<Context>('@ctx').then((ctx) => {
            removeQuota(ctx, Unit.AuthenticatedRequests, false).then((res) => {
              expect(res.status).to.equal(404);
            });
          });
        });

        describe('remove two quotas', () => {
          beforeEach(() => {
            cy.get<Context>('@ctx').then((ctx) => {
              ensureQuotaIsAdded(ctx, Unit.AuthenticatedRequests, true, 1);
              ensureQuotaIsAdded(ctx, Unit.ExecutionSeconds, true, 1);
            });
          });
          it('should remove a quota for each unit', () => {
            cy.get<Context>('@ctx').then((ctx) => {
              removeQuota(ctx, Unit.AuthenticatedRequests);
              removeQuota(ctx, Unit.ExecutionSeconds);
            });
          });
        });
      });
    });
  });

  describe('usage', () => {
    beforeEach(() => {
      cy.get<Context>('@ctx')
        .then((ctx) => {
          return [
            `${ctx.api.oidcBaseURL}/userinfo`,
            `${ctx.api.authBaseURL}/users/me`,
            `${ctx.api.mgmtBaseURL}/iam`,
            `${ctx.api.adminBaseURL}/instances/me`,
            // `${api.assetsBaseURL}/instance/policy/label/icon`,
            `${ctx.api.oauthBaseURL}/keys`,
            `${ctx.api.samlBaseURL}/certificate`,
          ];
        })
        .as('authenticatedUrls');
    });

    describe('authenticated requests', () => {
      beforeEach(() => {
        cy.get<Array<string>>('@authenticatedUrls').then((urls) => {
          cy.get<Context>('@ctx').then((ctx) => {
            ensureQuotaIsAdded(ctx, Unit.AuthenticatedRequests, true, urls.length);
            cy.task('runSQL', `TRUNCATE logstore.access;`);
          });
        });
      });

      it('authenticated requests are limited', () => {
        cy.get<Array<string>>('@authenticatedUrls').then((urls) => {
          cy.get<Context>('@ctx').then((ctx) => {
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
              url: urls[0],
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
    });

    describe('notifications', () => {
      const callURL = `http://${Cypress.env('WEBHOOK_HANDLER_HOST')}:${Cypress.env('WEBHOOK_HANDLER_PORT')}/do_something`;

      beforeEach(() => cy.task('resetWebhookEvents'));

      const amount = 100
      const percent = 10
      const usage = 25

      describe('without repetition', () => {
        beforeEach(() => {
          cy.get<Context>('@ctx').then((ctx) => {
            ensureQuotaIsAdded(ctx, Unit.AuthenticatedRequests, false, amount, [
              {
                callUrl: callURL,
                percent: percent,
                repeat: false,
              },
            ]);
            cy.task('runSQL', `TRUNCATE logstore.access;`);
          });
        });

        it('fires once with the expected payload', () => {
          cy.get<Array<string>>('@authenticatedUrls').then((urls) => {
            cy.get<Context>('@ctx').then((ctx) => {
              for (let i = 0; i < usage; i++) {
                cy.request({
                  url: urls[0],
                  method: 'GET',
                  auth: {
                    bearer: ctx.api.token,
                  },
                });
              }
            });
            cy.waitUntil(() =>
              cy.task<Array<ZITADELWebhookEvent>>('handledWebhookEvents').then((events) => {
                if (events.length != 1) {
                  return false;
                }
                return Cypress._.matches(<ZITADELWebhookEvent>{
                  callURL: callURL,
                  threshold: percent,
                  unit: 1,
                  usage: percent,
                })(events[0]);
              }),
            );
          });
        });
      });

      describe('with repetition', () => {
        beforeEach(() => {
          cy.get<Array<string>>('@authenticatedUrls').then((urls) => {
            cy.get<Context>('@ctx').then((ctx) => {
              ensureQuotaIsAdded(ctx, Unit.AuthenticatedRequests, false, amount, [
                {
                  callUrl: callURL,
                  percent: percent,
                  repeat: true,
                },
              ]);
              cy.task('runSQL', `TRUNCATE logstore.access;`);
            });
          });
        });

        it('fires repeatedly with the expected payloads', () => {
          cy.get<Array<string>>('@authenticatedUrls').then((urls) => {
            cy.get<Context>('@ctx').then((ctx) => {
              for (let i = 0; i < usage; i++) {
                cy.request({
                  url: urls[0],
                  method: 'GET',
                  auth: {
                    bearer: ctx.api.token,
                  },
                });
              }
            });
          });
          cy.waitUntil(() =>
            cy.task<Array<ZITADELWebhookEvent>>('handledWebhookEvents').then((events) => {
              if (events.length != 1) {
                return false;
              }
              for (let i = 0; i < events.length; i++) {
                const threshold = percent * (i + 1)
                if (
                  !Cypress._.matches(<ZITADELWebhookEvent>{
                    callURL: callURL,
                    threshold: threshold,
                    unit: 1,
                    usage: threshold,
                  })(events[i])
                ) {
                  return false;
                }
              }
              return true;
            }),
          );
        });
      });
    });
  });
});
