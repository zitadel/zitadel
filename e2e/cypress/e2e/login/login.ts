import { Before, Given, Then, When } from '@badeball/cypress-cucumber-preprocessor';
import { apiAuth } from 'support/api/apiauth';
import { ensureLoginPolicy } from 'support/api/policies';
import { ensureApplicationExists, ensureProjectExists } from 'support/api/projects';
import { API } from 'support/api/types';
import { ensureHumanUserExists } from 'support/api/users';
import { loginname } from 'support/login/users';

const testUser = 'testuser';
const testProject = 'testProject';
const testApp = 'testApp';

// TODO: Make this dynamic somehow
const testAppClientID = '189165202854445059@testproject';

Before(() => {
  apiAuth().as('api');
});

Given('A user with password {string} and verified email does exist', (pw: string) => {
  cy.get('@api').then((api: unknown) => {
    ensureHumanUserExists(<API>api, loginname(testUser, Cypress.env('ORGANIZATION')), pw);
  });
});
Given('an application with redirect uri {string} exists', (redirectUri: string) => {
  cy.get('@api').then((api: unknown) => {
    ensureProjectExists(<API>api, testProject).then((projectId) => {
      ensureApplicationExists(<API>api, projectId, testApp, [redirectUri]);
    });
  });
});
Given('login policy has values {string}', (policy: string) => {
  cy.get('@api').then((api: unknown) => {
    // TODO: Make that work
    // ensureLoginPolicy(<API>api, JSON.parse(policy));
  });
});
Given('a clear browser session', Cypress.session.clearAllSavedSessions);
Given('user navigates to authorize endpoint with redirect uri {string}', (redirectUri: string) => {
  cy.visit({
    url: 'http://localhost:8080/oauth/v2/authorize',
    qs: {
      scope: 'openid',
      prompt: 'login',
      response_type: 'code',
      client_id: encodeURI(testAppClientID),
      redirect_uri: encodeURI(redirectUri),
    },
  });
});

When('user enters loginname', () => {
  cy.get('#loginName').type(loginname(testUser, Cypress.env('ORGANIZATION')));
  cy.get('#submit-button').click();
});
When('user enters password {string}', (pw: string) => {
  cy.get('#password').type(pw);
  cy.get('#submit-button').click();
});

Then('user is redirected to {string}', (expectedUrl: string) => {
  cy.url().then((actualUrl: string) => {
    expect(actualUrl.startsWith(expectedUrl)).to.be.true;
  });
});
