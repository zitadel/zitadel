import { apiAuth } from '../../../support/api/apiauth';
import { ensureOIDCSettingsSet } from '../../../support/api/oidc-settings';

const smtpPath = `/settings?id=smtpprovider`;

beforeEach(() => {
  cy.context().as('ctx');
});

describe('instance notifications', () => {
  describe('smtp settings', () => {
    it(`should show SMTP provider settings`, () => {
      cy.visit(smtpPath);
      cy.contains('SMTP Settings');
    });
  });
});
