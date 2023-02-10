import { Apps, ensureProjectExists, ensureProjectResourceDoesntExist } from '../../support/api/projects';
import { Context } from 'support/commands';

const testProjectName = 'e2eprojectapplication';
const testAppName = 'e2eappundertest';

describe('applications', () => {
  beforeEach(() => {
    cy.context()
      .as('ctx')
      .then((ctx) => {
        ensureProjectExists(ctx.api, testProjectName).as('projectId');
      });
  });

  describe('add app', () => {
    beforeEach(`ensure it doesn't exist already`, () => {
      cy.get<Context>('@ctx').then((ctx) => {
        cy.get<string>('@projectId').then((projectId) => {
          ensureProjectResourceDoesntExist(ctx.api, projectId, Apps, testAppName);
          cy.visit(`/projects/${projectId}`);
        });
      });
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
      cy.shouldConfirmSuccess();
      const expectClientId = new RegExp(`^.*[0-9]+\\@${testProjectName}.*$`);
      cy.get('[data-e2e="client-id-copy"]').click();
      cy.contains('[data-e2e="client-id"]', expectClientId);
      cy.clipboardMatches(expectClientId);
    });
  });

  describe('edit app', () => {
    it('should configure an application to enable dev mode');
    it('should configure an application to put user roles and info inside id token');
  });
});
