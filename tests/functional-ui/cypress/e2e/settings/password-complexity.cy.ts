import { login, User } from '../../support/login/users';

// password-complexity.cy.ts
//
// Tests for the Password Complexity Policy settings page, including the
// history_count field introduced by the password-reuse-prevention feature.
//
// Instance-scope policy: /instance?id=complexity  (requires IAM admin login)
// Org-scope policy:      /org path, click Modify on the Password Complexity card
//
// The history-count value uses a +/- stepper mirroring the min-length row.
// The current value is rendered in span.history-count-input. Increment and
// decrement buttons carry .history-count-increment and .history-count-decrement.

describe('password complexity', () => {
  const instanceComplexityPath = '/instance?id=complexity';

  const setHistoryCountTo = (target: number) => {
    cy.get('.history-count-input')
      .invoke('text')
      .then((text) => {
        const current = parseInt(text.trim(), 10) || 0;
        const delta = target - current;
        const button = delta > 0 ? '.history-count-increment' : '.history-count-decrement';
        for (let i = 0; i < Math.abs(delta); i++) {
          cy.get(button).click();
        }
      });
  };

  describe('instance-scope policy (history count)', () => {
    beforeEach(() => {
      login(User.IAMAdminUser);
      cy.visit(instanceComplexityPath);
      // Wait for the form to be visible (cnsl-card wraps the inputs)
      cy.get('cnsl-password-complexity-policy').should('be.visible');
    });

    it('should display the history-count input', () => {
      cy.get('.history-count-input').should('be.visible');
    });

    it('should save history_count=3 and persist after reload', () => {
      setHistoryCountTo(3);
      cy.contains('button', 'Save').click();
      cy.shouldConfirmSuccess();
      cy.reload();
      cy.get('cnsl-password-complexity-policy').should('be.visible');
      cy.get('.history-count-input').should('have.text', '3');
      // Restore to 0 so other tests start clean
      setHistoryCountTo(0);
      cy.contains('button', 'Save').click();
      cy.shouldConfirmSuccess();
    });

    it('should save history_count=0 and persist after reload', () => {
      // First set to a non-zero value so the test is meaningful
      setHistoryCountTo(5);
      cy.contains('button', 'Save').click();
      cy.shouldConfirmSuccess();
      // Now set back to 0
      setHistoryCountTo(0);
      cy.contains('button', 'Save').click();
      cy.shouldConfirmSuccess();
      cy.reload();
      cy.get('cnsl-password-complexity-policy').should('be.visible');
      cy.get('.history-count-input').should('have.text', '0');
    });
  });

  describe('org-scope policy (history count via policy card)', () => {
    const orgPath = '/org';

    [User.OrgOwner].forEach((user) => {
      describe(`as user "${user}"`, () => {
        beforeEach(() => {
          login(user);
          cy.visit(orgPath);
          cy.contains('[data-e2e="policy-card"]', 'Password Complexity')
            .contains('button', 'Modify')
            .click({ force: true });
          cy.get('cnsl-password-complexity-policy').should('be.visible');
        });

        it('should display the history-count input', () => {
          cy.get('.history-count-input').should('be.visible');
        });

        // The existing stubs (no-body it() calls) are preserved below.
        it(`should restrict passwords that don't have the minimal length`);
        it(`should require passwords to contain a number if option is switched on`);
        it(`should not require passwords to contain a number if option is switched off`);
        it(`should require passwords to contain a symbol if option is switched on`);
        it(`should not require passwords to contain a symbol if option is switched off`);
        it(`should require passwords to contain a lowercase letter if option is switched on`);
        it(`should not require passwords to contain a lowercase letter if option is switched off`);
        it(`should require passwords to contain an uppercase letter if option is switched on`);
        it(`should not require passwords to contain an uppercase letter if option is switched off`);
      });
    });
  });
});
