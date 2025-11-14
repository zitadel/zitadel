import { ensureExternalLinksSettingsSet } from 'support/api/external-links-settings';
import { apiAuth } from '../../support/api/apiauth';

describe('external link settings', () => {
  const tosLink = '';
  const privacyPolicyLink = '';
  const helpLink = '';
  const supportEmail = '';
  const customLink = '';
  const customLinkText = '';
  const docsLink = 'https://zitadel.com/docs';

  beforeEach(`reset`, () => {
    apiAuth().then((apiCallProperties) => {
      ensureExternalLinksSettingsSet(apiCallProperties, tosLink, privacyPolicyLink, docsLink);
    });
  });

  describe('instance', () => {
    beforeEach(`visit`, () => {
      cy.visit(`/instance?id=privacypolicy`);
    });

    it(`should have default settings`, () => {
      cy.get('[formcontrolname="tosLink"]').should('value', tosLink);
      cy.get('[formcontrolname="privacyLink"]').should('value', privacyPolicyLink);
      cy.get('[formcontrolname="helpLink"]').should('value', helpLink);
      cy.get('[formcontrolname="supportEmail"]').should('value', supportEmail);
      cy.get('[formcontrolname="customLink"]').should('value', customLink);
      cy.get('[formcontrolname="customLinkText"]').should('value', customLinkText);
      cy.get('[formcontrolname="docsLink"]').should('value', docsLink);
    });

    it(`should update external links`, () => {
      cy.get('[formcontrolname="tosLink"]').clear().type('tosLink2');
      cy.get('[formcontrolname="privacyLink"]').clear().type('privacyLink2');
      cy.get('[formcontrolname="helpLink"]').clear().type('helpLink');
      cy.get('[formcontrolname="supportEmail"]').clear().type('support@example.com');
      cy.get('[formcontrolname="customLink"]').clear().type('customLink');
      cy.get('[formcontrolname="customLinkText"]').clear().type('customLinkText');
      cy.get('[formcontrolname="docsLink"]').clear().type('docsLink');
      cy.get('[data-e2e="save-button"]').click();
      cy.shouldConfirmSuccess();
    });

    it(`should return to default values`, () => {
      cy.get('[formcontrolname="tosLink"]').should('value', tosLink);
      cy.get('[formcontrolname="privacyLink"]').should('value', privacyPolicyLink);
      cy.get('[formcontrolname="helpLink"]').should('value', helpLink);
      cy.get('[formcontrolname="supportEmail"]').should('value', supportEmail);
      cy.get('[formcontrolname="customLink"]').should('value', customLink);
      cy.get('[formcontrolname="customLinkText"]').should('value', customLinkText);
      cy.get('[formcontrolname="docsLink"]').should('value', docsLink);
    });
  });

  describe('org', () => {
    beforeEach(`visit`, () => {
      cy.visit(`/org-settings?id=privacypolicy`);
    });

    it(`should have default settings`, () => {
      cy.get('[formcontrolname="tosLink"]').should('value', tosLink);
      cy.get('[formcontrolname="privacyLink"]').should('value', privacyPolicyLink);
      cy.get('[formcontrolname="helpLink"]').should('value', helpLink);
      cy.get('[formcontrolname="supportEmail"]').should('value', supportEmail);
      cy.get('[formcontrolname="customLink"]').should('value', customLink);
      cy.get('[formcontrolname="customLinkText"]').should('value', customLinkText);
      cy.get('[formcontrolname="docsLink"]').should('value', docsLink);
    });

    it(`should update external links`, () => {
      cy.get('[formcontrolname="tosLink"]').clear().type('tosLink2');
      cy.get('[formcontrolname="privacyLink"]').clear().type('privacyLink2');
      cy.get('[formcontrolname="helpLink"]').clear().type('helpLink');
      cy.get('[formcontrolname="supportEmail"]').clear().type('support@example.com');
      cy.get('[formcontrolname="customLink"]').clear().type('customLink');
      cy.get('[formcontrolname="customLinkText"]').clear().type('customLinkText');
      cy.get('[formcontrolname="docsLink"]').clear().type('docsLink');
      cy.get('[data-e2e="save-button"]').click();
      cy.shouldConfirmSuccess();
    });

    it(`should return to default values`, () => {
      cy.get('[data-e2e="reset-button"]').click();
      cy.get('[data-e2e="confirm-dialog-button"]').click();
      cy.get('[formcontrolname="tosLink"]').should('value', tosLink);
      cy.get('[formcontrolname="privacyLink"]').should('value', privacyPolicyLink);
      cy.get('[formcontrolname="helpLink"]').should('value', helpLink);
      cy.get('[formcontrolname="supportEmail"]').should('value', supportEmail);
      cy.get('[formcontrolname="customLink"]').should('value', customLink);
      cy.get('[formcontrolname="customLinkText"]').should('value', customLinkText);
      cy.get('[formcontrolname="docsLink"]').should('value', docsLink);
    });
  });
});
