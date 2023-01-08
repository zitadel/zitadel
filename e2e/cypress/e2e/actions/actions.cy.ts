import { ensureActionExists, setTriggerTypes } from 'support/api/actions';
import { apiAuth } from 'support/api/apiauth';
import { ensureProjectExists, ensureRoleExists } from 'support/api/projects';
import { API } from 'support/api/types';
import { ensureHumanUserExists, ensureUserDoesntExist, setMetadata } from 'support/api/users';
import { login, loginAsPredefinedUser, User } from 'support/login/users';

describe('actions', () => {
  const emailVerifiedScript = 'e2eSetEmailVerified',
    addGrantScript = 'e2eAddGrant',
    addMetadataScript = 'e2eAddMetadata',
    storeUsernameScript = 'e2eSetLastUsernameMD',
    projectName = 'e2eaction',
    roleKey = 'e2eactionrole',
    userFirstname = 'e2eFirstname',
    userLastname = 'e2eLastname',
    userEmail = 'e2e@zitadelaction.com',
    userPw = 'Password1!',
    specPath = 'cypress/e2e/actions';

  const loginUrl = `${Cypress.env('BACKEND_URL')}/ui/login`;

  beforeEach(() => {
    apiAuth().as('api');
  });

  describe('triggers', () => {
    describe('creation', () => {
      beforeEach(() => {
        cy.get<API>('@api').then((api) => {
          ensureUserDoesntExist(api, userEmail);
          cy.readFile(`${specPath}/${emailVerifiedScript}.js`).then((script) => {
            ensureActionExists(api, emailVerifiedScript, script).as('emailVerifiedId');
          });
        });
      });

      describe('pre', () => {
        beforeEach(() => {
          cy.get<API>('@api').then((api) => {
            cy.readFile(`${specPath}/${addMetadataScript}.js`).then((script) => {
              ensureActionExists(api, addMetadataScript, script).as('metadataId');
            });
          });
        });

        describe('internal', () => {
          beforeEach(() => {
            cy.get<API>('@api').then((api) => {
              cy.get<number>('@emailVerifiedId').then((emailVerifiedId) => {
                cy.get<number>('@metadataId').then((metadataId) => {
                  setTriggerTypes(api, 3, 2, [emailVerifiedId, metadataId]);
                });
              });
            });
          });

          it(`shouldn't prompt for email code and add metadata`, () => {
            register();
            cy.get('[data-e2e="email-is-verified"]');
            cy.get('[data-e2e="sidenav-element-metadata"]').click();
            cy.contains('tr', 'akey').contains('avalue');
          });
        });

        describe('external', () => {});
      });

      describe('post', () => {
        beforeEach(() => {
          cy.get<API>('@api').then((api) => {
            ensureProjectExists(api, projectName).then((projectId) => {
              ensureRoleExists(api, projectId, roleKey);
              cy.readFile<string>(`${specPath}/${addGrantScript}.js`).then((script) => {
                ensureActionExists(
                  api,
                  addGrantScript,
                  script.replace('<PROJECT_ID>', `${projectId}`).replace('<ROLE_KEY>', roleKey),
                ).as('addGrantId');
              });
            });
          });
        });

        describe('internal', () => {
          beforeEach(() => {
            cy.get<API>('@api').then((api) => {
              cy.get<number>('@emailVerifiedId').then((emailVerifiedId) => {
                cy.get<number>('@addGrantId').then((addGrantId) => {
                  setTriggerTypes(api, 3, 2, [emailVerifiedId]);
                  setTriggerTypes(api, 3, 3, [addGrantId]);
                });
              });
            });
          });
          it(`should add a grant when registering via UI`, () => {
            register();
            cy.get('[data-e2e="sidenav-element-grants"]').click();
            cy.contains('tr', roleKey);
          });
        });
        describe('external', () => {});
      });
    });

    describe('post authentication', () => {
      beforeEach(() => {
        cy.get<API>('@api').then((api) => {
          ensureHumanUserExists(api, userEmail, userPw, true).as('userId');
          cy.readFile<string>(`${specPath}/${storeUsernameScript}.js`).then((script) => {
            ensureActionExists(api, storeUsernameScript, script).as('storeUsernameId');
          });
        });
      });

      describe('internal', () => {
        beforeEach(() => {
          cy.get<API>('@api').then((api) => {
            cy.get<number>('@storeUsernameId').then((storeUsernameId) => {
              setTriggerTypes(api, 3, 1, [storeUsernameId]);
            });
          });
        });

        it('should store auth request states in metadata', () => {
          Cypress.session.clearAllSavedSessions();
          login(<User>userEmail, userPw);
          loginAsPredefinedUser(User.IAMAdminUser);
          cy.get('@userId').then((userId) => {
            cy.visit(`/users/${userId}?id=metadata`);
          });
          cy.contains('tr', 'last username used').contains(userEmail);
        });
      });
      describe('external', () => {});
    });
  });

  function register() {
    cy.intercept(
      {
        method: 'GET',
        url: `${Cypress.env('BACKEND_URL')}/oauth/v2/authorize*`,
      },
      (req) => {
        req.query['login_hint'] = userEmail;
        req.query['prompt'] = 'create';
      },
    ).as('authreq');
    cy.visit(loginUrl);
    cy.wait('@authreq');

    cy.get('#firstname').type(userFirstname);
    cy.get('#lastname').type(userLastname);
    cy.get('#register-password').type(userPw);
    cy.get('#register-password-confirmation').type(userPw);
    cy.get('#register-term-confirmation').check({ force: true });
    cy.get('#register-term-confirmation-privacy').check({ force: true });
    cy.get('form').submit();
    cy.get('#password').type(userPw);

    cy.intercept({
      method: 'POST',
      url: `${loginUrl}/password*`,
      times: 1,
    }).as('password');
    cy.get('form').submit();

    cy.wait('@password').then((interception) => {
      if (interception.response.body.indexOf('/ui/login/mfa/prompt') === -1) {
        return;
      }

      cy.contains('button', 'skip').click();
    });

    cy.contains('[data-e2e="top-view-title"]', `${userFirstname} ${userLastname}`);
    loginAsPredefinedUser(User.IAMAdminUser);
    cy.visit('/users');
    cy.contains('tr', userEmail).click();
  }
});
