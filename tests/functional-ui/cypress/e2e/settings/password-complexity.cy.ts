import { login, User } from '../../support/login/users';

// password-complexity.cy.ts
//
// Tests for the Password Complexity Policy settings page, including the
// history_count field introduced by the password-reuse-prevention feature.
//
// Instance-scope policy: /instance?id=complexity  (requires IAM admin login)
// Org-scope policy:      /org path, click Modify on the Password Complexity card
//
// The history-count input uses [(ngModel)]="historyCount" with class="history-count-input"
// (no formControlName). We target it by class.

describe('password complexity', () => {
  const instanceComplexityPath = '/instance?id=complexity';

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
      cy.get('.history-count-input').clear().type('3');
      cy.contains('button', 'Save').click();
      cy.shouldConfirmSuccess();
      cy.reload();
      cy.get('cnsl-password-complexity-policy').should('be.visible');
      cy.get('.history-count-input').should('have.value', '3');
      // Restore to 0 so other tests start clean
      cy.get('.history-count-input').clear().type('0');
      cy.contains('button', 'Save').click();
      cy.shouldConfirmSuccess();
    });

    it('should save history_count=0 and persist after reload', () => {
      // First set to a non-zero value so the test is meaningful
      cy.get('.history-count-input').clear().type('5');
      cy.contains('button', 'Save').click();
      cy.shouldConfirmSuccess();
      // Now set back to 0
      cy.get('.history-count-input').clear().type('0');
      cy.contains('button', 'Save').click();
      cy.shouldConfirmSuccess();
      cy.reload();
      cy.get('cnsl-password-complexity-policy').should('be.visible');
      cy.get('.history-count-input').should('have.value', '0');
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
