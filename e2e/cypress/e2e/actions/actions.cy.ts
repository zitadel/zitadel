import { ensureActionExists, setTriggerTypes } from 'support/api/actions';
import { apiAuth } from 'support/api/apiauth';
import { ensureDomainPolicy } from 'support/api/policies';
import { ensureProjectExists, ensureRoleExists } from 'support/api/projects';
import { API } from 'support/api/types';
import { ensureHumanUserExists, ensureUserDoesntExist } from 'support/api/users';
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
    preCreationEmail = 'precre@action.com',
    postCreationEmail = 'postcre@action.com',
    postAuthEmail = 'postauth@action.com',
    userPw = 'Password1!',
    specPath = 'cypress/e2e/actions';

  const loginUrl = `/ui/login`;

  beforeEach(() => {
    apiAuth()
      .as('api')
      .then((api) => {
        ensureDomainPolicy(api, false, false, false);
        cy.readFile(`${specPath}/${emailVerifiedScript}.js`).then((script) => {
          ensureActionExists(api, emailVerifiedScript, script)
            .as('emailVerifiedId')
            .then((emailVerifiedId) => {
              setTriggerTypes(api, 3, 2, [emailVerifiedId]);
            });
        });
      });
  });

  describe('triggers', () => {
    describe('creation', () => {
      describe('pre', () => {
        beforeEach(() => {
          cy.get<API>('@api').then((api) => {
            ensureUserDoesntExist(api, preCreationEmail);
            cy.get<number>('@emailVerifiedId').then((emailVerifiedId) => {
              cy.readFile(`${specPath}/${addMetadataScript}.js`).then((script) => {
                ensureActionExists(api, addMetadataScript, script).then((metadataId) => {
                  setTriggerTypes(api, 3, 2, [emailVerifiedId, metadataId]);
                });
              });
            });
          });
        });

        it(`shouldn't prompt for email code and add metadata`, () => {
          register(preCreationEmail);
          loginAsPredefinedUser(User.IAMAdminUser);
          cy.get('@registeredUserId').then((userId) => {
            cy.visit(`/users/${userId}?id=metadata`);
          });
          cy.contains('tr', 'akey').contains('avalue');
        });
      });

      describe('post', () => {
        beforeEach(() => {
          cy.get<API>('@api').then((api) => {
            cy.get<number>('@emailVerifiedId').then((emailVerifiedId) => {
              setTriggerTypes(api, 3, 2, [emailVerifiedId]);
            });
            ensureUserDoesntExist(api, postCreationEmail);
            ensureProjectExists(api, projectName).then((projectId) => {
              ensureRoleExists(api, projectId, roleKey);
              cy.readFile<string>(`${specPath}/${addGrantScript}.js`).then((script) => {
                ensureActionExists(
                  api,
                  addGrantScript,
                  script.replace('<PROJECT_ID>', `${projectId}`).replace('<ROLE_KEY>', roleKey),
                ).then((addGrantId) => {
                    setTriggerTypes(api, 3, 3, [addGrantId]);
                });
              });
            });
          });
        });

        it(`should add a grant when registering via UI`, () => {
          register(postCreationEmail);
          loginAsPredefinedUser(User.IAMAdminUser);
          cy.get('@registeredUserId').then((userId) => {
            cy.visit(`/users/${userId}?id=grants`);
          });
          cy.contains('tr', roleKey);
        });
      });
    });

    describe('post authentication', () => {
      beforeEach(() => {
        cy.get<API>('@api').then((api) => {
          ensureUserDoesntExist(api, postAuthEmail).as('userId');
          ensureHumanUserExists(api, postAuthEmail, true).as('userId');
          cy.readFile<string>(`${specPath}/${storeUsernameScript}.js`).then((script) => {
            ensureActionExists(api, storeUsernameScript, script).then((storeUsernameId) => {
              setTriggerTypes(api, 3, 1, [storeUsernameId]);
            });
          });
        });
      });

      it('should store the username in metadata after password authentication', () => {
        cy.get('@userId').then((userId) => {
          cy.log('user exists', userId);
        });
        login(postAuthEmail, userPw);
        loginAsPredefinedUser(User.IAMAdminUser);
        cy.get('@userId').then((userId) => {
          cy.visit(`/users/${userId}?id=metadata`);
        });
        cy.contains('tr', 'last username used').contains(postAuthEmail);
      });

      it('should store the username in metadata after mutlifactor authentication', () => {
        login(postAuthEmail, userPw);
        cy.visit('/users/me?id=security');
        cy.get('[data-e2e="add-factor"]').click();
        cy.get('[data-e2e="add-factor-otp"]').should('be.visible').click();
        cy.get('[data-e2e="otp-secret"]')
          .as('otpSecret')
          .then((secret) => {
            cy.task<string>('generateOTP', secret.text().trim()).then((token) => {
              cy.get('[data-e2e="otp-code-input"]').type(token, { force: true });
              cy.get('[data-e2e="save-otp-factor"]').click();
            });
          });
        login(postAuthEmail, userPw, () => {
          cy.task<string>('generateOTP').then((token) => {
            cy.get('#code').type(token);
            cy.get('#submit-button').click();
          });
        });
        loginAsPredefinedUser(User.IAMAdminUser);
        cy.get('@userId').then((userId) => {
          cy.visit(`/users/${userId}?id=metadata`);
        });
        cy.contains('tr', 'last username used').contains(postAuthEmail);
      });
    });
  });

  function register(email: string) {
    // We want to have a clean session but miss cypresses sesssion cache
    cy.session(Math.random().toString(), () => {
      cy.intercept(
        {
          method: 'GET',
          url: `${Cypress.env('BACKEND_URL')}/oauth/v2/authorize*`,
        },
        (req) => {
          req.query['prompt'] = 'create';
          req.query['login_hint'] = email;
          req.continue();
        },
      ).as('regAuthReq');
      cy.visit(loginUrl);
      cy.wait('@regAuthReq');
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

      cy.get('[data-e2e="user-id"]')
        .then(($el) => {
          return $el.text().trim();
        })
        .as('registeredUserId');
    });
  }
});
