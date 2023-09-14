import { apiAuth } from '../../../support/api/apiauth';
import { ensureOIDCSettingsSet } from '../../../support/api/oidc-settings';

beforeEach(() => {
  cy.context().as('ctx');
});

describe('oidc settings', () => {
  const smtpPath = `/settings?id=smtpprovider`;

  it(`should show SMTP provider settings`, () => {
    cy.visit(smtpPath);
    cy.contains('SMTP Settings');
  });
});
