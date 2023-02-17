import { getInstance } from 'support/api/instances';
import { ensureQuotaIsRemoved, removeQuota, Unit } from 'support/api/quota';
import { Context } from 'support/commands';

describe('api internationalization', () => {
  beforeEach(() => {
    cy.context()
      .as('ctx')
      .then((ctx) => {
        ensureQuotaIsRemoved(ctx, Unit.ExecutionSeconds);
      });
  });
  it('instance not found error should be translated', () => {
    cy.get<Context>('@ctx').then((ctx) => {
      removeQuota(ctx, Unit.ExecutionSeconds, false).then((res) => {
        expect(res.body.message).to.contain('Quota not found for this unit');
      });
      getInstance(ctx.system, "this ID clearly doesn't exist", false).then((res) => {
        expect(res.body.message).to.contain('Instance not found');
      });
    });
  });
});
