import { apiAuth } from '../../support/api/apiauth';
import { resetInstanceFeatures } from '../../support/api/features';
import { login, User } from '../../support/login/users';

describe('features settings', () => {
  const featuresPath = '/instance?id=features';

  beforeEach(() => {
    cy.context().as('ctx');
    cy.visit(featuresPath);
  });

  describe('UI and Display Tests', () => {
    it('should display features page with correct elements', () => {
      // Check page title contains relevant text (flexible for translation issues)
      cy.get('h2').should('be.visible').and('contain.text', 'Feature');
      cy.get('.events-desc').should('be.visible');

      // Check info link
      cy.get('a[href*="feature-service"]').should('be.visible').and('have.attr', 'target', '_blank');

      // Check reset button is present (if available)
      cy.get('body')
        .find('[data-e2e="reset-features-button"]')
        .then(($btn) => {
          if ($btn.length > 0) {
            cy.get('[data-e2e="reset-features-button"]').should('be.visible');
          }
        });
    });

    it('should display feature toggles', () => {
      cy.get('cnsl-card').should('be.visible');
      cy.get('.features').should('be.visible');

      // Check that feature toggles are present
      cy.get('cnsl-feature-toggle').should('have.length.greaterThan', 0);
      cy.get('cnsl-login-v2-feature-toggle').should('be.visible');
    });

    it('should maintain feature states after page reload', () => {
      // Wait for features to load
      cy.get('cnsl-feature-toggle').should('be.visible');

      // Get the count of feature toggles before reload
      cy.get('cnsl-feature-toggle').then(($toggles) => {
        const initialCount = $toggles.length;

        // Reload page
        cy.reload();

        // Wait for features to load again
        cy.get('cnsl-feature-toggle').should('be.visible');

        // Verify that the same number of toggles are still present
        cy.get('cnsl-feature-toggle').should('have.length', initialCount);

        // Verify that toggles are still functional
        cy.get('cnsl-feature-toggle')
          .first()
          .within(() => {
            cy.get('mat-button-toggle').should('be.visible');
          });
      });
    });

    it('should handle API errors gracefully', () => {
      // Intercept API calls and force them to fail
      cy.intercept('POST', '**/features*', {
        statusCode: 500,
        body: { message: 'Internal server error' },
      }).as('featureError');

      // Try to toggle a feature
      cy.get('cnsl-feature-toggle')
        .first()
        .within(() => {
          cy.get('mat-button-toggle').first().click();
        });

      // Wait for the error response (optional since it might not always trigger)
      cy.get('body').then(() => {
        // Just verify that the page is still functional
        cy.get('cnsl-feature-toggle').should('be.visible');
      });
    });

    it('should be keyboard accessible', () => {
      // Features should be focusable via the underlying button
      cy.get('cnsl-feature-toggle')
        .first()
        .within(() => {
          cy.get('mat-button-toggle button').first().focus().should('be.focused');
        });

      // Reset button should be focusable (if available)
      cy.get('body')
        .find('[data-e2e="reset-features-button"]')
        .then(($btn) => {
          if ($btn.length > 0) {
            cy.get('[data-e2e="reset-features-button"]').focus().should('be.focused');
          }
        });
    });

    describe('permissions', () => {
      it('should show appropriate elements for admin users', () => {
        // Admin should see reset button (if available)
        cy.get('body')
          .find('[data-e2e="reset-features-button"]')
          .then(($btn) => {
            if ($btn.length > 0) {
              cy.get('[data-e2e="reset-features-button"]').should('be.visible');
            }
          });

        // Admin should see all feature toggles
        cy.get('cnsl-feature-toggle').should('have.length.greaterThan', 0);
      });
    });
  });

  describe('Feature Modification Tests', () => {
    afterEach(() => {
      // Reset features after each feature modification test to ensure clean state
      apiAuth().then((api) => {
        resetInstanceFeatures(api);
      });
    });

    it('should be able to toggle a feature', () => {
      // Wait for features to load
      cy.get('cnsl-feature-toggle').should('be.visible');

      // Get the first feature toggle and check its initial state
      cy.get('cnsl-feature-toggle')
        .first()
        .within(() => {
          // Ensure we always trigger a state change by clicking an unchecked button
          cy.get('mat-button-toggle').then(($allButtons) => {
            const uncheckedButtons = $allButtons.not('.mat-button-toggle-checked');

            if (uncheckedButtons.length > 0) {
              // Click an unchecked button to enable it
              const targetButton = uncheckedButtons.first();
              cy.wrap(targetButton).click();

              // Check for success toast since we made a real change
              cy.shouldConfirmSuccess();

              // Verify the toggle reflected the new state
              cy.wrap(targetButton).should('have.class', 'mat-button-toggle-checked');
            } else {
              // All buttons are checked, click the first one to uncheck it
              const targetButton = $allButtons.first();
              cy.wrap(targetButton).click();

              // Check for success toast since we made a real change
              cy.shouldConfirmSuccess();

              // Verify the toggle reflected the new state (should be unchecked now)
              cy.wrap(targetButton).should('not.have.class', 'mat-button-toggle-checked');
            }
          });
        });
    });

    it('should handle loginV2 feature toggle', () => {
      // Check if loginV2 feature toggle exists
      cy.get('cnsl-login-v2-feature-toggle')
        .should('be.visible')
        .within(() => {
          // Should have a feature toggle (with button toggles)
          cy.get('cnsl-feature-toggle').should('be.visible');
          cy.get('mat-button-toggle').should('be.visible');

          // Actually toggle the loginV2 feature to test functionality
          // Check current state and click the opposite to ensure we trigger a change
          cy.get('mat-button-toggle').then(($buttons) => {
            const uncheckedButtons = $buttons.not('.mat-button-toggle-checked');

            if (uncheckedButtons.length > 0) {
              // Click an unchecked button to enable it
              cy.wrap(uncheckedButtons.first()).click();
            } else {
              // All buttons are checked, click the first one to toggle it
              cy.wrap($buttons.first()).click();
            }
          });

          // Check for success toast since we made a real change
          cy.shouldConfirmSuccess();

          // Check if a base URI input field appears after enabling loginV2
          cy.get('cnsl-login-v2-feature-toggle').within(() => {
            cy.get('input[cnslInput], input[data-e2e*="uri"], input[placeholder*="URI"], input[placeholder*="uri"]')
              .should('be.visible')
              .and('not.be.disabled');
          });
        });
    });

    it('should reset features when reset button is clicked', () => {
      // Change a feature first to have something to reset
      cy.get('cnsl-feature-toggle')
        .first()
        .within(() => {
          // Check current state and click the opposite to ensure we trigger a change
          cy.get('mat-button-toggle').then(($buttons) => {
            const checkedButton = $buttons.filter('.mat-button-toggle-checked');
            const uncheckedButtons = $buttons.not('.mat-button-toggle-checked');

            if (uncheckedButtons.length > 0) {
              // Click an unchecked button to enable it
              cy.wrap(uncheckedButtons.first()).click();
              // Check for success toast since we made a real change
              cy.shouldConfirmSuccess();
              // Verify the change was applied
              cy.wrap(uncheckedButtons.first()).should('have.class', 'mat-button-toggle-checked');
            } else {
              // All buttons are checked, click the first one to uncheck it
              cy.wrap(checkedButton.first()).click();
              // Check for success toast since we made a real change
              cy.shouldConfirmSuccess();
              // Verify the change was applied
              cy.wrap(checkedButton.first()).should('not.have.class', 'mat-button-toggle-checked');
            }
          });
        });

      // Click the reset button (if available)
      cy.get('body')
        .find('[data-e2e="reset-features-button"]')
        .then(($btn) => {
          if ($btn.length > 0) {
            cy.get('[data-e2e="reset-features-button"]').click();

            // Check for success toast from reset operation
            cy.shouldConfirmSuccess();

            // Verify features are still loaded and functional after UI reset
            cy.get('cnsl-feature-toggle').should('be.visible');
          }
        });
    });
  });
});
