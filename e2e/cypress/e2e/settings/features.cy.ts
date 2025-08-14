import { apiAuth } from '../../support/api/apiauth';
import { login, User } from '../../support/login/users';

describe('features settings', () => {
  const featuresPath = '/instance?id=features';

  beforeEach(() => {
    login(User.IAMAdminUser);
    cy.visit(featuresPath);
  });

  it('should display features page with correct elements', () => {
    // Check page title and description
    cy.get('h2').should('contain.text', 'Features');
    cy.get('.events-desc').should('be.visible');

    // Check info link
    cy.get('a[href*="feature-service"]').should('be.visible').and('have.attr', 'target', '_blank');

    // Check reset button is present
    cy.get('button').contains('Reset').should('be.visible').and('have.attr', 'color', 'warn');
  });

  it('should display feature toggles', () => {
    cy.get('cnsl-card').should('be.visible');
    cy.get('.features').should('be.visible');

    // Check that feature toggles are present
    cy.get('cnsl-feature-toggle').should('have.length.greaterThan', 0);
    cy.get('cnsl-login-v2-feature-toggle').should('be.visible');
  });

  it('should be able to toggle a feature', () => {
    // Wait for features to load
    cy.get('cnsl-feature-toggle').should('be.visible');

    // Get the first feature toggle and click it
    cy.get('cnsl-feature-toggle')
      .first()
      .within(() => {
        cy.get('mat-slide-toggle, mat-checkbox, button').first().click();
      });

    // Wait for the save operation
    cy.wait(1500);

    // Should show a toast message (success or error)
    cy.get('simple-snack-bar, .mat-snack-bar-container', { timeout: 5000 }).should('exist');
  });

  it('should handle loginV2 feature toggle', () => {
    cy.get('cnsl-login-v2-feature-toggle')
      .should('be.visible')
      .within(() => {
        // Should have a toggle
        cy.get('mat-slide-toggle, mat-checkbox').should('be.visible');

        // May have a base URI input field
        cy.get('body').then(() => {
          if (Cypress.$('input[type="text"]').length > 0) {
            cy.get('input[type="text"]').should('be.visible');
          }
        });
      });
  });

  it('should reset features when reset button is clicked', () => {
    // Click reset button
    cy.get('button').contains('Reset').click();

    // Wait for the reset operation
    cy.wait(1500);

    // Should show a toast message
    cy.get('simple-snack-bar, .mat-snack-bar-container', { timeout: 5000 }).should('exist');

    // Page should still show feature toggles
    cy.get('cnsl-feature-toggle').should('be.visible');
  });

  it('should maintain feature states after page reload', () => {
    // Wait for features to load
    cy.get('cnsl-feature-toggle').should('be.visible');

    // Store initial state of first toggle
    cy.get('cnsl-feature-toggle')
      .first()
      .within(() => {
        cy.get('mat-slide-toggle input, mat-checkbox input')
          .first()
          .then(($input) => {
            const initialState = $input.prop('checked');

            // Reload page
            cy.reload();

            // Wait for features to load again
            cy.get('cnsl-feature-toggle').should('be.visible');

            // Check that state is maintained
            cy.get('cnsl-feature-toggle')
              .first()
              .within(() => {
                cy.get('mat-slide-toggle input, mat-checkbox input').first().should('have.prop', 'checked', initialState);
              });
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
        cy.get('mat-slide-toggle, mat-checkbox, button').first().click();
      });

    // Wait for the error response
    cy.wait('@featureError');

    // Should show error toast
    cy.get('simple-snack-bar, .mat-snack-bar-container', { timeout: 5000 }).should('exist');
  });

  it('should be keyboard accessible', () => {
    // Features should be focusable
    cy.get('cnsl-feature-toggle')
      .first()
      .within(() => {
        cy.get('mat-slide-toggle, mat-checkbox, button').first().focus().should('be.focused');
      });

    // Reset button should be focusable
    cy.get('button').contains('Reset').focus().should('be.focused');
  });

  describe('specific feature tests', () => {
    const testableFeatures = ['userSchema', 'consoleUseV2UserApi', 'loginDefaultOrg', 'oidcTokenExchange'];

    testableFeatures.forEach((featureName) => {
      it(`should be able to toggle ${featureName} feature`, () => {
        // Look for the specific feature by its name/label
        cy.contains('cnsl-feature-toggle', featureName, { timeout: 10000 })
          .should('be.visible')
          .within(() => {
            cy.get('mat-slide-toggle, mat-checkbox').first().click();
          });

        // Wait for save operation
        cy.wait(1500);

        // Verify toast appears
        cy.get('simple-snack-bar, .mat-snack-bar-container', { timeout: 5000 }).should('exist');
      });
    });
  });

  describe('permissions', () => {
    it('should show appropriate elements for admin users', () => {
      // Admin should see reset button
      cy.get('button').contains('Reset').should('be.visible');

      // Admin should see all feature toggles
      cy.get('cnsl-feature-toggle').should('have.length.greaterThan', 0);
    });
  });
});
