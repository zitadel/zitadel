import { addQuota, ensureQuotaIsAdded, ensureQuotaIsRemoved, removeQuota, Unit } from 'support/api/quota';
import { createHumanUser, ensureUserDoesntExist } from 'support/api/users';
import { Context } from 'support/commands';
import { ZITADELWebhookEvent } from 'support/types';
import { textChangeRangeIsUnchanged } from 'typescript';

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
            `${ctx.api.oauthBaseURL}/keys`,
            `${ctx.api.samlBaseURL}/certificate`,
          ];
        })
        .as('authenticatedUrls');
    });

    describe('authenticated requests', () => {
      const testUserName = 'shouldNotBeCreated';
      beforeEach(() => {
        cy.get<Array<string>>('@authenticatedUrls').then((urls) => {
          cy.get<Context>('@ctx').then((ctx) => {
            ensureUserDoesntExist(ctx.api, testUserName);
            ensureQuotaIsAdded(ctx, Unit.AuthenticatedRequests, true, urls.length);
            cy.task('runSQL', `TRUNCATE logstore.access;`);
          });
        });
      });

      it('only authenticated requests are limited', () => {
        cy.get<Array<string>>('@authenticatedUrls').then((urls) => {
          cy.get<Context>('@ctx').then((ctx) => {
            const start = new Date();
            urls.forEach((url) => {
              cy.request({
                url: url,
                method: 'GET',
                auth: {
                  bearer: ctx.api.token,
                },
              });
            });
            expectCookieDoesntExist();
            const expiresMax = new Date();
            expiresMax.setMinutes(expiresMax.getMinutes() + 20);
            cy.request({
              url: urls[1],
              method: 'GET',
              auth: {
                bearer: ctx.api.token,
              },
              failOnStatusCode: false,
            }).then((res) => {
              expect(res.status).to.equal(429);
            });
            cy.getCookie('zitadel.quota.limiting').then((cookie) => {
              expect(cookie.value).to.equal('true');
              const cookieExpiry = new Date();
              cookieExpiry.setTime(cookie.expiry * 1000);
              expect(cookieExpiry).to.be.within(start, expiresMax);
            });
            createHumanUser(ctx.api, testUserName, false).then((res) => {
              expect(res.status).to.equal(429);
            });
            // visit limited console
            // cy.visit('/users/me');
            // cy.contains('#authenticated-requests-exhausted-dialog button', 'Continue').click();
            // const upgradeInstancePage = `https://example.com/instances/${ctx.instanceId}`;
            // cy.origin(upgradeInstancePage, { args: { upgradeInstancePage } }, ({ upgradeInstancePage }) => {
            //   cy.location('href').should('equal', upgradeInstancePage);
            // });
            // upgrade instance
            ensureQuotaIsRemoved(ctx, Unit.AuthenticatedRequests);
            // visit upgraded console again
            cy.visit('/users/me');
            cy.get('[data-e2e="top-view-title"]');
            expectCookieDoesntExist();
            createHumanUser(ctx.api, testUserName);
            expectCookieDoesntExist();
          });
        });
      });
    });

    describe('notifications', () => {
      const callURL = `http://${Cypress.env('WEBHOOK_HANDLER_HOST')}:${Cypress.env('WEBHOOK_HANDLER_PORT')}/do_something`;

      beforeEach(() => cy.task('resetWebhookEvents'));

      const amount = 100;
      const percent = 10;
      const usage = 35;

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

        it('fires at least once with the expected payload', () => {
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
                if (events.length < 1) {
                  return false;
                }
                return Cypress._.matches(<ZITADELWebhookEvent>{
                  sentStatus: 200,
                  payload: {
                    callURL: callURL,
                    threshold: percent,
                    unit: 1,
                    usage: percent,
                  },
                })(events[0]);
              }),
            );
          });
        });

        it('fires until the webhook returns a successful message', () => {
          cy.task('failWebhookEvents', 8);
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
            cy.waitUntil(
              () =>
                cy.task<Array<ZITADELWebhookEvent>>('handledWebhookEvents').then((events) => {
                  if (events.length != 9) {
                    return false;
                  }
                  return events.reduce<boolean>((a, b, i) => {
                    return !a
                      ? a
                      : i < 8
                      ? Cypress._.matches(<ZITADELWebhookEvent>{
                          sentStatus: 500,
                          payload: {
                            callURL: callURL,
                            threshold: percent,
                            unit: 1,
                            usage: percent,
                          },
                        })(b)
                      : Cypress._.matches(<ZITADELWebhookEvent>{
                          sentStatus: 200,
                          payload: {
                            callURL: callURL,
                            threshold: percent,
                            unit: 1,
                            usage: percent,
                          },
                        })(b);
                  }, true);
                }),
              { timeout: 60_000 },
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
              let foundExpected = 0;
              for (let i = 0; i < events.length; i++) {
                for (let expect = 10; expect <= 30; expect += 10) {
                  if (
                    Cypress._.matches(<ZITADELWebhookEvent>{
                      sentStatus: 200,
                      payload: {
                        callURL: callURL,
                        threshold: expect,
                        unit: 1,
                        usage: expect,
                      },
                    })(events[i])
                  ) {
                    foundExpected++;
                  }
                }
              }
              return foundExpected >= 3;
            }),
          );
        });
      });
    });
  });
});

function expectCookieDoesntExist() {
  cy.getCookie('zitadel.quota.limiting').then((cookie) => {
    expect(cookie).to.be.null;
  });
}
