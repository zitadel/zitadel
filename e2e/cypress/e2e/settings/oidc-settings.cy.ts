import { newTarget } from 'support/api/target';
import { ensureOIDCSettings } from '../../support/api/oidc-settings';

// TODO: As these are instance level settings,
// we should set a deterministic state before each test
describe('oidc settings', () => {
  const oidcSettingsPath = `/settings?id=oidc`;
  const accessTokenPrecondition = 1;
  const idTokenPrecondition = 2;
  const refreshTokenExpirationPrecondition = 7;
  const refreshTokenIdleExpirationPrecondition = 2;

  beforeEach(`ensure they are set`, () => {

    newTarget('e2eoidcsettings').then((target) => {
      ensureOIDCSettings(
        target,
        accessTokenPrecondition,
        idTokenPrecondition,
        refreshTokenExpirationPrecondition,
        refreshTokenIdleExpirationPrecondition,
      );
      cy.visit(oidcSettingsPath);
    });
  });

  it(`should update oidc settings`, () => {
    cy.get('[formcontrolname="accessTokenLifetime"]').should('value', accessTokenPrecondition).clear().type('2');
    cy.get('[formcontrolname="idTokenLifetime"]').should('value', idTokenPrecondition).clear().type('24');
    cy.get('[formcontrolname="refreshTokenExpiration"]')
      .should('value', refreshTokenExpirationPrecondition)
      .clear()
      .type('30');
    cy.get('[formcontrolname="refreshTokenIdleExpiration"]')
      .should('value', refreshTokenIdleExpirationPrecondition)
      .clear()
      .type('7');
    cy.get('[data-e2e="save-button"]').click();
    cy.shouldConfirmSuccess();
  });
});
