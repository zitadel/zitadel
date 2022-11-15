import { Apps, ensureProjectExists, ensureProjectResourceDoesntExist } from '../../support/api/projects';
import { apiAuth } from '../../support/api/apiauth';

describe('applications', () => {
  const testProjectName = 'e2eprojectapplication';
  const testAppName = 'e2eappundertest';

  beforeEach(() => {
    apiAuth()
      .as('api')
      .then((api) => {
        ensureProjectExists(api, testProjectName).as('projectId');
      });
  });

  describe('add app', function () {
    beforeEach(`ensure it doesn't exist already`, function () {
      ensureProjectResourceDoesntExist(this.api, this.projectId, Apps, testAppName);
      cy.visit(`/projects/${this.projectId}`);
    });

    it('add app', () => {
      cy.get('[data-e2e="app-card-add"]').should('be.visible').click();
      cy.get('[formcontrolname="name"]').focus().type(testAppName);
      cy.get('[for="WEB"]').click();
      cy.get('[data-e2e="continue-button-nameandtype"]').click();
      cy.get('[for="PKCE"]').should('be.visible').click();
      cy.get('[data-e2e="continue-button-authmethod"]').click();
      cy.get('[data-e2e="redirect-uris"] input').focus().type('http://localhost:3000/api/auth/callback/zitadel');
      cy.get('[data-e2e="postlogout-uris"] input').focus().type('http://localhost:3000');
      cy.get('[data-e2e="continue-button-redirecturis"]').click();
      cy.get('[data-e2e="create-button"]').click();
      cy.get('[id*=overlay]').should('exist');
      cy.get('.data-e2e-success');
      const expectClientId = new RegExp(`^.*[0-9]+\\@${testProjectName}.*$`);
      cy.get('[data-e2e="client-id-copy"]').click();
      cy.contains('[data-e2e="client-id"]', expectClientId);
      cy.clipboardMatches(expectClientId);
      cy.shouldNotExist({ selector: '.data-e2e-failure' });
    });
  });

  describe('edit app', () => {
    it('should configure an application to enable dev mode');
    it('should configure an application to put user roles and info inside id token');
  });
});
