import { ensureOIDCAppDoesntExist } from 'support/api/oidc-applications';
import { ensureProjectExists } from 'support/api/projects';
import { newTarget } from '../../support/api/target';

describe('applications', () => {
  const testProjectName = 'e2eprojectapplication';
  const testAppName = 'e2eappundertest';

  describe('add app', () => {
    beforeEach(`ensure it doesn't exist already`, () => {
      newTarget('e2eapplications').then((target) => {
        ensureProjectExists(target, testProjectName).then((projectId) => {
          ensureOIDCAppDoesntExist(target, projectId, testAppName);
          cy.visit(`/projects/${projectId}?org=${target.orgId}`);
        });
      });
    });

    it('add app', () => {
      cy.get('[data-e2e="app-card-add"]').should('be.visible').click();
      cy.get('[formcontrolname="name"]').focus().should('be.visible').type(testAppName);
      cy.get('[for="WEB"]').should('be.visible').click();
      cy.get('[data-e2e="continue-button-nameandtype"]').should('be.visible').click();
      cy.get('[for="PKCE"]').should('be.visible').click();
      cy.get('[data-e2e="continue-button-authmethod"]').should('be.visible').click();
      cy.get('[data-e2e="redirect-uris"] input')
        .focus()
        .should('be.visible')
        .type('http://localhost:3000/api/auth/callback/zitadel');
      cy.get('[data-e2e="postlogout-uris"] input').focus().should('be.visible').type('http://localhost:3000');
      cy.get('[data-e2e="continue-button-redirecturis"]').should('be.visible').click();
      cy.get('[data-e2e="create-button"]').should('be.visible').click();
      cy.get('[id*=overlay]').should('exist');
      cy.shouldConfirmSuccess();
      const expectClientId = new RegExp(`^.*[0-9]+\\@${testProjectName}.*$`);
      cy.get('[data-e2e="client-id-copy"]').should('be.visible').click();
      cy.contains('[data-e2e="client-id"]', expectClientId);
      cy.clipboardMatches(expectClientId);
    });

    describe('edit app', () => {
      it('should configure an application to enable dev mode');
      it('should configure an application to put user roles and info inside id token');
    });
  });
});
